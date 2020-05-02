package workspaceapi

import (
	"github.com/nkhang/pluto/internal/workspace"
)

type WorkspaceInfoResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func toWorkspaceInfoResponse(w workspace.Workspace) WorkspaceInfoResponse {
	return WorkspaceInfoResponse{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
	}
}
