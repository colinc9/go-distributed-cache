package network

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)
// Get all addresses (all tasks and their ports)
//		ListTasks, DescribeTasks, DescribeNetworkInterfaces
//		Magic ip: http://169.254.170.2/v2/metadata
//		http://169.254.169.254/latest/meta-data/
//		service discovery?
// Get own port and identify each one with metainfo
// set the myaddress and targetaddress
// poll?

func sdkDiscovery() (*ecs.DescribeTasksOutput, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}
	svc := ecs.NewFromConfig(cfg)
	input := &ecs.ListTasksInput{
		Cluster: aws.String("app-cluster"),
	}
	result, err := svc.ListTasks(ctx, input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	taskInput := &ecs.DescribeTasksInput{
		Tasks: result.TaskArns,
		Cluster: aws.String("app-cluster"),
	}
	output, err := svc.DescribeTasks(ctx, taskInput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

const (
	ecsMetadataUriEnvV3 = "ECS_CONTAINER_METADATA_URI"
)

type ContainerMetadataV3 struct {
	DockerID   string `json:"DockerId"`
	Name       string `json:"Name"`
	DockerName string `json:"DockerName"`
	Image      string `json:"Image"`
	ImageID    string `json:"ImageID"`
	Labels     struct {
		EcsCluster               string `json:"com.amazonaws.ecs.cluster"`
		EcsContainerName         string `json:"com.amazonaws.ecs.container-name"`
		EcsTaskArn               string `json:"com.amazonaws.ecs.task-arn"`
		EcsTaskDefinitionFamily  string `json:"com.amazonaws.ecs.task-definition-family"`
		EcsTaskDefinitionVersion string `json:"com.amazonaws.ecs.task-definition-version"`
	} `json:"Labels"`
	DesiredStatus string `json:"DesiredStatus"`
	KnownStatus   string `json:"KnownStatus"`
	Limits        struct {
		CPU    int `json:"CPU"`
		Memory int `json:"Memory"`
	} `json:"Limits"`
	CreatedAt time.Time `json:"CreatedAt"`
	StartedAt time.Time `json:"StartedAt,omitempty"`
	Type      string    `json:"Type"`
	Networks  []struct {
		NetworkMode   string   `json:"NetworkMode"`
		IPv4Addresses []string `json:"IPv4Addresses"`
	} `json:"Networks"`
}

type TaskMetadataV3 struct {
	Cluster       string                `json:"Cluster"`
	TaskARN       string                `json:"TaskARN"`
	Family        string                `json:"Family"`
	Revision      string                `json:"Revision"`
	DesiredStatus string                `json:"DesiredStatus"`
	KnownStatus   string                `json:"KnownStatus"`
	Containers    []ContainerMetadataV3 `json:"Containers"`
	Limits        struct {
		CPU    float64 `json:"CPU"`
		Memory int     `json:"Memory"`
	} `json:"Limits"`
	PullStartedAt time.Time `json:"PullStartedAt"`
	PullStoppedAt time.Time `json:"PullStoppedAt"`
}

// Retrieve ECS Task Metadata in V3 format
func GetTaskV3(ctx context.Context, client *http.Client) (*TaskMetadataV3, error) {
	metadataUrl := os.Getenv(ecsMetadataUriEnvV3)
	if metadataUrl == "" {
		return nil, fmt.Errorf("missing metadata uri in environment (%s)", ecsMetadataUriEnvV3)
	}

	taskMetadata := &TaskMetadataV3{}
	body, err := fetch(ctx, client, fmt.Sprintf("%s/task", metadataUrl))

	err = json.Unmarshal(body, &taskMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal into task metadata (v3): %w", err)
	}

	return taskMetadata, nil
}

//keep querying


// Retrieve ECS Container Metadata in V3 format
func GetContainerV3(ctx context.Context, client *http.Client) (*ContainerMetadataV3, error) {
	metadataUrl := os.Getenv(ecsMetadataUriEnvV3)
	if metadataUrl == "" {
		return nil, fmt.Errorf("missing metadata uri in environment (%s)", ecsMetadataUriEnvV3)
	}

	contaienrMetadata := &ContainerMetadataV3{}
	body, err := fetch(ctx, client, fmt.Sprintf("%s", metadataUrl))

	err = json.Unmarshal(body, &contaienrMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal into container metadata (v3): %w", err)
	}

	return contaienrMetadata, nil
}

func fetch(ctx context.Context, client *http.Client, metadataUrl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, metadataUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create metadata request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send metadata request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metadata response: %w", err)
	}

	return body, nil
}

