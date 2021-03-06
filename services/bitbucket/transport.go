package bitbucket

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/estafette/estafette-ci-api/clients/bitbucketapi"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

type Handler struct {
	service Service
}

func (h *Handler) Handle(c *gin.Context) {

	// https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html

	eventType := c.GetHeader("X-Event-Key")
	// h.prometheusInboundEventTotals.With(prometheus.Labels{"event": eventType, "source": "bitbucket"}).Inc()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Reading body from Bitbucket webhook failed")
		c.String(http.StatusInternalServerError, "Reading body from Bitbucket webhook failed")
		return
	}

	// unmarshal json body to check if installation is whitelisted
	var anyEvent bitbucketapi.AnyEvent
	err = json.Unmarshal(body, &anyEvent)
	if err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to BitbucketAnyEvent failed")
		return
	}

	// verify owner is whitelisted
	isWhitelisted, _ := h.service.IsWhitelistedOwner(anyEvent.Repository)
	if !isWhitelisted {
		c.Status(http.StatusUnauthorized)
		return
	}

	switch eventType {
	case "repo:push":

		// unmarshal json body
		var pushEvent bitbucketapi.RepositoryPushEvent
		err := json.Unmarshal(body, &pushEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to BitbucketRepositoryPushEvent failed")
			return
		}

		err = h.service.CreateJobForBitbucketPush(c.Request.Context(), pushEvent)

		if err != nil && !errors.Is(err, ErrNonCloneableEvent) && !errors.Is(err, ErrNoManifest) {
			c.String(http.StatusInternalServerError, "Oops, something went wrong!")
			return
		}
	case
		"repo:fork",
		"repo:transfer",
		"repo:created",
		"repo:commit_comment_created",
		"repo:commit_status_created",
		"repo:commit_status_updated",
		"issue:created",
		"issue:updated",
		"issue:comment_created",
		"pullrequest:created",
		"pullrequest:updated",
		"pullrequest:approved",
		"pullrequest:unapproved",
		"pullrequest:fulfilled",
		"pullrequest:rejected",
		"pullrequest:comment_created",
		"pullrequest:comment_updated",
		"pullrequest:comment_deleted":

	case "repo:updated":
		log.Debug().Str("event", eventType).Str("requestBody", string(body)).Msgf("Bitbucket webhook event of type '%v', logging request body", eventType)

		// unmarshal json body
		var repoUpdatedEvent bitbucketapi.RepoUpdatedEvent
		err := json.Unmarshal(body, &repoUpdatedEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to BitbucketRepoUpdatedEvent failed")
			return
		}

		if repoUpdatedEvent.IsValidRenameEvent() {
			log.Info().Msgf("Renaming repository from %v/%v/%v to %v/%v/%v", repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetOldRepoOwner(), repoUpdatedEvent.GetOldRepoName(), repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetNewRepoOwner(), repoUpdatedEvent.GetNewRepoName())
			err = h.service.Rename(c.Request.Context(), repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetOldRepoOwner(), repoUpdatedEvent.GetOldRepoName(), repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetNewRepoOwner(), repoUpdatedEvent.GetNewRepoName())
			if err != nil {
				log.Error().Err(err).Msgf("Failed renaming repository from %v/%v/%v to %v/%v/%v", repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetOldRepoOwner(), repoUpdatedEvent.GetOldRepoName(), repoUpdatedEvent.GetRepoSource(), repoUpdatedEvent.GetNewRepoOwner(), repoUpdatedEvent.GetNewRepoName())
				return
			}
		}

	case "repo:deleted":
		log.Debug().Str("event", eventType).Str("requestBody", string(body)).Msgf("Bitbucket webhook event of type '%v', logging request body", eventType)
		// unmarshal json body
		var repoDeletedEvent bitbucketapi.RepoDeletedEvent
		err := json.Unmarshal(body, &repoDeletedEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to BitbucketRepoDeletedEvent failed")
			return
		}

		log.Info().Msgf("Archiving repository %v/%v/%v", repoDeletedEvent.GetRepoSource(), repoDeletedEvent.GetRepoOwner(), repoDeletedEvent.GetRepoName())
		err = h.service.Archive(c.Request.Context(), repoDeletedEvent.GetRepoSource(), repoDeletedEvent.GetRepoOwner(), repoDeletedEvent.GetRepoName())
		if err != nil {
			log.Error().Err(err).Msgf("Failed archiving repository %v/%v/%v", repoDeletedEvent.GetRepoSource(), repoDeletedEvent.GetRepoOwner(), repoDeletedEvent.GetRepoName())
			return
		}

	default:
		log.Warn().Str("event", eventType).Msgf("Unsupported Bitbucket webhook event of type '%v'", eventType)
	}

	c.String(http.StatusOK, "Aye aye!")
}
