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
	persistenceProvider   *AWSS3PersistenceProvider
)

func initProviders() error {
	start := time.Now()
	fmt.Printf("Starting initialization at: %v\n", start)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}
	fmt.Printf("Config loaded, time elapsed: %v\n", time.Since(start))

	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")
	fmt.Printf("Environment variables - Bucket: %s, Key: %s\n", bucket, key)

	if bucket == "" || key == "" {
		return fmt.Errorf("BUCKET and KEY environment variables must be set")
	}

	pipelineStateProvider = &AWSPipelineStateProvider{cfg}
	persistenceProvider = &AWSS3PersistenceProvider{
		config: cfg,
		bucket: bucket,
		key:    key,
	}
	fmt.Printf("Successfully initialized providers, total init time: %v\n", time.Since(start))
	return nil
}

func updateProjectsStatus(stateProvider PipelineStateProvider, persistenceProvider PersistenceProvider) error {
	start := time.Now()
	fmt.Printf("Starting updateProjectsStatus at: %v\n", start)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	pipelineFetchStart := time.Now()
	fmt.Println("Fetching pipeline states...")
	pipelineStates, err := stateProvider.GetPipelineState()
	if err != nil {
		return fmt.Errorf("unable to get pipeline state: %v", err)
	}
	fmt.Printf("Pipeline states fetched in %v, got %d states\n",
		time.Since(pipelineFetchStart), len(pipelineStates))

	convertStart := time.Now()
	fmt.Println("Converting pipeline states to projects...")
	projects := Convert(pipelineStates)
	fmt.Printf("Conversion completed in %v, got %d projects\n",
		time.Since(convertStart), len(projects))

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation timed out")
	default:
		persistStart := time.Now()
		fmt.Println("Persisting projects to S3...")
		if err := persistenceProvider.PersistProjects(projects); err != nil {
			return fmt.Errorf("unable to persist projects data: %v", err)
		}
		fmt.Printf("Projects persisted in %v\n", time.Since(persistStart))
	}

	fmt.Printf("Total execution time: %v\n", time.Since(start))
	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	if err := initProviders(); err != nil {
		fmt.Printf("Initialization failed: %v\n", err)
		return err
	}

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
