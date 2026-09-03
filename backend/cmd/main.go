package main

import (
	"fmt"
	"log"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/database"
	"iot-rd-backend/internal/handler"
	"iot-rd-backend/internal/provider"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/router"
	"iot-rd-backend/internal/service"
)

func main() {
	log.Println("[INFO] Starting IoT Product R&D Control Center API Backend...")

	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Database and Auto-Migrate / Seed
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Database initialization failed: %v", err)
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	productRepo := repository.NewProductRepository(db)
	tradingRepo := repository.NewTradingRepository(db)
	projectItemRepo := repository.NewProjectItemRepository(db)
	versionRepo := repository.NewVersionRepository(db)
	revRepo := repository.NewRevisionRepository(db)
	compRepo := repository.NewComponentRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	mfgRepo := repository.NewManufacturerRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	bomRepo := repository.NewBOMRepository(db)
	protoRepo := repository.NewPrototypeRepository(db)
	rateRepo := repository.NewExchangeRateRepository(db)
	apiTokenRepo := repository.NewAPITokenRepository(db)
	fxProvider := provider.NewOpenExchangeRateProvider()

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, cfg)
	apiTokenService := service.NewAPITokenService(apiTokenRepo, activityRepo)
	rateService := service.NewExchangeRateService(rateRepo, fxProvider, activityRepo)
	productService := service.NewProductService(productRepo, tradingRepo, projectItemRepo, versionRepo, activityRepo)
	projectItemService := service.NewProjectItemService(projectItemRepo, productRepo, compRepo, activityRepo)
	versionService := service.NewVersionService(versionRepo, revRepo, productRepo, activityRepo)
	revService := service.NewRevisionService(revRepo, versionRepo, activityRepo)
	bomService := service.NewBOMService(bomRepo, revRepo, versionRepo, productRepo, compRepo, tradingRepo, projectItemRepo, activityRepo)
	productService.SetBOMService(bomService)
	productService.SetExchangeRateService(rateService)
	compService := service.NewComponentService(compRepo, activityRepo)
	catService := service.NewCategoryService(catRepo, activityRepo)
	mfgService := service.NewManufacturerService(mfgRepo, activityRepo)
	supplierService := service.NewSupplierService(supplierRepo, activityRepo)
	unitService := service.NewUnitService(unitRepo)
	userService := service.NewUserService(userRepo, activityRepo)
	roleService := service.NewRoleService(roleRepo, activityRepo)
	activityService := service.NewActivityService(activityRepo)
	uploadService := service.NewUploadService(cfg)
	protoService := service.NewPrototypeService(protoRepo, revRepo, versionRepo, productRepo, bomRepo, activityRepo)
	testingRepo := repository.NewTestingRepository(db)
	testingService := service.NewTestingService(testingRepo, protoRepo, activityRepo)
	changeRepo := repository.NewChangeRepository(db)
	changeService := service.NewChangeService(changeRepo, activityRepo)
	impRepo := repository.NewImportRepository(db)
	aiClassifier := provider.NewHeuristicRuleAIClassifier()
	importService := service.NewImportService(impRepo, productRepo, tradingRepo, compRepo, catRepo, mfgRepo, supplierRepo, activityRepo, aiClassifier, cfg)
	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo, activityRepo)
	importService.SetCustomerService(customerService)

	// 5. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	apiTokenHandler := handler.NewAPITokenHandler(apiTokenService)
	productHandler := handler.NewProductHandler(productService, projectItemService)
	versionHandler := handler.NewVersionHandler(versionService)
	revHandler := handler.NewRevisionHandler(revService)
	bomHandler := handler.NewBOMHandler(bomService)
	rateHandler := handler.NewExchangeRateHandler(rateService)
	compHandler := handler.NewComponentHandler(compService)
	catHandler := handler.NewCategoryHandler(catService)
	partnerHandler := handler.NewPartnerHandler(mfgService, supplierService)
	unitHandler := handler.NewUnitHandler(unitService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	activityHandler := handler.NewActivityHandler(activityService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	protoHandler := handler.NewPrototypeHandler(protoService)
	testingHandler := handler.NewTestingHandler(testingService)
	changeHandler := handler.NewChangeHandler(changeService)
	importHandler := handler.NewImportHandler(importService)
	customerHandler := handler.NewCustomerHandler(customerService)

	// 6. Setup Router
	r := router.SetupRouter(&router.RouterDependencies{
		Config:              cfg,
		AuthService:         authService,
		AuthHandler:         authHandler,
		APITokenService:     apiTokenService,
		APITokenHandler:     apiTokenHandler,
		ProductHandler:      productHandler,
		VersionHandler:      versionHandler,
		RevisionHandler:     revHandler,
		BOMHandler:          bomHandler,
		ExchangeRateHandler: rateHandler,
		CompHandler:         compHandler,
		CategoryHandler:     catHandler,
		PartnerHandler:      partnerHandler,
		CustomerHandler:     customerHandler,
		UnitHandler:         unitHandler,
		UserHandler:         userHandler,
		RoleHandler:         roleHandler,
		ActivityHandler:     activityHandler,
		UploadHandler:       uploadHandler,
		PrototypeHandler:    protoHandler,
		TestingHandler:      testingHandler,
		ChangeHandler:       changeHandler,
		ImportHandler:       importHandler,
	})

	// 7. Start HTTP Server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[INFO] Server listening on http://localhost%s (Environment: %s)", addr, cfg.AppEnv)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[FATAL] Server startup failed: %v", err)
	}
}
