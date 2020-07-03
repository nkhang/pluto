package taskapi

import (
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/task"
)

type Repository interface {
	CreateTask(request CreateTaskRequest) error
}

type repository struct {
	repository task.Repository
	imgRepo    image.Repository
}

func NewRepository(r task.Repository, ir image.Repository) *repository {
	return &repository{
		repository: r,
		imgRepo:    ir,
	}
}

func (r *repository) CreateTask(request CreateTaskRequest) error {
	imgs, err := r.imgRepo.GetAllImageByDataset(request.DatasetID)
	if err != nil {
		return err
	}
	truncated := truncate(imgs, request.Quantity)
	ids := make([]uint64, len(truncated))
	for i := range truncated {
		ids[i] = truncated[i].ID
	}
	return r.repository.CreateTask(request.Assigner, request.Labeler, request.Reviewer, request.DatasetID, ids)
}

func truncate(imgs []image.Image, s int) []image.Image {
	l := len(imgs)
	if l <= s {
		return imgs
	}
	return imgs[:s]
}
