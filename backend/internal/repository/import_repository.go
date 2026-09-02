package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ImportRepository struct {
	db *gorm.DB
}

func NewImportRepository(db *gorm.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

// DB returns the underlying gorm database connection for transactions.
func (r *ImportRepository) DB() *gorm.DB {
	return r.db
}

// CreateBatch creates a new import batch record.
func (r *ImportRepository) CreateBatch(batch *model.ImportBatch) error {
	return r.db.Create(batch).Error
}

// GetBatchByID retrieves an import batch by ID with its files.
func (r *ImportRepository) GetBatchByID(id string) (*model.ImportBatch, error) {
	var batch model.ImportBatch
	err := r.db.Preload("Files").Where("id = ? OR batch_code = ?", id, id).First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatches returns batches filtered by year and status.
func (r *ImportRepository) GetBatches(year string, status string) ([]model.ImportBatch, error) {
	var batches []model.ImportBatch
	query := r.db.Model(&model.ImportBatch{}).Preload("Files").Order("created_at DESC")
	if year != "" && year != "ALL" {
		query = query.Where("source_year = ?", year)
	}
	if status != "" && status != "ALL" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&batches).Error
	return batches, err
}

// UpdateBatch updates batch metadata and counters.
func (r *ImportRepository) UpdateBatch(batch *model.ImportBatch) error {
	return r.db.Save(batch).Error
}

// DeleteBatch deletes an import batch and cascades to files and staged rows.
func (r *ImportRepository) DeleteBatch(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Where("batch_id = ?", id).Delete(&model.ImportStagedRow{}).Error
		_ = tx.Where("batch_id = ?", id).Delete(&model.ImportFile{}).Error
		return tx.Where("id = ? OR batch_code = ?", id, id).Delete(&model.ImportBatch{}).Error
	})
}

// RecalculateBatchStats recalculates row counts and statuses for a batch.
func (r *ImportRepository) RecalculateBatchStats(batchID string) error {
	var total, valid, warning, errorRows, approved, rejected, imported int64
	
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ?", batchID).Count(&total)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND validation_status = ?", batchID, model.RowValidationValid).Count(&valid)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND validation_status = ?", batchID, model.RowValidationWarning).Count(&warning)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND validation_status = ?", batchID, model.RowValidationError).Count(&errorRows)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND review_decision = ?", batchID, model.RowDecisionApproved).Count(&approved)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND review_decision = ?", batchID, model.RowDecisionRejected).Count(&rejected)
	r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ? AND final_status = ?", batchID, "IMPORTED").Count(&imported)

	var filesCount int64
	r.db.Model(&model.ImportFile{}).Where("batch_id = ?", batchID).Count(&filesCount)

	updates := map[string]interface{}{
		"rows_count":    int(total),
		"valid_rows":    int(valid),
		"warning_rows":  int(warning),
		"error_rows":    int(errorRows),
		"approved_rows": int(approved),
		"rejected_rows": int(rejected),
		"imported_rows": int(imported),
		"files_count":   int(filesCount),
	}

	return r.db.Model(&model.ImportBatch{}).Where("id = ?", batchID).Updates(updates).Error
}

// CreateFile creates an import file record.
func (r *ImportRepository) CreateFile(file *model.ImportFile) error {
	return r.db.Create(file).Error
}

// GetFilesByBatchID retrieves all files belonging to a batch.
func (r *ImportRepository) GetFilesByBatchID(batchID string) ([]model.ImportFile, error) {
	var files []model.ImportFile
	err := r.db.Where("batch_id = ?", batchID).Order("created_at ASC").Find(&files).Error
	return files, err
}

// GetFileByID retrieves an import file by ID.
func (r *ImportRepository) GetFileByID(id string) (*model.ImportFile, error) {
	var file model.ImportFile
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// CreateStagedRows inserts a batch of staged rows into database.
func (r *ImportRepository) CreateStagedRows(rows []model.ImportStagedRow) error {
	if len(rows) == 0 {
		return nil
	}
	// Insert in chunks of 250 to avoid MySQL parameter count limits
	return r.db.CreateInBatches(rows, 250).Error
}

// GetStagedRows retrieves staged rows with filtering, search and server-side pagination.
func (r *ImportRepository) GetStagedRows(batchID string, status string, decision string, classification string, search string, page int, limit int) ([]model.ImportStagedRow, int64, error) {
	var rows []model.ImportStagedRow
	var total int64

	query := r.db.Model(&model.ImportStagedRow{}).Where("batch_id = ?", batchID)

	if status != "" && status != "ALL" {
		query = query.Where("validation_status = ?", status)
	}
	if decision != "" && decision != "ALL" {
		query = query.Where("review_decision = ?", decision)
	}
	if classification != "" && classification != "ALL" {
		query = query.Where("item_classification = ?", classification)
	}
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("normalized_name LIKE ? OR normalized_code LIKE ? OR normalized_category LIKE ? OR raw_data_json LIKE ?", s, s, s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	offset := (page - 1) * limit

	err := query.Order("row_number ASC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// GetAllStagedRowsInBatch returns all staged rows for batch-level processing.
func (r *ImportRepository) GetAllStagedRowsInBatch(batchID string) ([]model.ImportStagedRow, error) {
	var rows []model.ImportStagedRow
	err := r.db.Where("batch_id = ?", batchID).Order("row_number ASC").Find(&rows).Error
	return rows, err
}

// GetStagedRowByID retrieves a single staged row.
func (r *ImportRepository) GetStagedRowByID(id string) (*model.ImportStagedRow, error) {
	var row model.ImportStagedRow
	err := r.db.Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// UpdateStagedRow updates a staged row.
func (r *ImportRepository) UpdateStagedRow(row *model.ImportStagedRow) error {
	return r.db.Save(row).Error
}

// BulkUpdateRows updates multiple staged rows in batch.
func (r *ImportRepository) BulkUpdateRows(rows []model.ImportStagedRow) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i := range rows {
			if err := tx.Save(&rows[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetApprovedRowsByBatch returns all approved rows ready for atomic master dispatch.
func (r *ImportRepository) GetApprovedRowsByBatch(batchID string) ([]model.ImportStagedRow, error) {
	var rows []model.ImportStagedRow
	err := r.db.Where("batch_id = ? AND review_decision = ?", batchID, model.RowDecisionApproved).Order("row_number ASC").Find(&rows).Error
	return rows, err
}

// DeleteStagedRowsByBatch deletes all staged rows for a batch (e.g. before re-parsing).
func (r *ImportRepository) DeleteStagedRowsByBatch(batchID string) error {
	return r.db.Where("batch_id = ?", batchID).Delete(&model.ImportStagedRow{}).Error
}
