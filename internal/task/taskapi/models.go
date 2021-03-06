package taskapi

import (
	"github.com/nkhang/pluto/internal/dataset/datasetapi"
	"github.com/nkhang/pluto/internal/image/imageapi"
	"github.com/nkhang/pluto/internal/label/labelapi"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
)

type CreateTaskRequest struct {
	Title       string         `json:"title" form:"title"`
	Description string         `json:"description" form:"description"`
	DatasetID   uint64         `json:"dataset_id" form:"dataset_id" binding:"required"`
	Quantity    int            `json:"quantity" form:"quantity" binding:"required"`
	Assignees   []AssigneePair `json:"assignees" form:"assignees" binding:"required"`
}

const (
	SrcAllTasks uint32 = iota + 1
	SrcAssignerTasks
	SrcLabelingTasks
	SrcReviewingTasks
)

type GetTasksRequest struct {
	Source   uint32 `json:"src" form:"src" binding:"required"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type GetTaskResponse struct {
	Total int            `json:"total"`
	Tasks []TaskResponse `json:"tasks"`
}

type AssigneePair struct {
	Labeler  uint64 `json:"labeler" form:"labeler" binding:"required"`
	Reviewer uint64 `json:"reviewer" form:"reviewer" binding:"required"`
}

type GetTaskDetailsRequest struct {
	CurrentID uint64            `form:"current_id" json:"current_id"`
	PageSize  int               `form:"page_size" json:"page_size"`
	Status    task.DetailStatus `form:"status" json:"status"`
}

type UpdateTaskDetailRequest struct {
	Status task.DetailStatus `form:"status" json:"status"`
}

type NATSUpdateDetailRequest struct {
	TaskID   uint64            `json:"task"`
	DetailID uint64            `json:"task_detail"`
	Status   task.DetailStatus `form:"status" json:"status"`
}

type TaskDetailResponse struct {
	ID     uint64                 `json:"id"`
	Status int32                  `json:"status"`
	TaskID uint64                 `json:"task_id"`
	Image  imageapi.ImageResponse `json:"image"`
}

type TaskResponse struct {
	ID          uint64                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Project     ProjectObject              `json:"project"`
	Workspace   WorkspaceObject            `json:"workspace"`
	Assigner    uint64                     `json:"assigner"`
	Labeler     uint64                     `json:"labeler"`
	Reviewer    uint64                     `json:"reviewer"`
	Status      uint32                     `json:"status"`
	ImageCount  int                        `json:"image_count"`
	CreatedAt   int64                      `json:"created_at"`
	Dataset     datasetapi.DatasetResponse `json:"dataset"`
}

type ProjectObject struct {
	projectapi.ProjectBaseResponse
	ProjectManagers []uint64 `json:"project_managers"`
}

type WorkspaceObject struct {
	workspaceapi.WorkspaceBaseResponse
	Admin uint64 `json:"admin"`
}

type PushTaskMessage struct {
	Workspace workspaceapi.WorkspaceDetailResponse `json:"workspace"`
	Project   projectapi.ProjectResponse           `json:"project"`
	Dataset   datasetapi.DatasetResponse           `json:"dataset"`
	Tasks     []TaskResponse                       `json:"tasks"`
	Labels    []labelapi.LabelResponse             `json:"labels"`
}

func ToTaskDetailResponse(detail task.Detail) TaskDetailResponse {
	return TaskDetailResponse{
		ID:     detail.ID,
		Status: int32(detail.Status),
		TaskID: detail.TaskID,
		Image:  imageapi.ToImageResponse(detail.Image),
	}
}
