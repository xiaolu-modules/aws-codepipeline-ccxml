// +build integration

package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
)

func TestAWSGetPipelineState(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Errorf("TestAWSGetPipelineState() unable to load AWS config: %v", err)
	}

	pipelineStateProvider := AWSPipelineStateProvider{cfg}
	pipelineStates, err := pipelineStateProvider.GetPipelineState()
	if err != nil {
		t.Errorf("TestAWSGetPipelineState() unable to retrieve pipeline states: %v", err)
	}

	if len(pipelineStates) == 0 {
		t.Errorf("TestAWSGetPipelineState() doesn't return any state information")
	}
}
