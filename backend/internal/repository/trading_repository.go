package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TradingRepository struct {
	db *gorm.DB
}

func NewTradingRepository(db *gorm.DB) *TradingRepository {
	return &TradingRepository{db: db}
}

func (r *TradingRepository) FindByProductID(productID string) (*model.TradingProductDetail, error) {
	var detail model.TradingProductDetail
	err := r.db.Where("product_id = ?", productID).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *TradingRepository) Upsert(detail *model.TradingProductDetail) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"manufacturer_id", "manufacturer_name", "supplier_id", "supplier_name",
			"manufacturer_part_number", "supplier_part_number", "purchase_price",
			"selling_price", "currency", "lead_time_days", "warranty_information",
			"notes", "updated_at",
		}),
	}).Save(detail).Error
}

func (r *TradingRepository) Delete(productID string) error {
	return r.db.Where("product_id = ?", productID).Delete(&model.TradingProductDetail{}).Error
}
