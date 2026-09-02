package router

import (
	"net/http"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/handler"
	"iot-rd-backend/internal/middleware"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	Config          *config.Config
	AuthService     *service.AuthService
	AuthHandler     *handler.AuthHandler
	ProductHandler  *handler.ProductHandler
	VersionHandler  *handler.VersionHandler
	RevisionHandler *handler.RevisionHandler
	CompHandler     *handler.ComponentHandler
	CategoryHandler *handler.CategoryHandler
	PartnerHandler  *handler.PartnerHandler
	UnitHandler     *handler.UnitHandler
	UserHandler     *handler.UserHandler
	RoleHandler     *handler.RoleHandler
	ActivityHandler *handler.ActivityHandler
	UploadHandler       *handler.UploadHandler
	BOMHandler          *handler.BOMHandler
	ExchangeRateHandler *handler.ExchangeRateHandler
	PrototypeHandler    *handler.PrototypeHandler
	TestingHandler      *handler.TestingHandler
	ChangeHandler       *handler.ChangeHandler
	ImportHandler       *handler.ImportHandler
	APITokenService     *service.APITokenService
	APITokenHandler     *handler.APITokenHandler
}

func SetupRouter(deps *RouterDependencies) *gin.Engine {
	if deps.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())

	// Static file serving for uploads
	r.Static("/uploads", deps.Config.UploadPath)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "IoT R&D Control Center API", "version": "1.0.0"})
	})

	api := r.Group("/api/v1")
	{
		// Public Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", deps.AuthHandler.Login)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(deps.Config, deps.AuthService, deps.APITokenService))
		{
			// Auth Profile
			protected.GET("/auth/me", deps.AuthHandler.GetMe)

			// File Upload
			protected.POST("/upload", deps.UploadHandler.UploadFile)

			// Products (Phase 1 + Phase 2C Extensions)
			products := protected.Group("/products")
			{
				products.GET("", middleware.RequirePermission("products.view"), deps.ProductHandler.GetAll)
				products.GET("/:id", middleware.RequirePermission("products.view"), deps.ProductHandler.GetByID)
				products.GET("/:id/lifecycle", middleware.RequirePermission("products.view"), deps.VersionHandler.GetLifecycle)
				products.POST("", middleware.RequirePermission("products.create"), deps.ProductHandler.Create)
				products.PUT("/:id", middleware.RequirePermission("products.edit"), deps.ProductHandler.Update)
				products.DELETE("/:id", middleware.RequirePermission("products.delete"), deps.ProductHandler.Delete)

				// Trading Product Detail
				products.GET("/:id/trading-detail", middleware.RequirePermission("products.view"), deps.ProductHandler.GetTradingDetail)
				products.PUT("/:id/trading-detail", middleware.RequirePermission("products.edit"), deps.ProductHandler.UpdateTradingDetail)

				// Project Items Hierarchy
				products.GET("/:id/project-items", middleware.RequirePermission("products.view"), deps.ProductHandler.GetProjectItems)
				products.POST("/:id/project-items", middleware.RequirePermission("products.create"), deps.ProductHandler.CreateProjectItem)
				products.POST("/:id/project-items/reorder", middleware.RequirePermission("products.edit"), deps.ProductHandler.ReorderProjectItems)

				// Cost Summary Foundations (Phase 3)
				products.GET("/:id/cost-summary", middleware.RequirePermission("products.view"), deps.BOMHandler.GetProductCostSummary)
				products.GET("/:id/project-cost-summary", middleware.RequirePermission("products.view"), deps.BOMHandler.GetProjectCostSummary)
			}

			// Project Items Single Entity Operations
			protected.PUT("/project-items/:id", middleware.RequirePermission("products.edit"), deps.ProductHandler.UpdateProjectItem)
			protected.DELETE("/project-items/:id", middleware.RequirePermission("products.delete"), deps.ProductHandler.DeleteProjectItem)

			// Product Versions (Phase 2B)
			versions := protected.Group("/product-versions")
			{
				versions.GET("", middleware.RequirePermission("versions.view"), deps.VersionHandler.GetAll)
				versions.GET("/compare", middleware.RequirePermission("versions.view"), deps.VersionHandler.Compare)
				versions.GET("/:id", middleware.RequirePermission("versions.view"), deps.VersionHandler.GetByID)
				versions.POST("", middleware.RequirePermission("versions.create"), deps.VersionHandler.Create)
				versions.PUT("/:id", middleware.RequirePermission("versions.edit"), deps.VersionHandler.Update)
				versions.DELETE("/:id", middleware.RequirePermission("versions.delete"), deps.VersionHandler.Delete)
			}

			// Hardware Revisions (Phase 2B)
			revisions := protected.Group("/hardware-revisions")
			{
				revisions.GET("", middleware.RequirePermission("revisions.view"), deps.RevisionHandler.GetAll)
				revisions.GET("/:id", middleware.RequirePermission("revisions.view"), deps.RevisionHandler.GetByID)
				revisions.POST("", middleware.RequirePermission("revisions.create"), deps.RevisionHandler.Create)
				revisions.PUT("/:id", middleware.RequirePermission("revisions.edit"), deps.RevisionHandler.Update)
				revisions.DELETE("/:id", middleware.RequirePermission("revisions.delete"), deps.RevisionHandler.Delete)

				// Revision Engineering BOM (Phase 3)
				revisions.GET("/:id/bom", middleware.RequirePermission("bom.view"), deps.BOMHandler.GetByRevision)
				revisions.POST("/:id/bom", middleware.RequirePermission("bom.create"), deps.BOMHandler.Create)
			}

			// Engineering BOM Endpoints (Phase 3)
			boms := protected.Group("/boms")
			{
				boms.GET("/:id", middleware.RequirePermission("bom.view"), deps.BOMHandler.GetByID)
				boms.PUT("/:id", middleware.RequirePermission("bom.edit"), deps.BOMHandler.Update)
				boms.GET("/:id/cost-summary", middleware.RequirePermission("bom.view"), deps.BOMHandler.GetCostSummary)
				boms.POST("/:id/items", middleware.RequirePermission("bom.create"), deps.BOMHandler.AddItem)
				boms.POST("/:id/items/reorder", middleware.RequirePermission("bom.edit"), deps.BOMHandler.ReorderItems)
			}

			// BOM Single Item Operations
			protected.PUT("/bom-items/:id", middleware.RequirePermission("bom.edit"), deps.BOMHandler.UpdateItem)
			protected.DELETE("/bom-items/:id", middleware.RequirePermission("bom.delete"), deps.BOMHandler.DeleteItem)

			// Prototype & Engineering Build Management (Phase 4)
			prototypes := protected.Group("/prototypes")
			{
				prototypes.GET("/dashboard", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetDashboard)
				prototypes.GET("", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetAll)
				prototypes.POST("", middleware.RequirePermission("prototypes.create"), deps.PrototypeHandler.Create)
				prototypes.GET("/:id", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetByID)
				prototypes.PUT("/:id", middleware.RequirePermission("prototypes.edit"), deps.PrototypeHandler.Update)
				prototypes.DELETE("/:id", middleware.RequirePermission("prototypes.delete"), deps.PrototypeHandler.Delete)

				// BOM Snapshot
				prototypes.GET("/:id/bom-snapshot", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetBOMSnapshot)

				// Component Preparation
				prototypes.GET("/:id/components", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetComponents)
				prototypes.PUT("/:id/components/:componentId", middleware.RequirePermission("prototype.build.execute"), deps.PrototypeHandler.UpdateComponent)

				// Assembly Stages
				prototypes.GET("/:id/assembly", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetAssembly)
				prototypes.PUT("/:id/assembly/:stageId", middleware.RequirePermission("prototype.build.execute"), deps.PrototypeHandler.UpdateAssemblyStage)

				// Engineering Notes
				prototypes.GET("/:id/notes", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetNotes)
				prototypes.POST("/:id/notes", middleware.RequirePermission("prototype.notes.create"), deps.PrototypeHandler.CreateNote)
				prototypes.PUT("/:id/notes/:noteId", middleware.RequirePermission("prototype.notes.create"), deps.PrototypeHandler.UpdateNote)
				prototypes.DELETE("/:id/notes/:noteId", middleware.RequirePermission("prototype.notes.create"), deps.PrototypeHandler.DeleteNote)

				// Prototype Test Sessions (Phase 5)
				prototypes.GET("/:id/test-sessions", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetPrototypeTestSessions)
				prototypes.POST("/:id/test-sessions", middleware.RequirePermission("testing.execute"), deps.TestingHandler.CreatePrototypeTestSession)
			}

			// Revision Prototype Builds Shortcut Routes
			protected.GET("/revisions/:id/prototypes", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetByRevision)
			protected.GET("/hardware-revisions/:id/prototypes", middleware.RequirePermission("prototypes.view"), deps.PrototypeHandler.GetByRevision)

			// Testing & Validation Control (Phase 5)
			testingGroup := protected.Group("/testing")
			{
				testingGroup.GET("/dashboard", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetDashboard)
			}

			// Master Test Plans
			testPlans := protected.Group("/test-plans")
			{
				testPlans.GET("", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetTestPlans)
				testPlans.POST("", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.CreateTestPlan)
				testPlans.GET("/:id", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetTestPlan)
				testPlans.PUT("/:id", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.UpdateTestPlan)
				testPlans.DELETE("/:id", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.DeleteTestPlan)

				// Plan Versions
				testPlans.POST("/:id/versions", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.CreateTestPlanVersion)
			}

			// Test Plan Versions & Test Cases
			planVersions := protected.Group("/test-plan-versions")
			{
				planVersions.GET("/:id", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetTestPlanVersion)
				planVersions.PUT("/:id", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.UpdateTestPlanVersion)
				planVersions.POST("/:id/release", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.ReleaseTestPlanVersion)
				planVersions.POST("/:id/clone", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.CloneTestPlanVersion)

				// Test Cases on Version
				planVersions.POST("/:id/test-cases", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.CreateTestCase)
			}

			// Individual Test Cases
			testCases := protected.Group("/test-cases")
			{
				testCases.PUT("/:id", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.UpdateTestCase)
				testCases.DELETE("/:id", middleware.RequirePermission("testplans.manage"), deps.TestingHandler.DeleteTestCase)
			}

			// Prototype Test Sessions & Execution Runs
			testSessions := protected.Group("/test-sessions")
			{
				testSessions.GET("/:id", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetTestSession)
				testSessions.PUT("/:id", middleware.RequirePermission("testing.execute"), deps.TestingHandler.UpdateTestSession)
				testSessions.GET("/:id/executions", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetSessionExecutions)
				testSessions.POST("/:id/executions", middleware.RequirePermission("testing.execute"), deps.TestingHandler.RecordExecution)
				testSessions.POST("/:id/executions/:caseId/next-run", middleware.RequirePermission("testing.execute"), deps.TestingHandler.CreateNextRun)

				// Validation Decisions
				testSessions.GET("/:id/validation-decisions", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetValidationDecisions)
				testSessions.POST("/:id/validation-decisions", middleware.RequirePermission("testing.validate"), deps.TestingHandler.SubmitValidationDecision)
			}

			// Test Executions
			testExecutions := protected.Group("/test-executions")
			{
				testExecutions.POST("/:id/evidence", middleware.RequirePermission("testing.execute"), deps.TestingHandler.AddEvidence)
			}

			// Engineering Findings
			findings := protected.Group("/findings")
			{
				findings.GET("", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetFindings)
				findings.POST("", middleware.RequirePermission("findings.manage"), deps.TestingHandler.CreateFinding)
				findings.GET("/:id", middleware.RequirePermission("testing.view"), deps.TestingHandler.GetFinding)
				findings.PUT("/:id", middleware.RequirePermission("findings.manage"), deps.TestingHandler.UpdateFinding)
				findings.DELETE("/:id", middleware.RequirePermission("findings.manage"), deps.TestingHandler.DeleteFinding)
			}

			// Engineering Change Management (Phase 6)
			changesGroup := protected.Group("/changes")
			{
				changesGroup.GET("/dashboard", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetDashboard)
				changesGroup.GET("/traceability", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetTraceability)
			}

			// Engineering Change Requests (ECR)
			ecrGroup := protected.Group("/ecr")
			{
				ecrGroup.GET("", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetECRs)
				ecrGroup.POST("", middleware.RequirePermission("ecr.create"), deps.ChangeHandler.CreateECR)
				ecrGroup.GET("/:id", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetECR)
				ecrGroup.PUT("/:id", middleware.RequirePermission("ecr.edit"), deps.ChangeHandler.UpdateECR)
				ecrGroup.DELETE("/:id", middleware.RequirePermission("ecr.edit"), deps.ChangeHandler.DeleteECR)
				ecrGroup.POST("/:id/submit", middleware.RequirePermission("ecr.create"), deps.ChangeHandler.SubmitECR)
				ecrGroup.POST("/:id/items", middleware.RequirePermission("ecr.edit"), deps.ChangeHandler.AddChangeItem)
				ecrGroup.DELETE("/items/:id", middleware.RequirePermission("ecr.edit"), deps.ChangeHandler.DeleteChangeItem)
				ecrGroup.POST("/:id/impact", middleware.RequirePermission("ecr.review"), deps.ChangeHandler.SaveImpact)
				ecrGroup.POST("/:id/reviews", middleware.RequirePermission("ecr.review"), deps.ChangeHandler.SubmitReview)
				ecrGroup.POST("/:id/approvals", middleware.RequirePermission("ecr.approve"), deps.ChangeHandler.SubmitApproval)
				ecrGroup.POST("/:id/convert-to-eco", middleware.RequirePermission("eco.create"), deps.ChangeHandler.CreateECO)
			}

			// Engineering Change Orders (ECO)
			ecoGroup := protected.Group("/eco")
			{
				ecoGroup.GET("", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetECOs)
				ecoGroup.POST("", middleware.RequirePermission("eco.create"), deps.ChangeHandler.CreateECO)
				ecoGroup.GET("/:id", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetECO)
				ecoGroup.POST("/:id/approve", middleware.RequirePermission("eco.approve"), deps.ChangeHandler.ApproveECO)
				ecoGroup.GET("/:id/change-preview", middleware.RequirePermission("changes.view"), deps.ChangeHandler.GetChangePreview)
				ecoGroup.POST("/:id/implement", middleware.RequirePermission("eco.implement"), deps.ChangeHandler.ImplementECO)
				ecoGroup.POST("/:id/verify", middleware.RequirePermission("eco.verify"), deps.ChangeHandler.VerifyECO)
				ecoGroup.POST("/:id/close", middleware.RequirePermission("eco.approve"), deps.ChangeHandler.CloseECO)
			}

			// Components
			components := protected.Group("/components")
			{
				components.GET("", middleware.RequirePermission("components.view"), deps.CompHandler.GetAll)
				components.GET("/:id", middleware.RequirePermission("components.view"), deps.CompHandler.GetByID)
				components.POST("", middleware.RequirePermission("components.create"), deps.CompHandler.Create)
				components.PUT("/:id", middleware.RequirePermission("components.edit"), deps.CompHandler.Update)
				components.DELETE("/:id", middleware.RequirePermission("components.delete"), deps.CompHandler.Delete)
				components.POST("/:id/documents", middleware.RequirePermission("components.edit"), deps.CompHandler.AddDocument)
				components.DELETE("/:id/documents/:docId", middleware.RequirePermission("components.edit"), deps.CompHandler.DeleteDocument)
			}

			// Categories
			categories := protected.Group("/categories")
			{
				categories.GET("/tree", middleware.RequirePermission("categories.view"), deps.CategoryHandler.GetTree)
				categories.GET("", middleware.RequirePermission("categories.view"), deps.CategoryHandler.GetFlat)
				categories.GET("/:id", middleware.RequirePermission("categories.view"), deps.CategoryHandler.GetByID)
				categories.POST("", middleware.RequirePermission("categories.create"), deps.CategoryHandler.Create)
				categories.PUT("/:id", middleware.RequirePermission("categories.edit"), deps.CategoryHandler.Update)
				categories.DELETE("/:id", middleware.RequirePermission("categories.delete"), deps.CategoryHandler.Delete)
			}

			// Manufacturers
			manufacturers := protected.Group("/manufacturers")
			{
				manufacturers.GET("", middleware.RequirePermission("manufacturers.view"), deps.PartnerHandler.GetAllManufacturers)
				manufacturers.GET("/:id", middleware.RequirePermission("manufacturers.view"), deps.PartnerHandler.GetManufacturerByID)
				manufacturers.POST("", middleware.RequirePermission("manufacturers.create"), deps.PartnerHandler.CreateManufacturer)
				manufacturers.PUT("/:id", middleware.RequirePermission("manufacturers.edit"), deps.PartnerHandler.UpdateManufacturer)
				manufacturers.DELETE("/:id", middleware.RequirePermission("manufacturers.delete"), deps.PartnerHandler.DeleteManufacturer)
			}

			// Suppliers
			suppliers := protected.Group("/suppliers")
			{
				suppliers.GET("", middleware.RequirePermission("suppliers.view"), deps.PartnerHandler.GetAllSuppliers)
				suppliers.GET("/:id", middleware.RequirePermission("suppliers.view"), deps.PartnerHandler.GetSupplierByID)
				suppliers.POST("", middleware.RequirePermission("suppliers.create"), deps.PartnerHandler.CreateSupplier)
				suppliers.PUT("/:id", middleware.RequirePermission("suppliers.edit"), deps.PartnerHandler.UpdateSupplier)
				suppliers.DELETE("/:id", middleware.RequirePermission("suppliers.delete"), deps.PartnerHandler.DeleteSupplier)
			}

			// Units
			units := protected.Group("/units")
			{
				units.GET("", middleware.RequirePermission("units.view"), deps.UnitHandler.GetAll)
				units.GET("/:id", middleware.RequirePermission("units.view"), deps.UnitHandler.GetByID)
				units.POST("", middleware.RequirePermission("units.create"), deps.UnitHandler.Create)
				units.PUT("/:id", middleware.RequirePermission("units.edit"), deps.UnitHandler.Update)
				units.DELETE("/:id", middleware.RequirePermission("units.delete"), deps.UnitHandler.Delete)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.GetAll)
				users.GET("/:id", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.GetByID)
				users.POST("", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.Create)
				users.PUT("/:id", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.Update)
				users.PATCH("/:id/status", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.ToggleStatus)
				users.POST("/:id/reset-password", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.ResetPassword)
				users.DELETE("/:id", middleware.RequirePermission("administration.manageUsers"), deps.UserHandler.Delete)
			}

			// Roles & Permissions
			roles := protected.Group("/roles")
			{
				roles.GET("", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.GetAll)
				roles.GET("/:id", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.GetByID)
				roles.POST("", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.Create)
				roles.PUT("/:id", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.Update)
				roles.PUT("/:id/permissions", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.UpdatePermissions)
			}
			protected.GET("/permissions", middleware.RequirePermission("administration.manageRoles"), deps.RoleHandler.GetAllPermissions)

			// Activities
			protected.GET("/activity", middleware.RequirePermission("administration.auditLogs"), deps.ActivityHandler.GetAll)

			// Exchange Rates (Phase 3 Multi-Currency Engine)
			exchangeRates := protected.Group("/exchange-rates")
			{
				exchangeRates.GET("/latest", middleware.RequirePermission("products.view"), deps.ExchangeRateHandler.GetLatestRate)
				exchangeRates.POST("/refresh", middleware.RequirePermission("products.edit"), deps.ExchangeRateHandler.RefreshRate)
				exchangeRates.POST("/calculate-pricing", middleware.RequirePermission("products.view"), deps.ExchangeRateHandler.CalculatePricing)
			}

			// Cross-Cutting Platform: Product Data Import & Master Data Migration
			dataImport := protected.Group("/data-import")
			{
				dataImport.GET("/batches", middleware.RequirePermission("dataimport.view"), deps.ImportHandler.GetBatches)
				dataImport.POST("/batches", middleware.RequirePermission("dataimport.upload"), deps.ImportHandler.CreateBatch)
				dataImport.GET("/batches/:id", middleware.RequirePermission("dataimport.view"), deps.ImportHandler.GetBatchByID)
				dataImport.DELETE("/batches/:id", middleware.RequirePermission("dataimport.upload"), deps.ImportHandler.DeleteBatch)
				dataImport.POST("/batches/:id/files", middleware.RequirePermission("dataimport.upload"), deps.ImportHandler.UploadFiles)
				dataImport.POST("/batches/:id/parse", middleware.RequirePermission("dataimport.validate"), deps.ImportHandler.ParseBatch)
				dataImport.POST("/batches/:id/mapping", middleware.RequirePermission("dataimport.validate"), deps.ImportHandler.ApplyMapping)
				dataImport.POST("/batches/:id/ai-classify", middleware.RequirePermission("dataimport.ai"), deps.ImportHandler.RunAIClassification)
				dataImport.POST("/batches/:id/generate-codes", middleware.RequirePermission("dataimport.validate"), deps.ImportHandler.GenerateMissingCodes)
				dataImport.GET("/batches/:id/rows", middleware.RequirePermission("dataimport.view"), deps.ImportHandler.GetStagedRows)
				dataImport.POST("/batches/:id/approve", middleware.RequirePermission("dataimport.approve"), deps.ImportHandler.ApproveBatch)
				dataImport.POST("/batches/:id/commit", middleware.RequirePermission("dataimport.execute"), deps.ImportHandler.CommitBatch)
				dataImport.PUT("/rows/:rowId", middleware.RequirePermission("dataimport.review"), deps.ImportHandler.TriageRow)
			}

			// Developer & System Integrations: API Access Tokens
			apiTokens := protected.Group("/api-tokens")
			{
				apiTokens.GET("", deps.APITokenHandler.GetAll)
				apiTokens.POST("", deps.APITokenHandler.Generate)
				apiTokens.DELETE("/:id", deps.APITokenHandler.Revoke)
			}
		}
	}

	return r
}
