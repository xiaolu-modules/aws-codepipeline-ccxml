package main

import (
	"time"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// PipelineState captures the current state of a pipeline
type PipelineState struct {
	Name        string
	Created     time.Time
	Region      string
	StageStates []types.StageState
}

// PipelineStateProvider provides access to the current state of a pipeline
type PipelineStateProvider interface {
	// GetPipelineState returns the current state of a pipeline
	GetPipelineState() ([]PipelineState, error)
}

// AWSPipelineStateProvider provides access to the current state of a pipeline using the AWS API
type AWSPipelineStateProvider struct {
	config aws.Config
}

// GetPipelineState provides access to the current state of a pipeline using the AWS API
func (p *AWSPipelineStateProvider) GetPipelineState() ([]PipelineState, error) {
	client := codepipeline.NewFromConfig(p.config)

	resp, err := client.ListPipelines(context.Background(), &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, err
	}

	pipelineStates := make([]PipelineState, 0)

	for _, pipeline := range resp.Pipelines {
		 stageStates, err := client.GetPipelineState(context.Background(), &codepipeline.GetPipelineStateInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, err
		}

		pipelineStates = append(pipelineStates, PipelineState{
			Name:        *pipeline.Name,
			Created:     *pipeline.Created,
			Region:      p.config.Region,
			StageStates: stageStates.StageStates,
		})
	}

	return pipelineStates, nil
}
