package estafette

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/estafette/estafette-ci-api/api"
	"github.com/estafette/estafette-ci-api/clients/builderapi"
	"github.com/estafette/estafette-ci-api/clients/cloudstorage"
	"github.com/estafette/estafette-ci-api/clients/cockroachdb"
	contracts "github.com/estafette/estafette-ci-contracts"
	crypt "github.com/estafette/estafette-ci-crypt"
	manifest "github.com/estafette/estafette-ci-manifest"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"
)

// NewHandler returns a new estafette.Handler
func NewHandler(configFilePath string, config *api.APIConfig, encryptedConfig *api.APIConfig, cockroachDBClient cockroachdb.Client, cloudStorageClient cloudstorage.Client, ciBuilderClient builderapi.Client, buildService Service, warningHelper api.WarningHelper, secretHelper crypt.SecretHelper, githubJobVarsFunc func(context.Context, string, string, string) (string, string, error), bitbucketJobVarsFunc func(context.Context, string, string, string) (string, string, error), cloudsourceJobVarsFunc func(context.Context, string, string, string) (string, string, error)) Handler {

	return Handler{
		configFilePath:         configFilePath,
		config:                 config,
		encryptedConfig:        encryptedConfig,
		cockroachDBClient:      cockroachDBClient,
		cloudStorageClient:     cloudStorageClient,
		ciBuilderClient:        ciBuilderClient,
		buildService:           buildService,
		warningHelper:          warningHelper,
		secretHelper:           secretHelper,
		githubJobVarsFunc:      githubJobVarsFunc,
		bitbucketJobVarsFunc:   bitbucketJobVarsFunc,
		cloudsourceJobVarsFunc: cloudsourceJobVarsFunc,
	}
}

type Handler struct {
	configFilePath         string
	config                 *api.APIConfig
	encryptedConfig        *api.APIConfig
	cockroachDBClient      cockroachdb.Client
	cloudStorageClient     cloudstorage.Client
	ciBuilderClient        builderapi.Client
	buildService           Service
	warningHelper          api.WarningHelper
	secretHelper           crypt.SecretHelper
	githubJobVarsFunc      func(context.Context, string, string, string) (string, string, error)
	bitbucketJobVarsFunc   func(context.Context, string, string, string) (string, string, error)
	cloudsourceJobVarsFunc func(context.Context, string, string, string) (string, string, error)
}

func (h *Handler) GetPipelines(c *gin.Context) {

	pageNumber, pageSize, filters, sortings := api.GetQueryParameters(c)

	// filter on organizations / groups
	filters = api.SetPermissionsFilters(c, filters)

	response, err := api.GetPagedListResponse(
		func() ([]interface{}, error) {
			pipelines, err := h.cockroachDBClient.GetPipelines(c.Request.Context(), pageNumber, pageSize, filters, sortings, true)
			if err != nil {
				return nil, err
			}

			// convert typed array to interface array O(n)
			items := make([]interface{}, len(pipelines))
			for i := range pipelines {
				items[i] = pipelines[i]
			}

			return items, nil
		},
		func() (int, error) {
			return h.cockroachDBClient.GetPipelinesCount(c.Request.Context(), filters)
		},
		pageNumber,
		pageSize)

	if err != nil {
		log.Error().Err(err).Msg("Failed retrieving pipelines or count from db")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPipeline(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	filters := api.GetPipelineFilters(c)

	pipeline, err := h.cockroachDBClient.GetPipeline(c.Request.Context(), source, owner, repo, filters, true)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving pipeline for %v/%v/%v from db", source, owner, repo)
	}
	if pipeline == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

func (h *Handler) GetPipelineRecentBuilds(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	builds, err := h.cockroachDBClient.GetPipelineRecentBuilds(c.Request.Context(), source, owner, repo, true)
	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving recent builds for %v/%v/%v from db", source, owner, repo)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, builds)
}

func (h *Handler) GetPipelineBuilds(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	pageNumber, pageSize, filters, sortings := api.GetQueryParameters(c)

	response, err := api.GetPagedListResponse(
		func() ([]interface{}, error) {
			builds, err := h.cockroachDBClient.GetPipelineBuilds(c.Request.Context(), source, owner, repo, pageNumber, pageSize, filters, sortings, true)
			if err != nil {
				return nil, err
			}

			// convert typed array to interface array O(n)
			items := make([]interface{}, len(builds))
			for i := range builds {
				items[i] = builds[i]
			}

			return items, nil
		},
		func() (int, error) {
			return h.cockroachDBClient.GetPipelineBuildsCount(c.Request.Context(), source, owner, repo, filters)
		},
		pageNumber,
		pageSize)

	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving builds for %v/%v/%v from db", source, owner, repo)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPipelineBuild(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	revisionOrID := c.Param("revisionOrId")

	if len(revisionOrID) == 40 {

		build, err := h.cockroachDBClient.GetPipelineBuild(c.Request.Context(), source, owner, repo, revisionOrID, false)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, revisionOrID)
		}
		if build == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline not found"})
			return
		}

		c.JSON(http.StatusOK, build)
		return
	}

	id, err := strconv.Atoi(revisionOrID)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/builds/%v", source, owner, repo, revisionOrID)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	build, err := h.cockroachDBClient.GetPipelineBuildByID(c.Request.Context(), source, owner, repo, id, false)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, id)
	}
	if build == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline build not found"})
		return
	}

	// obfuscate all secrets
	build.Manifest, err = h.obfuscateSecrets(build.Manifest)
	if err != nil {
		log.Error().Err(err).Msgf("Failed obfuscating secrets")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}
	build.ManifestWithDefaults, err = h.obfuscateSecrets(build.ManifestWithDefaults)
	if err != nil {
		log.Error().Err(err).Msgf("Failed obfuscating secrets")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, build)
}

