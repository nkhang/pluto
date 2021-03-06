package imagefx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/project"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/image/imageapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/objectstorage"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideImageRepository(db *gorm.DB, cache cache.Cache) image.Repository {
	dbRepo := image.NewDBRepository(db)
	return image.NewRepository(dbRepo, cache)
}

func provideService(r image.Repository, s objectstorage.ObjectStorage,
	d dataset.Repository, p project.Repository) (pgin.Router, pgin.StandaloneRouter) {
	repository := imageapi.NewRepository(r, s, d, p)
	router := imageapi.NewService(repository)
	return router, router
}
