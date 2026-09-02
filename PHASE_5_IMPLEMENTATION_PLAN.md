# PHASE 5 IMPLEMENTATION PLAN
## IoT Product R&D Control Center
**Module**: Phase 5 — Testing & Validation Control  
**Document Type**: Engineering Execution Roadmap & Staging Plan  
**Status**: DRAFT IMPLEMENTATION PLAN (Awaiting Approval)  
**Date**: August 2026  

---

## 1. Execution Overview & Staging Order

Phase 5 will be executed across 7 strictly ordered stages to ensure zero regression of existing Phase 1–4 capabilities, mathematical accuracy of parametric tolerances, and database persistence fidelity.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       PHASE 5 EXECUTION STAGING ORDER                       │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
                 Stage 1: Database Schema & Relational Models
                                       ▼
                 Stage 2: Backend Repository & GORM Queries
                                       ▼
                 Stage 3: Business Logic Service Layer
                                       ▼
                 Stage 4: REST API Handlers & Router RBAC
                                       ▼
                 Stage 5: Frontend Store & Vue 2 Pinia Views
                                       ▼
                 Stage 6: Automated Testing & Runtime Verification
                                       ▼
                 Stage 7: Strict Forensic Runtime Audit
```

---

## 2. Detailed Staging Roadmap

### Stage 1: Database Schema & Relational Models
1. **Define GORM Models** in `backend/internal/model/testing.go`:
   - `TestPlan`: Master protocol template with code, name, version, category, and target revision.
   - `TestCase`: Parametric specification with procedure, expected limits, measurement unit, and evaluation type.
   - `PrototypeTestSession`: Test plan instance bound to a `PrototypeBuild` with `test_plan_snapshot_json`.
   - `TestExecution`: Execution runs (`run_number`, `status`, `executed_by`, `notes`).
   - `TestMeasurement`: Quantitative / qualitative measurement data with automated tolerance evaluation.
   - `TestEvidence`: File attachment records with mime-type and storage path.
   - `EngineeringFinding`: Defect / anomaly with severity, category, disposition, and `is_eco_candidate`.
   - `PrototypeValidationDecision`: Formal sign-off record (`VALIDATED`, `CONDITIONALLY_VALIDATED`, `REJECTED`).
2. **Update Database Auto-Migration** in `backend/internal/database/database.go`:
   - Register all 8 Phase 5 entities in `db.AutoMigrate(...)`.
   - Create compound indexes:
     - `idx_session_build_plan` on `prototype_test_sessions(prototype_build_id, test_plan_id)`.
     - `idx_exec_session_case_run` on `test_executions(test_session_id, test_case_id, run_number)`.
     - `idx_val_build_unique` on `prototype_validation_decisions(prototype_build_id)`.
3. **Update Seeder Data** in `backend/internal/database/seeder.go`:
   - Seed Phase 5 permissions (`testing.view`, `testplans.manage`, `testing.execute`, `findings.manage`, `testing.validate`).
   - Assign permissions to roles (`ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-HWENG`, `ROLE-FWENG`, `ROLE-QA`, `ROLE-VIEW`).
   - Seed realistic Master Test Plan Templates:
     - `TP-GW400X-EVT-01` (Smart Gateway EVT Comprehensive Hardware Verification Suite) with 12 parametric test cases across Functional, Power, Electrical, RF, and Thermal categories.
   - Seed initial test session for `PROTO-2024-001` with frozen snapshot, sample measurements, and finding records.

---

### Stage 2: Backend Repository Layer
1. **Create Repository** in `backend/internal/repository/testing_repository.go`:
   - Master Test Plans & Test Cases CRUD with pagination, category filter, and target revision lookup.
   - Prototype Test Sessions query by prototype ID with preloaded executions and snapshot data.
   - Test Executions query, creation of multi-run iterations (`run_number`), and measurement batching.
   - Test Evidence insertion and file linkage.
   - Engineering Findings filtered query (by prototype, severity, status, `is_eco_candidate`).
   - Validation Decisions retrieve and atomic commit.
   - Testing Dashboard KPIs aggregator (pass rate %, active builds under test, open critical/high findings, validation counts).

---

### Stage 3: Business Logic Service Layer
1. **Create Service** in `backend/internal/service/testing_service.go`:
   - `CreateTestPlan`: Version allocation, uniqueness validation, test case assembly.
   - `InitiatePrototypeTestSession`:
     - Validate product type = `RND` and prototype exists.
     - Freeze Test Plan and Test Cases into immutable JSON snapshot (`test_plan_snapshot_json`).
     - Auto-generate initial `TestExecution` records (`run_number = 1`, `status = 'NOT_STARTED'`).
     - Execute within atomic database transaction.
     - Log activity to `activity_logs`.
   - `RecordTestExecution`:
     - Parse numeric measurements against min/max bounds.
     - Auto-evaluate `PASSED` / `FAILED` status based on tolerance criteria.
     - Support multi-sample calculations (Min, Max, Mean).
     - If failed, offer automated creation of linked `EngineeringFinding`.
   - `UpdateEngineeringFinding`:
     - Validate severity levels (`CRITICAL`, `HIGH`, `MEDIUM`, `LOW`, `INFO`).
     - Handle disposition transitions (`FIXED_IN_REWORK`, `DEFERRED_TO_ECO`, `ACCEPTED_RISK`, `CLOSED`).
   - `SubmitValidationDecision`:
     - Guard: Block `VALIDATED` if any unresolved `CRITICAL` or `HIGH` findings exist.
     - Require mandatory condition notes if `CONDITIONALLY_VALIDATED`.
     - Update prototype build status to `COMPLETED` or `VALIDATED`.
     - Log formal validation audit event.

---

### Stage 4: REST API Handlers & Router Wiring
1. **Create Handlers** in `backend/internal/handler/testing_handler.go`:
   - Complete REST endpoints with standard JSON response envelopes (`{ success, message, data, meta, errors }`).
2. **Register Routes** in `backend/internal/router/router.go`:
   - `/api/v1/testing/dashboard` (GET) -> `testing.view`
   - `/api/v1/test-plans` (GET, POST) -> `testplans.manage`
   - `/api/v1/test-plans/:id` (GET, PUT, DELETE) -> `testplans.manage`
   - `/api/v1/test-plans/:id/cases` (GET, POST) -> `testplans.manage`
   - `/api/v1/prototypes/:id/test-sessions` (GET, POST) -> `testing.execute`
   - `/api/v1/test-sessions/:id` (GET) -> `testing.view`
   - `/api/v1/test-sessions/:id/executions/:caseId` (POST) -> `testing.execute`
   - `/api/v1/test-executions/:id/measurements` (POST) -> `testing.execute`
   - `/api/v1/test-executions/:id/evidence` (POST) -> `testing.execute`
   - `/api/v1/findings` (GET) -> `testing.view`
   - `/api/v1/findings/:id` (PUT, DELETE) -> `findings.manage`
   - `/api/v1/prototypes/:id/validation` (GET, POST) -> `testing.validate`
3. **Wire in `backend/cmd/main.go`**:
   - Inject repositories, services, and handlers.

---

### Stage 5: Frontend API Client, Pinia Store & UI Views
1. **Create API Client** in `src/api/testing.js`:
   - Axios wrapper for all `/api/v1/testing/...`, `/test-plans/...`, `/findings/...`, and prototype validation endpoints.
2. **Create Pinia Store** in `src/stores/testingStore.js`:
   - State: `dashboardStats`, `testPlans`, `activeSession`, `findings`, `validationDecision`, `loading`.
   - Actions: `fetchDashboard()`, `fetchTestPlans()`, `fetchSessionByPrototype(protoId)`, `recordExecution()`, `recordMeasurement()`, `createFinding()`, `submitValidation()`.
3. **Build Vue 2 UI Views**:
   - `src/views/testing/TestingDashboardView.vue` (`/testing`): Command center with KPI metrics, active test matrix, severity distribution chart, and failure stream.
   - `src/views/testing/TestPlanRegistryView.vue` (`/test-plans`): Master test protocol library with versioning, clone, and category filters.
   - `src/views/testing/TestPlanDetailView.vue` (`/test-plans/:id`): Interactive test case builder and tolerance limit editor.
   - `src/views/testing/EngineeringFindingsView.vue` (`/findings`): Cross-build defect tracker with severity filters and ECO candidate flags.
4. **Integrate into Prototype Detail Workspace (`PrototypeDetailView.vue`)**:
   - Add **Tab 7: Testing & Validation**.
   - Embed category-grouped interactive test matrix, measurement logging dialog, findings drawer, and validation sign-off banner.
5. **Update Navigation in `AppSidebar.vue`**:
   - Add `Testing & Validation (Phase 5)` navigation link group.

---

### Stage 6: Automated Testing & Verification
1. **Build Compilations**:
   - `go build ./...` in `backend/` ➔ 0 errors.
   - `npm run build` in root ➔ 0 errors.
2. **Create End-to-End Automated Verification Script** (`scratch/verify_phase5_fullstack.go`):
   - Test master test plan creation & test case limit validation.
   - Test test session initiation & immutable test plan snapshot freezing.
   - Test parametric measurement compliance ($Min \le Measured \le Max$).
   - Test multi-run iterations (Run 1 Fail $\rightarrow$ Run 2 Pass).
   - Test automatic and manual Engineering Finding creation.
   - Test Validation Decision governance (critical finding blocking rule).
   - Test RBAC permissions (viewer 403 checks, validation sign-off permissions).
   - Test activity logging integration in `activity_logs`.
   - Test server restart persistence in MySQL.

---

### Stage 7: Strict Forensic Runtime Audit
1. Execute strict forensic audit verifying 0 production file mutations during audit, 100% test pass rate, and zero regression of Phases 1–4.
2. Generate `PHASE_5_STRICT_FORENSIC_AUDIT.md`.

---

## 3. Regression Protection Strategy

To guarantee that Phase 5 does not impact existing Phase 1, Phase 2, Phase 2C, Phase 3, and Phase 4 systems:
1. **Isolated Table Names**: All Phase 5 tables use explicit prefixes (`test_plans`, `test_cases`, `prototype_test_sessions`, `test_executions`, `engineering_findings`, `prototype_validation_decisions`).
2. **Preserve Existing Models**: Zero modifications to `Product`, `ProductVersion`, `HardwareRevision`, `EngineeringBOM`, or `PrototypeBuild` core columns.
3. **Automated Regression Suite Execution**: Re-run `verify_phase2c.go`, `verify_phase3.go`, `verify_phase3_multicurrency.go`, and `verify_phase4b_fullstack.go` after every backend build.