func (h *Handler) CreatePipelineBuild(c *gin.Context) {

	if !api.RequestTokenIsValid(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid"})
		return
	}

	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)

	var buildCommand contracts.Build
	err := c.BindJSON(&buildCommand)
	if err != nil {
		errorMessage := fmt.Sprint("Binding CreatePipelineBuild body failed")
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// match source, owner, repo with values in binded release
	if buildCommand.RepoSource != c.Param("source") {
		errorMessage := fmt.Sprintf("RepoSource in path and post data do not match for pipeline %v/%v/%v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, email)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if buildCommand.RepoOwner != c.Param("owner") {
		errorMessage := fmt.Sprintf("RepoOwner in path and post data do not match for pipeline %v/%v/%v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, email)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if buildCommand.RepoName != c.Param("repo") {
		errorMessage := fmt.Sprintf("RepoName in path and post data do not match for pipeline %v/%v/%v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, email)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// check if version exists and is valid to re-run
	failedBuilds, err := h.cockroachDBClient.GetPipelineBuildsByVersion(c.Request.Context(), buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, []string{"failed", "canceled"}, 1, false)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build %v/%v/%v version %v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, email)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	nonFailedBuilds, err := h.cockroachDBClient.GetPipelineBuildsByVersion(c.Request.Context(), buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, []string{"succeeded", "running", "pending", "canceling"}, 1, false)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build %v/%v/%v version %v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, email)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	var failedBuild *contracts.Build
	// get first failed build
	if len(failedBuilds) > 0 {
		failedBuild = failedBuilds[0]
	}

	// ensure there's no succeeded or running builds
	hasNonFailedBuilds := len(nonFailedBuilds) > 0

	if failedBuild == nil {
		errorMessage := fmt.Sprintf("No failed build %v/%v/%v version %v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, email)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if hasNonFailedBuilds {
		errorMessage := fmt.Sprintf("Version %v of pipeline %v/%v/%v has builds that are succeeded or running ; only if all builds are failed the pipeline can be re-run", buildCommand.BuildVersion, buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// set trigger event to manual
	failedBuild.Events = []manifest.EstafetteEvent{
		manifest.EstafetteEvent{
			Manual: &manifest.EstafetteManualEvent{
				UserID: email,
			},
		},
	}

	// hand off to build service
	createdBuild, err := h.buildService.CreateBuild(c.Request.Context(), *failedBuild, false)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed creating build %v/%v/%v version %v for build command issued by %v", buildCommand.RepoSource, buildCommand.RepoOwner, buildCommand.RepoName, buildCommand.BuildVersion, email)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, createdBuild)
}

func (h *Handler) CancelPipelineBuild(c *gin.Context) {

	if !api.RequestTokenIsValid(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid"})
		return
	}

	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	revisionOrID := c.Param("revisionOrId")

	id, err := strconv.Atoi(revisionOrID)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/builds/%v", source, owner, repo, revisionOrID)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	// retrieve build
	build, err := h.cockroachDBClient.GetPipelineBuildByID(c.Request.Context(), source, owner, repo, id, false)
	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, revisionOrID)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Retrieving pipeline build failed"})
		return
	}
	if build == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline build not found"})
		return
	}
	if build.BuildStatus == "canceling" {
		// apparently cancel was already clicked, but somehow the job didn't update the status to canceled
		jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "build", build.RepoOwner, build.RepoName, build.ID)
		h.ciBuilderClient.CancelCiBuilderJob(c.Request.Context(), jobName)
		h.cockroachDBClient.UpdateBuildStatus(c.Request.Context(), build.RepoSource, build.RepoOwner, build.RepoName, id, "canceled")
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Canceled build by user %v", email)})
		return
	}

	if build.BuildStatus != "pending" && build.BuildStatus != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("Build with status %v cannot be canceled", build.BuildStatus)})
		return
	}

	// this build can be canceled, set status 'canceling' and cancel the build job
	jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "build", build.RepoOwner, build.RepoName, build.ID)
	cancelErr := h.ciBuilderClient.CancelCiBuilderJob(c.Request.Context(), jobName)
	buildStatus := "canceling"
	if build.BuildStatus == "pending" {
		// job might not have created a builder yet, so set status to canceled straightaway
		buildStatus = "canceled"
	}
	err = h.cockroachDBClient.UpdateBuildStatus(c.Request.Context(), build.RepoSource, build.RepoOwner, build.RepoName, id, buildStatus)
	if err != nil {
		log.Error().Err(err).Msgf("Failed updating build status for %v/%v/%v/builds/%v in db", source, owner, repo, revisionOrID)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed setting pipeline build status to canceling"})
		return
	}

	// canceling the job failed because it no longer existed we should set canceled status right after having set it to canceling
	if cancelErr != nil && build.BuildStatus == "running" {
		buildStatus = "canceled"
		err = h.cockroachDBClient.UpdateBuildStatus(c.Request.Context(), build.RepoSource, build.RepoOwner, build.RepoName, id, buildStatus)
		if err != nil {
			log.Error().Err(err).Msgf("Failed updating build status to canceled after setting it to canceling for %v/%v/%v/builds/%v in db", source, owner, repo, revisionOrID)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed setting pipeline build status to canceled"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Canceled build by user %v", email)})
}

func (h *Handler) GetPipelineBuildLogs(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	revisionOrID := c.Param("revisionOrId")

	var build *contracts.Build
	var err error
	if len(revisionOrID) == 40 {

		build, err = h.cockroachDBClient.GetPipelineBuild(c.Request.Context(), source, owner, repo, revisionOrID, false)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, revisionOrID)
		}
	} else {

		id, err := strconv.Atoi(revisionOrID)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed reading id from path parameter for %v/%v/%v/builds/%v", source, owner, repo, revisionOrID)
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
			return
		}

		build, err = h.cockroachDBClient.GetPipelineBuildByID(c.Request.Context(), source, owner, repo, id, false)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, id)
		}
	}

	if build == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline build not found"})
		return
	}

	buildLog, err := h.cockroachDBClient.GetPipelineBuildLogs(c.Request.Context(), source, owner, repo, build.RepoBranch, build.RepoRevision, build.ID, h.config.APIServer.ReadLogFromDatabase())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving build logs for %v/%v/%v/builds/%v/logs from db", source, owner, repo, revisionOrID)
	}
	if buildLog == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline build log not found"})
		return
	}

	if h.config.APIServer.ReadLogFromCloudStorage() {
		err := h.cloudStorageClient.GetPipelineBuildLogs(c.Request.Context(), *buildLog, strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip"), c.Writer)
		if err != nil {

			if errors.Is(err, cloudstorage.ErrLogNotExist) {
				log.Warn().Err(err).
					Msgf("Failed retrieving build logs for %v/%v/%v/builds/%v/logs from cloud storage", source, owner, repo, revisionOrID)
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound)})
			}

			log.Error().Err(err).
				Msgf("Failed retrieving build logs for %v/%v/%v/builds/%v/logs from cloud storage", source, owner, repo, revisionOrID)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
			return
		}
		c.Writer.Flush()
		return
	}

	c.JSON(http.StatusOK, buildLog.Steps)
}

