package main

import (
	"context"
	"testing"

	"github.com/estafette/estafette-ci-api/clients/bitbucketapi"
	"github.com/estafette/estafette-ci-api/clients/builderapi"
	"github.com/estafette/estafette-ci-api/clients/cloudsourceapi"
	"github.com/estafette/estafette-ci-api/clients/cloudstorage"
	"github.com/estafette/estafette-ci-api/clients/cockroachdb"
	"github.com/estafette/estafette-ci-api/clients/githubapi"
	"github.com/estafette/estafette-ci-api/clients/pubsubapi"
	"github.com/estafette/estafette-ci-api/clients/slackapi"
	"github.com/estafette/estafette-ci-api/config"
	"github.com/estafette/estafette-ci-api/helpers"
	"github.com/estafette/estafette-ci-api/services/bitbucket"
	"github.com/estafette/estafette-ci-api/services/cloudsource"
	"github.com/estafette/estafette-ci-api/services/estafette"
	"github.com/estafette/estafette-ci-api/services/github"
	"github.com/estafette/estafette-ci-api/services/pubsub"
	"github.com/estafette/estafette-ci-api/services/rbac"
	"github.com/estafette/estafette-ci-api/services/slack"
	crypt "github.com/estafette/estafette-ci-crypt"
)

func TestConfigureGinGonic(t *testing.T) {
	t.Run("DoesNotPanic", func(t *testing.T) {

		config := &config.APIConfig{
			Auth: &config.AuthConfig{
				JWT: &config.JWTConfig{
					Domain: "mydomain",
					Key:    "abc",
				},
			},
		}

		ctx := context.Background()

		cockroachdbClient := cockroachdb.MockClient{}
		cloudstorageClient := cloudstorage.MockClient{}
		builderapiClient := builderapi.MockClient{}
		estafetteService := estafette.MockService{}
		secretHelper := crypt.NewSecretHelper("abc", false)
		warningHelper := helpers.NewWarningHelper(secretHelper)
		githubapiClient := githubapi.MockClient{}
		bitbucketapiClient := bitbucketapi.MockClient{}
		cloudsourceapiClient := cloudsourceapi.MockClient{}
		pubsubapiclient := pubsubapi.MockClient{}
		slackapiClient := slackapi.MockClient{}

		bitbucketHandler := bitbucket.NewHandler(bitbucket.MockService{})
		githubHandler := github.NewHandler(github.MockService{})
		estafetteHandler := estafette.NewHandler("", config, config, cockroachdbClient, cloudstorageClient, builderapiClient, estafetteService, warningHelper, secretHelper, githubapiClient.JobVarsFunc(ctx), bitbucketapiClient.JobVarsFunc(ctx), cloudsourceapiClient.JobVarsFunc(ctx))

		rbacHandler := rbac.NewHandler(config, rbac.MockService{}, cockroachdbClient)
		pubsubHandler := pubsub.NewHandler(pubsubapiclient, estafetteService)
		slackHandler := slack.NewHandler(secretHelper, config, slackapiClient, cockroachdbClient, estafetteService, githubapiClient.JobVarsFunc(ctx), bitbucketapiClient.JobVarsFunc(ctx))
		cloudsourceHandler := cloudsource.NewHandler(pubsubapiclient, cloudsource.MockService{})

		// act
		_ = configureGinGonic(config, bitbucketHandler, githubHandler, estafetteHandler, rbacHandler, pubsubHandler, slackHandler, cloudsourceHandler)
	})
}