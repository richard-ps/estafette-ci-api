package cockroachdb

import (
	"context"
	"time"

	"github.com/estafette/estafette-ci-api/api"
	contracts "github.com/estafette/estafette-ci-contracts"
	manifest "github.com/estafette/estafette-ci-manifest"
	"github.com/opentracing/opentracing-go"
)

// NewTracingClient returns a new instance of a tracing Client.
func NewTracingClient(c Client) Client {
	return &tracingClient{c, "cockroachdb"}
}

type tracingClient struct {
	Client
	prefix string
}

func (c *tracingClient) Connect(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "Connect"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.Connect(ctx)
}

func (c *tracingClient) ConnectWithDriverAndSource(ctx context.Context, driverName, dataSourceName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "ConnectWithDriverAndSource"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.ConnectWithDriverAndSource(ctx, driverName, dataSourceName)
}

func (c *tracingClient) GetAutoIncrement(ctx context.Context, shortRepoSource, repoOwner, repoName string) (autoincrement int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetAutoIncrement"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetAutoIncrement(ctx, shortRepoSource, repoOwner, repoName)
}

func (c *tracingClient) InsertBuild(ctx context.Context, build contracts.Build, jobResources JobResources) (b *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertBuild"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertBuild(ctx, build, jobResources)
}

func (c *tracingClient) UpdateBuildStatus(ctx context.Context, repoSource, repoOwner, repoName string, buildID int, buildStatus string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateBuildStatus"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateBuildStatus(ctx, repoSource, repoOwner, repoName, buildID, buildStatus)
}

func (c *tracingClient) UpdateBuildResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, buildID int, jobResources JobResources) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateBuildResourceUtilization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateBuildResourceUtilization(ctx, repoSource, repoOwner, repoName, buildID, jobResources)
}

func (c *tracingClient) InsertRelease(ctx context.Context, release contracts.Release, jobResources JobResources) (r *contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertRelease"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertRelease(ctx, release, jobResources)
}

func (c *tracingClient) UpdateReleaseStatus(ctx context.Context, repoSource, repoOwner, repoName string, id int, releaseStatus string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateReleaseStatus"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateReleaseStatus(ctx, repoSource, repoOwner, repoName, id, releaseStatus)
}

func (c *tracingClient) UpdateReleaseResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, id int, jobResources JobResources) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateReleaseResourceUtilization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateReleaseResourceUtilization(ctx, repoSource, repoOwner, repoName, id, jobResources)
}

func (c *tracingClient) InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog, writeLogToDatabase bool) (buildlog contracts.BuildLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertBuildLog"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertBuildLog(ctx, buildLog, writeLogToDatabase)
}

func (c *tracingClient) InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog, writeLogToDatabase bool) (releaselog contracts.ReleaseLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertReleaseLog"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertReleaseLog(ctx, releaseLog, writeLogToDatabase)
}

func (c *tracingClient) UpsertComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpsertComputedPipeline"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpsertComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *tracingClient) UpdateComputedPipelinePermissions(ctx context.Context, pipeline contracts.Pipeline) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateComputedPipelinePermissions"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateComputedPipelinePermissions(ctx, pipeline)
}

func (c *tracingClient) UpdateComputedPipelineFirstInsertedAt(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateComputedPipelineFirstInsertedAt"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateComputedPipelineFirstInsertedAt(ctx, repoSource, repoOwner, repoName)
}

func (c *tracingClient) UpsertComputedRelease(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpsertComputedRelease"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpsertComputedRelease(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *tracingClient) UpdateComputedReleaseFirstInsertedAt(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateComputedReleaseFirstInsertedAt"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateComputedReleaseFirstInsertedAt(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *tracingClient) ArchiveComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "ArchiveComputedPipeline"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.ArchiveComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *tracingClient) UnarchiveComputedPipeline(ctx context.Context, repoSource, repoOwner, repoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UnarchiveComputedPipeline"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UnarchiveComputedPipeline(ctx, repoSource, repoOwner, repoName)
}

func (c *tracingClient) GetPipelines(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField, optimized bool) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelines"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelines(ctx, pageNumber, pageSize, filters, sortings, optimized)
}

