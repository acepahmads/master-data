# PRODUCT DATA IMPORT & MASTER DATA MIGRATION
## UI/UX Design System & Workspace Blueprint

**MODULE ROUTE:** `/data-import`  
**DESIGN AESTHETIC:** Clean Enterprise Engineering Grid (Slate / Navy / Brand / Emerald)  
**STATUS:** **DESIGN SPECIFICATION READY (GO)**  

---

## 1. Information Architecture & Navigation

The module lives in the main application under **PRODUCT & MASTER DATA** navigation:

```
PRODUCT & MASTER DATA
├── Products
├── Components
├── Categories
├── Manufacturers
├── Suppliers
├── Units (UoM)
└── Data Import & Migration  <-- NEW ROUTE (/data-import)
```

---

## 2. Import Dashboard & Batch Management (`/data-import`)

### Top KPI Executive Strip:
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ TOTAL BATCHES   │ ACTIVE 2026 ROWS│ VALIDATED (100%)│ WARNINGS/ERRORS │ AI SUGGESTIONS  │
│ 14 Batches      │ 12,480 Rows     │ 11,920 Valid    │ 560 Issues      │ 3,240 Generated │
│ 2016-2026 Total │ 2026 Migration  │ 95.5% Quality   │ Triage Required │ 94.2% Conf Avg  │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

### Year Scope Switcher:
- **`[● 2026 Active Scope]`** *(Default highlighted)*: Shows active import batches queued for Master Data deployment.
- **`[2016–2025 Historical Archive]`**: Shows read-only historical staging batches for data lineage and pricing audits.

### Batches Registry Table:
| Batch Code | Source File(s) | Import Year | Staged Rows | Valid | Warnings | Errors | AI Ready | Status | Action |
| :--- | :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `IMP-2026-0001` | `2026-Q1-Commercial.xlsx` (3 sheets) | **2026** | 1,450 | 1,380 | 52 | 18 | 410 | `READY_FOR_REVIEW` | `[Open Wizard]` |
| `IMP-2026-0002` | `2026-Sensors-Sourcing.xlsx` (1 sheet) | **2026** | 820 | 820 | 0 | 0 | 215 | `APPROVED` | `[Execute Import]` |
| `IMP-2025-ARCH` | `2025-Annual-Log.xlsx` (12 sheets) | **2025** | 18,200 | 17,900 | 200 | 100 | — | `HISTORICAL_ARCHIVE`| `[View Archive]` |

---

## 3. Multi-Step Import Wizard Workflow

```mermaid
graph LR
    S1[1. Upload Files] --> S2[2. Config]
    S2 --> S3[3. Sheets]
    S3 --> S4[4. Column Map]
    S4 --> S5[5. Normalize]
    S5 --> S6[6. Validate]
    S6 --> S7[7. Duplicates]
    S7 --> S8[8. AI Suggest]
    S8 --> S9[9. Human Review]
    S9 --> S10[10. Approval]
    S10 --> S11[11. Commit Master]
```

### Step-by-Step UX Interactions:

#### Step 1: Upload Files
- Drag-and-drop dropzone supporting multiple files (`.xlsx`, `.xls`, `.csv`).
- File list queue showing: File name, file size, SHA-256 hash, and upload progress bar.

#### Step 2: Source Configuration
- **Import Year Selector:** Default `2026`.
- **Source Category:** Commercial Sales / Procurement / Warehouse / Vendor Catalog.
- **Auditor Notes:** Textarea for migration rationale.

#### Step 3: Sheet Selection & Header Detection
- Interactive workbook tree showing each sheet, total row count, non-empty columns, and detected header row index (default row 1).
- Checkboxes to include/exclude specific sheets.

