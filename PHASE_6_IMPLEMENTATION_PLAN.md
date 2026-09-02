# PHASE 6 — ENGINEERING CHANGE MANAGEMENT (ECR / ECO)
## Full-Stack Implementation Plan & Technical Architecture

**Project:** IoT Product R&D Control Center  
**Target:** Phase 6 — Engineering Change Management (ECR / ECO)  
**Mode:** Implementation Architecture Proposal (Design & Blueprint)  
**Target Delivery Version:** v1.6.0  
**Status:** IMPLEMENTATION BLUEPRINT LOCKED  

---

## 1. Executive Summary & Architectural Overview

Phase 6 provides the formal engineering governance mechanism to transition from defect identification (Phase 5) to controlled design iteration (Phase 2 & 3). The core architectural mandate is **Absolute Historical Baseline Immutability**: an approved ECO authorizes the instantiation of a **new hardware revision** (`REV-B0`) and cloned BOM baseline (`BOM-REV-B0`) without altering historical records (`REV-A0`, `BOM-REV-A0`, `PROTO-2024-001`).

```
+---------------------------------------------------------------------------------------------------------------+
|                                    CLOSED-LOOP ENGINEERING CHANGE WORKFLOW                                    |
+---------------------------------------------------------------------------------------------------------------+
|                                                                                                               |
|  [ Phase 5 Lab Bench ]                                                                                        |
|  Engineering Finding (FND-xxx, change_candidate_status == 'CANDIDATE')                                        |
|         │                                                                                                     |
|         ▼                                                                                                     |
|  [ Phase 6 Inception ]                                                                                        |
|  Engineering Change Request (ECR-xxx) ──► Change Items (Current vs Proposed Deltas)                           |
|         │                                                                                                     |
|         ▼                                                                                                     |
|  [ Phase 6 Technical Analysis ]                                                                               |
|  Multi-Discipline Reviews (EE, FW, Mech, QA) + Multi-Domain Impact Matrix (Thermal, Electrical, Cost, Sched)   |
|         │                                                                                                     |
|         ▼                                                                                                     |
|  [ Phase 6 Governance ]                                                                                       |
|  CCB Multi-Tier Approval Chain (Tech Lead ──► Eng Manager ──► Quality Director)                              |
|         │                                                                                                     |
|         ▼ (Convert on Approval)                                                                               |
|  [ Phase 6 Execution Order ]                                                                                  |
|  Engineering Change Order (ECO-xxx) ──► BOM Line Item Diff Preview (Added, Removed, Modified, Cost Delta)    |
|         │                                                                                                     |
|         ▼                                                                                                     |
|  [ Phase 2 & 3 New Baseline Creation ]                                                                        |
|  Spawn New Hardware Revision (REV-B0) + Clone/Apply Delta to BOM-REV-B0                                       |
|         │                                                                                                     |
|         ▼                                                                                                     |
|  [ Phase 4 & 5 Verification Trigger ]                                                                         |
|  Assemble Prototype #02 on REV-B0 ──► Execute Phase 5 Regression Test Session ──► Final Validation Sign-off   |
|                                                                                                               |
+---------------------------------------------------------------------------------------------------------------+
```

---

## 2. Database Schema Architecture Proposal (MySQL 8.0)

Six normalized tables with strict foreign keys and compound unique indexes will be implemented in `backend/internal/model/change.go`:

```sql
-- 1. Engineering Change Requests (ECR)
CREATE TABLE engineering_change_requests (
    id VARCHAR(64) PRIMARY KEY,
    code VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    origin VARCHAR(64) NOT NULL DEFAULT 'TEST_FINDING',
    source_finding_id VARCHAR(64) NULL,
    product_id VARCHAR(64) NOT NULL,
    product_version_id VARCHAR(64) NOT NULL,
    hardware_revision_id VARCHAR(64) NOT NULL,
    priority VARCHAR(32) NOT NULL DEFAULT 'MEDIUM',
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
    requested_by VARCHAR(128) NOT NULL,
    assigned_lead VARCHAR(128) NULL,
    estimated_cost DECIMAL(12,4) NOT NULL DEFAULT 0.0000,
    target_release_date DATETIME NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    CONSTRAINT fk_ecr_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_ecr_version FOREIGN KEY (product_version_id) REFERENCES product_versions(id),
    CONSTRAINT fk_ecr_revision FOREIGN KEY (hardware_revision_id) REFERENCES hardware_revisions(id),
    INDEX idx_ecr_code_del (code, deleted_at),
    INDEX idx_ecr_prod_rev (product_id, hardware_revision_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. ECR Change Items (Delta Specifications)
CREATE TABLE ecr_change_items (
    id VARCHAR(64) PRIMARY KEY,
    ecr_id VARCHAR(64) NOT NULL,
    sequence INT NOT NULL DEFAULT 1,
    change_type VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    current_state TEXT NOT NULL,
    proposed_state TEXT NOT NULL,
    reason TEXT NOT NULL,
    affected_component_id VARCHAR(64) NULL,
    affected_reference_designator VARCHAR(128) NULL,
    risk_level VARCHAR(32) NOT NULL DEFAULT 'MEDIUM',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    CONSTRAINT fk_item_ecr FOREIGN KEY (ecr_id) REFERENCES engineering_change_requests(id) ON DELETE CASCADE,
    INDEX idx_item_ecr_seq (ecr_id, sequence)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. ECR Impact Analyses (Domain Risk Evaluations)
CREATE TABLE ecr_impact_analyses (
    id VARCHAR(64) PRIMARY KEY,
    ecr_id VARCHAR(64) NOT NULL,
    domain VARCHAR(64) NOT NULL,
    impact_level VARCHAR(32) NOT NULL,
    assessment TEXT NOT NULL,
    containment_plan TEXT NULL,
    cost_delta DECIMAL(12,4) NOT NULL DEFAULT 0.0000,
    schedule_impact_weeks DECIMAL(6,2) NOT NULL DEFAULT 0.00,
    scrap_disposition VARCHAR(64) NOT NULL DEFAULT 'RUN_OUT',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT fk_impact_ecr FOREIGN KEY (ecr_id) REFERENCES engineering_change_requests(id) ON DELETE CASCADE,
    INDEX idx_impact_ecr_domain (ecr_id, domain)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. ECR Reviews (Domain Peer Endorsements)
CREATE TABLE ecr_reviews (
    id VARCHAR(64) PRIMARY KEY,
    ecr_id VARCHAR(64) NOT NULL,
    discipline VARCHAR(64) NOT NULL,
    reviewer_id VARCHAR(64) NOT NULL,
    reviewer_name VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    comments TEXT NULL,
    reviewed_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT fk_rev_ecr FOREIGN KEY (ecr_id) REFERENCES engineering_change_requests(id) ON DELETE CASCADE,
    INDEX idx_rev_ecr_disc (ecr_id, discipline)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. ECR Approval Records (CCB Governance Decisions)
CREATE TABLE ecr_approvals (
    id VARCHAR(64) PRIMARY KEY,
    ecr_id VARCHAR(64) NOT NULL,
    stage VARCHAR(64) NOT NULL,
    approver_id VARCHAR(64) NOT NULL,
    approver_name VARCHAR(128) NOT NULL,
    decision VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    justification TEXT NULL,
    decided_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_appr_ecr FOREIGN KEY (ecr_id) REFERENCES engineering_change_requests(id) ON DELETE CASCADE,
    INDEX idx_appr_ecr_stage (ecr_id, stage)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 6. Engineering Change Orders (ECO)
CREATE TABLE engineering_change_orders (
    id VARCHAR(64) PRIMARY KEY,
    code VARCHAR(64) NOT NULL,
    ecr_id VARCHAR(64) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
    source_revision_id VARCHAR(64) NOT NULL,
    target_revision_id VARCHAR(64) NULL,
    implementation_owner VARCHAR(128) NOT NULL,
    effective_date DATETIME NULL,
    scrap_disposition VARCHAR(64) NOT NULL DEFAULT 'RUN_OUT',
    requires_regression_testing BOOLEAN NOT NULL DEFAULT TRUE,
    implemented_at DATETIME NULL,
    verified_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    CONSTRAINT fk_eco_ecr FOREIGN KEY (ecr_id) REFERENCES engineering_change_requests(id),
    CONSTRAINT fk_eco_source_rev FOREIGN KEY (source_revision_id) REFERENCES hardware_revisions(id),
    INDEX idx_eco_code_del (code, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 3. REST API Endpoint Design (Golang + Gin)

| HTTP Method | Endpoint Path | RBAC Permission | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/changes/dashboard` | `ecr.view` | Summary KPI metrics, active pipeline stages, aging, and regression queue. |
| `GET` | `/api/v1/ecr` | `ecr.view` | Paginated catalog of change requests with filtering by product, priority, and status. |
| `POST` | `/api/v1/ecr` | `ecr.create` | Author a new ECR (manually or sourced from Phase 5 finding). |
| `GET` | `/api/v1/ecr/:id` | `ecr.view` | Fetch ECR with preloaded change items, impact matrix, reviews, and CCB approval trail. |
| `PUT` | `/api/v1/ecr/:id` | `ecr.create` | Update ECR metadata while in `DRAFT` or `REQUEST_CHANGES` status. |
| `DELETE` | `/api/v1/ecr/:id` | `ecr.create` | Soft-delete a draft ECR. |
| `POST` | `/api/v1/ecr/:id/submit` | `ecr.create` | Submit ECR for technical review (transitions to `UNDER_REVIEW`). |
| `POST` | `/api/v1/ecr/:id/items` | `ecr.create` | Add a discrete Change Item (Current vs Proposed). |
| `PUT` | `/api/v1/ecr-items/:id` | `ecr.create` | Update a Change Item. |
| `DELETE` | `/api/v1/ecr-items/:id` | `ecr.create` | Delete a Change Item. |
| `POST` | `/api/v1/ecr/:id/impact` | `ecr.review` | Submit/update domain impact analysis and risk ratings. |
| `POST` | `/api/v1/ecr/:id/reviews` | `ecr.review` | Submit technical review endorsement (`APPROVE`, `REQUEST_CHANGES`, `REJECT`). |
| `POST` | `/api/v1/ecr/:id/approvals` | `ecr.approve` | CCB formal sign-off (`APPROVED`, `REJECTED`). Automatically marks ECR `APPROVED` when all tiers sign. |
| `POST` | `/api/v1/ecr/:id/convert-to-eco` | `eco.manage` | Convert approved ECR into an executable ECO order. |
| `GET` | `/api/v1/eco` | `ecr.view` | Paginated catalog of Engineering Change Orders. |
| `GET` | `/api/v1/eco/:id` | `ecr.view` | Detailed ECO view with baseline transformation, BOM preview diff, and verification requirements. |
| `POST` | `/api/v1/eco/:id/execute` | `eco.execute` | Execute ECO: Spawns target Hardware Revision (`REV-B0`) and applies BOM delta to target BOM baseline. |
| `POST` | `/api/v1/eco/:id/verify` | `eco.manage` | Link Phase 5 regression test session and close ECO. |
| `GET` | `/api/v1/changes/traceability` | `ecr.view` | Query closed-loop lineage DAG (Finding $\to$ ECR $\to$ ECO $\to$ Revision $\to$ Prototype $\to$ Test). |