func (c *tracingClient) GetPipelinesByRepoName(ctx context.Context, repoName string, optimized bool) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesByRepoName"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesByRepoName(ctx, repoName, optimized)
}

func (c *tracingClient) GetPipelinesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesCount(ctx, filters)
}

func (c *tracingClient) GetPipeline(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string, optimized bool) (pipeline *contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipeline"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipeline(ctx, repoSource, repoOwner, repoName, filters, optimized)
}

func (c *tracingClient) GetPipelineRecentBuilds(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (builds []*contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineRecentBuilds"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineRecentBuilds(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *tracingClient) GetPipelineBuilds(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField, optimized bool) (builds []*contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuilds"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuilds(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize, filters, sortings, optimized)
}

func (c *tracingClient) GetPipelineBuildsCount(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildsCount(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName, repoRevision string, optimized bool) (build *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuild"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuild(ctx, repoSource, repoOwner, repoName, repoRevision, optimized)
}

func (c *tracingClient) GetPipelineBuildByID(ctx context.Context, repoSource, repoOwner, repoName string, id int, optimized bool) (build *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildByID(ctx, repoSource, repoOwner, repoName, id, optimized)
}

func (c *tracingClient) GetLastPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (build *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetLastPipelineBuild"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetLastPipelineBuild(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *tracingClient) GetFirstPipelineBuild(ctx context.Context, repoSource, repoOwner, repoName string, optimized bool) (build *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFirstPipelineBuild"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFirstPipelineBuild(ctx, repoSource, repoOwner, repoName, optimized)
}

func (c *tracingClient) GetLastPipelineBuildForBranch(ctx context.Context, repoSource, repoOwner, repoName, branch string) (build *contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetLastPipelineBuildForBranch"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetLastPipelineBuildForBranch(ctx, repoSource, repoOwner, repoName, branch)
}

func (c *tracingClient) GetLastPipelineReleases(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string, pageSize int) (releases []*contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetLastPipelineReleases"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetLastPipelineReleases(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction, pageSize)
}

func (c *tracingClient) GetFirstPipelineRelease(ctx context.Context, repoSource, repoOwner, repoName, releaseName, releaseAction string) (release *contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFirstPipelineRelease"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFirstPipelineRelease(ctx, repoSource, repoOwner, repoName, releaseName, releaseAction)
}

func (c *tracingClient) GetPipelineBuildsByVersion(ctx context.Context, repoSource, repoOwner, repoName, buildVersion string, statuses []string, limit uint64, optimized bool) (builds []*contracts.Build, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildsByVersion"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildsByVersion(ctx, repoSource, repoOwner, repoName, buildVersion, statuses, limit, optimized)
}

func (c *tracingClient) GetPipelineBuildLogs(ctx context.Context, repoSource, repoOwner, repoName, repoBranch, repoRevision, buildID string, readLogFromDatabase bool) (buildlog *contracts.BuildLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildLogs"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildLogs(ctx, repoSource, repoOwner, repoName, repoBranch, repoRevision, buildID, readLogFromDatabase)
}

func (c *tracingClient) GetPipelineBuildLogsPerPage(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber int, pageSize int) (buildlogs []*contracts.BuildLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildLogsPerPage"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildLogsPerPage(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize)
}

func (c *tracingClient) GetPipelineBuildMaxResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName string, lastNRecords int) (jobresources JobResources, count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildMaxResourceUtilization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildMaxResourceUtilization(ctx, repoSource, repoOwner, repoName, lastNRecords)
}

func (c *tracingClient) GetPipelineReleases(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (releases []*contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleases"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleases(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetPipelineReleasesCount(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleasesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleasesCount(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineRelease(ctx context.Context, repoSource, repoOwner, repoName string, id int) (release *contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineRelease"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineRelease(ctx, repoSource, repoOwner, repoName, id)
}

func (c *tracingClient) GetPipelineLastReleasesByName(ctx context.Context, repoSource, repoOwner, repoName, releaseName string, actions []string) (releases []contracts.Release, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineLastReleasesByName"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineLastReleasesByName(ctx, repoSource, repoOwner, repoName, releaseName, actions)
}

func (c *tracingClient) GetPipelineReleaseLogs(ctx context.Context, repoSource, repoOwner, repoName string, id int, readLogFromDatabase bool) (releaselog *contracts.ReleaseLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleaseLogs"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleaseLogs(ctx, repoSource, repoOwner, repoName, id, readLogFromDatabase)
}

func (c *tracingClient) GetPipelineReleaseLogsPerPage(ctx context.Context, repoSource, repoOwner, repoName string, pageNumber int, pageSize int) (releaselogs []*contracts.ReleaseLog, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleaseLogsPerPage"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleaseLogsPerPage(ctx, repoSource, repoOwner, repoName, pageNumber, pageSize)
}

func (c *tracingClient) GetPipelineReleaseMaxResourceUtilization(ctx context.Context, repoSource, repoOwner, repoName, targetName string, lastNRecords int) (jobresources JobResources, count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleaseMaxResourceUtilization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleaseMaxResourceUtilization(ctx, repoSource, repoOwner, repoName, targetName, lastNRecords)
}

func (c *tracingClient) GetBuildsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetBuildsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetBuildsCount(ctx, filters)
}

func (c *tracingClient) GetReleasesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetReleasesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetReleasesCount(ctx, filters)
}

