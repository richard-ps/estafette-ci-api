package slackapi

import (
	"context"
	"time"

	"github.com/estafette/estafette-ci-api/config"
	"github.com/estafette/estafette-ci-api/helpers"
	"github.com/go-kit/kit/metrics"
)

// NewMetricsClient returns a new instance of a metrics Client.
func NewMetricsClient(c Client, requestCount metrics.Counter, requestLatency metrics.Histogram) Client {
	return &metricsClient{c, requestCount, requestLatency}
}

type metricsClient struct {
	Client
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (c *metricsClient) GetUserProfile(ctx context.Context, userID string) (profile *UserProfile, err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(c.requestCount, c.requestLatency, "GetUserProfile", begin)
	}(time.Now())

	return c.Client.GetUserProfile(ctx, userID)
}

func (c *metricsClient) RefreshConfig(config *config.APIConfig) {
	c.Client.RefreshConfig(config)
}