func (h *Handler) TailPipelineBuildLogs(c *gin.Context) {

	owner := c.Param("owner")
	repo := c.Param("repo")
	id := c.Param("revisionOrId")

	jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "build", owner, repo, id)

	logChannel := make(chan contracts.TailLogLine, 50)

	go h.ciBuilderClient.TailCiBuilderJobLogs(c.Request.Context(), jobName, logChannel)

	ticker := time.NewTicker(5 * time.Second)

	// ensure openresty doesn't buffer this response but sends the chunks rightaway
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		select {
		case ll, ok := <-logChannel:
			if !ok {
				c.SSEvent("close", true)
				return false
			}
			c.SSEvent("log", ll)
		case <-ticker.C:
			c.SSEvent("ping", true)
		}
		return true
	})
}

func (h *Handler) PostPipelineBuildLogs(c *gin.Context) {

	// ensure the request has the correct claims
	claims := jwt.ExtractClaims(c)
	job := claims["job"].(string)
	if job == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid or has invalid claim"})
		return
	}

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	revisionOrID := c.Param("revisionOrId")

	var buildLog contracts.BuildLog
	err := c.Bind(&buildLog)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed binding v2 logs for %v/%v/%v/%v", source, owner, repo, revisionOrID)
		c.String(http.StatusInternalServerError, "Oops, something went wrong")
		return
	}

	if len(revisionOrID) != 40 {

		_, err := strconv.Atoi(revisionOrID)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed reading id from path parameter for %v/%v/%v/builds/%v", source, owner, repo, revisionOrID)
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
			return
		}

		buildLog.BuildID = revisionOrID
	}

	insertedBuildLog, err := h.cockroachDBClient.InsertBuildLog(c.Request.Context(), buildLog, h.config.APIServer.WriteLogToDatabase())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed inserting logs for %v/%v/%v/%v", source, owner, repo, revisionOrID)
		c.String(http.StatusInternalServerError, "Oops, something went wrong")
		return
	}

	if h.config.APIServer.WriteLogToCloudStorage() {
		err = h.cloudStorageClient.InsertBuildLog(c.Request.Context(), insertedBuildLog)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed inserting logs into cloudstorage for %v/%v/%v/%v", source, owner, repo, revisionOrID)
		}
	}

	c.String(http.StatusOK, "Aye aye!")
}

func (h *Handler) GetPipelineBuildWarnings(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	revisionOrID := c.Param("revisionOrId")

	id, err := strconv.Atoi(revisionOrID)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/builds/%v", source, owner, repo, revisionOrID)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	build, err := h.cockroachDBClient.GetPipelineBuildByID(c.Request.Context(), source, owner, repo, id, false)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving build for %v/%v/%v/builds/%v from db", source, owner, repo, id)
	}
	if build == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline build not found"})
		return
	}

	warnings, err := h.warningHelper.GetManifestWarnings(build.ManifestObject, build.GetFullRepoPath())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed getting warnings for %v/%v/%v/builds/%v manifest", source, owner, repo, revisionOrID)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed getting warnings for manifest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"warnings": warnings})
}

