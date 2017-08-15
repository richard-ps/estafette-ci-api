package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	// channel for passing push events to handler that creates ci-builder job
	bitbucketPushEvents = make(chan BitbucketRepositoryPushEvent, 100)
)

func bitbucketWebhookHandler(w http.ResponseWriter, r *http.Request) {

	// https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html

	eventType := r.Header.Get("X-Event-Key")
	webhookTotal.With(prometheus.Labels{"event": eventType, "source": "bitbucket"}).Inc()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Reading body from Bitbucket webhook failed")
		http.Error(w, "Reading body from Bitbucket webhook failed", 500)
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body from Bitbucket webhook failed")
		http.Error(w, "Deserializing body from Bitbucket webhook failed", 500)
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Interface("headers", r.Header).
		Interface("body", b).
		Msgf("Received webhook event of type '%v' from Bitbucket...", eventType)

	switch eventType {
	case "repo:push":
		handleBitbucketPush(body)

	case "repo:fork":
	case "repo:updated":
	case "repo:transfer":
	case "repo:commit_comment_created":
	case "repo:commit_status_created":
	case "repo:commit_status_updated":
	case "issue:created":
	case "issue:updated":
	case "issue:comment_created":
	case "pullrequest:created":
	case "pullrequest:updated":
	case "pullrequest:approved":
	case "pullrequest:unapproved":
	case "pullrequest:fulfilled":
	case "pullrequest:rejected":
	case "pullrequest:comment_created":
	case "pullrequest:comment_updated":
	case "pullrequest:comment_deleted":
		log.Debug().Str("event", eventType).Msgf("Not implemented Bitbucket webhook event of type '%v'", eventType)

	default:
		log.Warn().Str("event", eventType).Msgf("Unsupported Bitbucket webhook event of type '%v'", eventType)
	}

	fmt.Fprintf(w, "Aye aye!")
}

func handleBitbucketPush(body []byte) {

	// unmarshal json body
	var pushEvent BitbucketRepositoryPushEvent
	err := json.Unmarshal(body, &pushEvent)
	if err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Deserializing body to BitbucketRepositoryPushEvent failed")
		return
	}

	log.Debug().Interface("pushEvent", pushEvent).Msgf("Deserialized Bitbucket push event for repository %v", pushEvent.Repository.FullName)

	// test making api calls for bitbucket app in the background
	bitbucketPushEvents <- pushEvent
}
