package builderapi

import (
	"context"

	"github.com/estafette/estafette-ci-api/api"
	contracts "github.com/estafette/estafette-ci-contracts"
	batchv1 "k8s.io/api/batch/v1"
)

// NewLoggingClient returns a new instance of a logging Client.
func NewLoggingClient(c Client) Client {
	return &loggingClient{c, "builderapi"}
}

type loggingClient struct {
	Client
	prefix string
}

func (c *loggingClient) CreateCiBuilderJob(ctx context.Context, params CiBuilderParams) (job *batchv1.Job, err error) {
	defer func() { api.HandleLogError(c.prefix, "CreateCiBuilderJob", err) }()

	return c.Client.CreateCiBuilderJob(ctx, params)
}

func (c *loggingClient) RemoveCiBuilderJob(ctx context.Context, jobName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RemoveCiBuilderJob", err) }()

	return c.Client.RemoveCiBuilderJob(ctx, jobName)
}

func (c *loggingClient) CancelCiBuilderJob(ctx context.Context, jobName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "CancelCiBuilderJob", err) }()

	return c.Client.CancelCiBuilderJob(ctx, jobName)
}

func (c *loggingClient) RemoveCiBuilderConfigMap(ctx context.Context, configmapName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RemoveCiBuilderConfigMap", err) }()

	return c.Client.RemoveCiBuilderConfigMap(ctx, configmapName)
}

func (c *loggingClient) RemoveCiBuilderSecret(ctx context.Context, secretName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RemoveCiBuilderSecret", err) }()

	return c.Client.RemoveCiBuilderSecret(ctx, secretName)
}

func (c *loggingClient) TailCiBuilderJobLogs(ctx context.Context, jobName string, logChannel chan contracts.TailLogLine) (err error) {
	defer func() { api.HandleLogError(c.prefix, "TailCiBuilderJobLogs", err) }()

	return c.Client.TailCiBuilderJobLogs(ctx, jobName, logChannel)
}

func (c *loggingClient) GetJobName(ctx context.Context, jobType, repoOwner, repoName, id string) string {
	return c.Client.GetJobName(ctx, jobType, repoOwner, repoName, id)
}
