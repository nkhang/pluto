package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nkhang/pluto/internal/task/taskapi"

	"github.com/nkhang/pluto/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/tool/toolapi"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/pgin"
)

type params struct {
	fx.In

	GormDB           *gorm.DB
	Router           *gin.Engine
	ToolRepository   toolapi.Repository
	WorkspaceService pgin.StandaloneRouter `name:"WorkspaceService"`
	ProjectService   pgin.StandaloneRouter `name:"ProjectService"`
	ToolService      pgin.StandaloneRouter `name:"ToolService"`
	TaskService      pgin.StandaloneRouter `name:"TaskService"`
	ImageService     pgin.StandaloneRouter `name:"ImageService"`
	TaskServiceIns   *taskapi.Service      `name:"TaskService"`
}

func initializer(l fx.Lifecycle, p params) {
	migrate(p.GormDB)
	router := p.Router.Group("/pluto/api/v1")
	p.ImageService.RegisterStandalone(router.Group("/images"))
	p.TaskServiceIns.RegisterInternal(router.Group(""))
	if viper.GetBool("service.authen") {
		router.Use(pgin.ApplyVerifyToken())
	}
	p.ToolService.RegisterStandalone(router.Group("/tools"))
	p.ProjectService.RegisterStandalone(router.Group("/projects"))
	p.WorkspaceService.RegisterStandalone(router.Group("/workspaces"))
	p.TaskService.RegisterStandalone(router.Group("/tasks"))
	l.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := viper.GetInt("service.port")
				go func() {
					addr := fmt.Sprintf(":%d", port)
					err := p.Router.Run(addr)
					if err != nil {
						logger.Panic(err)
					}
				}()
				logger.Infof("Server is running at port %d", port)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Stop Server")
				return nil
			},
		},
	)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&tool.Tool{})
	db.FirstOrCreate(&tool.Tool{Name: "RECTANGLE"}, "name = ?", "RECTANGLE")
	db.FirstOrCreate(&tool.Tool{Name: "POINT"}, "name = ?", "POINT")
	db.FirstOrCreate(&tool.Tool{Name: "POLYLINE"}, "name = ?", "POLYLINE")
	db.FirstOrCreate(&tool.Tool{Name: "POLYGON"}, "name = ?", "POLYGON")
	db.AutoMigrate(&dataset.Dataset{})
	db.AutoMigrate(&label.Label{})
	db.AutoMigrate(&project.Project{})
	db.AutoMigrate(&project.Permission{})
	db.AutoMigrate(&workspace.Workspace{})
	db.AutoMigrate(&workspace.Permission{})
	db.AutoMigrate(&image.Image{})
	db.AutoMigrate(&task.Task{})
	db.AutoMigrate(&task.Detail{})
	db.AutoMigrate(&task.Detail{TaskID: 1})
	db.AutoMigrate(&task.Detail{TaskID: 2})
	db.AutoMigrate(&task.Detail{TaskID: 3})
	db.AutoMigrate(&task.Detail{TaskID: 4})
	db.AutoMigrate(&task.Detail{TaskID: 5})
	db.AutoMigrate(&task.Detail{TaskID: 6})
	db.AutoMigrate(&task.Detail{TaskID: 7})
	db.AutoMigrate(&task.Detail{TaskID: 8})
	db.AutoMigrate(&task.Detail{TaskID: 9})
}
