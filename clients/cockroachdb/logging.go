package cockroachdb

import (
	"context"
	"time"

	"github.com/estafette/estafette-ci-api/api"
	contracts "github.com/estafette/estafette-ci-contracts"
	manifest "github.com/estafette/estafette-ci-manifest"
)

// NewLoggingClient returns a new instance of a logging Client.
func NewLoggingClient(c Client) Client {
	return &loggingClient{c, "cockroachdb"}
}

type loggingClient struct {
	Client
	prefix string
}

func (c *loggingClient) Connect(ctx context.Context) (err error) {
	defer func() { api.HandleLogError(c.prefix, "Connect", err) }()

	return c.Client.Connect(ctx)
}

func (c *loggingClient) ConnectWithDriverAndSource(ctx context.Context, driverName, dataSourceName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "ConnectWithDriverAndSource", err) }()

	return c.Client.ConnectWithDriverAndSource(ctx, driverName, dataSourceName)
}

func (c *loggingClient) GetAutoIncrement(ctx context.Context, shortRepoSource, repoOwner, repoName string) (autoincrement int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetAutoIncrement", err) }()

	return c.Client.GetAutoIncrement(ctx, shortRepoSource, repoOwner, repoName)
}

func (c *loggingClient) InsertBuild(ctx context.Context, build contracts.Build, jobResources JobResources) (b *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertBuild", err) }()

	return c.Client.InsertBuild(ctx, build, jobResources)
}

func (c *loggingClient) UpdateBuildStatus(ctx context.Context, repoSource, repoOwner, repoName string, buildID int, buildStatus string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateBuildStatus", err) }()

	return c.Client.UpdateBuildStatus(ctx, repoSource, repoOwner, repoName, buildID, buildStatus)
}

func (c *loggingClient) UpdateBuildResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, buildID int, jobResources JobResources) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateBuildResourceUtilization", err) }()

	return c.Client.UpdateBuildResourceUtilization(ctx, repoSource, repoOwner, repoName, buildID, jobResources)
}

func (c *loggingClient) InsertRelease(ctx context.Context, release contracts.Release, jobResources JobResources) (r *contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertRelease", err) }()

	return c.Client.InsertRelease(ctx, release, jobResources)
}

func (c *loggingClient) UpdateReleaseStatus(ctx context.Context, repoSource, repoOwner, repoName string, id int, releaseStatus string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateReleaseStatus", err) }()

	return c.Client.UpdateReleaseStatus(ctx, repoSource, repoOwner, repoName, id, releaseStatus)
}

func (c *loggingClient) UpdateReleaseResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, id int, jobResources JobResources) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateReleaseResourceUtilization", err) }()

	return c.Client.UpdateReleaseResourceUtilization(ctx, repoSource, repoOwner, repoName, id, jobResources)
}

func (c *loggingClient) InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog, writeLogToDatabase bool) (buildlog contracts.BuildLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertBuildLog", err) }()

	return c.Client.InsertBuildLog(ctx, buildLog, writeLogToDatabase)
}

func (c *loggingClient) InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog, writeLogToDatabase bool) (releaselog contracts.ReleaseLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertReleaseLog", err) }()

	return c.Client.InsertReleaseLog(ctx, releaseLog, writeLogToDatabase)
}

func (c *loggingClient) UpsertComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpsertComputedPipeline", err) }()

	return c.Client.UpsertComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *loggingClient) UpdateComputedPipelinePermissions(ctx context.Context, pipeline contracts.Pipeline) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateComputedPipelinePermissions", err) }()

	return c.Client.UpdateComputedPipelinePermissions(ctx, pipeline)
}

func (c *loggingClient) UpdateComputedPipelineFirstInsertedAt(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateComputedPipelineFirstInsertedAt", err) }()

	return c.Client.UpdateComputedPipelineFirstInsertedAt(ctx, repoSource, repoOwner, repoName)
}