func (h *Handler) GetPipelineReleases(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	pageNumber, pageSize, filters, sortings := api.GetQueryParameters(c)

	response, err := api.GetPagedListResponse(
		func() ([]interface{}, error) {
			releases, err := h.cockroachDBClient.GetPipelineReleases(c.Request.Context(), source, owner, repo, pageNumber, pageSize, filters, sortings)
			if err != nil {
				return nil, err
			}

			// convert typed array to interface array O(n)
			items := make([]interface{}, len(releases))
			for i := range releases {
				items[i] = releases[i]
			}

			return items, nil
		},
		func() (int, error) {
			return h.cockroachDBClient.GetPipelineReleasesCount(c.Request.Context(), source, owner, repo, filters)
		},
		pageNumber,
		pageSize)

	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving releases for %v/%v/%v from db", source, owner, repo)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreatePipelineRelease(c *gin.Context) {

	if !api.RequestTokenIsValid(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid"})
		return
	}

	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)

	var releaseCommand contracts.Release
	err := c.BindJSON(&releaseCommand)
	if err != nil {
		errorMessage := fmt.Sprint("Binding CreatePipelineRelease body failed")
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// match source, owner, repo with values in binded release
	if releaseCommand.RepoSource != c.Param("source") {
		errorMessage := fmt.Sprintf("RepoSource in path and post data do not match for pipeline %v/%v/%v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if releaseCommand.RepoOwner != c.Param("owner") {
		errorMessage := fmt.Sprintf("RepoOwner in path and post data do not match for pipeline %v/%v/%v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if releaseCommand.RepoName != c.Param("repo") {
		errorMessage := fmt.Sprintf("RepoName in path and post data do not match for pipeline %v/%v/%v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	filters := api.GetPipelineFilters(c)

	pipeline, err := h.cockroachDBClient.GetPipeline(c.Request.Context(), releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, filters, false)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving pipeline %v/%v/%v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}
	if pipeline == nil {
		errorMessage := fmt.Sprintf("No pipeline %v/%v/%v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// check if version exists and is valid to release
	builds, err := h.cockroachDBClient.GetPipelineBuildsByVersion(c.Request.Context(), releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.ReleaseVersion, []string{"succeeded"}, 1, false)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build %v/%v/%v version %v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.ReleaseVersion)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	var build *contracts.Build
	// get first build
	if len(builds) > 0 {
		build = builds[0]
	}

	if build == nil {
		errorMessage := fmt.Sprintf("No build %v/%v/%v version %v for release command", releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.ReleaseVersion)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}
	if build.BuildStatus != "succeeded" {
		errorMessage := fmt.Sprintf("Build %v for pipeline %v/%v/%v has status %v for release command; only succeeded pipelines are allowed to be released", releaseCommand.ReleaseVersion, releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, build.BuildStatus)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// check if release target exists
	releaseTargetExists := false
	actionExists := false
	for _, releaseTarget := range build.ReleaseTargets {
		if releaseTarget.Name == releaseCommand.Name {
			if len(releaseTarget.Actions) == 0 && releaseCommand.Action == "" {
				actionExists = true
			} else if len(releaseTarget.Actions) > 0 {
				for _, a := range releaseTarget.Actions {
					if a.Name == releaseCommand.Action {
						actionExists = true
						break
					}
				}
			}

			releaseTargetExists = true
			break
		}
	}
	if !releaseTargetExists {
		errorMessage := fmt.Sprintf("Build %v for pipeline %v/%v/%v has no release %v for release command", releaseCommand.ReleaseVersion, releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.Name)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// check if action is defined
	if !actionExists {
		errorMessage := fmt.Sprintf("Build %v for pipeline %v/%v/%v has no action %v for release action", releaseCommand.ReleaseVersion, releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.Action)
		log.Error().Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	// create release object and hand off to build service
	createdRelease, err := h.buildService.CreateRelease(c.Request.Context(), contracts.Release{
		Name:           releaseCommand.Name,
		Action:         releaseCommand.Action,
		RepoSource:     releaseCommand.RepoSource,
		RepoOwner:      releaseCommand.RepoOwner,
		RepoName:       releaseCommand.RepoName,
		ReleaseVersion: releaseCommand.ReleaseVersion,
		Groups:         build.Groups,
		Organizations:  build.Organizations,

		// set trigger event to manual
		Events: []manifest.EstafetteEvent{
			{
				Manual: &manifest.EstafetteManualEvent{
					UserID: email,
				},
			},
		},
	}, *build.ManifestObject, build.RepoBranch, build.RepoRevision, true)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed creating release %v for pipeline %v/%v/%v version %v for release command issued by %v", releaseCommand.Name, releaseCommand.RepoSource, releaseCommand.RepoOwner, releaseCommand.RepoName, releaseCommand.ReleaseVersion, email)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, createdRelease)
}

func (h *Handler) CancelPipelineRelease(c *gin.Context) {

	if !api.RequestTokenIsValid(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid"})
		return
	}

	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	idValue := c.Param("id")

	id, err := strconv.Atoi(idValue)
	if err != nil {
		log.Error().Err(err).Msgf("Failed reading id from path parameter for %v/%v/%v/%v", source, owner, repo, idValue)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	release, err := h.cockroachDBClient.GetPipelineRelease(c.Request.Context(), source, owner, repo, id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving release for %v/%v/%v/%v from db", source, owner, repo, id)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Retrieving pipeline release failed"})
		return
	}
	if release == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline release not found"})
		return
	}
	if release.ReleaseStatus == "canceling" {
		jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "release", release.RepoOwner, release.RepoName, release.ID)
		h.ciBuilderClient.CancelCiBuilderJob(c.Request.Context(), jobName)
		h.cockroachDBClient.UpdateReleaseStatus(c.Request.Context(), release.RepoSource, release.RepoOwner, release.RepoName, id, "canceled")
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Canceled release by user %v", email)})
		return
	}
	if release.ReleaseStatus != "pending" && release.ReleaseStatus != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("Release with status %v cannot be canceled", release.ReleaseStatus)})
		return
	}

	// this release can be canceled, set status 'canceling' and cancel the release job
	jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "release", release.RepoOwner, release.RepoName, release.ID)
	cancelErr := h.ciBuilderClient.CancelCiBuilderJob(c.Request.Context(), jobName)
	releaseStatus := "canceling"
	if release.ReleaseStatus == "pending" {
		// job might not have created a builder yet, so set status to canceled straightaway
		releaseStatus = "canceled"
	}
	err = h.cockroachDBClient.UpdateReleaseStatus(c.Request.Context(), release.RepoSource, release.RepoOwner, release.RepoName, id, releaseStatus)
	if err != nil {
		log.Error().Err(err).Msgf("Failed updating release status for %v/%v/%v/builds/%v in db", source, owner, repo, id)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed setting pipeline release status to canceling"})
		return
	}

	// canceling the job failed because it no longer existed we should set canceled status right after having set it to canceling
	if cancelErr != nil && release.ReleaseStatus == "running" {
		releaseStatus = "canceled"
		err = h.cockroachDBClient.UpdateReleaseStatus(c.Request.Context(), release.RepoSource, release.RepoOwner, release.RepoName, id, releaseStatus)
		if err != nil {
			log.Error().Err(err).Msgf("Failed updating release status to canceled after setting it to canceling for %v/%v/%v/builds/%v in db", source, owner, repo, id)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed setting pipeline release status to canceled"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Canceled release by user %v", email)})
}

func (h *Handler) GetPipelineRelease(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	idValue := c.Param("id")

	id, err := strconv.Atoi(idValue)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/%v", source, owner, repo, idValue)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	release, err := h.cockroachDBClient.GetPipelineRelease(c.Request.Context(), source, owner, repo, id)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving release for %v/%v/%v/%v from db", source, owner, repo, id)
	}
	if release == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline release not found"})
		return
	}

	c.JSON(http.StatusOK, release)
}

func (h *Handler) GetPipelineReleaseLogs(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	idValue := c.Param("id")

	id, err := strconv.Atoi(idValue)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/%v", source, owner, repo, idValue)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	releaseLog, err := h.cockroachDBClient.GetPipelineReleaseLogs(c.Request.Context(), source, owner, repo, id, h.config.APIServer.ReadLogFromDatabase())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed retrieving release logs for %v/%v/%v/%v from db", source, owner, repo, id)
	}
	if releaseLog == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline release log not found"})
		return
	}

	if h.config.APIServer.ReadLogFromCloudStorage() {
		err := h.cloudStorageClient.GetPipelineReleaseLogs(c.Request.Context(), *releaseLog, strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip"), c.Writer)
		if err != nil {

			if errors.Is(err, cloudstorage.ErrLogNotExist) {
				log.Warn().Err(err).
					Msgf("Failed retrieving release logs for %v/%v/%v/%v from cloud storage", source, owner, repo, id)
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound)})
			}

			log.Error().Err(err).
				Msgf("Failed retrieving release logs for %v/%v/%v/%v from cloud storage", source, owner, repo, id)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
			return
		}
		c.Writer.Flush()
		return
	}

	c.JSON(http.StatusOK, releaseLog.Steps)
}

func (h *Handler) TailPipelineReleaseLogs(c *gin.Context) {

	owner := c.Param("owner")
	repo := c.Param("repo")
	id := c.Param("id")

	jobName := h.ciBuilderClient.GetJobName(c.Request.Context(), "release", owner, repo, id)

	logChannel := make(chan contracts.TailLogLine, 50)

	go h.ciBuilderClient.TailCiBuilderJobLogs(c.Request.Context(), jobName, logChannel)

	ticker := time.NewTicker(5 * time.Second)

	// ensure openresty doesn't buffer this response but sends the chunks rightaway
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		select {
		case ll, ok := <-logChannel:
			if !ok {
				c.SSEvent("close", true)
				return false
			}
			c.SSEvent("log", ll)
		case <-ticker.C:
			c.SSEvent("ping", true)
		}
		return true
	})
}