#### Step 4: Column Mapping Matrix
- Visual mapping grid:
  - `Source Excel Header` (e.g. *"Nama Barang"*, *"Harga Beli"*, *"Merk"*, *"Satuan"*)
  - `Target Master Field` (Dropdown: *Product Name*, *Purchase Price*, *Manufacturer*, *UoM*, *Ignore Column*)
  - `Sample Data Preview` (First 3 rows of actual cell values)
  - `Smart Auto-Map Button`: Maps known Indonesian & English terminology automatically.

#### Step 5: Normalization Engine
- Real-time preview of transformation rules:
  - Currency: `"Rp 4.500.000"` $\to$ `4500000` (`IDR`)
  - Unit: `"pcs"`, `"Buah"`, `"unit"` $\to$ `PCS`
  - String clean: Removes leading/trailing whitespace, non-printable characters.

#### Step 6: Validation Inspector
- Filterable status tabs: `[All (1,450)]`, `[Valid (1,380)]`, `[Warnings (52)]`, `[Errors (18)]`.
- Inline badges explaining validation diagnostics (e.g. *"Missing Category"*, *"Negative Price"*, *"Invalid Date Format"*).

#### Step 7: Duplicate Detection Inspector
- Visual match pill:
  - `[NEW RECORD]` (Emerald) — Zero collisions in Master database.
  - `[POSSIBLE DUPLICATE]` (Amber) — High name similarity (85%+) with existing record.
  - `[EXACT MASTER MATCH]` (Blue) — MPN or Code matches existing Master. Options: *Skip*, *Update Specs*, or *Create Alternative SKU*.

#### Step 8: AI Co-Pilot Classification Suggestion
- Displayed as a purple-tinted pill:
  - `AI: pH Sensor (PRODUCT / TRADING) • 94% Confidence`
  - Action buttons: `[✓ Accept]`, `[✎ Edit]`, `[✕ Ignore]`.

#### Step 9: Human Review & Triage Workspace (Split Screen)
- **Left Panel:** High-density interactive grid with quick approve/reject toggles.
- **Right Inspector Drawer:**
  - Original Excel row snapshot.
  - Normalized values.
  - AI reasoning payload.
  - Target destination selector (*Product Master* vs *Component Master* vs *Service Item*).

#### Step 10: Import Approval & Executive Sign-off
- Summary modal showing exact counts of items to be inserted into each master table.
- Formally signed by authenticated user (`Alex Chen, Super Admin`).

#### Step 11: Atomic Master Data Dispatch
- Real-time animated progress bar committing rows in database transactions.
- Provides downloadable migration receipt and direct links to newly created Product / Component master records.

---

## 4. High-Density Review Table Column Layout

To maintain full consistency with the zero-horizontal-scroll design system:

| Col # | Column Title | Width | Content |
| :---: | :--- | :---: | :--- |
| 1 | **Row** | `w-12` | `#142` (Source sheet row index) |
| 2 | **Source Name** | `min-w-[160px]` | Raw string from Excel |
| 3 | **Normalized Name** | `min-w-[160px]` | Clean canonical title |
| 4 | **Item Classification** | `w-28` | Pill: `PRODUCT` / `COMPONENT` / `SERVICE` |
| 5 | **Product Archetype** | `w-20` | `TRADING` / `PROJECT` / `RND` (if Product) |
| 6 | **Category** | `w-24` | Master category badge |
| 7 | **Cost / Currency** | `w-24` | `4,500,000 IDR` |
| 8 | **Validation** | `w-20` | `● Valid` / `▲ Warning` / `✕ Error` |
| 9 | **Duplicate** | `w-24` | `New` / `Match: PRD-001` |
| 10| **AI Suggestion** | `w-32` | `Sensors & Nodes (94%)` |
| 11| **Decision** | `w-24` | `[Approve]` `[Reject]` `[Edit]` |

---

## 5. Mobile & Tablet Responsiveness

- On screens $< 1024\text{px}$, the review grid collapses into expandable cards with the right inspector opening as a full-screen drawer.
- The top wizard navigation renders as a compact breadcrumb indicator with active step name.