func (c *loggingClient) UpsertComputedRelease(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpsertComputedRelease", err) }()

	return c.Client.UpsertComputedRelease(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *loggingClient) UpdateComputedReleaseFirstInsertedAt(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateComputedReleaseFirstInsertedAt", err) }()

	return c.Client.UpdateComputedReleaseFirstInsertedAt(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *loggingClient) ArchiveComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "ArchiveComputedPipeline", err) }()

	return c.Client.ArchiveComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *loggingClient) UnarchiveComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UnarchiveComputedPipeline", err) }()

	return c.Client.UnarchiveComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *loggingClient) GetPipelines(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField, optimized bool) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelines", err) }()

	return c.Client.GetPipelines(ctx, pageNumber, pageSize, filters, sortings, optimized)
}

func (c *loggingClient) GetPipelinesByRepoName(ctx context.Context, repoName string, optimized bool) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesByRepoName", err) }()

	return c.Client.GetPipelinesByRepoName(ctx, repoName, optimized)
}

func (c *loggingClient) GetPipelinesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesCount", err) }()

	return c.Client.GetPipelinesCount(ctx, filters)
}

func (c *loggingClient) GetPipeline(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string, optimized bool) (pipeline *contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipeline", err) }()

	return c.Client.GetPipeline(ctx, repoSource, repoOwner, repoName, filters, optimized)
}

