package main

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/estafette/estafette-ci-crypt"

	"github.com/alecthomas/kingpin"
	"github.com/estafette/estafette-ci-api/bitbucket"
	"github.com/estafette/estafette-ci-api/cockroach"
	"github.com/estafette/estafette-ci-api/estafette"
	"github.com/estafette/estafette-ci-api/github"
	"github.com/estafette/estafette-ci-api/slack"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	prometheusMetricsAddress = kingpin.Flag("metrics-listen-address", "The address to listen on for Prometheus metrics requests.").Default(":9001").String()
	prometheusMetricsPath    = kingpin.Flag("metrics-path", "The path to listen for Prometheus metrics requests.").Default("/metrics").String()

	apiAddress = kingpin.Flag("api-listen-address", "The address to listen on for api HTTP requests.").Default(":5000").String()

	githubAppPrivateKeyPath    = kingpin.Flag("github-app-privatey-key-path", "The path to the pem file for the private key of the Github App.").Default("/github-app-key/private-key.pem").String()
	githubAppID                = kingpin.Flag("github-app-id", "The Github App id.").Envar("GITHUB_APP_ID").String()
	githubAppOAuthClientID     = kingpin.Flag("github-app-oauth-client-id", "The OAuth client id for the Github App.").Envar("GITHUB_APP_OAUTH_CLIENT_ID").String()
	githubAppOAuthClientSecret = kingpin.Flag("github-app-oauth-client-secret", "The OAuth client secret for the Github App.").Envar("GITHUB_APP_OAUTH_CLIENT_SECRET").String()
	githubWebhookSecret        = kingpin.Flag("github-webhook-secret", "The secret to verify webhook authenticity.").Envar("GITHUB_WEBHOOK_SECRET").String()

	bitbucketAPIKey         = kingpin.Flag("bitbucket-api-key", "The api key for Bitbucket.").Envar("BITBUCKET_API_KEY").String()
	bitbucketAppOAuthKey    = kingpin.Flag("bitbucket-app-oauth-key", "The OAuth key for the Bitbucket App.").Envar("BITBUCKET_APP_OAUTH_KEY").String()
	bitbucketAppOAuthSecret = kingpin.Flag("bitbucket-app-oauth-secret", "The OAuth secret for the Bitbucket App.").Envar("BITBUCKET_APP_OAUTH_SECRET").String()

	estafetteCiServerBaseURL = kingpin.Flag("estafette-ci-server-base-url", "The base url of this api server.").Envar("ESTAFETTE_CI_SERVER_BASE_URL").String()
	estafetteCiAPIKey        = kingpin.Flag("estafette-ci-api-key", "An api key for estafette itself to use until real oauth is supported.").Envar("ESTAFETTE_CI_API_KEY").String()

	slackAppClientID          = kingpin.Flag("slack-app-client-id", "The Slack App id for accessing Slack API.").Envar("SLACK_APP_CLIENT_ID").String()
	slackAppClientSecret      = kingpin.Flag("slack-app-client-secret", "The Slack App secret for accessing Slack API.").Envar("SLACK_APP_CLIENT_ID").String()
	slackAppVerificationToken = kingpin.Flag("slack-app-verification-token", "The token used to verify incoming Slack webhook events.").Envar("SLACK_APP_VERIFICATION_TOKEN").String()
	slackAppOAuthAccessToken  = kingpin.Flag("slack-app-oauth-access-token", "The OAuth access token for the Slack App.").Envar("SLACK_APP_OAUTH_ACCESS_TOKEN").String()

	secretDecryptionKey = kingpin.Flag("secret-decryption-key", "The AES-256 key used to decrypt secrets that have been encrypted with it.").Envar("SECRET_DECRYPTION_KEY").String()

	cockroachDatabase       = kingpin.Flag("cockroach-database", "CockroachDB database.").Envar("COCKROACH_DATABASE").String()
	cockroachHost           = kingpin.Flag("cockroach-host", "CockroachDB host.").Envar("COCKROACH_HOST").String()
	cockroachInsecure       = kingpin.Flag("cockroach-insecure", "CockroachDB insecure connection.").Envar("COCKROACH_INSECURE").Bool()
	cockroachCertificateDir = kingpin.Flag("cockroach-certs-dir", "CockroachDB certificate directory.").Envar("COCKROACH_CERTS_DIR").String()
	cockroachPort           = kingpin.Flag("cockroach-port", "CockroachDB port.").Envar("COCKROACH_PORT").Int()
	cockroachUser           = kingpin.Flag("cockroach-user", "CockroachDB user.").Envar("COCKROACH_USER").String()
	cockroachPassword       = kingpin.Flag("cockroach-password", "CockroachDB password.").Envar("COCKROACH_PASSWORD").String()

	// prometheusInboundEventTotals is the prometheus timeline serie that keeps track of inbound events
	prometheusInboundEventTotals = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "estafette_ci_api_inbound_event_totals",
			Help: "Total of inbound events.",
		},
		[]string{"event", "source"},
	)

	// prometheusOutboundAPICallTotals is the prometheus timeline serie that keeps track of outbound api calls
	prometheusOutboundAPICallTotals = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "estafette_ci_api_outbound_api_call_totals",
			Help: "Total of outgoing api calls.",
		},
		[]string{"target"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(prometheusInboundEventTotals)
	prometheus.MustRegister(prometheusOutboundAPICallTotals)
}