func (c *tracingClient) GetBuildsDuration(ctx context.Context, filters map[api.FilterType][]string) (duration time.Duration, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetBuildsDuration"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetBuildsDuration(ctx, filters)
}

func (c *tracingClient) GetFirstBuildTimes(ctx context.Context, filters map[api.FilterType][]string) (times []time.Time, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFirstBuildTimes"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFirstBuildTimes(ctx, filters)
}

func (c *tracingClient) GetFirstReleaseTimes(ctx context.Context, filters map[api.FilterType][]string) (times []time.Time, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFirstReleaseTimes"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFirstReleaseTimes(ctx, filters)
}

func (c *tracingClient) GetPipelineBuildsDurations(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (durations []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildsDurations"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildsDurations(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineReleasesDurations(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (durations []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleasesDurations"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleasesDurations(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineBuildsCPUUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildsCPUUsageMeasurements"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildsCPUUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineReleasesCPUUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleasesCPUUsageMeasurements"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleasesCPUUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineBuildsMemoryUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineBuildsMemoryUsageMeasurements"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineBuildsMemoryUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetPipelineReleasesMemoryUsageMeasurements(ctx context.Context, repoSource, repoOwner, repoName string, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineReleasesMemoryUsageMeasurements"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineReleasesMemoryUsageMeasurements(ctx, repoSource, repoOwner, repoName, filters)
}

func (c *tracingClient) GetLabelValues(ctx context.Context, labelKey string) (labels []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetLabelValues"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetLabelValues(ctx, labelKey)
}

func (c *tracingClient) GetFrequentLabels(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (measurements []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFrequentLabels"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFrequentLabels(ctx, pageNumber, pageSize, filters)
}

func (c *tracingClient) GetFrequentLabelsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetFrequentLabelsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetFrequentLabelsCount(ctx, filters)
}

func (c *tracingClient) GetPipelinesWithMostBuilds(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (pipelines []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesWithMostBuilds"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesWithMostBuilds(ctx, pageNumber, pageSize, filters)
}

func (c *tracingClient) GetPipelinesWithMostBuildsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesWithMostBuildsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesWithMostBuildsCount(ctx, filters)
}

func (c *tracingClient) GetPipelinesWithMostReleases(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (pipelines []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesWithMostReleases"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesWithMostReleases(ctx, pageNumber, pageSize, filters)
}

func (c *tracingClient) GetPipelinesWithMostReleasesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelinesWithMostReleasesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelinesWithMostReleasesCount(ctx, filters)
}

func (c *tracingClient) GetTriggers(ctx context.Context, triggerType, identifier, event string) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetTriggers(ctx, triggerType, identifier, event)
}

func (c *tracingClient) GetGitTriggers(ctx context.Context, gitEvent manifest.EstafetteGitEvent) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetGitTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetGitTriggers(ctx, gitEvent)
}

