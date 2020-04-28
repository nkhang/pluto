package label

import (
	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/pkg/gorm"
)

const (
	fieldProjectID = "project_id"
)

type Label struct {
	gorm.Model
	Name      string
	ProjectID uint64
	ToolID    uint64
	Tool      tool.Tool `gorm:"foreignkey:ToolID;association_foreignkey:ID"`
}
