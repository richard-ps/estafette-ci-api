package api

import (
	"encoding/json"
	"math"
	"testing"

	crypt "github.com/estafette/estafette-ci-crypt"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFromFile(t *testing.T) {

	t.Run("ReturnsConfigWithoutErrors", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		_, err := configReader.ReadConfigFromFile("test-config.yaml", true)

		assert.Nil(t, err)
	})

	t.Run("ReturnsGithubConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		githubConfig := config.Integrations.Github

		assert.Equal(t, "/github-app-key/private-key.pem", githubConfig.PrivateKeyPath)
		assert.Equal(t, "15", githubConfig.AppID)
		assert.Equal(t, "asdas2342", githubConfig.ClientID)
		assert.Equal(t, "this is my secret", githubConfig.ClientSecret)
		assert.Equal(t, 2, len(githubConfig.WhitelistedInstallations))
		assert.Equal(t, 15, githubConfig.WhitelistedInstallations[0])
		assert.Equal(t, 83, githubConfig.WhitelistedInstallations[1])

		assert.Equal(t, 2, len(githubConfig.InstallationOrganizations))
		assert.Equal(t, 15, githubConfig.InstallationOrganizations[0].Installation)
		assert.Equal(t, 1, len(githubConfig.InstallationOrganizations[0].Organizations))
		assert.Equal(t, "Estafette", githubConfig.InstallationOrganizations[0].Organizations[0].Name)
		assert.Equal(t, 83, githubConfig.InstallationOrganizations[1].Installation)
		assert.Equal(t, 1, len(githubConfig.InstallationOrganizations[1].Organizations))
		assert.Equal(t, "Estafette", githubConfig.InstallationOrganizations[1].Organizations[0].Name)
	})

	t.Run("ReturnsBitbucketConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		bitbucketConfig := config.Integrations.Bitbucket

		assert.Equal(t, "sd9ewiwuejkwejkewk", bitbucketConfig.APIKey)
		assert.Equal(t, "2390w3e90jdsk", bitbucketConfig.AppOAuthKey)
		assert.Equal(t, "this is my secret", bitbucketConfig.AppOAuthSecret)
		assert.Equal(t, 1, len(bitbucketConfig.WhitelistedOwners))
		assert.Equal(t, "estafette", bitbucketConfig.WhitelistedOwners[0])

		assert.Equal(t, 1, len(bitbucketConfig.OwnerOrganizations))
		assert.Equal(t, "estafette", bitbucketConfig.OwnerOrganizations[0].Owner)
		assert.Equal(t, 1, len(bitbucketConfig.OwnerOrganizations[0].Organizations))
		assert.Equal(t, "Estafette", bitbucketConfig.OwnerOrganizations[0].Organizations[0].Name)
	})

	t.Run("ReturnsCloudsourceConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		cloudsourceConfig := config.Integrations.CloudSource

		assert.Equal(t, "estafette", cloudsourceConfig.WhitelistedProjects[0])
		assert.Equal(t, 1, len(cloudsourceConfig.ProjectOrganizations))
		assert.Equal(t, "estafette", cloudsourceConfig.ProjectOrganizations[0].Project)
		assert.Equal(t, 1, len(cloudsourceConfig.ProjectOrganizations[0].Organizations))
		assert.Equal(t, "Estafette", cloudsourceConfig.ProjectOrganizations[0].Organizations[0].Name)
	})

	t.Run("ReturnsSlackConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		slackConfig := config.Integrations.Slack

		assert.Equal(t, "d9ew90weoijewjke", slackConfig.ClientID)
		assert.Equal(t, "this is my secret", slackConfig.ClientSecret)
		assert.Equal(t, "this is my secret", slackConfig.AppVerificationToken)
		assert.Equal(t, "this is my secret", slackConfig.AppOAuthAccessToken)
	})

	t.Run("ReturnsPrometheusConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		prometheusConfig := config.Integrations.Prometheus

		assert.Equal(t, "http://prometheus-server.monitoring.svc.cluster.local", prometheusConfig.ServerURL)
		assert.Equal(t, 10, prometheusConfig.ScrapeIntervalSeconds)
	})

	t.Run("ReturnsBigQueryConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		bigqueryConfig := config.Integrations.BigQuery

		assert.Equal(t, true, bigqueryConfig.Enable)
		assert.Equal(t, "my-gcp-project", bigqueryConfig.ProjectID)
		assert.Equal(t, "my-dataset", bigqueryConfig.Dataset)
	})

	t.Run("ReturnsCloudStorageConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		cloudStorageConfig := config.Integrations.CloudStorage

		assert.Equal(t, "my-gcp-project", cloudStorageConfig.ProjectID)
		assert.Equal(t, "my-bucket", cloudStorageConfig.Bucket)
		assert.Equal(t, "logs", cloudStorageConfig.LogsDirectory)
	})

	t.Run("ReturnsAPIServerConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		apiServerConfig := config.APIServer

		assert.Equal(t, "https://ci.estafette.io/", apiServerConfig.BaseURL)
		assert.Equal(t, "http://estafette-ci-api.estafette.svc.cluster.local/", apiServerConfig.ServiceURL)
		assert.Equal(t, 2, len(apiServerConfig.LogWriters))
		assert.Equal(t, "database", apiServerConfig.LogWriters[0])
		assert.Equal(t, "cloudstorage", apiServerConfig.LogWriters[1])
		assert.Equal(t, "database", apiServerConfig.LogReader)
	})

	t.Run("ReturnsAuthConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		authConfig := config.Auth

		assert.Equal(t, "ci.estafette.io", authConfig.JWT.Domain)
		assert.Equal(t, "this is my secret", authConfig.JWT.Key)

		assert.Equal(t, 3, len(authConfig.Organizations))
		assert.Equal(t, "Org A", authConfig.Organizations[0].Name)
		assert.Equal(t, 1, len(authConfig.Organizations[0].OAuthProviders))
		assert.Equal(t, "google", authConfig.Organizations[0].OAuthProviders[0].Name)
		assert.Equal(t, "abcdasa", authConfig.Organizations[0].OAuthProviders[0].ClientID)
		assert.Equal(t, "asdsddsfdfs", authConfig.Organizations[0].OAuthProviders[0].ClientSecret)
		assert.Equal(t, ".+@estafette\\.io", authConfig.Organizations[0].OAuthProviders[0].AllowedIdentitiesRegex)

		assert.Equal(t, "Org B", authConfig.Organizations[1].Name)
		assert.Equal(t, 1, len(authConfig.Organizations[1].OAuthProviders))

		assert.Equal(t, "Org C", authConfig.Organizations[2].Name)
		assert.Equal(t, 1, len(authConfig.Organizations[2].OAuthProviders))

		assert.Equal(t, 2, len(authConfig.Administrators))
		assert.Equal(t, "admin1@server.com", authConfig.Administrators[0])
		assert.Equal(t, "admin2@server.com", authConfig.Administrators[1])
	})

	t.Run("ReturnsJobsConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		jobsConfig := config.Jobs

		assert.Equal(t, "estafette-ci-jobs", jobsConfig.Namespace)
		assert.Equal(t, 0.1, jobsConfig.MinCPUCores)
		assert.Equal(t, 3.5, jobsConfig.MaxCPUCores)
		assert.Equal(t, 1.0, jobsConfig.CPURequestRatio)
		assert.Equal(t, 64*math.Pow(2, 10)*math.Pow(2, 10), jobsConfig.MinMemoryBytes)                 // 64Mi
		assert.Equal(t, 12*math.Pow(2, 10)*math.Pow(2, 10)*math.Pow(2, 10), jobsConfig.MaxMemoryBytes) // 12Gi
		assert.Equal(t, 1.25, jobsConfig.MemoryRequestRatio)
	})

	t.Run("ReturnsDatabaseConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		databaseConfig := config.Database

		assert.Equal(t, "estafette_ci_api", databaseConfig.DatabaseName)
		assert.Equal(t, "cockroachdb-public.estafette.svc.cluster.local", databaseConfig.Host)
		assert.Equal(t, true, databaseConfig.Insecure)
		assert.Equal(t, "/cockroachdb-certificates/cockroachdb.crt", databaseConfig.CertificateDir)
		assert.Equal(t, 26257, databaseConfig.Port)
		assert.Equal(t, "myuser", databaseConfig.User)
		assert.Equal(t, "this is my secret", databaseConfig.Password)
	})

	t.Run("ReturnsManifestPreferences", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, err := configReader.ReadConfigFromFile("test-config.yaml", true)

		if !assert.Nil(t, err) {
			return
		}

		if !assert.NotNil(t, config.ManifestPreferences) {
			return
		}

		assert.Equal(t, 1, len(config.ManifestPreferences.LabelRegexes))
		assert.Equal(t, "api|web|library|container", config.ManifestPreferences.LabelRegexes["type"])
		assert.Equal(t, 2, len(config.ManifestPreferences.BuilderOperatingSystems))
		assert.Equal(t, "linux", config.ManifestPreferences.BuilderOperatingSystems[0])
		assert.Equal(t, "windows", config.ManifestPreferences.BuilderOperatingSystems[1])
		assert.Equal(t, 2, len(config.ManifestPreferences.BuilderTracksPerOperatingSystem))
		assert.Equal(t, 3, len(config.ManifestPreferences.BuilderTracksPerOperatingSystem["linux"]))
		assert.Equal(t, 3, len(config.ManifestPreferences.BuilderTracksPerOperatingSystem["windows"]))
	})

	t.Run("ReturnsCatalogConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		catalogConfig := config.Catalog

		assert.Equal(t, 2, len(catalogConfig.Filters))
		assert.Equal(t, "type", catalogConfig.Filters[0])
		assert.Equal(t, "team", catalogConfig.Filters[1])
	})

	t.Run("ReturnsCredentialsConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		credentialsConfig := config.Credentials

		assert.Equal(t, 9, len(credentialsConfig))
		assert.Equal(t, "container-registry-extensions", credentialsConfig[0].Name)
		assert.Equal(t, "container-registry", credentialsConfig[0].Type)
		assert.Equal(t, "extensions", credentialsConfig[0].AdditionalProperties["repository"])
		assert.Equal(t, "slack-webhook-estafette", credentialsConfig[6].Name)
		assert.Equal(t, "slack-webhook", credentialsConfig[6].Type)
		assert.Equal(t, "estafette", credentialsConfig[6].AdditionalProperties["workspace"])
	})

	t.Run("ReturnsTrustedImagesConfig", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		trustedImagesConfig := config.TrustedImages

		assert.Equal(t, 8, len(trustedImagesConfig))
		assert.Equal(t, "extensions/docker", trustedImagesConfig[0].ImagePath)
		assert.True(t, trustedImagesConfig[0].RunDocker)
		assert.Equal(t, 1, len(trustedImagesConfig[0].InjectedCredentialTypes))
		assert.Equal(t, "container-registry", trustedImagesConfig[0].InjectedCredentialTypes[0])

		assert.Equal(t, "multiple-git-sources-test", trustedImagesConfig[7].ImagePath)
		assert.False(t, trustedImagesConfig[7].RunDocker)
		assert.Equal(t, 3, len(trustedImagesConfig[7].InjectedCredentialTypes))
		assert.Equal(t, "bitbucket-api-token", trustedImagesConfig[7].InjectedCredentialTypes[0])
		assert.Equal(t, "github-api-token", trustedImagesConfig[7].InjectedCredentialTypes[1])
		assert.Equal(t, "cloudsource-api-token", trustedImagesConfig[7].InjectedCredentialTypes[2])
	})

	t.Run("ReturnsRegistryMirror", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		registryMirrorConfig := config.RegistryMirror

		assert.NotNil(t, registryMirrorConfig)
		assert.Equal(t, "https://mirror.gcr.io", *registryMirrorConfig)
	})

	t.Run("ReturnsDockerDaemonMTU", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		dockerDaemonMTUConfig := config.DockerDaemonMTU

		assert.NotNil(t, dockerDaemonMTUConfig)
		assert.Equal(t, "1360", *dockerDaemonMTUConfig)
	})

	t.Run("ReturnsDockerDaemonBIP", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		dockerDaemonBIPConfig := config.DockerDaemonBIP

		assert.NotNil(t, dockerDaemonBIPConfig)
		assert.Equal(t, "192.168.1.1/24", *dockerDaemonBIPConfig)
	})

	t.Run("ReturnsDockerNetwork", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		dockerNetworkConfig := config.DockerNetwork

		if assert.NotNil(t, dockerNetworkConfig) {
			assert.Equal(t, "estafette", dockerNetworkConfig.Name)
			assert.Equal(t, "192.168.2.1/24", dockerNetworkConfig.Subnet)
			assert.Equal(t, "192.168.2.1", dockerNetworkConfig.Gateway)
		}
	})

	t.Run("AllowsCredentialConfigWithComplexAdditionalPropertiesToBeJSONMarshalled", func(t *testing.T) {

		configReader := NewConfigReader(crypt.NewSecretHelper("SazbwMf3NZxVVbBqQHebPcXCqrVn3DDp", false))

		// act
		config, _ := configReader.ReadConfigFromFile("test-config.yaml", true)

		credentialsConfig := config.Credentials

		bytes, err := json.Marshal(credentialsConfig[2])

		assert.Nil(t, err)
		assert.Equal(t, "{\"name\":\"gke-estafette-production\",\"type\":\"kubernetes-engine\",\"additionalProperties\":{\"cluster\":\"production-europe-west2\",\"defaults\":{\"autoscale\":{\"min\":2},\"container\":{\"repository\":\"estafette\"},\"namespace\":\"estafette\",\"sidecars\":[{\"image\":\"estafette/openresty-sidecar:1.13.6.1-alpine\",\"type\":\"openresty\"}]},\"project\":\"estafette-production\",\"region\":\"europe-west2\",\"serviceAccountKeyfile\":\"{}\"}}", string(bytes))
	})
}

