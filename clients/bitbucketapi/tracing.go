package bitbucketapi

import (
	"context"

	"github.com/estafette/estafette-ci-api/api"
	"github.com/opentracing/opentracing-go"
)

// NewTracingClient returns a new instance of a tracing Client.
func NewTracingClient(c Client) Client {
	return &tracingClient{c, "bitbucketapi"}
}

type tracingClient struct {
	Client
	prefix string
}

func (c *tracingClient) GetAccessToken(ctx context.Context) (accesstoken AccessToken, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetAccessToken"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetAccessToken(ctx)
}

func (c *tracingClient) GetAuthenticatedRepositoryURL(ctx context.Context, accesstoken AccessToken, htmlURL string) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetAuthenticatedRepositoryURL"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetAuthenticatedRepositoryURL(ctx, accesstoken, htmlURL)
}

func (c *tracingClient) GetEstafetteManifest(ctx context.Context, accesstoken AccessToken, event RepositoryPushEvent) (valid bool, manifest string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetEstafetteManifest"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetEstafetteManifest(ctx, accesstoken, event)
}

func (c *tracingClient) JobVarsFunc(ctx context.Context) func(ctx context.Context, repoSource, repoOwner, repoName string) (token string, url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "JobVarsFunc"))
	defer func() { api.FinishSpan(span) }()

	return c.Client.JobVarsFunc(ctx)
}
