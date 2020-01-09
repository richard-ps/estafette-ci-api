package github

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/estafette/estafette-ci-api/clients/githubapi"
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

	// https://developer.github.com/webhooks/
	eventType := c.GetHeader("X-Github-Event")
	// h.prometheusInboundEventTotals.With(prometheus.Labels{"event": eventType, "source": "github"}).Inc()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Reading body from Github webhook failed")
		c.String(http.StatusInternalServerError, "Reading body from Github webhook failed")
		return
	}

	// verify hmac signature
	hasValidSignature, err := h.service.HasValidSignature(c.Request.Context(), body, c.GetHeader("X-Hub-Signature"))
	if err != nil {
		log.Error().Err(err).Msg("Verifying signature from Github webhook failed")
		c.String(http.StatusInternalServerError, "Verifying signature from Github webhook failed")
		return
	}
	if !hasValidSignature {
		log.Warn().Msg("Signature from Github webhook is invalid")
		c.String(http.StatusBadRequest, "Signature from Github webhook is invalid")
		return
	}

	// unmarshal json body to check if installation is whitelisted
	var anyEvent githubapi.AnyEvent
	err = json.Unmarshal(body, &anyEvent)
	if err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to GithubAnyEvent failed")
		return
	}

	// verify installation id is whitelisted
	if !h.service.IsWhitelistedInstallation(c.Request.Context(), anyEvent.Installation) {
		c.Status(http.StatusUnauthorized)
		return
	}

	switch eventType {
	case "push": // Any Git push to a Repository, including editing tags or branches. Commits via API actions that update references are also counted. This is the default event.

		// unmarshal json body
		var pushEvent githubapi.PushEvent
		err := json.Unmarshal(body, &pushEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to GithubPushEvent failed")
			return
		}

		err = h.service.CreateJobForGithubPush(c.Request.Context(), pushEvent)

		if err != nil && !errors.Is(err, ErrNonCloneableEvent) && !errors.Is(err, ErrNoManifest) {
			c.String(http.StatusInternalServerError, "Oops, something went wrong!")
			return
		}

	case
		"commit_comment",                        // Any time a Commit is commented on.
		"create",                                // Any time a Branch or Tag is created.
		"delete",                                // Any time a Branch or Tag is deleted.
		"deployment",                            // Any time a Repository has a new deployment created from the API.
		"deployment_status",                     // Any time a deployment for a Repository has a status update from the API.
		"fork",                                  // Any time a Repository is forked.
		"gollum",                                // Any time a Wiki page is updated.
		"installation",                          // Any time a GitHub App is installed or uninstalled.
		"installation_repositories",             // Any time a repository is added or removed from an installation.
		"issue_comment",                         // Any time a comment on an issue is created, edited, or deleted.
		"issues",                                // Any time an Issue is assigned, unassigned, labeled, unlabeled, opened, edited, milestoned, demilestoned, closed, or reopened.
		"label",                                 // Any time a Label is created, edited, or deleted.
		"marketplace_purchase",                  // Any time a user purchases, cancels, or changes their GitHub Marketplace plan.
		"member",                                // Any time a User is added or removed as a collaborator to a Repository, or has their permissions modified.
		"membership",                            // Any time a User is added or removed from a team. Organization hooks only.
		"milestone",                             // Any time a Milestone is created, closed, opened, edited, or deleted.
		"organization",                          // Any time a user is added, removed, or invited to an Organization. Organization hooks only.
		"org_block",                             // Any time an organization blocks or unblocks a user. Organization hooks only.
		"page_build",                            // Any time a Pages site is built or results in a failed build.
		"project_card",                          // Any time a Project Card is created, edited, moved, converted to an issue, or deleted.
		"project_column",                        // Any time a Project Column is created, edited, moved, or deleted.
		"project",                               // Any time a Project is created, edited, closed, reopened, or deleted.
		"public",                                // Any time a Repository changes from private to public.
		"pull_request_review_comment",           // Any time a comment on a pull request's unified diff is created, edited, or deleted (in the Files Changed tab).
		"pull_request_review",                   // Any time a pull request review is submitted, edited, or dismissed.
		"pull_request",                          // Any time a pull request is assigned, unassigned, labeled, unlabeled, opened, edited, closed, reopened, or synchronized (updated due to a new push in the branch that the pull request is tracking). Also any time a pull request review is requested, or a review request is removed.
		"release",                               // Any time a Release is published in a Repository.
		"status",                                // Any time a Repository has a status update from the API
		"team",                                  // Any time a team is created, deleted, modified, or added to or removed from a repository. Organization hooks only
		"team_add",                              // Any time a team is added or modified on a Repository.
		"watch",                                 // Any time a User stars a Repository.
		"integration_installation_repositories": // ?

	case "repository": // Any time a Repository is created, deleted (organization hooks only), made public, or made private.
		log.Debug().Str("event", eventType).Str("requestBody", string(body)).Msgf("Github webhook event of type '%v', logging request body", eventType)

		// unmarshal json body
		var repositoryEvent githubapi.RepositoryEvent
		err := json.Unmarshal(body, &repositoryEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to GithubRepositoryEvent failed")
			return
		}

		if repositoryEvent.IsValidRenameEvent() {
			log.Info().Msgf("Renaming repository from %v/%v/%v to %v/%v/%v", repositoryEvent.GetRepoSource(), repositoryEvent.GetOldRepoOwner(), repositoryEvent.GetOldRepoName(), repositoryEvent.GetRepoSource(), repositoryEvent.GetNewRepoOwner(), repositoryEvent.GetNewRepoName())
			err = h.service.Rename(c.Request.Context(), repositoryEvent.GetRepoSource(), repositoryEvent.GetOldRepoOwner(), repositoryEvent.GetOldRepoName(), repositoryEvent.GetRepoSource(), repositoryEvent.GetNewRepoOwner(), repositoryEvent.GetNewRepoName())
			if err != nil {
				log.Error().Err(err).Msgf("Failed renaming repository from %v/%v/%v to %v/%v/%v", repositoryEvent.GetRepoSource(), repositoryEvent.GetOldRepoOwner(), repositoryEvent.GetOldRepoName(), repositoryEvent.GetRepoSource(), repositoryEvent.GetNewRepoOwner(), repositoryEvent.GetNewRepoName())
				return
			}
		}

	default:
		log.Warn().Str("event", eventType).Msgf("Unsupported Github webhook event of type '%v'", eventType)
	}

	c.String(http.StatusOK, "Aye aye!")
}
