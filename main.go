package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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
	start := time.Now()
	fmt.Printf("Starting updateProjectsStatus at: %v\n", start)
	
	// 跟踪获取pipeline状态的时间
	pipelineFetchStart := time.Now()
	fmt.Println("Fetching pipeline states...")
	pipelineStates, err := stateProvider.GetPipelineState()
	if err != nil {
		fmt.Printf("Error getting pipeline state: %v\n", err)
		return fmt.Errorf("unable to get state pipeline state: %v", err)
	}
	fmt.Printf("Pipeline states fetched in %v, got %d states\n", 
		time.Since(pipelineFetchStart), len(pipelineStates))

	// 跟踪转换时间
	convertStart := time.Now()
	fmt.Println("Converting pipeline states to projects...")
	projects := Convert(pipelineStates)
	fmt.Printf("Conversion completed in %v, got %d projects\n", 
		time.Since(convertStart), len(projects))
	
	// 跟踪持久化时间
	persistStart := time.Now()
	fmt.Println("Persisting projects to S3...")
	err = persistenceProvider.PersistProjects(projects)
	if err != nil {
		fmt.Printf("Error persisting projects: %v\n", err)
		return fmt.Errorf("unable to persist projects data: %v", err)
	}
	fmt.Printf("Projects persisted in %v\n", time.Since(persistStart))

	fmt.Printf("Total execution time: %v\n", time.Since(start))
	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	start := time.Now()
	fmt.Printf("Starting request handling at: %v\n", start)
	fmt.Printf("Received event: %s\n", string(event))

	err := updateProjectsStatus(pipelineStateProvider, persistenceProvider)
	if err != nil {
		fmt.Printf("Error updating projects status after %v: %v\n", 
			time.Since(start), err)
		return err
	}

	fmt.Printf("Request completed successfully in %v\n", time.Since(start))
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