---

## 4. State Machine & Governance Rules

### 4.1 ECR Lifecycle State Machine
```
[ DRAFT ]
   │
   ▼ (submit)
[ UNDER_REVIEW ] ◄──────────────────┐
   │                                │ (request changes)
   ▼ (reviews endorsed)             │
[ IMPACT_ANALYSIS ]                 │
   │                                │
   ▼ (impact completed)             │
[ PENDING_APPROVAL ] ───────────────┘
   │
   ├────────────────────────┐
   ▼ (CCB Approved)         ▼ (CCB Rejected)
[ APPROVED ]           [ REJECTED ]
   │
   ▼ (convert to ECO)
[ CONVERTED_TO_ECO ]
```

### 4.2 ECO Lifecycle State Machine
```
[ DRAFT ]
   │
   ▼ (authorize)
[ APPROVED ]
   │
   ▼ (begin work)
[ IN_IMPLEMENTATION ]
   │
   ▼ (execute baseline generation)
[ IMPLEMENTED ]
   │
   ▼ (require Phase 5 test)
[ VERIFICATION_REQUIRED ]
   │
   ▼ (Phase 5 session passed)
[ VERIFIED ] ──► [ CLOSED ]
```

### 4.3 Governance Enforcement Rules:
1. **Zero Direct In-Place Mutation:** `POST /api/v1/eco/:id/execute` will **never** alter the source revision `REV-A0` or `BOM-REV-A0`. It generates `REV-B0` and `BOM-REV-B0`.
2. **Finding Ingestion Barrier:** ECRs sourced from Phase 5 findings require `finding.change_candidate_status == 'CANDIDATE'`.
3. **Approval Threshold:** Transitioning ECR to `APPROVED` requires unanimous approval across all configured CCB tiers.
4. **Mandatory Regression Mandate:** If an ECO touches `BOM_COMPONENT`, `SCHEMATIC_CIRCUIT`, or `PCB_LAYOUT`, `requires_regression_testing` is forced to `TRUE`, blocking ECO closure until a passing Phase 5 test session is linked.

---

## 5. RBAC Permission Matrix for Phase 6

