package repository

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iot-rd-backend/internal/model"
)

type ExchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) FindRate(baseCurrency, quoteCurrency string) (*model.ExchangeRate, error) {
	base := strings.ToUpper(strings.TrimSpace(baseCurrency))
	quote := strings.ToUpper(strings.TrimSpace(quoteCurrency))

	if base == quote {
		return &model.ExchangeRate{
			ID:            "EXR-" + base + "-" + quote,
			BaseCurrency:  base,
			QuoteCurrency: quote,
			Rate:          1.0,
			Provider:      "System Parity",
			Source:        "Direct 1:1",
			FetchedAt:     time.Now(),
			EffectiveDate: time.Now().Format("2006-01-02"),
		}, nil
	}

	var rate model.ExchangeRate
	err := r.db.Where("base_currency = ? AND quote_currency = ?", base, quote).First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *ExchangeRateRepository) UpsertRate(rate *model.ExchangeRate) error {
	rate.BaseCurrency = strings.ToUpper(strings.TrimSpace(rate.BaseCurrency))
	rate.QuoteCurrency = strings.ToUpper(strings.TrimSpace(rate.QuoteCurrency))
	if rate.ID == "" {
		rate.ID = "EXR-" + rate.BaseCurrency + "-" + rate.QuoteCurrency
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"rate", "provider", "source", "fetched_at", "effective_date", "updated_at",
		}),
	}).Create(rate).Error
}

func (r *ExchangeRateRepository) FindAll() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	err := r.db.Find(&rates).Error
	return rates, err
}