func TestWriteLogToDatabase(t *testing.T) {

	t.Run("ReturnsTrueIfLogWritersIsEmpty", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{},
		}

		// act
		result := config.WriteLogToDatabase()

		assert.True(t, result)
	})

	t.Run("ReturnsTrueIfLogWritersContainsDatabase", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{
				"cloudstorage",
				"database",
			},
		}

		// act
		result := config.WriteLogToDatabase()

		assert.True(t, result)
	})

	t.Run("ReturnsFalseIfLogWritersDoesNotContainDatabase", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{
				"cloudstorage",
			},
		}

		// act
		result := config.WriteLogToDatabase()

		assert.False(t, result)
	})
}

func TestWriteLogToCloudStorage(t *testing.T) {

	t.Run("ReturnsFalseIfLogWritersIsEmpty", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{},
		}

		// act
		result := config.WriteLogToCloudStorage()

		assert.False(t, result)
	})

	t.Run("ReturnsTrueIfLogWritersContainsCloudStorage", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{
				"cloudstorage",
				"database",
			},
		}

		// act
		result := config.WriteLogToCloudStorage()

		assert.True(t, result)
	})

	t.Run("ReturnsFalseIfLogWritersDoesNotContainCloudStorage", func(t *testing.T) {

		config := APIServerConfig{
			LogWriters: []string{
				"database",
			},
		}

		// act
		result := config.WriteLogToCloudStorage()

		assert.False(t, result)
	})
}

