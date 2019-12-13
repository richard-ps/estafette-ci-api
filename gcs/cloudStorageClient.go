package gcs

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"cloud.google.com/go/storage"
	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/opentracing/opentracing-go"
)

// CloudStorageClient is the interface for connecting to google cloud storage
type CloudStorageClient interface {
	InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog) (err error)
	InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog) (err error)
	GetPipelineBuildLogs(ctx context.Context, buildLog contracts.BuildLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error)
	GetPipelineReleaseLogs(ctx context.Context, releaseLog contracts.ReleaseLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error)
}

type cloudStorageClientImpl struct {
	client *storage.Client
	config *config.CloudStorageConfig
}

// NewCloudStorageClient returns new CloudStorageClient
func NewCloudStorageClient(config *config.CloudStorageConfig) (CloudStorageClient, error) {

	if config == nil {
		return &cloudStorageClientImpl{
			config: config,
		}, nil
	}

	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &cloudStorageClientImpl{
		client: storageClient,
		config: config,
	}, nil
}

func (impl *cloudStorageClientImpl) InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog) (err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudStorageClient::InsertBuildLog")
	defer span.Finish()

	logPath := impl.getBuildLogPath(buildLog)

	return foundation.Retry(func() error {
		return impl.insertLog(ctx, logPath, buildLog.Steps)
	})
}

func (impl *cloudStorageClientImpl) InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog) (err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudStorageClient::InsertReleaseLog")
	defer span.Finish()

	logPath := impl.getReleaseLogPath(releaseLog)

	return foundation.Retry(func() error {
		return impl.insertLog(ctx, logPath, releaseLog.Steps)
	})
}

func (impl *cloudStorageClientImpl) insertLog(ctx context.Context, path string, steps []*contracts.BuildLogStep) (err error) {

	bucket := impl.client.Bucket(impl.config.Bucket)

	// marshal json
	jsonBytes, err := json.Marshal(steps)
	if err != nil {
		return err
	}

	// create writer for cloud storage object
	logObject := bucket.Object(path)

	// don't allow overwrites, return when file already exists
	_, err = logObject.Attrs(ctx)
	if err == nil {
		// log file already exists, return
		return nil
	}
	if err != nil && err != storage.ErrObjectNotExist {
		// some other error happened, return it
		return err
	}

	// object doesn't exist, okay to write it
	writer := logObject.NewWriter(ctx)
	if writer == nil {
		return fmt.Errorf("Writer for logobject %v is nil", path)
	}

	// write compressed bytes
	gz, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
	if err != nil {
		_ = writer.Close()
		return err
	}
	_, err = gz.Write(jsonBytes)
	if err != nil {
		_ = writer.Close()
		return err
	}
	err = gz.Close()
	if err != nil {
		_ = writer.Close()
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (impl *cloudStorageClientImpl) GetPipelineBuildLogs(ctx context.Context, buildLog contracts.BuildLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudStorageClient::GetPipelineBuildLogs")
	defer span.Finish()

	logPath := impl.getBuildLogPath(buildLog)

	return impl.getLog(ctx, logPath, acceptGzipEncoding, responseWriter)
}

func (impl *cloudStorageClientImpl) GetPipelineReleaseLogs(ctx context.Context, releaseLog contracts.ReleaseLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudStorageClient::GetPipelineReleaseLogs")
	defer span.Finish()

	logPath := impl.getReleaseLogPath(releaseLog)

	return impl.getLog(ctx, logPath, acceptGzipEncoding, responseWriter)
}

func (impl *cloudStorageClientImpl) getLog(ctx context.Context, path string, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudStorageClient::getLog")
	defer span.Finish()

	span.SetTag("log-path", path)
	span.SetTag("accept-gzip-encoding", acceptGzipEncoding)

	bucket := impl.client.Bucket(impl.config.Bucket)

	// create reader for cloud storage object
	logObject := bucket.Object(path).ReadCompressed(true)
	reader, err := logObject.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	// create source reader to either copy compressed bytes or decompress them first
	sourceReader := io.Reader(reader)
	if acceptGzipEncoding {
		responseWriter.Header().Set("Content-Encoding", "gzip")
		responseWriter.Header().Set("Vary", "Accept-Encoding")
	} else {
		gzr, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer gzr.Close()
		sourceReader = io.Reader(gzr)
	}

	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	writtenBytes, err := io.Copy(responseWriter, sourceReader)
	if err != nil {
		return err
	}

	responseWriter.Header().Set("Content-Length", fmt.Sprint(writtenBytes))

	return nil
}

func (impl *cloudStorageClientImpl) getBuildLogPath(buildLog contracts.BuildLog) (logPath string) {

	logPath = path.Join(impl.config.LogsDirectory, buildLog.RepoSource, buildLog.RepoOwner, buildLog.RepoName, "builds", fmt.Sprintf("%v.log", buildLog.ID))

	return logPath
}

func (impl *cloudStorageClientImpl) getReleaseLogPath(releaseLog contracts.ReleaseLog) (logPath string) {

	logPath = path.Join(impl.config.LogsDirectory, releaseLog.RepoSource, releaseLog.RepoOwner, releaseLog.RepoName, "releases", fmt.Sprintf("%v.log", releaseLog.ID))

	return logPath
}
