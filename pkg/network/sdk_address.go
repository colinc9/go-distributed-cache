package network

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func PeriodicSdkDiscovery() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				getAddresses()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func getAddresses() {
	tasks, err := getTasks()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	c := &http.Client{Timeout: time.Duration(1) * time.Second}
	currTask, err := getTaskV4(context.Background(), c)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	revision, err := strconv.ParseInt(currTask.Revision, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	currIp := currTask.Networks[0].IPv4Addresses[0]
	MyAddress = currIp
	log.Printf("tasks metadata: %+v ", tasks)
	var targets []string
	for _, task := range tasks.Tasks {
		if task.Version == revision {
			attachments := task.Attachments
			for _, attachment := range attachments {
				if *attachment.Type == "ElasticNetworkInterface" && *attachment.Status == "ATTACHED"{
					details := attachment.Details
					for _, detail := range details {
						if *detail.Name == "privateIPv4Address" {
							if currIp != *detail.Value {
								targets = append(targets, *detail.Value)
							}
							break
						}
					}
					break
				}
			}

		}
	}
	TargetAddress = targets
}

func getTasks() (*ecs.DescribeTasksOutput, error) {
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
	ecsMetadataUriEnvV4 = "ECS_CONTAINER_METADATA_URI_V4"
)

type Limits struct {
	CPU    float64 `json:"CPU"`
	Memory int     `json:"Memory"`
}

type ContainerMetadataV4 struct {
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
	DesiredStatus string    `json:"DesiredStatus"`
	KnownStatus   string    `json:"KnownStatus"`
	Limits        Limits    `json:"Limits"`
	CreatedAt     time.Time `json:"CreatedAt"`
	StartedAt     time.Time `json:"StartedAt"`
	Type          string    `json:"Type"`
	ContainerARN  string    `json:"ContainerARN"`
	LogDriver     string    `json:"LogDriver"`
	LogOptions    struct {
		AwsLogsCreateGroup string `json:"awslogs-create-group"`
		AwsLogsGroup       string `json:"awslogs-group"`
		AwsLogsStream      string `json:"awslogs-stream"`
		AwsRegion          string `json:"awslogs-region"`
	} `json:"LogOptions"`
	Networks []struct {
		NetworkMode              string   `json:"NetworkMode"`
		IPv4Addresses            []string `json:"IPv4Addresses"`
		AttachmentIndex          int      `json:"AttachmentIndex"`
		IPv4SubnetCIDRBlock      string   `json:"IPv4SubnetCIDRBlock"`
		MACAddress               string   `json:"MACAddress"`
		DomainNameServers        []string `json:"DomainNameServers"`
		DomainNameSearchList     []string `json:"DomainNameSearchList"`
		PrivateDNSName           string   `json:"PrivateDNSName"`
		SubnetGatewayIpv4Address string   `json:"SubnetGatewayIpv4Address"`
	} `json:"Networks"`
}

type TaskMetadataV4 struct {
	Cluster          string                `json:"Cluster"`
	TaskARN          string                `json:"TaskARN"`
	Family           string                `json:"Family"`
	Revision         string                `json:"Revision"`
	DesiredStatus    string                `json:"DesiredStatus"`
	KnownStatus      string                `json:"KnownStatus"`
	Limits           Limits                `json:"Limits"`
	PullStartedAt    time.Time             `json:"PullStartedAt"`
	PullStoppedAt    time.Time             `json:"PullStoppedAt"`
	AvailabilityZone string                `json:"AvailabilityZone"`
	LaunchType       string                `json:"LaunchType"`
	Containers       []ContainerMetadataV4 `json:"Containers"`
	Networks []struct {
		NetworkMode              string   `json:"NetworkMode"`
		IPv4Addresses            []string `json:"IPv4Addresses"`
		AttachmentIndex          int      `json:"AttachmentIndex"`
		IPv4SubnetCIDRBlock      string   `json:"IPv4SubnetCIDRBlock"`
		MACAddress               string   `json:"MACAddress"`
		DomainNameServers        []string `json:"DomainNameServers"`
		DomainNameSearchList     []string `json:"DomainNameSearchList"`
		PrivateDNSName           string   `json:"PrivateDNSName"`
		SubnetGatewayIpv4Address string   `json:"SubnetGatewayIpv4Address"`
	} `json:"Networks"`
}

// Retrieve ECS Task Metadata in V4 format
func getTaskV4(ctx context.Context, client *http.Client) (*TaskMetadataV4, error) {
	metadataUrl := os.Getenv(ecsMetadataUriEnvV4)
	if metadataUrl == "" {
		return nil, fmt.Errorf("missing metadata uri in environment (%s)", ecsMetadataUriEnvV4)
	}

	taskMetadata := &TaskMetadataV4{}
	body, err := fetch(ctx, client, fmt.Sprintf("%s/task", metadataUrl))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve task metadata v4: %w", err)
	}

	err = json.Unmarshal(body, &taskMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal into task metadata v4: %w", err)
	}

	return taskMetadata, nil
}

// Retrieve ECS Container Metadata in V4 format
func getContainerV4(ctx context.Context, client *http.Client) (*ContainerMetadataV4, error) {
	metadataUrl := os.Getenv(ecsMetadataUriEnvV4)
	if metadataUrl == "" {
		return nil, fmt.Errorf("missing metadata uri in environment (%s)", ecsMetadataUriEnvV4)
	}

	containerMetadata := &ContainerMetadataV4{}
	body, err := fetch(ctx, client, metadataUrl)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve container metadata v4: %w", err)
	}

	err = json.Unmarshal(body, &containerMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal into container metadata v4: %w", err)
	}

	return containerMetadata, nil
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