func TestReadLogFromDatabase(t *testing.T) {

	t.Run("ReturnsTrueIfLogReaderIsEmpty", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "",
		}

		// act
		result := config.ReadLogFromDatabase()

		assert.True(t, result)
	})

	t.Run("ReturnsTrueIfLogReaderEqualsDatabase", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "database",
		}

		// act
		result := config.ReadLogFromDatabase()

		assert.True(t, result)
	})

	t.Run("ReturnsFalseIfLogReaderDoesNotEqualDatabase", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "cloudstorage",
		}

		// act
		result := config.ReadLogFromDatabase()

		assert.False(t, result)
	})
}

func TestReadLogFromCloudStorage(t *testing.T) {

	t.Run("ReturnsFalseIfLogReaderIsEmpty", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "",
		}

		// act
		result := config.ReadLogFromCloudStorage()

		assert.False(t, result)
	})

	t.Run("ReturnsTrueIfLogReaderEqualsCloudStorage", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "cloudstorage",
		}

		// act
		result := config.ReadLogFromCloudStorage()

		assert.True(t, result)
	})

	t.Run("ReturnsFalseIfLogReaderDoesNotCloudStorage", func(t *testing.T) {

		config := APIServerConfig{
			LogReader: "database",
		}

		// act
		result := config.ReadLogFromCloudStorage()

		assert.False(t, result)
	})
}

// // ReadLogFromDatabase indicates if logReader config is database
// func (c *APIServerConfig) ReadLogFromDatabase() bool {
// 	return c.LogReader == "database"
// }

// // ReadLogFromCloudStorage indicates if logReader config is cloudstorage
// func (c *APIServerConfig) ReadLogFromCloudStorage() bool {
// 	return c.LogReader == "cloudstorage"
// }
