package workspaceapi

import "github.com/nkhang/pluto/internal/workspace"

type Repository interface {
	GetByID(id uint64) (WorkspaceInfoResponse, error)
}

type repository struct {
	workspaceRepository workspace.Repository
}

func NewRepository(r workspace.Repository) *repository {
	return &repository{r}
}

func (r *repository) GetByID(id uint64) (WorkspaceInfoResponse, error) {
	w, err := r.workspaceRepository.Get(id)
	if err != nil {
		return WorkspaceInfoResponse{}, err
	}
	return toWorkspaceInfoResponse(w), nil
}
