package builderapi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ericchiang/k8s"
	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ericchiang/k8s/apis/resource"
	"github.com/estafette/estafette-ci-api/clients/dockerhubapi"
	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
	crypt "github.com/estafette/estafette-ci-crypt"
	manifest "github.com/estafette/estafette-ci-manifest"
	"github.com/rs/zerolog/log"
)

// Client is the interface for running kubernetes commands specific to this application
type Client interface {
	CreateCiBuilderJob(ctx context.Context, params CiBuilderParams) (job *batchv1.Job, err error)
	RemoveCiBuilderJob(ctx context.Context, jobName string) (err error)
	CancelCiBuilderJob(ctx context.Context, jobName string) (err error)
	RemoveCiBuilderConfigMap(ctx context.Context, configmapName string) (err error)
	RemoveCiBuilderSecret(ctx context.Context, secretName string) (err error)
	TailCiBuilderJobLogs(ctx context.Context, jobName string, logChannel chan contracts.TailLogLine) (err error)
	GetJobName(ctx context.Context, jobType, repoOwner, repoName, id string) (jobname string)
	GetBuilderConfig(ctx context.Context, params CiBuilderParams, jobName string) (config contracts.BuilderConfig, err error)
	RefreshConfig(config *config.APIConfig, encryptedConfig *config.APIConfig)
}

// NewClient returns a new estafette.Client
func NewClient(config config.APIConfig, encryptedConfig config.APIConfig, secretHelper crypt.SecretHelper, kubeClient *k8s.Client, dockerHubClient dockerhubapi.Client) Client {

	return &client{
		kubeClient:      kubeClient,
		dockerHubClient: dockerHubClient,
		config:          config,
		encryptedConfig: encryptedConfig,
		secretHelper:    secretHelper,
	}
}

type client struct {
	kubeClient      *k8s.Client
	dockerHubClient dockerhubapi.Client
	config          config.APIConfig
	encryptedConfig config.APIConfig
	secretHelper    crypt.SecretHelper
}

func (c *client) RefreshConfig(config *config.APIConfig, encryptedConfig *config.APIConfig) {
	c.config = *config
	c.encryptedConfig = *encryptedConfig
}

