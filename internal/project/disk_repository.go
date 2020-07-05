package project

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/pkg/logger"

	"github.com/nkhang/pluto/pkg/errors"
)

type DBRepository interface {
	Get(pID uint64) (Project, error)
	GetByWorkspaceID(wID uint64, offset, limit int) ([]Project, int, error)
	CountAllProject(workspaceID uint64) (int, error)
	GetProjectPermissions(pID uint64) ([]Permission, error)
	GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error)
	GetPermission(userID, projectID uint64) (Permission, error)
	CreateProject(wID uint64, title, desc, uid string) (Project, error)
	CreatePermission(projectID, userID uint64, role Role) (Permission, error)
}

type dbRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) Get(pID uint64) (Project, error) {
	var p Project
	result := r.db.First(&p, pID)
	if result.RecordNotFound() {
		return Project{}, errors.ProjectNotFound.NewWithMessage("project not found")
	}
	if err := result.Error; err != nil {
		return Project{}, errors.ProjectQueryError.Wrap(err, "query project error")
	}
	return p, nil
}

func (r *dbRepository) GetByWorkspaceID(wID uint64, offset, limit int) ([]Project, int, error) {
	var projects = make([]Project, 0)
	var total int
	err := r.db.Model(&Project{}).
		Where(fieldWorkspaceID+" = ?", wID).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&projects).Error
	if err != nil {
		return nil, 0, errors.ProjectQueryError.NewWithMessage("error getting project of workspace")
	}
	return projects, total, nil
}

func (r *dbRepository) CountAllProject(workspaceID uint64) (int, error) {
	var count int
	err := r.db.Model(Project{}).Count(&count).Error
	if err != nil {
		return 0, errors.ProjectCreatingError.Wrap(err, "cannot count all projects")
	}
	return count, nil
}

func (r *dbRepository) CreateProject(wID uint64, title, desc, uid string) (Project, error) {
	rand.Seed(time.Now().Unix())
	index := rand.Int()
	var p = Project{
		WorkspaceID: wID,
		Title:       title,
		Description: desc,
		Dir:         uid,
		Thumbnail:   defaultImage,
		Color:       Color[index%len(Color)],
	}
	err := r.db.Create(&p).Error
	if err != nil {
		return Project{}, errors.ProjectCreatingError.Wrap(err, "cannot create project")
	}
	return p, nil
}

func (r *dbRepository) GetProjectPermissions(pID uint64) ([]Permission, error) {
	var perms = make([]Permission, 0)
	err := r.db.Where("project_id = ?", pID).Find(&perms).Error
	if err != nil {
		return nil, errors.ProjectPermissionQueryError.Wrap(err, "cannot query project permissions for project")
	}
	return perms, nil
}

func (r *dbRepository) GetUserPermissions(userID uint64, role Role, offset, limit int) ([]Permission, int, error) {
	var perms = make([]Permission, 0)
	var total int
	db := r.db.Model(Permission{}).
		Where("user_id = ?", userID)
	if role != 0 {
		db = db.Where("role = ?", role)
	}
	db = db.Count(&total)
	if offset != 0 || limit != 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err := db.Preload("Project").Find(&perms).Error
	if err != nil {
		logger.Error(err)
		return nil, 0, errors.ProjectPermissionQueryError.Wrap(err, "cannot query project permissions for project")
	}
	return perms, total, nil
}

func (r *dbRepository) CreatePermission(projectID, userID uint64, role Role) (Permission, error) {
	perm := Permission{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	err := r.db.Create(&perm).Error
	if err != nil {
		return perm, errors.ProjectPermissionCreatingError.Wrap(err, "cannot create project permission")
	}
	return perm, nil
}

func (r *dbRepository) GetPermission(userID, projectID uint64) (Permission, error) {
	perm := Permission{
		UserID:    userID,
		ProjectID: projectID,
	}
	db := r.db.Where(&perm).First(&perm)
	if db.RecordNotFound() {
		return Permission{}, errors.ProjectPermissionNotFound.NewWithMessage("project permission not found")
	}
	if err := db.Error; err != nil {
		return Permission{}, errors.ProjectPermissionQueryError.NewWithMessage("error query project permission")
	}
	return perm, nil
}
