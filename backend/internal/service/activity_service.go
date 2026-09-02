package service

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) GetAll(p repository.PaginationParams, module, entityType, entityID string) ([]model.ActivityLog, repository.PaginationMeta, error) {
	activities, total, err := s.repo.FindAll(p, module, entityType, entityID)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return activities, meta, nil
}
