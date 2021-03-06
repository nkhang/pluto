package workspaceapi

import (
	"encoding/json"

	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/paging"
)

type Repository interface {
	GetByID(id uint64) (WorkspaceDetailResponse, error)
	GetByUserID(userID uint64, request GetByUserIDRequest) (GetByUserResponse, error)
	CreateWorkspace(admin uint64, p CreateWorkspaceRequest) (WorkspaceDetailResponse, error)
	UpdateWorkspace(id uint64, request UpdateWorkspaceRequest) (WorkspaceDetailResponse, error)
	DeleteWorkspace(id uint64) error
}

type repository struct {
	workspaceRepository workspace.Repository
	projectRepo         project.Repository
}

func NewRepository(workspaceRepo workspace.Repository, projectRepo project.Repository) *repository {
	return &repository{
		workspaceRepository: workspaceRepo,
		projectRepo:         projectRepo,
	}
}

func (r *repository) GetByID(id uint64) (WorkspaceDetailResponse, error) {
	w, err := r.workspaceRepository.Get(id)
	if err != nil {
		return WorkspaceDetailResponse{}, err
	}
	return r.convertResponse(
		w), nil
}

func (r *repository) GetByUserID(userID uint64, request GetByUserIDRequest) (GetByUserResponse, error) {
	offset, limit := paging.Parse(request.Page, request.PageSize)
	var (
		workspaces []workspace.Workspace
		total      int
		err        error
	)
	switch request.Source {
	case 1:
		workspaces, total, err = r.workspaceRepository.GetByUserID(userID, workspace.Any, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	case 2:
		workspaces, total, err = r.workspaceRepository.GetByUserID(userID, workspace.Admin, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	case 3:
		workspaces, total, err = r.workspaceRepository.GetByUserID(userID, workspace.Member, offset, limit)
		if err != nil {
			return GetByUserResponse{}, err
		}
	default:
		return GetByUserResponse{}, errors.BadRequest.NewWithMessage("unsupported src")
	}
	responses := make([]WorkspaceDetailResponse, len(workspaces))
	for i := range workspaces {
		responses[i] = r.convertResponse(workspaces[i])
	}
	return GetByUserResponse{
		Total:      total,
		Workspaces: responses,
	}, nil
}

func (r *repository) CreateWorkspace(admin uint64, p CreateWorkspaceRequest) (WorkspaceDetailResponse, error) {
	w, err := r.workspaceRepository.Create(admin, p.Title, p.Description, p.Color)
	if err != nil {
		return WorkspaceDetailResponse{}, err
	}
	err = r.workspaceRepository.CreatePermission(w.ID, p.Members, workspace.Member)
	if err != nil {
		logger.Errorf("permission has not been created for workspace %d", w.ID)
		err2 := r.workspaceRepository.DeleteWorkspace(w.ID)
		if err2 != nil {
			return WorkspaceDetailResponse{}, err2
		}
		return WorkspaceDetailResponse{}, err
	}
	response := r.convertResponse(w)
	return response, nil

}

func (r *repository) convertResponse(w workspace.Workspace) WorkspaceDetailResponse {
	var projectCount int
	projects, err := r.projectRepo.GetByWorkspaceID(w.ID)
	if err == nil {
		projectCount = len(projects)
	} else {
		logger.Errorf("[WORKSPACE-API] - cannot get all projects by workspace. err %v", err)
	}
	var admin uint64
	perms, permissionCount, err := r.workspaceRepository.GetPermission(w.ID, workspace.Any, 0, 0)
	if err != nil {
		logger.Errorf("[WORKSPACE-API] - cannot get all permissions, and admin in workspace %d", w.ID)
		permissionCount = 0
	} else {
		for _, perm := range perms {
			if perm.Role == workspace.Admin {
				admin = perm.UserID
				break
			}
		}
	}
	return WorkspaceDetailResponse{
		WorkspaceBaseResponse: ToWorkspaceInfoResponse(w),
		ProjectCount:          projectCount,
		MemberCount:           permissionCount,
		Admin:                 admin,
	}
}

func (r *repository) UpdateWorkspace(id uint64, request UpdateWorkspaceRequest) (WorkspaceDetailResponse, error) {
	var changes = make(map[string]interface{})
	b, _ := json.Marshal(&request)
	_ = json.Unmarshal(b, &changes)
	logger.Info(changes)
	w, err := r.workspaceRepository.UpdateWorkspace(id, changes)
	if err != nil {
		return WorkspaceDetailResponse{}, err
	}
	return r.convertResponse(w), nil
}

func (r *repository) DeleteWorkspace(id uint64) error {
	return r.workspaceRepository.DeleteWorkspace(id)
}