func (h *Handler) PostPipelineReleaseLogs(c *gin.Context) {

	// ensure the request has the correct claims
	claims := jwt.ExtractClaims(c)
	job := claims["job"].(string)
	if job == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid or has invalid claim"})
		return
	}

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")
	idValue := c.Param("id")

	id, err := strconv.Atoi(idValue)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed reading id from path parameter for %v/%v/%v/%v", source, owner, repo, idValue)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": "Path parameter id is not of type integer"})
		return
	}

	var releaseLog contracts.ReleaseLog
	err = c.Bind(&releaseLog)
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed binding release logs for %v/%v/%v/%v", source, owner, repo, id)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "Failed binding release logs from body"})
		return
	}

	insertedReleaseLog, err := h.cockroachDBClient.InsertReleaseLog(c.Request.Context(), releaseLog, h.config.APIServer.WriteLogToDatabase())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed inserting release logs for %v/%v/%v/%v", source, owner, repo, id)
		c.String(http.StatusInternalServerError, "Oops, something went wrong")
		return
	}

	if h.config.APIServer.WriteLogToCloudStorage() {
		err = h.cloudStorageClient.InsertReleaseLog(c.Request.Context(), insertedReleaseLog)
		if err != nil {
			log.Error().Err(err).
				Msgf("Failed inserting release logs into cloud storage for %v/%v/%v/%v", source, owner, repo, id)
		}
	}

	c.String(http.StatusOK, "Aye aye!")
}

func (h *Handler) GetFrequentLabels(c *gin.Context) {

	pageNumber, pageSize, filters, _ := api.GetQueryParameters(c)

	// filter on organizations / groups
	filters = api.SetPermissionsFilters(c, filters)

	response, err := api.GetPagedListResponse(
		func() ([]interface{}, error) {
			labels, err := h.cockroachDBClient.GetFrequentLabels(c.Request.Context(), pageNumber, pageSize, filters)
			if err != nil {
				return nil, err
			}

			// convert typed array to interface array O(n)
			items := make([]interface{}, len(labels))
			for i := range labels {
				items[i] = labels[i]
			}

			return items, nil
		},
		func() (int, error) {
			return h.cockroachDBClient.GetFrequentLabelsCount(c.Request.Context(), filters)
		},
		pageNumber,
		pageSize)

	if err != nil {
		log.Error().Err(err).Msg("Failed retrieving frequent labels from db")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPipelineStatsBuildsDurations(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	durations, err := h.cockroachDBClient.GetPipelineBuildsDurations(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build durations from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"durations": durations,
	})
}

func (h *Handler) GetPipelineStatsReleasesDurations(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	durations, err := h.cockroachDBClient.GetPipelineReleasesDurations(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving releases durations from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"durations": durations,
	})
}

func (h *Handler) GetPipelineStatsBuildsCPUUsageMeasurements(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	measurements, err := h.cockroachDBClient.GetPipelineBuildsCPUUsageMeasurements(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build cpu usage measurements from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
	})
}

func (h *Handler) GetPipelineStatsReleasesCPUUsageMeasurements(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	measurements, err := h.cockroachDBClient.GetPipelineReleasesCPUUsageMeasurements(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving release cpu usage measurements from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
	})
}

func (h *Handler) GetPipelineStatsBuildsMemoryUsageMeasurements(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	measurements, err := h.cockroachDBClient.GetPipelineBuildsMemoryUsageMeasurements(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build memory usage measurements from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
	})
}

func (h *Handler) GetPipelineStatsReleasesMemoryUsageMeasurements(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	// get filters (?filter[last]=100)
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")
	filters[api.FilterLast] = api.GetLastFilter(c, 100)

	measurements, err := h.cockroachDBClient.GetPipelineReleasesMemoryUsageMeasurements(c.Request.Context(), source, owner, repo, filters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving release memory usage measurements from db for %v/%v/%v", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
	})
}

func (h *Handler) GetPipelineWarnings(c *gin.Context) {

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	filters := api.GetPipelineFilters(c)

	pipeline, err := h.cockroachDBClient.GetPipeline(c.Request.Context(), source, owner, repo, filters, false)
	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving pipeline for %v/%v/%v from db", source, owner, repo)
	}
	if pipeline == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusText(http.StatusNotFound), "message": "Pipeline not found"})
		return
	}

	warnings := []contracts.Warning{}

	// get filters (?filter[last]=100)
	buildsFilters := map[api.FilterType][]string{}

	// only use durations of successful builds
	buildsFilters[api.FilterStatus] = api.GetStatusFilter(c, "succeeded")

	// get last 25 builds
	buildsFilters[api.FilterLast] = api.GetLastFilter(c, 25)

	durations, err := h.cockroachDBClient.GetPipelineBuildsDurations(c.Request.Context(), source, owner, repo, buildsFilters)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed retrieving build durations from db for pipeline %v/%v/%v warnings", source, owner, repo)
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	if len(durations) > 0 {
		// pick the item at half of the length
		medianIndex := len(durations)/2 - 1
		if medianIndex < 0 {
			medianIndex = 0
		}
		duration := durations[medianIndex]["duration"].(time.Duration)
		durationInSeconds := duration.Seconds()

		if durationInSeconds > 300.0 {
			warnings = append(warnings, contracts.Warning{
				Status:  "danger",
				Message: fmt.Sprintf("The [median build time](/pipelines/%v/%v/%v/statistics?last=25) of this pipeline is **%v**. This is too slow, please optimize your build speed by using smaller images or running less intensive steps to ensure it finishes at least within 5 minutes, but preferably within 2 minutes.", source, owner, repo, duration),
			})
		} else if durationInSeconds > 120.0 {
			warnings = append(warnings, contracts.Warning{
				Status:  "warning",
				Message: fmt.Sprintf("The [median build time](/pipelines/%v/%v/%v/statistics?last=25) of this pipeline is **%v**. This is a bit too slow, please optimize your build speed by using smaller images or running less intensive steps to ensure it finishes within 2 minutes.", source, owner, repo, duration),
			})
		}
	}

	manifestWarnings, err := h.warningHelper.GetManifestWarnings(pipeline.ManifestObject, pipeline.GetFullRepoPath())
	if err != nil {
		log.Error().Err(err).
			Msgf("Failed getting warnings for %v/%v/%v manifest", source, owner, repo)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": "Failed getting warnings for manifest"})
		return
	}
	warnings = append(warnings, manifestWarnings...)

	c.JSON(http.StatusOK, gin.H{"warnings": warnings})
}

