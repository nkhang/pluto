package dataset

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository DbRepository

type repository struct {
	dbRepo    DbRepository
	cacheRepo cache.Cache
}

func NewRepository(d DbRepository, c cache.Cache) *repository {
	return &repository{
		dbRepo:    d,
		cacheRepo: c,
	}
}

func (r *repository) Get(dID uint64) (d Dataset, err error) {
	k := rediskey.DatasetByID(dID)
	err = r.cacheRepo.Get(k, &d)
	if err == nil {
		logger.Infof("cache hit getting dataset [%d]", dID)
		return d, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss getting dataset [%d]", dID)
	} else {
		logger.Errorf("error getting dataset [%d] from cache", dID)
	}
	d, err = r.dbRepo.Get(dID)
	if err != nil {
		logger.Error("error getting dataset [%d] from database", dID)
		return Dataset{}, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &d)
		if err != nil {
			logger.Error("error set dataset [%d] to cache", dID)
		}
	}()
	return d, nil
}

func (r *repository) GetByProject(pID uint64) ([]Dataset, error) {
	var ds = make([]Dataset, 0)
	k := rediskey.DatasetByProject(pID)
	err := r.cacheRepo.Get(k, &ds)
	if err == nil {
		logger.Infof("cache hit getting datasets of project [%d]", pID)
		return ds, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss getting datasets of project [%d]", pID)
	} else {
		logger.Errorf("error getting datasets of project [%d] from cache", pID)
	}
	ds, err = r.dbRepo.GetByProject(pID)
	if err != nil {
		logger.Error("error getting datasets of projects [%d] from database", pID)
		return nil, err
	}
	go func() {
		err := r.cacheRepo.Set(k, &ds)
		if err != nil {
			logger.Error("error set datasets of projects [%d] to cache", pID)
		}
	}()
	return ds, nil
}