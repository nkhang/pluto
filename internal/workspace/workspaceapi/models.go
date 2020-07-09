package workspaceapi

import (
	"github.com/nkhang/pluto/internal/workspace"
)

type GetByUserIDRequest struct {
	UserID   uint64 `form:"user_id" binding:"required"`
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Source   int    `form:"src" binding:"required"`
}

type CreateWorkspaceRequest struct {
	Title       string   `form:"title" json:"title"`
	Description string   `form:"description" json:"description"`
	Color       string   `form:"color" json:"color"`
	Admin       uint64   `form:"admin" json:"admin" binding:"required"`
	Members     []uint64 `form:"members" json:"members"`
}

type UpdateWorkspaceRequest struct {
	Title       string `form:"title" json:"title,omitempty"`
	Description string `form:"description" json:"description,omitempty"`
}

type GetByUserResponse struct {
	Total      int                 `json:"total"`
	Workspaces []WorkspaceResponse `json:"workspaces"`
}

type WorkspaceResponse struct {
	ID           uint64 `json:"id"`
	Updated      int64  `json:"updated"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	ProjectCount int    `json:"project_count"`
	MemberCount  int    `json:"member_count"`
	Admin        uint64 `json:"admin"`
}

func toWorkspaceInfoResponse(w workspace.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
	}
}
