package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	pipelineStateProvider *AWSPipelineStateProvider
	persistenceProvider *AWSS3PersistenceProvider
)

func init() {
	// 在初始化阶段创建 AWS 客户端
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	
	// 获取环境变量
	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")
	if bucket == "" || key == "" {
		log.Fatalf("BUCKET and KEY environment variables must be set")
	}

	// 初始化 providers
	pipelineStateProvider = &AWSPipelineStateProvider{cfg}
	persistenceProvider = &AWSS3PersistenceProvider{
		cfg,
		bucket,
		key,
	}
}

func updateProjectsStatus(stateProvider PipelineStateProvider, persistenceProvider PersistenceProvider) error {
	log.Printf("Starting updateProjectsStatus")
	pipelineStates, err := stateProvider.GetPipelineState()
	if err != nil {
		log.Printf("Error getting pipeline state: %v", err)
		return fmt.Errorf("unable to get state pipeline state: %v", err)
	}
	log.Printf("Successfully retrieved pipeline states")

	err = persistenceProvider.PersistProjects(Convert(pipelineStates))
	if err != nil {
		log.Printf("Error persisting projects: %v", err)
		return fmt.Errorf("unable to persist projects data: %v", err)
	}
	log.Printf("Successfully persisted projects")

	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	log.Printf("Received event: %s", string(event))

	err := updateProjectsStatus(pipelineStateProvider, persistenceProvider)
	if err != nil {
		log.Printf("Error updating projects status: %v", err)
		return err
	}

	log.Printf("Successfully processed pipeline state update")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
