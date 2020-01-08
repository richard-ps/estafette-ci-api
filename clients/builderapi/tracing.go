package builderapi

import (
	"context"

	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	"github.com/estafette/estafette-ci-api/helpers"
	contracts "github.com/estafette/estafette-ci-contracts"
	"github.com/opentracing/opentracing-go"
)

// NewTracingClient returns a new instance of a tracing Client.
func NewTracingClient(c Client) Client {
	return &tracingClient{c, "builderapi"}
}

type tracingClient struct {
	Client
	prefix string
}

func (c *tracingClient) CreateCiBuilderJob(ctx context.Context, params CiBuilderParams) (job *batchv1.Job, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "CreateCiBuilderJob"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.CreateCiBuilderJob(ctx, params)
}

func (c *tracingClient) RemoveCiBuilderJob(ctx context.Context, jobName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "RemoveCiBuilderJob"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.RemoveCiBuilderJob(ctx, jobName)
}

func (c *tracingClient) CancelCiBuilderJob(ctx context.Context, jobName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "CancelCiBuilderJob"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.CancelCiBuilderJob(ctx, jobName)
}

func (c *tracingClient) RemoveCiBuilderConfigMap(ctx context.Context, configmapName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "RemoveCiBuilderConfigMap"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.RemoveCiBuilderConfigMap(ctx, configmapName)
}

func (c *tracingClient) RemoveCiBuilderSecret(ctx context.Context, secretName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "RemoveCiBuilderSecret"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.RemoveCiBuilderSecret(ctx, secretName)
}

func (c *tracingClient) TailCiBuilderJobLogs(ctx context.Context, jobName string, logChannel chan contracts.TailLogLine) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "TailCiBuilderJobLogs"))
	defer func() { helpers.FinishSpanWithError(span, err) }()

	return c.Client.TailCiBuilderJobLogs(ctx, jobName, logChannel)
}

func (c *tracingClient) GetJobName(ctx context.Context, jobType, repoOwner, repoName, id string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "GetJobName"))
	defer func() { helpers.FinishSpan(span) }()

	return c.Client.GetJobName(ctx, jobType, repoOwner, repoOwner, id)
}

func (c *tracingClient) GetBuilderConfig(ctx context.Context, params CiBuilderParams, jobName string) contracts.BuilderConfig {
	span, ctx := opentracing.StartSpanFromContext(ctx, helpers.GetSpanName(c.prefix, "GetBuilderConfig"))
	defer func() { helpers.FinishSpan(span) }()

	return c.Client.GetBuilderConfig(ctx, params, jobName)
}