func main() {

	// parse command line parameters
	kingpin.Parse()

	// configure json logging
	initLogging()

	// define channels and waitgroup to gracefully shutdown the application
	sigs := make(chan os.Signal, 1)                                    // Create channel to receive OS signals
	stop := make(chan struct{})                                        // Create channel to receive stop signal
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // Register the sigs channel to receieve SIGTERM
	wg := &sync.WaitGroup{}                                            // Goroutines can add themselves to this to be waited on so that they finish

	// start prometheus
	go startPrometheus()

	// handle api requests
	srv := handleRequests(stop, wg)

	// wait for graceful shutdown to finish
	<-sigs // Wait for signals (this hangs until a signal arrives)
	log.Debug().Msg("Shutting down...")

	// shut down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Graceful server shutdown failed")
	}

	log.Debug().Msg("Stopping goroutines...")
	close(stop) // Tell goroutines to stop themselves

	log.Debug().Msg("Awaiting waitgroup...")
	wg.Wait() // Wait for all to be stopped

	log.Info().Msg("Server gracefully stopped")
}

func startPrometheus() {
	log.Debug().
		Str("port", *prometheusMetricsAddress).
		Str("path", *prometheusMetricsPath).
		Msg("Serving Prometheus metrics...")

	http.Handle(*prometheusMetricsPath, promhttp.Handler())

	if err := http.ListenAndServe(*prometheusMetricsAddress, nil); err != nil {
		log.Fatal().Err(err).Msg("Starting Prometheus listener failed")
	}
}

func initLogging() {

	// log as severity for stackdriver logging to recognize the level
	zerolog.LevelFieldName = "severity"

	// set some default fields added to all logs
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "estafette-ci-api").
		Str("version", version).
		Logger()

	// use zerolog for any logs sent via standard log library
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// log startup message
	log.Info().
		Str("branch", branch).
		Str("revision", revision).
		Str("buildDate", buildDate).
		Str("goVersion", goVersion).
		Msg("Starting estafette-ci-api...")
}

func createRouter() *gin.Engine {

	// run gin in release mode and other defaults
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger
	gin.DisableConsoleColor()

	// Creates a router without any middleware by default
	router := gin.New()

	// Logging middleware
	router.Use(ZeroLogMiddleware())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Gzip middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// liveness and readiness
	router.GET("/liveness", func(c *gin.Context) {
		c.String(200, "I'm alive!")
	})
	router.GET("/readiness", func(c *gin.Context) {
		c.String(200, "I'm ready!")
	})

	return router
}

