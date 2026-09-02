package database

import (
	"database/sql"
	"fmt"
	"log"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.DBDriver == "mysql" {
		// Attempt to create database if not exists using standard database/sql
		rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
		
		rawDB, rawErr := sql.Open("mysql", rootDSN)
		if rawErr == nil {
			_, _ = rawDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName))
			_ = rawDB.Close()
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})

		if err != nil {
			log.Printf("[WARNING] Could not connect to MySQL at %s:%s (Error: %v). Falling back to SQLite database for continuous operation...", cfg.DBHost, cfg.DBPort, err)
			db, err = gorm.Open(sqlite.Open("iot_rd_master.db"), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Warn),
			})
		}
	} else {
		db, err = gorm.Open(sqlite.Open("iot_rd_master.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.APIToken{},
		&model.Product{},
		&model.TradingProductDetail{},
		&model.ProjectItem{},
		&model.ComponentCategory{},
		&model.Manufacturer{},
		&model.Supplier{},
		&model.Unit{},
		&model.Component{},
		&model.ComponentSpecification{},
		&model.ComponentDocument{},
		&model.ProductVersion{},
		&model.HardwareRevision{},
		&model.EngineeringBOM{},
		&model.BOMItem{},
		&model.ExchangeRate{},
		&model.ActivityLog{},
		&model.PrototypeBuild{},
		&model.PrototypeComponentPreparation{},
		&model.PrototypeAssemblyStage{},
		&model.PrototypeEngineeringNote{},
		&model.TestPlan{},
		&model.TestPlanVersion{},
		&model.TestCase{},
		&model.PrototypeTestSession{},
		&model.TestExecution{},
		&model.TestMeasurement{},
		&model.TestEvidence{},
		&model.EngineeringFinding{},
		&model.PrototypeValidationDecision{},
		&model.EngineeringChangeRequest{},
		&model.ECRChangeItem{},
		&model.ECRImpactAnalysis{},
		&model.ECRReview{},
		&model.ECRApproval{},
		&model.EngineeringChangeOrder{},
		&model.ImportBatch{},
		&model.ImportFile{},
		&model.ImportStagedRow{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}

	// Safe backfill for existing products to guarantee RND type
	_ = db.Exec("UPDATE products SET product_type = 'RND' WHERE product_type IS NULL OR product_type = '';")

	// Create compound indexes if not present
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_prod_version_num ON product_versions(product_id, version_number, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_ver_rev_code ON hardware_revisions(product_version_id, code, deleted_at);")
	_ = db.Exec("ALTER TABLE prototype_builds DROP INDEX idx_prototype_builds_code;")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_proto_code_del ON prototype_builds(code, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_rev_build_num ON prototype_builds(hardware_revision_id, build_number, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tp_code_del ON test_plans(code, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tp_ver_num ON test_plan_versions(test_plan_id, version_number, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_sess_code_del ON prototype_test_sessions(session_code, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_fnd_code_del ON engineering_findings(code, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_ecr_code_del ON engineering_change_requests(code, deleted_at);")
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_eco_code_del ON engineering_change_orders(code, deleted_at);")

	DB = db
	log.Println("[INFO] Database connected and schema auto-migrated successfully.")

	// Run Seed Data
	SeedAll(db)
	SeedPhase2(db)
	SeedPhase2C(db)
	SeedPhase3(db)
	SeedMultiCurrency(db)
	SeedPhase4(db)
	SeedPhase5(db)
	SeedPhase6(db)
	SeedDataImport(db)

	return db, nil
}