// CreateCiBuilderJob creates an estafette-ci-builder job in Kubernetes to run the estafette build
func (c *client) CreateCiBuilderJob(ctx context.Context, ciBuilderParams CiBuilderParams) (job *batchv1.Job, err error) {

	// create job name of max 63 chars
	id := strconv.Itoa(ciBuilderParams.BuildID)
	if ciBuilderParams.JobType == "release" {
		id = strconv.Itoa(ciBuilderParams.ReleaseID)
	}

	jobName := c.GetJobName(ctx, ciBuilderParams.JobType, ciBuilderParams.RepoOwner, ciBuilderParams.RepoName, id)

	log.Info().Msgf("Creating job %v...", jobName)

	// extend builder config to parameterize the builder and replace all other envvars to improve security
	localBuilderConfig, err := c.GetBuilderConfig(ctx, ciBuilderParams, jobName)
	if err != nil {
		return
	}

	builderConfigPathName := "BUILDER_CONFIG_PATH"
	builderConfigPathValue := "/configs/builder-config.json"
	builderConfigJSONBytes, err := json.Marshal(localBuilderConfig)
	if err != nil {
		return
	}
	builderConfigValue := string(builderConfigJSONBytes)
	builderConfigValue, newKey, err := c.secretHelper.ReencryptAllEnvelopes(builderConfigValue, ciBuilderParams.GetFullRepoPath(), false)
	if err != nil {
		return
	}

	estafetteLogFormatName := "ESTAFETTE_LOG_FORMAT"
	estafetteLogFormatValue := os.Getenv("ESTAFETTE_LOG_FORMAT")
	jaegerServiceNameName := "JAEGER_SERVICE_NAME"
	jaegerServiceNameValue := "estafette-ci-builder"
	jaegerAgentHostName := "JAEGER_AGENT_HOST"
	jaegerAgentHostFieldPath := "status.hostIP"
	jaegerSamplerTypeName := "JAEGER_SAMPLER_TYPE"
	jaegerSamplerTypeValue := "const"
	jaegerSamplerParamName := "JAEGER_SAMPLER_PARAM"
	jaegerSamplerParamValue := "1"
	podNameName := "POD_NAME"
	podNameMetadataNameFieldPath := "metadata.name"
	environmentVariables := []*corev1.EnvVar{
		&corev1.EnvVar{
			Name:  &builderConfigPathName,
			Value: &builderConfigPathValue,
		},
		&corev1.EnvVar{
			Name:  &estafetteLogFormatName,
			Value: &estafetteLogFormatValue,
		},
		&corev1.EnvVar{
			Name:  &jaegerServiceNameName,
			Value: &jaegerServiceNameValue,
		},
		&corev1.EnvVar{
			Name: &jaegerAgentHostName,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: &jaegerAgentHostFieldPath,
				},
			},
		},
		&corev1.EnvVar{
			Name:  &jaegerSamplerTypeName,
			Value: &jaegerSamplerTypeValue,
		},
		&corev1.EnvVar{
			Name:  &jaegerSamplerTypeName,
			Value: &jaegerSamplerTypeValue,
		},
		&corev1.EnvVar{
			Name:  &jaegerSamplerParamName,
			Value: &jaegerSamplerParamValue,
		},
		&corev1.EnvVar{
			Name: &podNameName,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: &podNameMetadataNameFieldPath,
				},
			},
		},
	}

	// forward all envars prefixed with JAEGER_ to builder job
	for _, e := range os.Environ() {
		kvPair := strings.SplitN(e, "=", 2)

		if len(kvPair) == 2 {
			envvarName := kvPair[0]
			envvarValue := kvPair[1]

			if strings.HasPrefix(envvarName, "JAEGER_") && envvarName != "JAEGER_SERVICE_NAME" && envvarName != "JAEGER_AGENT_HOST" && envvarName != "JAEGER_SAMPLER_TYPE" && envvarName != "JAEGER_SAMPLER_PARAM" && envvarValue != "" {
				environmentVariables = append(environmentVariables, &corev1.EnvVar{
					Name:  &envvarName,
					Value: &envvarValue,
				})
			}
		}
	}

	// define resource request and limit values from job resources struct, so we can autotune later on
	cpuRequest := fmt.Sprintf("%f", ciBuilderParams.JobResources.CPURequest)
	cpuLimit := fmt.Sprintf("%f", ciBuilderParams.JobResources.CPULimit)
	memoryRequest := fmt.Sprintf("%.0f", ciBuilderParams.JobResources.MemoryRequest)
	memoryLimit := fmt.Sprintf("%.0f", ciBuilderParams.JobResources.MemoryLimit)

	// other job config
	containerName := "estafette-ci-builder"
	repository := "estafette/estafette-ci-builder"
	tag := ciBuilderParams.Track
	image := fmt.Sprintf("%v:%v", repository, tag)
	imagePullPolicy := "Always"
	digest, err := c.dockerHubClient.GetDigestCached(ctx, repository, tag)
	if err == nil && digest.Digest != "" {
		image = fmt.Sprintf("%v@%v", repository, digest.Digest)
		imagePullPolicy = "IfNotPresent"
	}
	restartPolicy := "Never"
	privileged := true

	preemptibleAffinityWeight := int32(10)
	preemptibleAffinityKey := "cloud.google.com/gke-preemptible"
	preemptibleAffinityOperator := "In"

	operatingSystemAffinityKey := "beta.kubernetes.io/os"
	operatingSystemAffinityOperator := "In"
	operatingSystemAffinityValue := ciBuilderParams.OperatingSystem

	// create configmap for builder config
	builderConfigConfigmapName := jobName
	builderConfigVolumeName := "app-configs"
	builderConfigVolumeMountPath := "/configs"
	configmap := &corev1.ConfigMap{
		Metadata: &metav1.ObjectMeta{
			Name:      &builderConfigConfigmapName,
			Namespace: &c.config.Jobs.Namespace,
			Labels: map[string]string{
				"createdBy": "estafette",
				"jobType":   ciBuilderParams.JobType,
			},
		},
		Data: map[string]string{
			"builder-config.json": builderConfigValue,
		},
	}

	err = c.kubeClient.Create(context.Background(), configmap)
	if err != nil {
		return
	}

	log.Info().Msgf("Configmap %v is created", builderConfigConfigmapName)

	// create secret for decryption key secret
	decryptionKeySecretName := jobName
	decryptionKeySecretVolumeName := "app-secret"
	decryptionKeySecretVolumeMountPath := "/secrets"
	secret := &corev1.Secret{
		Metadata: &metav1.ObjectMeta{
			Name:      &decryptionKeySecretName,
			Namespace: &c.config.Jobs.Namespace,
			Labels: map[string]string{
				"createdBy": "estafette",
				"jobType":   ciBuilderParams.JobType,
			},
		},
		Data: map[string][]byte{
			"secretDecryptionKey": []byte(newKey),
		},
	}

	err = c.kubeClient.Create(context.Background(), secret)
	if err != nil {
		return
	}

	log.Info().Msgf("Secret %v is created", decryptionKeySecretName)

	affinity := &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []*corev1.PreferredSchedulingTerm{
				&corev1.PreferredSchedulingTerm{
					Weight: &preemptibleAffinityWeight,
					Preference: &corev1.NodeSelectorTerm{
						MatchExpressions: []*corev1.NodeSelectorRequirement{
							&corev1.NodeSelectorRequirement{
								Key:      &preemptibleAffinityKey,
								Operator: &preemptibleAffinityOperator,
								Values:   []string{"true"},
							},
						},
					},
				},
			},
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []*corev1.NodeSelectorTerm{
					&corev1.NodeSelectorTerm{
						MatchExpressions: []*corev1.NodeSelectorRequirement{
							&corev1.NodeSelectorRequirement{
								Key:      &operatingSystemAffinityKey,
								Operator: &operatingSystemAffinityOperator,
								Values:   []string{operatingSystemAffinityValue},
							},
						},
					},
				},
			},
		},
	}

	tolerations := []*corev1.Toleration{}

	if ciBuilderParams.JobType == "release" {
		// keep off of preemptibles
		preemptibleAffinityOperator := "DoesNotExist"

		affinity = &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []*corev1.NodeSelectorTerm{
						&corev1.NodeSelectorTerm{
							MatchExpressions: []*corev1.NodeSelectorRequirement{
								&corev1.NodeSelectorRequirement{
									Key:      &preemptibleAffinityKey,
									Operator: &preemptibleAffinityOperator,
								},
							},
						},
						&corev1.NodeSelectorTerm{
							MatchExpressions: []*corev1.NodeSelectorRequirement{
								&corev1.NodeSelectorRequirement{
									Key:      &operatingSystemAffinityKey,
									Operator: &operatingSystemAffinityOperator,
									Values:   []string{operatingSystemAffinityValue},
								},
							},
						},
					},
				},
			},
		}
	}

	volumes := []*corev1.Volume{
		&corev1.Volume{
			Name: &builderConfigVolumeName,
			VolumeSource: &corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: &corev1.LocalObjectReference{
						Name: &jobName,
					},
				},
			},
		},
		&corev1.Volume{
			Name: &decryptionKeySecretVolumeName,
			VolumeSource: &corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: &jobName,
				},
			},
		},
	}
	volumeMounts := []*corev1.VolumeMount{
		&corev1.VolumeMount{
			Name:      &builderConfigVolumeName,
			MountPath: &builderConfigVolumeMountPath,
		},
		&corev1.VolumeMount{
			Name:      &decryptionKeySecretVolumeName,
			MountPath: &decryptionKeySecretVolumeMountPath,
		},
	}

	if ciBuilderParams.OperatingSystem == "windows" {
		// use emptydir volume in order to be able to have docker daemon on host mount path into internal container
		workingDirectoryVolumeName := "working-directory"
		volumes = append(volumes, &corev1.Volume{
			Name: &workingDirectoryVolumeName,
			VolumeSource: &corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})

		workingDirectoryVolumeMountPath := "C:/estafette-work"
		volumeMounts = append(volumeMounts, &corev1.VolumeMount{
			Name:      &workingDirectoryVolumeName,
			MountPath: &workingDirectoryVolumeMountPath,
		})

		// windows builds uses docker-outside-docker, for which the hosts docker socket needs to be mounted into the ci-builder container
		dockerSocketVolumeName := "docker-socket"
		dockerSocketVolumeHostPath := `\\.\pipe\docker_engine`
		volumes = append(volumes, &corev1.Volume{
			Name: &dockerSocketVolumeName,
			VolumeSource: &corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: &dockerSocketVolumeHostPath,
				},
			},
		})

		dockerSocketVolumeMountPath := `\\.\pipe\docker_engine`
		volumeMounts = append(volumeMounts, &corev1.VolumeMount{
			Name:      &dockerSocketVolumeName,
			MountPath: &dockerSocketVolumeMountPath,
		})

		// in order not to have to install the docker cli into the ci-builder container it's mounted from the host as well
		dockerCLIVolumeName := "docker-cli"
		dockerCLIVolumeHostPath := `C:/Program Files/Docker`
		volumes = append(volumes, &corev1.Volume{
			Name: &dockerCLIVolumeName,
			VolumeSource: &corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: &dockerCLIVolumeHostPath,
				},
			},
		})

		dockerCLIVolumeMountPath := `C:/Program Files/Docker`
		volumeMounts = append(volumeMounts, &corev1.VolumeMount{
			Name:      &dockerCLIVolumeName,
			MountPath: &dockerCLIVolumeMountPath,
		})

		// docker in kubernetes on windows is still at 18.09.7, which has api version 1.39
		// todo - use auto detect for the docker api version
		dockerAPIVersionName := "DOCKER_API_VERSION"
		dockerAPIVersionValue := "1.39"
		environmentVariables = append(environmentVariables,
			&corev1.EnvVar{
				Name:  &dockerAPIVersionName,
				Value: &dockerAPIVersionValue,
			},
		)

		podUIDName := "POD_UID"
		podUIDFieldPath := "metadata.uid"
		environmentVariables = append(environmentVariables,
			&corev1.EnvVar{
				Name: &podUIDName,
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: &podUIDFieldPath,
					},
				},
			},
		)

		// this is the path on the host mounted into any of the stage containers; with docker-outside-docker the daemon can't see paths inside the ci-builder container
		estafetteWorkdirName := "ESTAFETTE_WORKDIR"
		estafetteWorkdirValue := "c:/var/lib/kubelet/pods/$(POD_UID)/volumes/kubernetes.io~empty-dir/" + workingDirectoryVolumeName
		environmentVariables = append(environmentVariables,
			&corev1.EnvVar{
				Name:  &estafetteWorkdirName,
				Value: &estafetteWorkdirValue,
			},
		)

		tolerationEffect := "NoSchedule"
		tolerationKey := "node.kubernetes.io/os"
		tolerationOperator := "Equal"
		tolerationValue := "windows"
		tolerations = append(tolerations, &corev1.Toleration{
			Effect:   &tolerationEffect,
			Key:      &tolerationKey,
			Operator: &tolerationOperator,
			Value:    &tolerationValue,
		})
	}

	job = &batchv1.Job{
		Metadata: &metav1.ObjectMeta{
			Name:      &jobName,
			Namespace: &c.config.Jobs.Namespace,
			Labels: map[string]string{
				"createdBy": "estafette",
				"jobType":   ciBuilderParams.JobType,
			},
		},
		Spec: &batchv1.JobSpec{
			Template: &corev1.PodTemplateSpec{
				Metadata: &metav1.ObjectMeta{
					Labels: map[string]string{
						"createdBy": "estafette",
						"jobType":   ciBuilderParams.JobType,
					},
				},
				Spec: &corev1.PodSpec{
					Containers: []*corev1.Container{
						&corev1.Container{
							Name:            &containerName,
							Image:           &image,
							ImagePullPolicy: &imagePullPolicy,
							Args: []string{
								"--run-as-job",
							},
							Env: environmentVariables,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privileged,
							},
							Resources: &corev1.ResourceRequirements{
								Requests: map[string]*resource.Quantity{
									"cpu":    &resource.Quantity{String_: &cpuRequest},
									"memory": &resource.Quantity{String_: &memoryRequest},
								},
								Limits: map[string]*resource.Quantity{
									"cpu":    &resource.Quantity{String_: &cpuLimit},
									"memory": &resource.Quantity{String_: &memoryLimit},
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					RestartPolicy: &restartPolicy,

					Volumes: volumes,

					Affinity: affinity,

					Tolerations: tolerations,
				},
			},
		},
	}

	// "error":"unregistered type *v1.Job",
	err = c.kubeClient.Create(context.Background(), job)

	if err != nil {
		return
	}

	log.Info().Msgf("Job %v is created", jobName)

	return
}

// RemoveCiBuilderJob waits for a job to finish and then removes it
func (c *client) RemoveCiBuilderJob(ctx context.Context, jobName string) (err error) {

	log.Info().Msgf("Deleting job %v...", jobName)

	// check if job is finished
	var job batchv1.Job
	err = c.kubeClient.Get(context.Background(), c.config.Jobs.Namespace, jobName, &job)
	if err != nil {
		log.Error().Err(err).
			Str("jobName", jobName).
			Msgf("Get call for job %v failed", jobName)
	}

	if err != nil || *job.Status.Succeeded != 1 {
		log.Debug().Str("jobName", jobName).Msgf("Job is not done yet, watching for job %v to succeed", jobName)

		// watch for job updates
		var job batchv1.Job
		watcher, err := c.kubeClient.Watch(context.Background(), c.config.Jobs.Namespace, &job, k8s.Timeout(time.Duration(300)*time.Second))
		defer watcher.Close()

		if err != nil {
			log.Error().Err(err).
				Str("jobName", jobName).
				Msgf("Watcher call for job %v failed", jobName)
		} else {
			// wait for job to succeed
			for {
				job := new(batchv1.Job)
				event, err := watcher.Next(job)
				if err != nil {
					log.Error().Err(err)
					break
				}

				if event == k8s.EventModified && *job.Metadata.Name == jobName && *job.Status.Succeeded == 1 {
					break
				}
			}
		}
	}

	// delete job
	err = c.kubeClient.Delete(context.Background(), &job)
	if err != nil {
		log.Error().Err(err).
			Str("jobName", jobName).
			Msgf("Deleting job %v failed", jobName)
		return
	}

	log.Info().Msgf("Job %v is deleted", jobName)

	c.RemoveCiBuilderConfigMap(ctx, jobName)
	c.RemoveCiBuilderSecret(ctx, jobName)

	return
}

func (c *client) RemoveCiBuilderConfigMap(ctx context.Context, configmapName string) (err error) {

	// check if configmap exists
	var configmap corev1.ConfigMap
	err = c.kubeClient.Get(context.Background(), c.config.Jobs.Namespace, configmapName, &configmap)
	if err != nil {
		log.Error().Err(err).
			Str("configmap", configmapName).
			Msgf("Get call for configmap %v failed", configmapName)
		return
	}

	// delete configmap
	err = c.kubeClient.Delete(context.Background(), &configmap)
	if err != nil {
		log.Error().Err(err).
			Str("configmap", configmapName).
			Msgf("Deleting configmap %v failed", configmapName)
		return
	}

	log.Info().Msgf("Configmap %v is deleted", configmapName)

	return
}

func (c *client) RemoveCiBuilderSecret(ctx context.Context, secretName string) (err error) {

	// check if secret exists
	var secret corev1.Secret
	err = c.kubeClient.Get(context.Background(), c.config.Jobs.Namespace, secretName, &secret)
	if err != nil {
		log.Error().Err(err).
			Str("secret", secretName).
			Msgf("Get call for secret %v failed", secretName)
		return
	}

	// delete secret
	err = c.kubeClient.Delete(context.Background(), &secret)
	if err != nil {
		log.Error().Err(err).
			Str("secret", secretName).
			Msgf("Deleting secret %v failed", secretName)
		return
	}

	log.Info().Msgf("Secret %v is deleted", secretName)

	return
}

// CancelCiBuilderJob removes a job and its pods to cancel a build/release
func (c *client) CancelCiBuilderJob(ctx context.Context, jobName string) (err error) {

	log.Info().Msgf("Canceling job %v...", jobName)

	// check if job is finished
	var job batchv1.Job
	err = c.kubeClient.Get(context.Background(), c.config.Jobs.Namespace, jobName, &job)
	if err != nil {
		log.Error().Err(err).
			Str("jobName", jobName).
			Msgf("Get call for job %v failed", jobName)
		return
	}

	// delete job
	err = c.kubeClient.Delete(context.Background(), &job)
	if err != nil {
		log.Error().Err(err).
			Str("jobName", jobName).
			Msgf("Canceling job %v failed", jobName)
		return
	}

	log.Info().Msgf("Job %v is canceled", jobName)

	c.RemoveCiBuilderConfigMap(ctx, jobName)
	c.RemoveCiBuilderSecret(ctx, jobName)

	return
}

// TailCiBuilderJobLogs tails logs of a running job
func (c *client) TailCiBuilderJobLogs(ctx context.Context, jobName string, logChannel chan contracts.TailLogLine) (err error) {

	// close channel so api handler can finish it's response
	defer close(logChannel)

	labels := new(k8s.LabelSelector)
	labels.Eq("job-name", jobName)

	log.Debug().Msgf("TailCiBuilderJobLogs - listing pods with job-name=%v namespace=%v", jobName, c.config.Jobs.Namespace)

	var pods corev1.PodList
	if err := c.kubeClient.List(context.Background(), c.config.Jobs.Namespace, &pods, labels.Selector()); err != nil {
		return err
	}

	log.Debug().Msgf("TailCiBuilderJobLogs - retrieved %v pods", len(pods.Items))

	for _, pod := range pods.Items {

		if *pod.Status.Phase == "Pending" {

			log.Debug().Msg("TailCiBuilderJobLogs - pod is pending, waiting for running state...")

			// watch for pod to go into Running state (or out of Pending state)
			var pendingPod corev1.Pod
			watcher, err := c.kubeClient.Watch(context.Background(), c.config.Jobs.Namespace, &pendingPod, k8s.Timeout(time.Duration(300)*time.Second))

			if err != nil {
				return err
			}

			// wait for pod to change Phase to succeed
			defer watcher.Close()
			for {
				watchedPod := new(corev1.Pod)
				event, err := watcher.Next(watchedPod)
				if err != nil {
					return err
				}

				if event == k8s.EventModified && *watchedPod.Metadata.Name == *pod.Metadata.Name && *watchedPod.Status.Phase != "Pending" {
					pod = watchedPod
					break
				}
			}
		}

		if *pod.Status.Phase != "Running" {
			log.Warn().Msgf("Post %v for job %v has unsupported phase %v", *pod.Metadata.Name, jobName, *pod.Status.Phase)
		}

		log.Debug().Msg("TailCiBuilderJobLogs - pod has running state...")

		// follow logs from pod
		url := fmt.Sprintf("%v/api/v1/namespaces/%v/pods/%v/log?follow=true", c.kubeClient.Endpoint, c.config.Jobs.Namespace, *pod.Metadata.Name)

		log.Debug().Msgf("TailCiBuilderJobLogs - k8s api url: %v", url)

		ct := "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error().Err(err).Msgf("Failed generating request for retrieving logs from pod %v for job %v", *pod.Metadata.Name, jobName)
			return err
		}
		if c.kubeClient.SetHeaders != nil {
			if err := c.kubeClient.SetHeaders(req.Header); err != nil {
				log.Error().Err(err).Msgf("Failed setting request headers for retrieving logs from pod %v for job %v", *pod.Metadata.Name, jobName)
				return err
			}
		}
		req = req.WithContext(context.Background())

		req.Header.Set("Accept", ct)

		resp, err := c.kubeClient.Client.Do(req)
		if err != nil {
			log.Error().Err(err).Msgf("Failed performing request for retrieving logs from pod %v for job %v", *pod.Metadata.Name, jobName)
			return err
		}

		if resp.StatusCode/100 != 2 {
			errorMessage := fmt.Sprintf("Request for retrieving logs from pod %v for job %v has status code %v", *pod.Metadata.Name, jobName, resp.StatusCode)
			log.Error().Msg(errorMessage)
			return fmt.Errorf(errorMessage)
		}

		log.Debug().Msgf("TailCiBuilderJobLogs - streaming logs")

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Warn().Err(err).Msgf("Error while reading lines from logs from pod %v for job %v", *pod.Metadata.Name, jobName)
			}

			// only forward if it's a json object with property 'tailLogLine'
			var zeroLogLine ZeroLogLine
			err = json.Unmarshal(line, &zeroLogLine)
			if err == nil {
				if zeroLogLine.TailLogLine != nil {
					logChannel <- *zeroLogLine.TailLogLine
				}
			} else {
				log.Error().Err(err).Str("line", string(line)).Msgf("Tailed log from pod %v for job %v is not of type json", *pod.Metadata.Name, jobName)
			}
		}
	}

	return
}

