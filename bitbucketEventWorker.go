package main

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// BitbucketWorker processes events pushed to channels
type BitbucketWorker interface {
	ListenToBitbucketPushEventChannel()
	Stop()
	CreateJobForBitbucketPush(BitbucketRepositoryPushEvent)
}

type bitbucketWorkerImpl struct {
	WaitGroup   *sync.WaitGroup
	QuitChannel chan bool
}

func newBitbucketWorker(waitGroup *sync.WaitGroup) BitbucketWorker {
	return &bitbucketWorkerImpl{
		WaitGroup:   waitGroup,
		QuitChannel: make(chan bool)}
}

func (w *bitbucketWorkerImpl) ListenToBitbucketPushEventChannel() {
	go func() {
		// handle github push events via channels
		log.Debug().Msg("Listening to Bitbucket push events channel...")
		for {
			select {
			case pushEvent := <-bitbucketPushEvents:
				w.WaitGroup.Add(1)
				w.CreateJobForBitbucketPush(pushEvent)
				w.WaitGroup.Done()
			case <-w.QuitChannel:
				log.Info().Msg("Stopping Bitbucket worker...")
				return
			}
		}
	}()
}

func (w *bitbucketWorkerImpl) Stop() {
	go func() {
		w.QuitChannel <- true
	}()
}

func (w *bitbucketWorkerImpl) CreateJobForBitbucketPush(pushEvent BitbucketRepositoryPushEvent) {

	// check to see that it's a cloneable event
	if len(pushEvent.Push.Changes) == 0 || pushEvent.Push.Changes[0].New == nil || pushEvent.Push.Changes[0].New.Type != "branch" || len(pushEvent.Push.Changes[0].New.Target.Hash) == 0 {
		return
	}

	// get authenticated url for the repository
	bbClient := newBitbucketAPIClient(*bitbucketAPIKey, *bitbucketAppOAuthKey, *bitbucketAppOAuthSecret)
	authenticatedRepositoryURL, accessToken, err := bbClient.GetAuthenticatedRepositoryURL(pushEvent.Repository.Links.HTML.Href)
	if err != nil {
		log.Error().Err(err).
			Msg("Retrieving authenticated repository failed")
		return
	}

	// create ci builder client
	ciBuilderClient, err := newCiBuilderClient()
	if err != nil {
		log.Error().Err(err).Msg("Initializing ci builder client failed")
		return
	}

	// define ci builder params
	ciBuilderParams := CiBuilderParams{
		RepoFullName: pushEvent.Repository.FullName,
		RepoURL:      authenticatedRepositoryURL,
		RepoBranch:   pushEvent.Push.Changes[0].New.Name,
		RepoRevision: pushEvent.Push.Changes[0].New.Target.Hash,
		EnvironmentVariables: map[string]string{"ESTAFETTE_BITBUCKET_API_TOKEN": accessToken.AccessToken},
	}

	// create ci builder job
	_, err = ciBuilderClient.CreateCiBuilderJob(ciBuilderParams)
	if err != nil {
		log.Error().Err(err).
			Str("fullname", ciBuilderParams.RepoFullName).
			Str("url", ciBuilderParams.RepoURL).
			Str("branch", ciBuilderParams.RepoBranch).
			Str("revision", ciBuilderParams.RepoRevision).
			Msgf("Created estafette-ci-builder job for Bitbucket repository %v revision %v failed", ciBuilderParams.RepoFullName, ciBuilderParams.RepoRevision)

		return
	}

	log.Debug().
		Str("fullname", ciBuilderParams.RepoFullName).
		Str("url", ciBuilderParams.RepoURL).
		Str("branch", ciBuilderParams.RepoBranch).
		Str("revision", ciBuilderParams.RepoRevision).
		Msgf("Created estafette-ci-builder job for Bitbucket repository %v revision %v", ciBuilderParams.RepoFullName, ciBuilderParams.RepoRevision)
}