func (h *Handler) GetCatalogFilters(c *gin.Context) {

	if h.config == nil || h.config.Catalog == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, h.config.Catalog.Filters)
}

func (h *Handler) GetCatalogFilterValues(c *gin.Context) {

	labelKey := c.DefaultQuery("filter[labels]", "type")

	labels, err := h.cockroachDBClient.GetLabelValues(c.Request.Context(), labelKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed retrieving label values from db")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (h *Handler) GetStatsPipelinesCount(c *gin.Context) {

	// get filters (?filter[status]=running,succeeded&filter[since]=1w
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c)
	filters[api.FilterSince] = api.GetSinceFilter(c)

	pipelinesCount, err := h.cockroachDBClient.GetPipelinesCount(c.Request.Context(), filters)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed retrieving pipelines count from db")
	}

	c.JSON(http.StatusOK, gin.H{
		"count": pipelinesCount,
	})
}

func (h *Handler) GetStatsReleasesCount(c *gin.Context) {

	// get filters (?filter[status]=running,succeeded&filter[since]=1w
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c)
	filters[api.FilterSince] = api.GetSinceFilter(c)

	releasesCount, err := h.cockroachDBClient.GetReleasesCount(c.Request.Context(), filters)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed retrieving releases count from db")
	}

	c.JSON(http.StatusOK, gin.H{
		"count": releasesCount,
	})
}

func (h *Handler) GetStatsBuildsCount(c *gin.Context) {

	// get filters (?filter[status]=running,succeeded&filter[since]=1w
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c)
	filters[api.FilterSince] = api.GetSinceFilter(c)

	buildsCount, err := h.cockroachDBClient.GetBuildsCount(c.Request.Context(), filters)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed retrieving builds count from db")
	}

	c.JSON(http.StatusOK, gin.H{
		"count": buildsCount,
	})
}

