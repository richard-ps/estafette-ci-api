package cloudstorage

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

var (
	// ErrLogNotExist is returned when a log cannot be found
	ErrLogNotExist = errors.New("The log does not exist")
)

// Client is the interface for connecting to google cloud storage
type Client interface {
	InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog) (err error)
	InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog) (err error)
	GetPipelineBuildLogs(ctx context.Context, buildLog contracts.BuildLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error)
	GetPipelineReleaseLogs(ctx context.Context, releaseLog contracts.ReleaseLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error)
	Rename(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error)
}

// NewClient returns new cloudstorage.Client
func NewClient(config *config.CloudStorageConfig, storageClient *storage.Client) Client {

	if config == nil {
		return &client{
			config: config,
		}
	}

	return &client{
		client: storageClient,
		config: config,
	}
}

type client struct {
	client *storage.Client
	config *config.CloudStorageConfig
}

func (c *client) InsertBuildLog(ctx context.Context, buildLog contracts.BuildLog) (err error) {

	logPath := c.getBuildLogPath(buildLog)

	return foundation.Retry(func() error {
		return c.insertLog(ctx, logPath, buildLog.Steps)
	})
}

func (c *client) InsertReleaseLog(ctx context.Context, releaseLog contracts.ReleaseLog) (err error) {

	logPath := c.getReleaseLogPath(releaseLog)

	return foundation.Retry(func() error {
		return c.insertLog(ctx, logPath, releaseLog.Steps)
	})
}

func (c *client) insertLog(ctx context.Context, path string, steps []*contracts.BuildLogStep) (err error) {

	bucket := c.client.Bucket(c.config.Bucket)

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

func (c *client) GetPipelineBuildLogs(ctx context.Context, buildLog contracts.BuildLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	logPath := c.getBuildLogPath(buildLog)

	return c.getLog(ctx, logPath, acceptGzipEncoding, responseWriter)
}

func (c *client) GetPipelineReleaseLogs(ctx context.Context, releaseLog contracts.ReleaseLog, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	logPath := c.getReleaseLogPath(releaseLog)

	return c.getLog(ctx, logPath, acceptGzipEncoding, responseWriter)
}

func (c *client) getLog(ctx context.Context, path string, acceptGzipEncoding bool, responseWriter http.ResponseWriter) (err error) {

	bucket := c.client.Bucket(c.config.Bucket)

	// create reader for cloud storage object
	logObject := bucket.Object(path).ReadCompressed(true)
	reader, err := logObject.NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return ErrLogNotExist
		}

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

func (c *client) getBuildLogPath(buildLog contracts.BuildLog) (logPath string) {

	logDirectory := c.getLogDirectory(buildLog.RepoSource, buildLog.RepoOwner, buildLog.RepoName, "builds")

	logPath = path.Join(logDirectory, fmt.Sprintf("%v.log", buildLog.ID))

	return logPath
}

func (c *client) getReleaseLogPath(releaseLog contracts.ReleaseLog) (logPath string) {

	logDirectory := c.getLogDirectory(releaseLog.RepoSource, releaseLog.RepoOwner, releaseLog.RepoName, "releases")

	logPath = path.Join(logDirectory, fmt.Sprintf("%v.log", releaseLog.ID))

	return logPath
}

func (c *client) getLogDirectory(repoSource, repoOwner, repoName, logType string) (logDirectory string) {
	logDirectory = path.Join(c.config.LogsDirectory, repoSource, repoOwner, repoName, logType)

	return logDirectory
}

func (c *client) Rename(ctx context.Context, fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName string) (err error) {

	log.Info().Msgf("Renaming cloud storage logs from %v/%v/%v to %v/%v/%v for bucket %v", fromRepoSource, fromRepoOwner, fromRepoName, toRepoSource, toRepoOwner, toRepoName, c.config.Bucket)

	bucket := c.client.Bucket(c.config.Bucket)

	// list all build log files in old location, rename to new location
	fromBuildLogDirectory := c.getLogDirectory(fromRepoSource, fromRepoOwner, fromRepoName, "builds")
	toBuildLogDirectory := c.getLogDirectory(toRepoSource, toRepoOwner, toRepoName, "builds")

	err = c.renameFilesInDirectory(ctx, bucket, fromBuildLogDirectory, toBuildLogDirectory)
	if err != nil {
		return err
	}

	// list all release log files in old location, rename to new location
	fromReleaseLogDirectory := c.getLogDirectory(fromRepoSource, fromRepoOwner, fromRepoName, "releases")
	toReleaseLogDirectory := c.getLogDirectory(toRepoSource, toRepoOwner, toRepoName, "releases")

	err = c.renameFilesInDirectory(ctx, bucket, fromReleaseLogDirectory, toReleaseLogDirectory)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) renameFilesInDirectory(ctx context.Context, bucket *storage.BucketHandle, fromLogFileDirectory, toLogFileDirectory string) (err error) {

	query := &storage.Query{
		Prefix:    fromLogFileDirectory,
		Delimiter: "/",
	}
	query.SetAttrSelection([]string{"Name"})

	it := bucket.Objects(ctx, query)

	log.Info().Interface("bucket", *bucket).Interface("query", *query).Interface("iterator", *it).Msgf("Renaming cloud storage logs in bucket %v from directory %v to %v", bucket, fromLogFileDirectory, toLogFileDirectory)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fromLogFilePath := attrs.Name
		toLogFilePath := strings.Replace(fromLogFilePath, fromLogFileDirectory, toLogFileDirectory, 1)

		err = c.renameFile(ctx, bucket, fromLogFilePath, toLogFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *client) renameFile(ctx context.Context, bucket *storage.BucketHandle, fromLogFilePath, toLogFilePath string) (err error) {

	log.Debug().Interface("bucket", *bucket).Msgf("Renaming cloud storage log in bucket %v from path %v to %v", bucket, fromLogFilePath, toLogFilePath)

	src := bucket.Object(fromLogFilePath)
	dst := bucket.Object(toLogFilePath)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}

	return nil
}
