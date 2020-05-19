package estafette

import (
	"context"
	"time"

	"github.com/estafette/estafette-ci-api/clients/builderapi"
	"github.com/estafette/estafette-ci-api/helpers"
	contracts "github.com/estafette/estafette-ci-contracts"
	manifest "github.com/estafette/estafette-ci-manifest"
	"github.com/go-kit/kit/metrics"
)

// NewMetricsService returns a new instance of a metrics Service.
func NewMetricsService(s Service, requestCount metrics.Counter, requestLatency metrics.Histogram) Service {
	return &metricsService{s, requestCount, requestLatency}
}

type metricsService struct {
	Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (s *metricsService) CreateBuild(ctx context.Context, build contracts.Build, waitForJobToStart bool) (b *contracts.Build, err error) {
	defer func(begin time.Time) { helpers.UpdateMetrics(s.requestCount, s.requestLatency, "CreateBuild", begin) }(time.Now())

	return s.Service.CreateBuild(ctx, build, waitForJobToStart)
}

func (s *metricsService) FinishBuild(ctx context.Context, repoSource, repoOwner, repoName string, buildID int, buildStatus string) (err error) {
	defer func(begin time.Time) { helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FinishBuild", begin) }(time.Now())

	return s.Service.FinishBuild(ctx, repoSource, repoOwner, repoName, buildID, buildStatus)
}

func (s *metricsService) CreateRelease(ctx context.Context, release contracts.Release, mft manifest.EstafetteManifest, repoBranch, repoRevision string, waitForJobToStart bool) (r *contracts.Release, err error) {
	defer func(begin time.Time) { helpers.UpdateMetrics(s.requestCount, s.requestLatency, "CreateRelease", begin) }(time.Now())

	return s.Service.CreateRelease(ctx, release, mft, repoBranch, repoRevision, waitForJobToStart)
}

func (s *metricsService) FinishRelease(ctx context.Context, repoSource, repoOwner, repoName string, releaseID int, releaseStatus string) (err error) {
	defer func(begin time.Time) { helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FinishRelease", begin) }(time.Now())

	return s.Service.FinishRelease(ctx, repoSource, repoOwner, repoName, releaseID, releaseStatus)
}

func (s *metricsService) FireGitTriggers(ctx context.Context, gitEvent manifest.EstafetteGitEvent) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FireGitTriggers", begin)
	}(time.Now())

	return s.Service.FireGitTriggers(ctx, gitEvent)
}

func (s *metricsService) FirePipelineTriggers(ctx context.Context, build contracts.Build, event string) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FirePipelineTriggers", begin)
	}(time.Now())

	return s.Service.FirePipelineTriggers(ctx, build, event)
}

func (s *metricsService) FireReleaseTriggers(ctx context.Context, release contracts.Release, event string) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FireReleaseTriggers", begin)
	}(time.Now())

	return s.Service.FireReleaseTriggers(ctx, release, event)
}

func (s *metricsService) FirePubSubTriggers(ctx context.Context, pubsubEvent manifest.EstafettePubSubEvent) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FirePubSubTriggers", begin)
	}(time.Now())

	return s.Service.FirePubSubTriggers(ctx, pubsubEvent)
}

func (s *metricsService) FireCronTriggers(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "FireCronTriggers", begin)
	}(time.Now())

	return s.Service.FireCronTriggers(ctx)
}

func (s *metricsService) Rename(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func(begin time.Time) { helpers.UpdateMetrics(s.requestCount, s.requestLatency, "Rename", begin) }(time.Now())

	return s.Service.Rename(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (s *metricsService) UpdateBuildStatus(ctx context.Context, event builderapi.CiBuilderEvent) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "UpdateBuildStatus", begin)
	}(time.Now())

	return s.Service.UpdateBuildStatus(ctx, event)
}

func (s *metricsService) UpdateJobResources(ctx context.Context, event builderapi.CiBuilderEvent) (err error) {
	defer func(begin time.Time) {
		helpers.UpdateMetrics(s.requestCount, s.requestLatency, "UpdateJobResources", begin)
	}(time.Now())

	return s.Service.UpdateJobResources(ctx, event)
}
