# PRODUCT DATA IMPORT & MASTER DATA MIGRATION
## Detailed Technical Implementation Plan

**PROJECT:** Product & Master Data Control Center  
**TARGET RELEASE:** v1.1.0 Platform Capability  
**CLASSIFICATION:** Cross-Cutting Data Migration Engine (Non-Phase 7)  
**STATUS:** **READY FOR FUTURE IMPLEMENTATION SIGN-OFF**  

---

## 1. Implementation Milestones

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ IMP-M1: STAGING SCHEMA & GOLANG EXCELIZE STREAM PARSER                     │
│ DDL migrations for imp_batches, imp_files, imp_staged_rows + parser daemon │
├─────────────────────────────────────────────────────────────────────────────┤
│ IMP-M2: NORMALIZATION & RULE VALIDATION PIPELINE                            │
│ Column mapping, currency/date/UoM clean, validation rule engine             │
├─────────────────────────────────────────────────────────────────────────────┤
│ IMP-M3: MULTI-LAYER DUPLICATE DETECTION & AI CO-PILOT ADAPTER               │
│ Exact & fuzzy matching engine + optional LLM classification service         │
├─────────────────────────────────────────────────────────────────────────────┤
│ IMP-M4: HUMAN REVIEW WORKSPACE & FRONTEND WIZARD UI                         │
│ Multi-step wizard, review grid, row inspector, partial approval UI          │
├─────────────────────────────────────────────────────────────────────────────┤
│ IMP-M5: ATOMIC MASTER DATA DISPATCHER & AUDIT PROVENANCE COMMIT             │
│ Transactional dispatcher to Products, Components, and Services + receipt    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Target Database Schema (MySQL 8.0 DDL)

```sql
-- 1. Import Batches Table
CREATE TABLE IF NOT EXISTS imp_batches (
    id VARCHAR(36) PRIMARY KEY,
    batch_code VARCHAR(32) NOT NULL UNIQUE,
    import_year VARCHAR(4) NOT NULL DEFAULT '2026',
    source_type VARCHAR(64) NOT NULL DEFAULT 'EXCEL_WORKBOOK',
    source_notes TEXT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'UPLOADED',
    files_count INT NOT NULL DEFAULT 0,
    total_rows INT NOT NULL DEFAULT 0,
    valid_rows INT NOT NULL DEFAULT 0,
    warning_rows INT NOT NULL DEFAULT 0,
    error_rows INT NOT NULL DEFAULT 0,
    approved_rows INT NOT NULL DEFAULT 0,
    rejected_rows INT NOT NULL DEFAULT 0,
    uploaded_by VARCHAR(128) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_batch_year (import_year),
    INDEX idx_batch_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Import Files Table
CREATE TABLE IF NOT EXISTS imp_files (
    id VARCHAR(36) PRIMARY KEY,
    batch_id VARCHAR(36) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    stored_filepath VARCHAR(512) NOT NULL,
    file_sha256 VARCHAR(64) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    sheet_count INT NOT NULL DEFAULT 1,
    sheet_manifest JSON NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'UPLOADED',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (batch_id) REFERENCES imp_batches(id) ON DELETE CASCADE,
    INDEX idx_file_batch (batch_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Staged Rows Table
CREATE TABLE IF NOT EXISTS imp_staged_rows (
    id VARCHAR(36) PRIMARY KEY,
    batch_id VARCHAR(36) NOT NULL,
    file_id VARCHAR(36) NOT NULL,
    sheet_name VARCHAR(128) NOT NULL,
    source_row_index INT NOT NULL,
    raw_data_snapshot JSON NOT NULL,
    normalized_data JSON NULL,
    item_classification VARCHAR(32) NOT NULL DEFAULT 'PRODUCT', -- PRODUCT, COMPONENT, SPARE_PART, ACCESSORY, SERVICE
    product_type VARCHAR(32) NULL, -- TRADING, PROJECT, RND (when PRODUCT)
    target_category_id VARCHAR(36) NULL,
    target_mfg_id VARCHAR(36) NULL,
    target_supplier_id VARCHAR(36) NULL,
    unit_cost DECIMAL(15,4) NULL,
    currency VARCHAR(8) NOT NULL DEFAULT 'IDR',
    validation_status VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- VALID, WARNING, ERROR
    validation_messages JSON NULL,
    duplicate_status VARCHAR(32) NOT NULL DEFAULT 'NEW_RECORD', -- NEW_RECORD, POSSIBLE_DUPLICATE, MATCHED_EXISTING
    matched_master_id VARCHAR(36) NULL,
    matched_master_type VARCHAR(32) NULL,
    ai_suggested_type VARCHAR(32) NULL,
    ai_confidence DECIMAL(5,2) NULL,
    ai_payload JSON NULL,
    review_decision VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED, EDITED
    reviewed_by VARCHAR(128) NULL,
    target_master_table VARCHAR(64) NULL,
    target_master_id VARCHAR(36) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (batch_id) REFERENCES imp_batches(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES imp_files(id) ON DELETE CASCADE,
    INDEX idx_staged_batch_val (batch_id, validation_status),
    INDEX idx_staged_batch_dec (batch_id, review_decision)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 3. Golang Backend Architecture (`backend/internal/dataimport`)

```
backend/internal/dataimport/
├── handlers/
│   └── import_handler.go      # REST endpoints for batch upload, wizard & commit
├── services/
│   ├── parser_service.go      # excelize/v2 streaming row extractor
│   ├── normalizer_service.go  # Currency, UoM, string sanitization engine
│   ├── validator_service.go   # Rule validation matrix
│   ├── duplicate_service.go   # Multi-layer duplicate detection
│   ├── ai_service.go          # Optional classification suggestion co-pilot
│   └── migration_service.go   # Atomic dispatcher to production Master Data
└── models/
    └── import_models.go       # Staging structs, DTOs & payloads
