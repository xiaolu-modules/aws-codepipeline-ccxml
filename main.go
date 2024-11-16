package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	pipelineStateProvider *AWSPipelineStateProvider
	persistenceProvider *AWSS3PersistenceProvider
)

func init() {
	start := time.Now()
	fmt.Printf("Starting initialization at: %v\n", start)
	
	// 在初始化阶段创建 AWS 客户端
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("unable to load SDK config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Config loaded, time elapsed: %v\n", time.Since(start))

	// 获取环境变量
	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")
	fmt.Printf("Environment variables - Bucket: %s, Key: %s\n", bucket, key)
	
	if bucket == "" || key == "" {
		fmt.Println("BUCKET and KEY environment variables must be set")
		os.Exit(1)
	}

	// 初始化 providers
	pipelineStateProvider = &AWSPipelineStateProvider{cfg}
	persistenceProvider = &AWSS3PersistenceProvider{
		cfg,
		bucket,
		key,
	}
	fmt.Printf("Successfully initialized providers, total init time: %v\n", time.Since(start))
}

func updateProjectsStatus(stateProvider PipelineStateProvider, persistenceProvider PersistenceProvider) error {
	fmt.Println("Starting updateProjectsStatus")
	
	fmt.Println("Fetching pipeline states...")
	pipelineStates, err := stateProvider.GetPipelineState()
	if err != nil {
		fmt.Printf("Error getting pipeline state: %v\n", err)
		return fmt.Errorf("unable to get state pipeline state: %v", err)
	}
	fmt.Printf("Successfully retrieved %d pipeline states\n", len(pipelineStates))

	fmt.Println("Converting and persisting projects...")
	projects := Convert(pipelineStates)
	fmt.Printf("Converted to %d projects\n", len(projects))
	
	err = persistenceProvider.PersistProjects(projects)
	if err != nil {
		fmt.Printf("Error persisting projects: %v\n", err)
		return fmt.Errorf("unable to persist projects data: %v", err)
	}
	fmt.Println("Successfully persisted projects")

	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	fmt.Printf("Handling request with context: %v\n", ctx)
	fmt.Printf("Received event: %s\n", string(event))

	fmt.Println("Starting projects status update...")
	err := updateProjectsStatus(pipelineStateProvider, persistenceProvider)
	if err != nil {
		fmt.Printf("Error updating projects status: %v\n", err)
		return err
	}

	fmt.Println("Successfully completed request processing")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