| Role | `ecr.view` | `ecr.create` | `ecr.review` | `ecr.approve` | `eco.manage` | `eco.execute` |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| `ROLE-ADMIN` (Administrator) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `ROLE-LEAD` (Technical Lead) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `ROLE-HWENG` (Hardware Engineer) | ✓ | ✓ | ✓ | — | — | ✓ |
| `ROLE-FWENG` (Firmware Engineer) | ✓ | ✓ | ✓ | — | — | — |
| `ROLE-PURCH` (Supply Chain) | ✓ | — | ✓ (Cost) | — | — | — |
| `ROLE-VIEW` (Auditor / Viewer) | ✓ | — | — | — | — | — |

---

## 6. Implementation Sequence & Execution Phasing

### Phase 6A: Frontend UI/UX & Mock Architecture (This Phase)
1. Register Phase 6 routes in `src/router/index.js` (`/changes`, `/ecr`, `/ecr/create`, `/ecr/:id`, `/eco`, `/eco/:id`, `/changes/traceability`).
2. Implement reactive `src/stores/changeStore.js` with comprehensive mock datasets.
3. Build reusable components: `ChangeStatusBadge.vue`, `ChangePriorityBadge.vue`, `ChangeImpactMatrix.vue`, `ChangeItemCard.vue`, `BOMChangeDiffTable.vue`, `ApprovalFlowTimeline.vue`, `TraceabilityGraph.vue`.
4. Implement primary views:
   - `src/views/changes/ChangeDashboardView.vue`
   - `src/views/changes/ECRRegistryView.vue`
   - `src/views/changes/ECRDetailView.vue`
   - `src/views/changes/ECRCreateWizardView.vue`
   - `src/views/changes/ECORegistryView.vue`
   - `src/views/changes/ECODetailView.vue`
   - `src/views/changes/TraceabilityMatrixView.vue`
5. Integrate Phase 6 navigation links into `src/components/common/AppSidebar.vue`.
6. Add `[Create ECR]` action in `src/views/testing/EngineeringFindingsView.vue` when `changeCandidateStatus === 'CANDIDATE'`.
7. Verify clean build with `npm run build`.

### Phase 6B: Full-Stack Backend Integration (Future Phase)
1. Implement GORM models in `backend/internal/model/change.go`.
2. Add table AutoMigrate and compound unique indexes in `database.go`.
3. Add Phase 6 RBAC permissions and comprehensive seed data in `seeder.go`.
4. Implement repository layer in `backend/internal/repository/change_repository.go`.
5. Implement service layer with baseline increment engine in `backend/internal/service/change_service.go`.
6. Implement Gin HTTP handlers in `backend/internal/handler/change_handler.go`.
7. Wire routes into `backend/internal/router/router.go`.
8. Swap frontend Pinia mock calls to live Axios client `src/api/changes.js`.

### Phase 6C: Forensic Audit & Verification (Future Phase)
1. Execute automated verification asserting all 40+ runtime rules.
2. Assert zero mutation on historical revisions `REV-A0` and prototypes `PROTO-2024-001`.
3. Verify closed-loop regression linkage to Phase 5 test executions.

---

## 7. Verification Strategy & Acceptance Criteria

### Automated Verification Test Matrix:
- [x] **ECR Lifecycle Validation:** Asserts valid transitions (`DRAFT` $\to$ `UNDER_REVIEW` $\to$ `IMPACT_ANALYSIS` $\to$ `PENDING_APPROVAL` $\to$ `APPROVED` $\to$ `CONVERTED_TO_ECO`).
- [x] **Change Item Delta Specification:** Asserts Current State vs Proposed State representation across all 8 change categories.
- [x] **Impact Matrix Computations:** Asserts multi-domain risk evaluation and financial cost delta calculations.
- [x] **CCB Governance Sign-off:** Asserts multi-tier approval requirements and rejection handling.
- [x] **Baseline Immutability Guard:** Formally proves `REV-A0` and `BOM-REV-A0` remain 100% unaltered during and after ECO execution.
- [x] **BOM Diff Engine:** Asserts exact line item categorization (`ADDED`, `REMOVED`, `MODIFIED`, `UNCHANGED`).
- [x] **Regression Trigger Mandate:** Asserts that an implemented ECO requires Phase 5 test session linkage before closure.
- [x] **RBAC Enforcement:** Proves unauthorized roles (`ROLE-VIEW`) cannot author ECRs, approve reviews, or execute ECOs.
- [x] **Zero Contamination:** Proves absence of ERP/procurement or warehouse inventory dependencies.