```

---

## 4. REST API Endpoints Specification

| Method | Endpoint | Description | Required Role |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/import/batches` | List import batches with year filter (`?year=2026`) | Viewer+ |
| `POST` | `/api/v1/import/batches` | Create new batch container | Uploader+ |
| `POST` | `/api/v1/import/batches/:id/files` | Multipart upload of XLS/XLSX workbooks | Uploader+ |
| `POST` | `/api/v1/import/batches/:id/parse` | Execute excelize parsing on selected sheets | Reviewer+ |
| `POST` | `/api/v1/import/batches/:id/map-columns` | Save header-to-master mapping definition | Reviewer+ |
| `POST` | `/api/v1/import/batches/:id/validate` | Run normalization, validation & duplicate rules | Reviewer+ |
| `POST` | `/api/v1/import/batches/:id/ai-classify` | Request AI classification suggestions | Reviewer+ |
| `GET` | `/api/v1/import/batches/:id/rows` | Fetch staged rows with filters (`?status=error`) | Viewer+ |
| `PUT` | `/api/v1/import/batches/:id/rows/:rowId` | Triage / edit specific staged row | Reviewer+ |
| `POST` | `/api/v1/import/batches/:id/approve` | Formally approve validated rows | Approver+ |
| `POST` | `/api/v1/import/batches/:id/commit` | Execute transactional insertion to Master Data | Executor+ |

---

## 5. Frontend Component Architecture (`src/views/import/`)

```
src/
├── views/import/
│   ├── ImportDashboardView.vue       # Batches list, year filter & KPI summary
│   ├── ImportWizardView.vue          # 11-step interactive import wizard
│   ├── ImportReviewWorkspace.vue     # High-density triage grid & inspector drawer
│   └── ImportReceiptView.vue         # Migration receipt & lineage audit view
├── components/import/
│   ├── ColumnMappingGrid.vue         # Interactive source-to-target field mapper
│   ├── RowInspectorDrawer.vue        # Right-hand split inspector with raw snapshot
│   ├── AiSuggestionBadge.vue         # AI classification chip with confidence %
│   └── DuplicateMatchModal.vue       # Side-by-side comparison modal
└── stores/
    └── importStore.js                # Pinia store managing staging state & wizard steps
```

---

## 6. Verification & Forensic Quality Gates

Before opening the module for production data migration:
1. **Streaming Load Test:** Verify parsing of a $20,000$ row Excel workbook within $< 4.0\text{s}$ with $< 80\text{MB}$ memory consumption.
2. **Formula Injection Sanitization Test:** Verify that malicious formula payloads (`=cmd|' /C ...'!A0`) are sanitized without execution.
3. **Partial Approval Integrity Test:** Verify that a batch with $900$ valid rows and $100$ error rows imports exactly $900$ items while leaving $100$ staged for correction.
4. **Data Lineage Traceability Test:** Verify that every generated Product Master record has a valid provenance envelope referencing its exact Excel source row.
