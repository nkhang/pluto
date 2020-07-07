package workspace

import (
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/pkg/gorm"
)

type Role int32

const (
	Any   Role = 0
	Admin Role = iota + 1
	Member
)

type Workspace struct {
	gorm.Model
	Title       string
	Description string
	Projects    []project.Project
	Perm        []Permission
}

type Permission struct {
	gorm.Model
	WorkspaceID uint64
	Workspace   Workspace `gorm:"association_save_reference:false"`
	Role        Role
	UserID      uint64
}

func (Permission) TableName() string {
	return "workspace_permissions"
}
