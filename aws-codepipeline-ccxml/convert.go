package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// Convert the pipeline states to Projects
func Convert(pipelineStates []PipelineState) []Project {
	projects := make([]Project, 0)

	for _, pipeline := range pipelineStates {
		lastBuildStatus := LastBuildStatusSuccess
		activity := ActivitySleeping
		var lastBuildTime time.Time

		// 检查所有阶段的状态
		for _, stage := range pipeline.StageStates {
			stageStatus := buildLastBuildStatus(stage)
			if stageStatus == LastBuildStatusFailure {
				lastBuildStatus = LastBuildStatusFailure
			}

			stageActivity := buildActivity(stage)
			if stageActivity == ActivityBuilding {
				activity = ActivityBuilding
			}

			// 获取最新的构建时间
			stageTime := getStageTime(pipeline.Created, stage)
			if stageTime.After(lastBuildTime) {
				lastBuildTime = stageTime
			}
		}

		projects = append(projects, Project{
			Name:            pipeline.Name,
			LastBuildStatus: lastBuildStatus,
			Activity:        activity,
			LastBuildTime:   lastBuildTime.Format(time.RFC3339),
		})
	}

	return projects
}

func buildLastBuildStatus(stage types.StageState) LastBuildStatus {
	if stage.LatestExecution == nil {
		return LastBuildStatusUnknown
	}

	switch stage.LatestExecution.Status {
	case types.StageExecutionStatusFailed:
		return LastBuildStatusFailure
	case types.StageExecutionStatusSucceeded:
		return LastBuildStatusSuccess
	}

	// assume Success as no easy way to work out previous state
	return LastBuildStatusSuccess
}

func buildActivity(stage types.StageState) Activity {
	if stage.LatestExecution != nil && stage.LatestExecution.Status == types.StageExecutionStatusInProgress {
		return ActivityBuilding
	}

	return ActivitySleeping
}

func getStageTime(created time.Time, stage types.StageState) time.Time {
	if stage.ActionStates == nil || len(stage.ActionStates) == 0 ||
		stage.ActionStates[0].LatestExecution == nil || stage.ActionStates[0].LatestExecution.LastStatusChange == nil {
		return created
	}
	return *stage.ActionStates[0].LatestExecution.LastStatusChange
}