func handleRequests(stopChannel <-chan struct{}, waitGroup *sync.WaitGroup) *http.Server {

	githubAPIClient := github.NewGithubAPIClient(*githubAppPrivateKeyPath, *githubAppID, *githubAppOAuthClientID, *githubAppOAuthClientSecret, prometheusOutboundAPICallTotals)
	bitbucketAPIClient := bitbucket.NewBitbucketAPIClient(*bitbucketAPIKey, *bitbucketAppOAuthKey, *bitbucketAppOAuthSecret, prometheusOutboundAPICallTotals)
	slackAPIClient := slack.NewSlackAPIClient(*slackAppClientID, *slackAppClientSecret, *slackAppOAuthAccessToken, prometheusOutboundAPICallTotals)
	secretHelper := crypt.NewSecretHelper(*secretDecryptionKey)
	cockroachDBClient := cockroach.NewCockroachDBClient(*cockroachDatabase, *cockroachHost, *cockroachInsecure, *cockroachCertificateDir, *cockroachPort, *cockroachUser, *cockroachPassword, prometheusOutboundAPICallTotals)
	ciBuilderClient, err := estafette.NewCiBuilderClient(*estafetteCiServerBaseURL, *estafetteCiAPIKey, *secretDecryptionKey, prometheusOutboundAPICallTotals)
	if err != nil {
		log.Fatal().Err(err).Msg("Creating new CiBuilderClient has failed")
	}

	// set up database
	err = cockroachDBClient.Connect()
	if err != nil {
		log.Warn().Err(err).Msg("Failed connecting to CockroachDB")
	}
	err = cockroachDBClient.MigrateSchema()
	if err != nil {
		log.Warn().Err(err).Msg("Failed migrating schema of CockroachDB")
	}

	// listen to channels for push events
	githubPushEvents := make(chan github.PushEvent, 100)
	githubEventWorker := github.NewGithubEventWorker(stopChannel, waitGroup, githubAPIClient, ciBuilderClient, githubPushEvents)
	githubEventWorker.ListenToEventChannels()

	bitbucketPushEvents := make(chan bitbucket.RepositoryPushEvent, 100)
	bitbucketEventWorker := bitbucket.NewBitbucketEventWorker(stopChannel, waitGroup, bitbucketAPIClient, ciBuilderClient, bitbucketPushEvents)
	bitbucketEventWorker.ListenToEventChannels()

	slackEvents := make(chan slack.SlashCommand, 100)
	slackEventWorker := slack.NewSlackEventWorker(stopChannel, waitGroup, slackAPIClient, slackEvents)
	slackEventWorker.ListenToEventChannels()

	estafetteCiBuilderEvents := make(chan estafette.CiBuilderEvent, 100)
	estafetteBuildJobLogs := make(chan cockroach.BuildJobLogs, 100)
	estafetteEventWorker := estafette.NewEstafetteEventWorker(stopChannel, waitGroup, ciBuilderClient, cockroachDBClient, estafetteCiBuilderEvents, estafetteBuildJobLogs)
	estafetteEventWorker.ListenToEventChannels()

	// listen to http calls
	log.Debug().
		Str("port", *apiAddress).
		Msg("Serving api calls...")

	// create and init router
	router := createRouter()

	githubEventHandler := github.NewGithubEventHandler(githubPushEvents, *githubWebhookSecret, prometheusInboundEventTotals)
	router.POST("/events/github", githubEventHandler.Handle)

	bitbucketEventHandler := bitbucket.NewBitbucketEventHandler(bitbucketPushEvents, prometheusInboundEventTotals)
	router.POST("/events/bitbucket", bitbucketEventHandler.Handle)

	slackEventHandler := slack.NewSlackEventHandler(secretHelper, *slackAppVerificationToken, slackEvents, prometheusInboundEventTotals)
	router.POST("/events/slack/slash", slackEventHandler.Handle)

	estafetteEventHandler := estafette.NewEstafetteEventHandler(*estafetteCiAPIKey, estafetteCiBuilderEvents, estafetteBuildJobLogs, prometheusInboundEventTotals)
	router.POST("/events/estafette/ci-builder", estafetteEventHandler.Handle)

	// instantiate servers instead of using router.Run in order to handle graceful shutdown
	srv := &http.Server{
		Addr:           *apiAddress,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Starting gin router failed")
		}
	}()

	return srv
}
