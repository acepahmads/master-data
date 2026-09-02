package repository

import (
	"fmt"
	"time"

	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) FindAll(p PaginationParams, module, entityType, entityID string) ([]model.ActivityLog, int64, error) {
	var activities []model.ActivityLog
	var total int64

	query := r.db.Model(&model.ActivityLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	if entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order("created_at desc").Find(&activities).Error

	now := time.Now()
	for i := range activities {
		diff := now.Sub(activities[i].CreatedAt)
		if diff < time.Minute {
			activities[i].Timestamp = "Just now"
		} else if diff < time.Hour {
			activities[i].Timestamp = fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		} else if diff < 24*time.Hour {
			activities[i].Timestamp = fmt.Sprintf("%d hours ago", int(diff.Hours()))
		} else {
			activities[i].Timestamp = fmt.Sprintf("%d days ago", int(diff.Hours()/24))
		}
		activities[i].Date = activities[i].CreatedAt.Format("2006-01-02 15:04")
	}

	return activities, total, err
}

func (r *ActivityRepository) Create(log *model.ActivityLog) error {
	return r.db.Create(log).Error
}