func (c *loggingClient) GetPipelineRecentBuilds(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (builds []*contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineRecentBuilds", err) }()

	return c.Client.GetPipelineRecentBuilds(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *loggingClient) GetPipelineBuilds(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField, optimized bool) (builds []*contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuilds", err) }()

	return c.Client.GetPipelineBuilds(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize, filters, sortings, optimized)
}

func (c *loggingClient) GetPipelineBuildsCount(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildsCount", err) }()

	return c.Client.GetPipelineBuildsCount(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName, repoRevision string, optimized bool) (build *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuild", err) }()

	return c.Client.GetPipelineBuild(ctx, repoSource, repoOwner, repoName, repoRevision, optimized)
}

func (c *loggingClient) GetPipelineBuildByID(ctx context.Context, repoSource, repoOwner, repoName string, id int, optimized bool) (build *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildByID", err) }()

	return c.Client.GetPipelineBuildByID(ctx, repoSource, repoOwner, repoName, id, optimized)
}

func (c *loggingClient) GetLastPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (build *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetLastPipelineBuild", err) }()

	return c.Client.GetLastPipelineBuild(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *loggingClient) GetFirstPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (build *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFirstPipelineBuild", err) }()

	return c.Client.GetFirstPipelineBuild(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *loggingClient) GetLastPipelineBuildForBranch(ctx context.Context, repoSource, repoOwner, repoName, branch string) (build *contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetLastPipelineBuildForBranch", err) }()

	return c.Client.GetLastPipelineBuildForBranch(ctx, repoSource, repoOwner, repoName, branch)
}

func (c *loggingClient) GetLastPipelineReleases(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string, pageSize int) (releases []*contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetLastPipelineReleases", err) }()

	return c.Client.GetLastPipelineReleases(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction, pageSize)
}

func (c *loggingClient) GetFirstPipelineRelease(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (release *contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFirstPipelineRelease", err) }()

	return c.Client.GetFirstPipelineRelease(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *loggingClient) GetPipelineBuildsByVersion(ctx context.Context, repoSource, repoOwner, repoName, buildVersion string, statuses []string, limit uint64, optimized bool) (builds []*contracts.Build, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildsByVersion", err) }()

	return c.Client.GetPipelineBuildsByVersion(ctx, repoSource, repoOwner, repoName, buildVersion, statuses, limit, optimized)
}

func (c *loggingClient) GetPipelineBuildLogs(ctx context.Context, repoSource, repoOwner, repoName, repoBranch, repoRevision, buildID string, readLogFromDatabase bool) (buildlog *contracts.BuildLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildLogs", err) }()

	return c.Client.GetPipelineBuildLogs(ctx, repoSource, repoOwner, repoName, repoBranch, repoRevision, buildID, readLogFromDatabase)
}

func (c *loggingClient) GetPipelineBuildLogsPerPage(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber int, pageSize int) (buildlogs []*contracts.BuildLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildLogsPerPage", err) }()

	return c.Client.GetPipelineBuildLogsPerPage(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize)
}

func (c *loggingClient) GetPipelineBuildMaxResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, lastNRecords int) (jobresources JobResources, count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildMaxResourceUtilization", err) }()

	return c.Client.GetPipelineBuildMaxResourceUtilization(ctx, repoSource, repoOwner, repoName, lastNRecords)
}

func (c *loggingClient) GetPipelineReleases(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (releases []*contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleases", err) }()

	return c.Client.GetPipelineReleases(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetPipelineReleasesCount(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleasesCount", err) }()

	return c.Client.GetPipelineReleasesCount(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineRelease(ctx context.Context, repoSource, repoOwner, repoName string, id int) (release *contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineRelease", err) }()

	return c.Client.GetPipelineRelease(ctx, repoSource, repoOwner, repoName, id)
}

func (c *loggingClient) GetPipelineLastReleasesByName(ctx context.Context, repoSource, repoOwner, repoName, releaseName string, actions []string) (releases []contracts.Release, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineLastReleasesByName", err) }()

	return c.Client.GetPipelineLastReleasesByName(ctx, repoSource, repoOwner, repoName, releaseName, actions)
}

func (c *loggingClient) GetPipelineReleaseLogs(ctx context.Context, repoSource, repoOwner, repoName string, id int, readLogFromDatabase bool) (releaselog *contracts.ReleaseLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleaseLogs", err) }()

	return c.Client.GetPipelineReleaseLogs(ctx, repoSource, repoOwner, repoName, id, readLogFromDatabase)
}

func (c *loggingClient) GetPipelineReleaseLogsPerPage(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber int, pageSize int) (releaselogs []*contracts.ReleaseLog, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleaseLogsPerPage", err) }()

	return c.Client.GetPipelineReleaseLogsPerPage(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize)
}

func (c *loggingClient) GetPipelineReleaseMaxResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName, targetName string, lastNRecords int) (jobresources JobResources, count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleaseMaxResourceUtilization", err) }()

	return c.Client.GetPipelineReleaseMaxResourceUtilization(ctx, repoSource, repoOwner, repoName, targetName, lastNRecords)
}

func (c *loggingClient) GetBuildsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetBuildsCount", err) }()

	return c.Client.GetBuildsCount(ctx, filters)
}

func (c *loggingClient) GetReleasesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetReleasesCount", err) }()

	return c.Client.GetReleasesCount(ctx, filters)
}

func (c *loggingClient) GetBuildsDuration(ctx context.Context, filters map[api.FilterType][]string) (duration time.Duration, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetBuildsDuration", err) }()

	return c.Client.GetBuildsDuration(ctx, filters)
}

func (c *loggingClient) GetFirstBuildTimes(ctx context.Context) (times []time.Time, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFirstBuildTimes", err) }()

	return c.Client.GetFirstBuildTimes(ctx)
}

func (c *loggingClient) GetFirstReleaseTimes(ctx context.Context) (times []time.Time, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFirstReleaseTimes", err) }()

	return c.Client.GetFirstReleaseTimes(ctx)
}

func (c *loggingClient) GetPipelineBuildsDurations(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (durations []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildsDurations", err) }()

	return c.Client.GetPipelineBuildsDurations(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineReleasesDurations(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (durations []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleasesDurations", err) }()

	return c.Client.GetPipelineReleasesDurations(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineBuildsCPUUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildsCPUUsageMeasurements", err) }()

	return c.Client.GetPipelineBuildsCPUUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineReleasesCPUUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleasesCPUUsageMeasurements", err) }()

	return c.Client.GetPipelineReleasesCPUUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineBuildsMemoryUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineBuildsMemoryUsageMeasurements", err) }()

	return c.Client.GetPipelineBuildsMemoryUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetPipelineReleasesMemoryUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineReleasesMemoryUsageMeasurements", err) }()

	return c.Client.GetPipelineReleasesMemoryUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *loggingClient) GetLabelValues(ctx context.Context, labelKey string) (labels []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetLabelValues", err) }()

	return c.Client.GetLabelValues(ctx, labelKey)
}

func (c *loggingClient) GetFrequentLabels(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFrequentLabels", err) }()

	return c.Client.GetFrequentLabels(ctx, pageNumber, pageSize, filters)
}

func (c *loggingClient) GetFrequentLabelsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetFrequentLabelsCount", err) }()

	return c.Client.GetFrequentLabelsCount(ctx, filters)
}

func (c *loggingClient) GetPipelinesWithMostBuilds(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (pipelines []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesWithMostBuilds", err) }()

	return c.Client.GetPipelinesWithMostBuilds(ctx, pageNumber, pageSize, filters)
}

func (c *loggingClient) GetPipelinesWithMostBuildsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesWithMostBuildsCount", err) }()

	return c.Client.GetPipelinesWithMostBuildsCount(ctx, filters)
}

func (c *loggingClient) GetPipelinesWithMostReleases(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (pipelines []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesWithMostReleases", err) }()

	return c.Client.GetPipelinesWithMostReleases(ctx, pageNumber, pageSize, filters)
}

func (c *loggingClient) GetPipelinesWithMostReleasesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelinesWithMostReleasesCount", err) }()

	return c.Client.GetPipelinesWithMostReleasesCount(ctx, filters)
}

func (c *loggingClient) GetTriggers(ctx context.Context, triggerType, identifier, event string) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetTriggers", err) }()

	return c.Client.GetTriggers(ctx, triggerType, identifier, event)
}

func (c *loggingClient) GetGitTriggers(ctx context.Context, gitEvent manifest.EstafetteGitEvent) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGitTriggers", err) }()

	return c.Client.GetGitTriggers(ctx, gitEvent)
}

func (c *loggingClient) GetPipelineTriggers(ctx context.Context, build contracts.Build, event string) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPipelineTriggers", err) }()

	return c.Client.GetPipelineTriggers(ctx, build, event)
}

func (c *loggingClient) GetReleaseTriggers(ctx context.Context, release contracts.Release, event string) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetReleaseTriggers", err) }()

	return c.Client.GetReleaseTriggers(ctx, release, event)
}

func (c *loggingClient) GetPubSubTriggers(ctx context.Context, pubsubEvent manifest.EstafettePubSubEvent) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetPubSubTriggers", err) }()

	return c.Client.GetPubSubTriggers(ctx, pubsubEvent)
}

func (c *loggingClient) GetCronTriggers(ctx context.Context) (pipelines []*contracts.Pipeline, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCronTriggers", err) }()

	return c.Client.GetCronTriggers(ctx)
}

func (c *loggingClient) Rename(ctx context.Context, shortFromRepoSource, fromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "Rename", err) }()

	return c.Client.Rename(ctx, shortFromRepoSource, fromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameBuildVersion(ctx context.Context, shortFromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameBuildVersion", err) }()

	return c.Client.RenameBuildVersion(ctx, shortFromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameBuilds(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameBuilds", err) }()

	return c.Client.RenameBuilds(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameBuildLogs(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameBuildLogs", err) }()

	return c.Client.RenameBuildLogs(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameReleases(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameReleases", err) }()

	return c.Client.RenameReleases(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameReleaseLogs(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameReleaseLogs", err) }()

	return c.Client.RenameReleaseLogs(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameComputedPipelines(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameComputedPipelines", err) }()

	return c.Client.RenameComputedPipelines(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) RenameComputedReleases(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "RenameComputedReleases", err) }()

	return c.Client.RenameComputedReleases(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *loggingClient) InsertUser(ctx context.Context, user contracts.User) (u *contracts.User, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertUser", err) }()

	return c.Client.InsertUser(ctx, user)
}

func (c *loggingClient) UpdateUser(ctx context.Context, user contracts.User) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateUser", err) }()

	return c.Client.UpdateUser(ctx, user)
}

func (c *loggingClient) DeleteUser(ctx context.Context, user contracts.User) (err error) {
	defer func() { api.HandleLogError(c.prefix, "DeleteUser", err) }()

	return c.Client.DeleteUser(ctx, user)
}

func (c *loggingClient) GetUserByIdentity(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetUserByIdentity", err, ErrUserNotFound) }()

	return c.Client.GetUserByIdentity(ctx, identity)
}

func (c *loggingClient) GetUserByID(ctx context.Context, id string, filters map[api.FilterType][]string) (user *contracts.User, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetUserByID", err, ErrUserNotFound) }()

	return c.Client.GetUserByID(ctx, id, filters)
}

func (c *loggingClient) GetUsers(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (users []*contracts.User, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetUsers", err) }()

	return c.Client.GetUsers(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetUsersCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetUsersCount", err) }()

	return c.Client.GetUsersCount(ctx, filters)
}

func (c *loggingClient) InsertGroup(ctx context.Context, group contracts.Group) (g *contracts.Group, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertGroup", err) }()

	return c.Client.InsertGroup(ctx, group)
}

func (c *loggingClient) UpdateGroup(ctx context.Context, group contracts.Group) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateGroup", err) }()

	return c.Client.UpdateGroup(ctx, group)
}

func (c *loggingClient) DeleteGroup(ctx context.Context, group contracts.Group) (err error) {
	defer func() { api.HandleLogError(c.prefix, "DeleteGroup", err) }()

	return c.Client.DeleteGroup(ctx, group)
}

func (c *loggingClient) GetGroupByIdentity(ctx context.Context, identity contracts.GroupIdentity) (group *contracts.Group, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroupByIdentity", err, ErrGroupNotFound) }()

	return c.Client.GetGroupByIdentity(ctx, identity)
}

func (c *loggingClient) GetGroupByID(ctx context.Context, id string, filters map[api.FilterType][]string) (group *contracts.Group, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroupByID", err, ErrGroupNotFound) }()

	return c.Client.GetGroupByID(ctx, id, filters)
}

func (c *loggingClient) GetGroups(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (groups []*contracts.Group, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroups", err) }()

	return c.Client.GetGroups(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetGroupsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroupsCount", err) }()

	return c.Client.GetGroupsCount(ctx, filters)
}

func (c *loggingClient) InsertOrganization(ctx context.Context, organization contracts.Organization) (o *contracts.Organization, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertOrganization", err) }()

	return c.Client.InsertOrganization(ctx, organization)
}

func (c *loggingClient) UpdateOrganization(ctx context.Context, organization contracts.Organization) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateOrganization", err) }()

	return c.Client.UpdateOrganization(ctx, organization)
}

func (c *loggingClient) DeleteOrganization(ctx context.Context, organization contracts.Organization) (err error) {
	defer func() { api.HandleLogError(c.prefix, "DeleteOrganization", err) }()

	return c.Client.DeleteOrganization(ctx, organization)
}

func (c *loggingClient) GetOrganizationByIdentity(ctx context.Context, identity contracts.OrganizationIdentity) (organization *contracts.Organization, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetOrganizationByIdentity", err, ErrOrganizationNotFound) }()

	return c.Client.GetOrganizationByIdentity(ctx, identity)
}

func (c *loggingClient) GetOrganizationByID(ctx context.Context, id string) (organization *contracts.Organization, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetOrganizationByID", err, ErrOrganizationNotFound) }()

	return c.Client.GetOrganizationByID(ctx, id)
}

func (c *loggingClient) GetOrganizationByName(ctx context.Context, name string) (organization *contracts.Organization, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetOrganizationByName", err, ErrOrganizationNotFound) }()

	return c.Client.GetOrganizationByName(ctx, name)
}

func (c *loggingClient) GetOrganizations(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (organizations []*contracts.Organization, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroups", err) }()

	return c.Client.GetOrganizations(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetOrganizationsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetGroupsCount", err) }()

	return c.Client.GetOrganizationsCount(ctx, filters)
}

func (c *loggingClient) InsertClient(ctx context.Context, client contracts.Client) (cl *contracts.Client, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertClient", err) }()

	return c.Client.InsertClient(ctx, client)
}

func (c *loggingClient) UpdateClient(ctx context.Context, client contracts.Client) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateClient", err) }()

	return c.Client.UpdateClient(ctx, client)
}

func (c *loggingClient) DeleteClient(ctx context.Context, client contracts.Client) (err error) {
	defer func() { api.HandleLogError(c.prefix, "DeleteClient", err) }()

	return c.Client.DeleteClient(ctx, client)
}

func (c *loggingClient) GetClientByClientID(ctx context.Context, clientID string) (client *contracts.Client, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetClientByClientID", err, ErrClientNotFound) }()

	return c.Client.GetClientByClientID(ctx, clientID)
}

func (c *loggingClient) GetClientByID(ctx context.Context, id string) (client *contracts.Client, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetClientByID", err, ErrClientNotFound) }()

	return c.Client.GetClientByID(ctx, id)
}

func (c *loggingClient) GetClients(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (clients []*contracts.Client, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetClients", err) }()

	return c.Client.GetClients(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetClientsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetClientsCount", err) }()

	return c.Client.GetClientsCount(ctx, filters)
}

func (c *loggingClient) InsertCatalogEntity(ctx context.Context, catalogEntity contracts.CatalogEntity) (insertedCatalogEntity *contracts.CatalogEntity, err error) {
	defer func() { api.HandleLogError(c.prefix, "InsertCatalogEntity", err) }()

	return c.Client.InsertCatalogEntity(ctx, catalogEntity)
}

func (c *loggingClient) UpdateCatalogEntity(ctx context.Context, catalogEntity contracts.CatalogEntity) (err error) {
	defer func() { api.HandleLogError(c.prefix, "UpdateCatalogEntity", err) }()

	return c.Client.UpdateCatalogEntity(ctx, catalogEntity)
}

func (c *loggingClient) DeleteCatalogEntity(ctx context.Context, id string) (err error) {
	defer func() { api.HandleLogError(c.prefix, "DeleteCatalogEntity", err) }()

	return c.Client.DeleteCatalogEntity(ctx, id)
}

func (c *loggingClient) GetCatalogEntityByID(ctx context.Context, id string) (catalogEntity *contracts.CatalogEntity, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityByID", err) }()

	return c.Client.GetCatalogEntityByID(ctx, id)
}

func (c *loggingClient) GetCatalogEntities(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (catalogEntities []*contracts.CatalogEntity, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntities", err) }()

	return c.Client.GetCatalogEntities(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetCatalogEntitiesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntitiesCount", err) }()

	return c.Client.GetCatalogEntitiesCount(ctx, filters)
}

func (c *loggingClient) GetCatalogEntityParentKeys(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (keys []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityParentKeys", err) }()

	return c.Client.GetCatalogEntityParentKeys(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetCatalogEntityParentKeysCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityParentKeysCount", err) }()

	return c.Client.GetCatalogEntityParentKeysCount(ctx, filters)
}

func (c *loggingClient) GetCatalogEntityParentValues(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (values []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityParentValues", err) }()

	return c.Client.GetCatalogEntityParentValues(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetCatalogEntityParentValuesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityParentValuesCount", err) }()

	return c.Client.GetCatalogEntityParentValuesCount(ctx, filters)
}

func (c *loggingClient) GetCatalogEntityKeys(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (keys []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityKeys", err) }()

	return c.Client.GetCatalogEntityKeys(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetCatalogEntityKeysCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityKeysCount", err) }()

	return c.Client.GetCatalogEntityKeysCount(ctx, filters)
}

func (c *loggingClient) GetCatalogEntityValues(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (values []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityValues", err) }()

	return c.Client.GetCatalogEntityValues(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *loggingClient) GetCatalogEntityValuesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityValuesCount", err) }()

	return c.Client.GetCatalogEntityValuesCount(ctx, filters)
}

func (c *loggingClient) GetCatalogEntityLabels(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (labels []map[string]interface{}, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityLabels", err) }()

	return c.Client.GetCatalogEntityLabels(ctx, pageNumber, pageSize, filters)
}

func (c *loggingClient) GetCatalogEntityLabelsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	defer func() { api.HandleLogError(c.prefix, "GetCatalogEntityLabelsCount", err) }()

	return c.Client.GetCatalogEntityLabelsCount(ctx, filters)
}
