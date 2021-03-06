package projectfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project/projectapi/permissionapi"
	"github.com/nkhang/pluto/internal/project/projectapi/statsapi"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/workspace/workspaceapi"
	"github.com/nkhang/pluto/pkg/annotation"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/project/projectapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideProjectDBRepository(db *gorm.DB) project.DBRepository {
	return project.NewDiskRepository(db)
}

func provideRepository(r project.DBRepository, c cache.Cache, t task.Repository, d dataset.Repository, i image.Repository) project.Repository {
	return project.NewRepository(r, c, t, d, i)
}

func provideAPIRepository(r project.Repository, dr dataset.Repository, wr workspaceapi.Repository, ann annotation.Service) projectapi.Repository {
	return projectapi.NewRepository(r, dr, wr, ann)
}

func provideStatsAPIRepo(d dataset.Repository, t task.Repository, i image.Repository, s annotation.Service, l label.Repository) statsapi.Repository {
	return statsapi.NewRepository(d, t, i, s, l)
}

type params struct {
	fx.In

	Repository        projectapi.Repository
	StatAPIRepo       statsapi.Repository
	ProjectRepo       project.Repository
	ProjectAPI        projectapi.Repository
	AnnotationService annotation.Service
	DatasetRouter     pgin.Router `name:"DatasetService"`
	TaskRouter        pgin.Router `name:"TaskService"`
	LabelRouter       pgin.Router `name:"LabelService"`
}

func provideService(p params) (pgin.Router, pgin.StandaloneRouter) {
	permRepo := permissionapi.NewProjectPermissionAPIRepository(p.ProjectRepo, p.ProjectAPI, p.AnnotationService)
	permService := permissionapi.NewService(permRepo, p.ProjectRepo)
	statService := statsapi.NewService(p.StatAPIRepo)
	service := projectapi.NewService(p.Repository, p.ProjectRepo,
		permService, p.TaskRouter, p.DatasetRouter, p.LabelRouter, statService)
	return service, service
}