func (c *tracingClient) GetPipelineTriggers(ctx context.Context, build contracts.Build, event string) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPipelineTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPipelineTriggers(ctx, build, event)
}

func (c *tracingClient) GetReleaseTriggers(ctx context.Context, release contracts.Release, event string) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetReleaseTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetReleaseTriggers(ctx, release, event)
}

func (c *tracingClient) GetPubSubTriggers(ctx context.Context, pubsubEvent manifest.EstafettePubSubEvent) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetPubSubTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetPubSubTriggers(ctx, pubsubEvent)
}

func (c *tracingClient) GetCronTriggers(ctx context.Context) (pipelines []*contracts.Pipeline, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCronTriggers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCronTriggers(ctx)
}

func (c *tracingClient) Rename(ctx context.Context, shortFromRepoSource, fromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "Rename"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.Rename(ctx, shortFromRepoSource, fromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameBuildVersion(ctx context.Context, shortFromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameBuildVersion"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameBuildVersion(ctx, shortFromRepoSource, fromRepoOwner, fromRepoName, shortToRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameBuilds(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameBuilds"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameBuilds(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameBuildLogs(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameBuildLogs"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameBuildLogs(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameReleases(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameReleases"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameReleases(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameReleaseLogs(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameReleaseLogs"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameReleaseLogs(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameComputedPipelines(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameComputedPipelines"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameComputedPipelines(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) RenameComputedReleases(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "RenameComputedReleases"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.RenameComputedReleases(ctx, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName)
}

func (c *tracingClient) InsertUser(ctx context.Context, user contracts.User) (u *contracts.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertUser"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertUser(ctx, user)
}

func (c *tracingClient) UpdateUser(ctx context.Context, user contracts.User) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateUser"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateUser(ctx, user)
}

func (c *tracingClient) DeleteUser(ctx context.Context, user contracts.User) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "DeleteUser"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.DeleteUser(ctx, user)
}

func (c *tracingClient) GetUserByIdentity(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetUserByIdentity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetUserByIdentity(ctx, identity)
}

func (c *tracingClient) GetUserByID(ctx context.Context, id string, filters map[api.FilterType][]string) (user *contracts.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetUserByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetUserByID(ctx, id, filters)
}

func (c *tracingClient) GetUsers(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (users []*contracts.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetUsers"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetUsers(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetUsersCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetUsersCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetUsersCount(ctx, filters)
}

func (c *tracingClient) InsertGroup(ctx context.Context, group contracts.Group) (g *contracts.Group, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertGroup"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertGroup(ctx, group)
}

func (c *tracingClient) UpdateGroup(ctx context.Context, group contracts.Group) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateGroup"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateGroup(ctx, group)
}

func (c *tracingClient) DeleteGroup(ctx context.Context, group contracts.Group) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "DeleteGroup"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.DeleteGroup(ctx, group)
}

func (c *tracingClient) GetGroupByIdentity(ctx context.Context, identity contracts.GroupIdentity) (group *contracts.Group, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetGroupByIdentity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetGroupByIdentity(ctx, identity)
}

func (c *tracingClient) GetGroupByID(ctx context.Context, id string, filters map[api.FilterType][]string) (group *contracts.Group, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetGroupByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetGroupByID(ctx, id, filters)
}

func (c *tracingClient) GetGroups(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (groups []*contracts.Group, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetGroups"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetGroups(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetGroupsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetGroupsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetGroupsCount(ctx, filters)
}

func (c *tracingClient) InsertOrganization(ctx context.Context, organization contracts.Organization) (o *contracts.Organization, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertOrganization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertOrganization(ctx, organization)
}

func (c *tracingClient) UpdateOrganization(ctx context.Context, organization contracts.Organization) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateOrganization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateOrganization(ctx, organization)
}

func (c *tracingClient) DeleteOrganization(ctx context.Context, organization contracts.Organization) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "DeleteOrganization"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.DeleteOrganization(ctx, organization)
}

func (c *tracingClient) GetOrganizationByIdentity(ctx context.Context, identity contracts.OrganizationIdentity) (organization *contracts.Organization, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetOrganizationByIdentity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetOrganizationByIdentity(ctx, identity)
}

func (c *tracingClient) GetOrganizationByID(ctx context.Context, id string) (organization *contracts.Organization, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetOrganizationByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetOrganizationByID(ctx, id)
}

func (c *tracingClient) GetOrganizationByName(ctx context.Context, name string) (organization *contracts.Organization, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetOrganizationByName"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetOrganizationByName(ctx, name)
}

func (c *tracingClient) GetOrganizations(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (organizations []*contracts.Organization, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetOrganizations"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetOrganizations(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetOrganizationsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetOrganizationsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetOrganizationsCount(ctx, filters)
}

func (c *tracingClient) InsertClient(ctx context.Context, client contracts.Client) (cl *contracts.Client, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertClient"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertClient(ctx, client)
}

func (c *tracingClient) UpdateClient(ctx context.Context, client contracts.Client) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateClient"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateClient(ctx, client)
}

func (c *tracingClient) DeleteClient(ctx context.Context, client contracts.Client) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "DeleteClient"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.DeleteClient(ctx, client)
}

func (c *tracingClient) GetClientByClientID(ctx context.Context, clientID string) (client *contracts.Client, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetClientByClientID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetClientByClientID(ctx, clientID)
}

func (c *tracingClient) GetClientByID(ctx context.Context, id string) (client *contracts.Client, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetClientByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetClientByID(ctx, id)
}

func (c *tracingClient) GetClients(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (clients []*contracts.Client, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetClients"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetClients(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetClientsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetClientsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetClientsCount(ctx, filters)
}

func (c *tracingClient) InsertCatalogEntity(ctx context.Context, catalogEntity contracts.CatalogEntity) (insertedCatalogEntity *contracts.CatalogEntity, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "InsertCatalogEntity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.InsertCatalogEntity(ctx, catalogEntity)
}

func (c *tracingClient) UpdateCatalogEntity(ctx context.Context, catalogEntity contracts.CatalogEntity) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "UpdateCatalogEntity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.UpdateCatalogEntity(ctx, catalogEntity)
}

func (c *tracingClient) DeleteCatalogEntity(ctx context.Context, id string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "DeleteCatalogEntity"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.DeleteCatalogEntity(ctx, id)
}

func (c *tracingClient) GetCatalogEntityByID(ctx context.Context, id string) (catalogEntity *contracts.CatalogEntity, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityByID"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityByID(ctx, id)
}

func (c *tracingClient) GetCatalogEntities(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (catalogEntities []*contracts.CatalogEntity, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntities"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntities(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetCatalogEntitiesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntitiesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntitiesCount(ctx, filters)
}

func (c *tracingClient) GetCatalogEntityParentKeys(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (keys []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityParentKeys"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityParentKeys(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetCatalogEntityParentKeysCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityParentKeysCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityParentKeysCount(ctx, filters)
}

func (c *tracingClient) GetCatalogEntityParentValues(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (values []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityParentValues"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityParentValues(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetCatalogEntityParentValuesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityParentValuesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityParentValuesCount(ctx, filters)
}

func (c *tracingClient) GetCatalogEntityKeys(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (keys []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityKeys"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityKeys(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetCatalogEntityKeysCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityKeysCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityKeysCount(ctx, filters)
}

func (c *tracingClient) GetCatalogEntityValues(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string, sortings []api.OrderField) (values []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityValues"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityValues(ctx, pageNumber, pageSize, filters, sortings)
}

func (c *tracingClient) GetCatalogEntityValuesCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityValuesCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityValuesCount(ctx, filters)
}

func (c *tracingClient) GetCatalogEntityLabels(ctx context.Context, pageNumber, pageSize int, filters map[api.FilterType][]string) (labels []map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityLabels"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityLabels(ctx, pageNumber, pageSize, filters)
}

func (c *tracingClient) GetCatalogEntityLabelsCount(ctx context.Context, filters map[api.FilterType][]string) (count int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, api.GetSpanName(c.prefix, "GetCatalogEntityLabelsCount"))
	defer func() { api.FinishSpanWithError(span, err) }()

	return c.Client.GetCatalogEntityLabelsCount(ctx, filters)
}
