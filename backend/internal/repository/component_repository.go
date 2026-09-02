package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ComponentRepository struct {
	db *gorm.DB
}

func NewComponentRepository(db *gorm.DB) *ComponentRepository {
	return &ComponentRepository{db: db}
}

type ComponentFilters struct {
	CategoryID     string
	ManufacturerID string
	SupplierID     string
	Status         string
	MountingType   string
	RoHS           string
	UnitID         string
}

func (r *ComponentRepository) FindAll(p PaginationParams, filters ComponentFilters) ([]model.Component, int64, error) {
	var components []model.Component
	var total int64

	query := r.db.Model(&model.Component{}).
		Preload("Category").
		Preload("Manufacturer").
		Preload("Supplier").
		Preload("Unit").
		Preload("Specs").
		Preload("Documents")

	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("part_number LIKE ? OR internal_code LIKE ? OR name LIKE ? OR short_description LIKE ?", s, s, s, s)
	}
	if filters.CategoryID != "" && filters.CategoryID != "All" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}
	if filters.ManufacturerID != "" && filters.ManufacturerID != "All" {
		query = query.Where("manufacturer_id = ?", filters.ManufacturerID)
	}
	if filters.SupplierID != "" && filters.SupplierID != "All" {
		query = query.Where("supplier_id = ?", filters.SupplierID)
	}
	if filters.Status != "" && filters.Status != "All" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.MountingType != "" {
		query = query.Where("mounting_type = ?", filters.MountingType)
	}
	if filters.RoHS == "true" {
		query = query.Where("ro_hs_compliant = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&components).Error

	// Populate flat names for UI convenience
	for i := range components {
		r.populateComputedFields(&components[i])
	}

	return components, total, err
}

func (r *ComponentRepository) FindByID(id string) (*model.Component, error) {
	var component model.Component
	err := r.db.Where("id = ? OR part_number = ? OR internal_code = ?", id, id, id).
		Preload("Category").
		Preload("Manufacturer").
		Preload("Supplier").
		Preload("Unit").
		Preload("Specs", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order asc, created_at asc")
		}).
		Preload("Documents").
		First(&component).Error
	if err != nil {
		return nil, err
	}

	r.populateComputedFields(&component)
	return &component, nil
}

func (r *ComponentRepository) populateComputedFields(c *model.Component) {
	if c.Category != nil {
		c.CategoryName = c.Category.Name
	}
	if c.Manufacturer != nil {
		c.ManufacturerName = c.Manufacturer.Name
	}
	if c.Supplier != nil {
		c.SupplierName = c.Supplier.Name
	}
	if c.Unit != nil {
		c.UnitCode = c.Unit.Code
	}
}

func (r *ComponentRepository) Create(component *model.Component) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(component).Error; err != nil {
			return err
		}
		// Increment manufacturer/supplier/category component counts
		if component.ManufacturerID != "" {
			tx.Model(&model.Manufacturer{}).Where("id = ?", component.ManufacturerID).
				UpdateColumn("total_components", gorm.Expr("total_components + 1"))
		}
		if component.SupplierID != "" {
			tx.Model(&model.Supplier{}).Where("id = ?", component.SupplierID).
				UpdateColumn("total_components", gorm.Expr("total_components + 1"))
		}
		if component.CategoryID != "" {
			tx.Model(&model.ComponentCategory{}).Where("id = ?", component.CategoryID).
				UpdateColumn("total_components", gorm.Expr("total_components + 1"))
		}
		return nil
	})
}

func (r *ComponentRepository) Update(component *model.Component) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update specs
		_ = tx.Where("component_id = ?", component.ID).Delete(&model.ComponentSpecification{}).Error
		if len(component.Specs) > 0 {
			for i := range component.Specs {
				component.Specs[i].ComponentID = component.ID
			}
			if err := tx.Create(&component.Specs).Error; err != nil {
				return err
			}
		}

		return tx.Save(component).Error
	})
}

func (r *ComponentRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var comp model.Component
		if err := tx.Where("id = ? OR part_number = ?", id, id).First(&comp).Error; err != nil {
			return err
		}

		if comp.ManufacturerID != "" {
			tx.Model(&model.Manufacturer{}).Where("id = ?", comp.ManufacturerID).
				UpdateColumn("total_components", gorm.Expr("GREATEST(total_components - 1, 0)"))
		}
		if comp.SupplierID != "" {
			tx.Model(&model.Supplier{}).Where("id = ?", comp.SupplierID).
				UpdateColumn("total_components", gorm.Expr("GREATEST(total_components - 1, 0)"))
		}
		if comp.CategoryID != "" {
			tx.Model(&model.ComponentCategory{}).Where("id = ?", comp.CategoryID).
				UpdateColumn("total_components", gorm.Expr("GREATEST(total_components - 1, 0)"))
		}

		_ = tx.Where("component_id = ?", comp.ID).Delete(&model.ComponentSpecification{}).Error
		_ = tx.Where("component_id = ?", comp.ID).Delete(&model.ComponentDocument{}).Error
		return tx.Delete(&comp).Error
	})
}

func (r *ComponentRepository) AddDocument(doc *model.ComponentDocument) error {
	return r.db.Create(doc).Error
}

func (r *ComponentRepository) DeleteDocument(docID string) error {
	return r.db.Where("id = ?", docID).Delete(&model.ComponentDocument{}).Error
}

func (r *ComponentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Component{}).Count(&count).Error
	return count, err
}
