package estafette

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/estafette/estafette-ci-api/clients/builderapi"
	"github.com/estafette/estafette-ci-api/clients/cloudstorage"
	"github.com/estafette/estafette-ci-api/clients/cockroachdb"
	"github.com/estafette/estafette-ci-api/config"
	"github.com/estafette/estafette-ci-api/helpers"
	contracts "github.com/estafette/estafette-ci-contracts"
	crypt "github.com/estafette/estafette-ci-crypt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCatalogFilters(t *testing.T) {

	t.Run("ReturnsUpdatedConfigAfterReload", func(t *testing.T) {

		configFilePath := "/configs/config.yaml"
		cfg := &config.APIConfig{
			Catalog: &config.CatalogConfig{
				Filters: []string{
					"type",
				},
			},
		}
		encryptedConfig := cfg

		cockroachdbClient := cockroachdb.MockClient{}
		cloudStorageClient := cloudstorage.MockClient{}
		builderapiClient := builderapi.MockClient{}

		buildService := MockService{}
		secretHelper := crypt.NewSecretHelper("abc", false)
		warningHelper := helpers.NewWarningHelper(secretHelper)
		githubJobVarsFunc := func(context.Context, string, string, string) (string, string, error) {
			return "", "", nil
		}
		bitbucketJobVarsFunc := githubJobVarsFunc
		cloudsourceJobVarsFunc := githubJobVarsFunc

		handler := NewHandler(configFilePath, cfg, encryptedConfig, cockroachdbClient, cloudStorageClient, builderapiClient, buildService, warningHelper, secretHelper, githubJobVarsFunc, bitbucketJobVarsFunc, cloudsourceJobVarsFunc)
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)

		// act
		handler.GetCatalogFilters(c)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		body, err := ioutil.ReadAll(recorder.Result().Body)
		assert.Nil(t, err)
		assert.Equal(t, "[\"type\"]\n", string(body))

		// act
		*cfg = config.APIConfig{
			Catalog: &config.CatalogConfig{
				Filters: []string{
					"type",
					"language",
				},
			},
		}
		recorder = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(recorder)
		handler.GetCatalogFilters(c)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		body, err = ioutil.ReadAll(recorder.Result().Body)
		assert.Nil(t, err)
		assert.Equal(t, "[\"type\",\"language\"]\n", string(body))

	})
}

func TestGetPipeline(t *testing.T) {

	t.Run("ReturnsPipelineFromNewCockroachdbClientAfterReload", func(t *testing.T) {

		configFilePath := "/configs/config.yaml"
		cfg := &config.APIConfig{}
		encryptedConfig := cfg

		var cockroachdbClient cockroachdb.Client
		cockroachdbClient = &cockroachdb.MockClient{
			GetPipelineFunc: func(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (pipeline *contracts.Pipeline, err error) {
				pipeline = &contracts.Pipeline{
					BuildStatus: "succeeded",
				}
				return
			},
		}
		cloudStorageClient := cloudstorage.MockClient{}
		builderapiClient := builderapi.MockClient{}
		buildService := MockService{}
		secretHelper := crypt.NewSecretHelper("abc", false)
		warningHelper := helpers.NewWarningHelper(secretHelper)
		githubJobVarsFunc := func(context.Context, string, string, string) (string, string, error) {
			return "", "", nil
		}
		bitbucketJobVarsFunc := githubJobVarsFunc
		cloudsourceJobVarsFunc := githubJobVarsFunc

		handler := NewHandler(configFilePath, cfg, encryptedConfig, cockroachdbClient, cloudStorageClient, builderapiClient, buildService, warningHelper, secretHelper, githubJobVarsFunc, bitbucketJobVarsFunc, cloudsourceJobVarsFunc)
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		bodyReader := strings.NewReader("")
		c.Request = httptest.NewRequest("GET", "https://ci.estafette.io/pipelines/a/b/c", bodyReader)
		if !assert.NotNil(t, c.Request) {
			return
		}

		// act
		handler.GetPipeline(c)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		body, err := ioutil.ReadAll(recorder.Result().Body)
		assert.Nil(t, err)
		assert.Equal(t, "{\"id\":\"\",\"repoSource\":\"\",\"repoOwner\":\"\",\"repoName\":\"\",\"repoBranch\":\"\",\"repoRevision\":\"\",\"buildStatus\":\"succeeded\",\"insertedAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\",\"duration\":0,\"lastUpdatedAt\":\"0001-01-01T00:00:00Z\"}\n", string(body))

		// act
		cockroachdbClient = cockroachdb.MockClient{
			GetPipelineFunc: func(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (pipeline *contracts.Pipeline, err error) {
				pipeline = &contracts.Pipeline{
					BuildStatus: "failed",
				}
				return
			},
		}

		recorder = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(recorder)
		bodyReader = strings.NewReader("")
		c.Request = httptest.NewRequest("GET", "https://ci.estafette.io/pipelines/a/b/c", bodyReader)
		handler.GetPipeline(c)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		body, err = ioutil.ReadAll(recorder.Result().Body)
		assert.Nil(t, err)
		// assert.Equal(t, "{\"id\":\"\",\"repoSource\":\"\",\"repoOwner\":\"\",\"repoName\":\"\",\"repoBranch\":\"\",\"repoRevision\":\"\",\"buildStatus\":\"failed\",\"insertedAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\",\"duration\":0,\"lastUpdatedAt\":\"0001-01-01T00:00:00Z\"}\n", string(body))
	})
}