// GetJobName returns the job name for a build or release job
func (c *client) GetJobName(ctx context.Context, jobType, repoOwner, repoName, id string) string {

	// create job name of max 63 chars
	maxJobNameLength := 63

	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	fullRepoName := re.ReplaceAllString(fmt.Sprintf("%v/%v", repoOwner, repoName), "-")

	maxRepoNameLength := maxJobNameLength - len(jobType) - 1 - len(id) - 1
	if len(fullRepoName) > maxRepoNameLength {
		fullRepoName = fullRepoName[:maxRepoNameLength]
	}

	return strings.ToLower(fmt.Sprintf("%v-%v-%v", jobType, fullRepoName, id))
}

// GetJobName returns the job name for a build or release job
func (c *client) GetBuilderConfig(ctx context.Context, ciBuilderParams CiBuilderParams, jobName string) (contracts.BuilderConfig, error) {

	// retrieve stages to filter trusted images and credentials
	stages := ciBuilderParams.Manifest.Stages
	if ciBuilderParams.JobType == "release" {

		releaseExists := false
		for _, r := range ciBuilderParams.Manifest.Releases {
			if r.Name == ciBuilderParams.ReleaseName {
				releaseExists = true
				stages = r.Stages
			}
		}
		if !releaseExists {
			stages = []*manifest.EstafetteStage{}
		}
	}

	// get configured credentials
	credentials := c.encryptedConfig.Credentials

	// add dynamic github api token credential
	if token, ok := ciBuilderParams.EnvironmentVariables["ESTAFETTE_GITHUB_API_TOKEN"]; ok {

		encryptedTokenEnvelope, err := c.secretHelper.EncryptEnvelope(token, crypt.DefaultPipelineWhitelist)
		if err != nil {
			return contracts.BuilderConfig{}, err
		}

		credentials = append(credentials, &contracts.CredentialConfig{
			Name: "github-api-token",
			Type: "github-api-token",
			AdditionalProperties: map[string]interface{}{
				"token": encryptedTokenEnvelope,
			},
		})
	}

	// add dynamic bitbucket api token credential
	if token, ok := ciBuilderParams.EnvironmentVariables["ESTAFETTE_BITBUCKET_API_TOKEN"]; ok {

		encryptedTokenEnvelope, err := c.secretHelper.EncryptEnvelope(token, crypt.DefaultPipelineWhitelist)
		if err != nil {
			return contracts.BuilderConfig{}, err
		}

		credentials = append(credentials, &contracts.CredentialConfig{
			Name: "bitbucket-api-token",
			Type: "bitbucket-api-token",
			AdditionalProperties: map[string]interface{}{
				"token": encryptedTokenEnvelope,
			},
		})
	}

	// add dynamic cloudsource api token credential
	if token, ok := ciBuilderParams.EnvironmentVariables["ESTAFETTE_CLOUDSOURCE_API_TOKEN"]; ok {

		encryptedTokenEnvelope, err := c.secretHelper.EncryptEnvelope(token, crypt.DefaultPipelineWhitelist)
		if err != nil {
			return contracts.BuilderConfig{}, err
		}

		credentials = append(credentials, &contracts.CredentialConfig{
			Name: "cloudsource-api-token",
			Type: "cloudsource-api-token",
			AdditionalProperties: map[string]interface{}{
				"token": encryptedTokenEnvelope,
			},
		})
	}

	// filter to only what's needed by the build/release job
	trustedImages := contracts.FilterTrustedImages(c.encryptedConfig.TrustedImages, stages, ciBuilderParams.GetFullRepoPath())
	credentials = contracts.FilterCredentials(credentials, trustedImages, ciBuilderParams.GetFullRepoPath())

	// add container-registry credentials to allow private registry images to be used in stages
	credentials = contracts.AddCredentialsIfNotPresent(credentials, contracts.FilterCredentialsByPipelinesWhitelist(contracts.GetCredentialsByType(c.encryptedConfig.Credentials, "container-registry"), ciBuilderParams.GetFullRepoPath()))

	localBuilderConfig := contracts.BuilderConfig{
		Credentials:     credentials,
		TrustedImages:   trustedImages,
		RegistryMirror:  c.config.RegistryMirror,
		DockerDaemonMTU: c.config.DockerDaemonMTU,
		DockerDaemonBIP: c.config.DockerDaemonBIP,
		DockerNetwork:   c.config.DockerNetwork,
	}

	localBuilderConfig.Action = &ciBuilderParams.JobType
	localBuilderConfig.Track = &ciBuilderParams.Track
	localBuilderConfig.Git = &contracts.GitConfig{
		RepoSource:   ciBuilderParams.RepoSource,
		RepoOwner:    ciBuilderParams.RepoOwner,
		RepoName:     ciBuilderParams.RepoName,
		RepoBranch:   ciBuilderParams.RepoBranch,
		RepoRevision: ciBuilderParams.RepoRevision,
	}
	if ciBuilderParams.Manifest.Version.SemVer != nil {
		versionParams := manifest.EstafetteVersionParams{
			AutoIncrement: ciBuilderParams.AutoIncrement,
			Branch:        ciBuilderParams.RepoBranch,
			Revision:      ciBuilderParams.RepoRevision,
		}
		patchWithLabel := ciBuilderParams.Manifest.Version.SemVer.GetPatchWithLabel(versionParams)
		label := ciBuilderParams.Manifest.Version.SemVer.GetLabel(versionParams)
		localBuilderConfig.BuildVersion = &contracts.BuildVersionConfig{
			Version:       ciBuilderParams.VersionNumber,
			Major:         &ciBuilderParams.Manifest.Version.SemVer.Major,
			Minor:         &ciBuilderParams.Manifest.Version.SemVer.Minor,
			Patch:         &patchWithLabel,
			Label:         &label,
			AutoIncrement: &ciBuilderParams.AutoIncrement,
		}
	} else {
		localBuilderConfig.BuildVersion = &contracts.BuildVersionConfig{
			Version:       ciBuilderParams.VersionNumber,
			AutoIncrement: &ciBuilderParams.AutoIncrement,
		}
	}

	localBuilderConfig.Manifest = &ciBuilderParams.Manifest

	localBuilderConfig.JobName = &jobName
	localBuilderConfig.CIServer = &contracts.CIServerConfig{
		BaseURL:          c.config.APIServer.BaseURL,
		BuilderEventsURL: strings.TrimRight(c.config.APIServer.ServiceURL, "/") + "/api/commands",
		PostLogsURL:      strings.TrimRight(c.config.APIServer.ServiceURL, "/") + fmt.Sprintf("/api/pipelines/%v/%v/%v/builds/%v/logs", ciBuilderParams.RepoSource, ciBuilderParams.RepoOwner, ciBuilderParams.RepoName, ciBuilderParams.BuildID),
		APIKey:           c.config.Auth.APIKey,
	}

	if ciBuilderParams.ReleaseID > 0 {
		localBuilderConfig.CIServer.PostLogsURL = strings.TrimRight(c.config.APIServer.ServiceURL, "/") + fmt.Sprintf("/api/pipelines/%v/%v/%v/releases/%v/logs", ciBuilderParams.RepoSource, ciBuilderParams.RepoOwner, ciBuilderParams.RepoName, ciBuilderParams.ReleaseID)
	}

	if *localBuilderConfig.Action == "build" {
		localBuilderConfig.BuildParams = &contracts.BuildParamsConfig{
			BuildID: ciBuilderParams.BuildID,
		}
	}
	if *localBuilderConfig.Action == "release" {
		localBuilderConfig.ReleaseParams = &contracts.ReleaseParamsConfig{
			ReleaseName:   ciBuilderParams.ReleaseName,
			ReleaseID:     ciBuilderParams.ReleaseID,
			ReleaseAction: ciBuilderParams.ReleaseAction,
			TriggeredBy:   ciBuilderParams.ReleaseTriggeredBy,
		}
	}

	localBuilderConfig.Events = make([]*manifest.EstafetteEvent, 0)
	for _, e := range ciBuilderParams.TriggeredByEvents {
		localBuilderConfig.Events = append(localBuilderConfig.Events, &e)
	}

	return localBuilderConfig, nil
}