func (h *Handler) GetStatsMostBuilds(c *gin.Context) {

	pageNumber, pageSize, filters, _ := api.GetQueryParameters(c)

	pipelines, err := h.cockroachDBClient.GetPipelinesWithMostBuilds(c.Request.Context(), pageNumber, pageSize, filters)
	if err != nil {
		errorMessage := "Failed retrieving pipelines with most builds from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}
	pipelinesCount, err := h.cockroachDBClient.GetPipelinesWithMostBuildsCount(c.Request.Context(), filters)
	if err != nil {
		errorMessage := "Failed retrieving pipelines count from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	response := contracts.ListResponse{
		Pagination: contracts.Pagination{
			Page:       pageNumber,
			Size:       pageSize,
			TotalItems: pipelinesCount,
			TotalPages: int(math.Ceil(float64(pipelinesCount) / float64(pageSize))),
		},
	}

	response.Items = make([]interface{}, len(pipelines))
	for i := range pipelines {
		response.Items[i] = pipelines[i]
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStatsMostReleases(c *gin.Context) {

	pageNumber, pageSize, filters, _ := api.GetQueryParameters(c)

	pipelines, err := h.cockroachDBClient.GetPipelinesWithMostReleases(c.Request.Context(), pageNumber, pageSize, filters)
	if err != nil {
		errorMessage := "Failed retrieving pipelines with most builds from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}
	pipelinesCount, err := h.cockroachDBClient.GetPipelinesWithMostReleasesCount(c.Request.Context(), filters)
	if err != nil {
		errorMessage := "Failed retrieving pipelines count from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	response := contracts.ListResponse{
		Pagination: contracts.Pagination{
			Page:       pageNumber,
			Size:       pageSize,
			TotalItems: pipelinesCount,
			TotalPages: int(math.Ceil(float64(pipelinesCount) / float64(pageSize))),
		},
	}

	response.Items = make([]interface{}, len(pipelines))
	for i := range pipelines {
		response.Items[i] = pipelines[i]
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStatsBuildsDuration(c *gin.Context) {

	// get filters (?filter[status]=running,succeeded&filter[since]=1w
	filters := map[api.FilterType][]string{}
	filters[api.FilterStatus] = api.GetStatusFilter(c)
	filters[api.FilterSince] = api.GetSinceFilter(c)

	buildsDuration, err := h.cockroachDBClient.GetBuildsDuration(c.Request.Context(), filters)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed retrieving builds duration from db")
	}

	c.JSON(http.StatusOK, gin.H{
		"duration": buildsDuration,
	})
}

func (h *Handler) GetStatsBuildsAdoption(c *gin.Context) {

	buildTimes, err := h.cockroachDBClient.GetFirstBuildTimes(c.Request.Context())
	if err != nil {
		errorMessage := "Failed retrieving first build times from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"datetimes": buildTimes,
	})
}

func (h *Handler) GetStatsReleasesAdoption(c *gin.Context) {

	releaseTimes, err := h.cockroachDBClient.GetFirstReleaseTimes(c.Request.Context())
	if err != nil {
		errorMessage := "Failed retrieving first release times from db"
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "message": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"datetimes": releaseTimes,
	})
}

func (h *Handler) GetConfig(c *gin.Context) {

	configBytes, err := yaml.Marshal(h.encryptedConfig)
	if err != nil {
		log.Error().Err(err).Msgf("Failed marshalling encrypted config")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// obfuscate all secrets
	configString, err := h.obfuscateSecrets(string(configBytes))
	if err != nil {
		log.Error().Err(err).Msgf("Failed obfuscating secrets")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// add extra whitespace after each top-level item
	addWhitespaceRegex := regexp.MustCompile(`\n([a-z])`)
	configString = addWhitespaceRegex.ReplaceAllString(configString, "\n\n$1")

	c.JSON(http.StatusOK, gin.H{"config": configString})
}

func (h *Handler) GetConfigCredentials(c *gin.Context) {

	configBytes, err := yaml.Marshal(h.encryptedConfig.Credentials)
	if err != nil {
		log.Error().Err(err).Msgf("Failed marshalling encrypted config")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// obfuscate all secrets
	configString, err := h.obfuscateSecrets(string(configBytes))
	if err != nil {
		log.Error().Err(err).Msgf("Failed obfuscating secrets")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// add extra whitespace after each top-level item
	addWhitespaceRegex := regexp.MustCompile(`\n([a-z])`)
	configString = addWhitespaceRegex.ReplaceAllString(configString, "\n\n$1")

	c.JSON(http.StatusOK, gin.H{"config": configString})
}

func (h *Handler) GetConfigTrustedImages(c *gin.Context) {

	configBytes, err := yaml.Marshal(h.encryptedConfig.TrustedImages)
	if err != nil {
		log.Error().Err(err).Msgf("Failed marshalling encrypted config")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// obfuscate all secrets
	configString, err := h.obfuscateSecrets(string(configBytes))
	if err != nil {
		log.Error().Err(err).Msgf("Failed obfuscating secrets")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// add extra whitespace after each top-level item
	addWhitespaceRegex := regexp.MustCompile(`\n([a-z])`)
	configString = addWhitespaceRegex.ReplaceAllString(configString, "\n\n$1")

	c.JSON(http.StatusOK, gin.H{"config": configString})
}

func (h *Handler) GetManifestTemplates(c *gin.Context) {

	configFiles, err := ioutil.ReadDir(filepath.Dir(h.configFilePath))
	if err != nil {
		log.Error().Err(err).Msgf("Failed listing config files directory")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	templates := []interface{}{}
	for _, f := range configFiles {

		configfileName := f.Name()

		// check if it's a manifest template
		re := regexp.MustCompile(`^manifest-(.+)\.tmpl`)
		match := re.FindStringSubmatch(configfileName)

		if len(match) == 2 {

			// read template file
			templateFilePath := filepath.Join(filepath.Dir(h.configFilePath), configfileName)
			data, err := ioutil.ReadFile(templateFilePath)
			if err != nil {
				log.Error().Err(err).Msgf("Failed reading template file %v", templateFilePath)
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
				return
			}

			placeholderRegex := regexp.MustCompile(`{{\.([a-zA-Z0-9]+)}}`)
			placeholderMatches := placeholderRegex.FindAllStringSubmatch(string(data), -1)

			// reduce and deduplicate [["{{.Application}}","Application"],["{{.Team}}","Team"],["{{.ProjectName}}","ProjectName"],["{{.ProjectName}}","ProjectName"]] to ["Application","Team","ProjectName"]
			placeholders := []string{}
			for _, m := range placeholderMatches {
				if len(m) == 2 && !api.StringArrayContains(placeholders, m[1]) {
					placeholders = append(placeholders, m[1])
				}
			}

			templateData := map[string]interface{}{
				"template":     match[1],
				"placeholders": placeholders,
			}

			templates = append(templates, templateData)
		}
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *Handler) GenerateManifest(c *gin.Context) {

	var aux struct {
		Template     string            `json:"template"`
		Placeholders map[string]string `json:"placeholders,omitempty"`
	}

	err := c.BindJSON(&aux)
	if err != nil {
		errorMessage := fmt.Sprint("Binding GenerateManifest body failed")
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	templateFilePath := filepath.Join(filepath.Dir(h.configFilePath), fmt.Sprintf("manifest-%v.tmpl", aux.Template))
	data, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed reading template file %v", templateFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	tmpl, err := template.New(".estafette.yaml").Parse(string(data))
	if err != nil {
		log.Error().Err(err).Msgf("Failed parsing template file %v", templateFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, aux.Placeholders)
	if err != nil {
		log.Error().Err(err).Msgf("Failed rendering template file %v", templateFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manifest": renderedTemplate.String()})
}

func (h *Handler) ValidateManifest(c *gin.Context) {

	var aux struct {
		Manifest string `json:"manifest"`
	}

	err := c.BindJSON(&aux)
	if err != nil {
		errorMessage := fmt.Sprint("Binding ValidateManifest body failed")
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	_, err = manifest.ReadManifest(h.config.ManifestPreferences, aux.Manifest, true)
	status := "succeeded"
	errorString := ""
	if err != nil {
		status = "failed"
		errorString = err.Error()
	}

	c.JSON(http.StatusOK, gin.H{"status": status, "errors": errorString})
}

func (h *Handler) EncryptSecret(c *gin.Context) {

	var aux struct {
		Base64Encode      bool   `json:"base64"`
		DoubleEncrypt     bool   `json:"double"`
		PipelineWhitelist string `json:"pipelineWhitelist"`
		Value             string `json:"value"`
	}

	err := c.BindJSON(&aux)
	if err != nil {
		errorMessage := fmt.Sprint("Binding EncryptSecret body failed")
		log.Error().Err(err).Msg(errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": errorMessage})
		return
	}

	value := aux.Value

	// trim any whitespace and newlines at beginning and end of string
	value = strings.TrimSpace(value)

	if aux.Base64Encode {
		// base64 encode for use in a kubernetes secret
		value = base64.StdEncoding.EncodeToString([]byte(value))
	}

	encryptedString, err := h.secretHelper.EncryptEnvelope(value, aux.PipelineWhitelist)
	if err != nil {
		log.Error().Err(err).Msg("Failed encrypting secret")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	if aux.DoubleEncrypt {
		encryptedString, err = h.secretHelper.EncryptEnvelope(encryptedString, crypt.DefaultPipelineWhitelist)
		if err != nil {
			log.Error().Err(err).Msg("Failed encrypting secret")
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"secret": encryptedString})
}

func (h *Handler) PostCronEvent(c *gin.Context) {

	// ensure the user has administrator role
	if !api.RequestTokenHasRole(c, api.RoleCronTrigger) {
		c.JSON(http.StatusForbidden, gin.H{"code": http.StatusText(http.StatusForbidden), "message": "JWT is invalid or user does not have cron-trigger role"})
		return
	}

	err := h.buildService.FireCronTriggers(c.Request.Context())

	if err != nil {
		log.Error().Err(err).Msg("Failed firing cron triggers")
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hey Cron, here's a tock for your tick"})
}

func (h *Handler) CopyLogsToCloudStorage(c *gin.Context) {

	// ensure the user has administrator role
	if !api.RequestTokenHasRole(c, api.RoleLogMigrator) {
		c.JSON(http.StatusForbidden, gin.H{"code": http.StatusText(http.StatusForbidden), "message": "JWT is invalid or user does not have log-migrator role"})
		return
	}

	pageNumber, pageSize, filters, _ := api.GetQueryParameters(c)

	searchValue := "builds"
	if search, ok := filters[api.FilterSearch]; ok && len(search) > 0 && search[0] != "" {
		searchValue = search[0]
	}

	source := c.Param("source")
	owner := c.Param("owner")
	repo := c.Param("repo")

	if searchValue == "builds" {
		buildLogs, err := h.cockroachDBClient.GetPipelineBuildLogsPerPage(c.Request.Context(), source, owner, repo, pageNumber, pageSize)
		if err != nil {
			log.Error().Err(err).Int("pageNumber", pageNumber).Int("pageSize", pageSize).Msgf("Failed retrieving build logs for %v/%v/%v", source, owner, repo)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "error": err})
			return
		}

		errors := make(chan error, len(buildLogs))

		var wg sync.WaitGroup
		wg.Add(len(buildLogs))

		for _, bl := range buildLogs {
			go func(ctx context.Context, bl contracts.BuildLog) {
				defer wg.Done()

				err = h.cloudStorageClient.InsertBuildLog(c.Request.Context(), bl)
				if err != nil {
					errors <- err
				}
			}(c.Request.Context(), *bl)
		}

		// wait for all parallel runs to finish
		wg.Wait()

		// return error if any of them have been generated
		close(errors)
		for e := range errors {
			log.Error().Err(err).Int("pageNumber", pageNumber).Int("pageSize", pageSize).Msgf("Failed inserting build logs for %v/%v/%v into cloud storage", source, owner, repo)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "error": e})
			return
		}

		c.String(http.StatusOK, strconv.Itoa(len(buildLogs)))
		return

	} else if searchValue == "releases" {
		releaseLogs, err := h.cockroachDBClient.GetPipelineReleaseLogsPerPage(c.Request.Context(), source, owner, repo, pageNumber, pageSize)
		if err != nil {
			log.Error().Err(err).Int("pageNumber", pageNumber).Int("pageSize", pageSize).Msgf("Failed retrieving release logs for %v/%v/%v", source, owner, repo)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "error": err})
			return
		}

		errors := make(chan error, len(releaseLogs))

		var wg sync.WaitGroup
		wg.Add(len(releaseLogs))

		for _, rl := range releaseLogs {
			go func(ctx context.Context, rl contracts.ReleaseLog) {
				defer wg.Done()

				err = h.cloudStorageClient.InsertReleaseLog(c.Request.Context(), rl)
				if err != nil {
					errors <- err
				}
			}(c.Request.Context(), *rl)
		}

		// wait for all parallel runs to finish
		wg.Wait()

		// return error if any of them have been generated
		close(errors)
		for e := range errors {
			log.Error().Err(err).Int("pageNumber", pageNumber).Int("pageSize", pageSize).Msgf("Failed inserting release logs for %v/%v/%v into cloud storage", source, owner, repo)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusText(http.StatusInternalServerError), "error": e})
			return
		}

		c.String(http.StatusOK, strconv.Itoa(len(releaseLogs)))
		return
	}

	c.String(http.StatusOK, "Aye aye!")
}

func (h *Handler) obfuscateSecrets(input string) (string, error) {

	r, err := regexp.Compile(`estafette\.secret\(([a-zA-Z0-9.=_-]+)\)`)
	if err != nil {
		return "", err
	}

	// obfuscate all secrets
	return r.ReplaceAllLiteralString(input, "***"), nil
}

func (h *Handler) Commands(c *gin.Context) {

	// ensure the request has the correct claims
	claims := jwt.ExtractClaims(c)
	job := claims["job"].(string)
	if job == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusText(http.StatusUnauthorized), "message": "JWT is invalid or has invalid claim"})
		return
	}

	eventType := c.GetHeader("X-Estafette-Event")
	log.Debug().Msgf("X-Estafette-Event is set to %v", eventType)
	// h.prometheusInboundEventTotals.With(prometheus.Labels{"event": eventType, "source": "estafette"}).Inc()

	eventJobname := c.GetHeader("X-Estafette-Event-Job-Name")
	log.Debug().Msgf("X-Estafette-Event-Job-Name is set to %v", eventJobname)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Reading body from Estafette 'build finished' event failed")
		c.String(http.StatusInternalServerError, "Reading body from Estafette 'build finished' event failed")
		return
	}

	log.Debug().Msgf("Read body for /api/commands for job %v", eventJobname)

	switch eventType {
	case
		"builder:nomanifest",
		"builder:running",
		"builder:succeeded",
		"builder:failed",
		"builder:canceled":

		// unmarshal json body
		var ciBuilderEvent builderapi.CiBuilderEvent
		err = json.Unmarshal(body, &ciBuilderEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to CiBuilderEvent failed")
			return
		}

		log.Debug().Interface("ciBuilderEvent", ciBuilderEvent).Msgf("Unmarshaled body of /api/commands event %v for job %v", eventType, eventJobname)

		err := h.buildService.UpdateBuildStatus(c.Request.Context(), ciBuilderEvent)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed updating build status for job %v to %v, not removing the job", eventJobname, ciBuilderEvent.BuildStatus)
			log.Error().Err(err).Interface("ciBuilderEvent", ciBuilderEvent).Msg(errorMessage)
			c.String(http.StatusInternalServerError, errorMessage)
			return
		}

	case "builder:clean":

		// unmarshal json body
		var ciBuilderEvent builderapi.CiBuilderEvent
		err = json.Unmarshal(body, &ciBuilderEvent)
		if err != nil {
			log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to CiBuilderEvent failed")
			return
		}

		log.Debug().Interface("ciBuilderEvent", ciBuilderEvent).Msgf("Unmarshaled body of /api/commands event %v for job %v", eventType, eventJobname)

		if ciBuilderEvent.BuildStatus != "canceled" {
			go func(eventJobname string) {
				err = h.ciBuilderClient.RemoveCiBuilderJob(c.Request.Context(), eventJobname)
				if err != nil {
					errorMessage := fmt.Sprintf("Failed removing job %v for event %v", eventJobname, eventType)
					log.Error().Err(err).Interface("ciBuilderEvent", ciBuilderEvent).Msg(errorMessage)
				}
			}(eventJobname)
		} else {
			log.Info().Msgf("Job %v is already removed by cancellation, no need to remove for event %v", eventJobname, eventType)
		}

		go func(ctx context.Context, ciBuilderEvent builderapi.CiBuilderEvent) {
			err := h.buildService.UpdateJobResources(c.Request.Context(), ciBuilderEvent)
			if err != nil {
				log.Error().Err(err).Msgf("Failed updating max cpu and memory from prometheus for pod %v", ciBuilderEvent.PodName)
			}
		}(c.Request.Context(), ciBuilderEvent)

	default:
		log.Warn().Str("event", eventType).Msgf("Unsupported Estafette event of type '%v'", eventType)
	}

	c.String(http.StatusOK, "Aye aye!")
}
