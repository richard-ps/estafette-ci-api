package cloudsourceapi

import (
	"context"

	"github.com/estafette/estafette-ci-api/helpers"
	"github.com/opentracing/opentracing-go"
)

// NewTracingClient returns a new instance of a tracing Client.
func NewTracingClient(c Client) Client {
	return &tracingClient{c, "cloudsourceapi"}
}

type tracingClient struct {
	Client
	prefix string
}

func (c *tracingClient) GetAccessToken(ctx context.Context) (accesstoken AccessToken, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "GetAccessToken"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.GetAccessToken(ctx)
}

func (c *tracingClient) GetAuthenticatedRepositoryURL(ctx context.Context, accesstoken AccessToken, htmlURL string) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "GetAuthenticatedRepositoryURL"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.GetAuthenticatedRepositoryURL(ctx, accesstoken, htmlURL)
}

func (c *tracingClient) GetEstafetteManifest(ctx context.Context, accesstoken AccessToken, notification PubSubNotification) (valid bool, manifest string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "GetEstafetteManifest"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.GetEstafetteManifest(ctx, accesstoken, notification)
}

func (c *tracingClient) JobVarsFunc(ctx context.Context) func(context.Context, string, string, string) (string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "JobVarsFunc"))
	defer func() { helpers.FinishSpan(span) }()

	return c.Client.JobVarsFunc(ctx)
}
