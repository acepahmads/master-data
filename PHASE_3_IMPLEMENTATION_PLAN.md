# PHASE 3 — ENGINEERING BOM & COST FOUNDATION IMPLEMENTATION PLAN

## IoT Product R&D Control Center
**Date:** August 31, 2026  
**Status:** DRAFT & READY FOR APPROVAL  

---

## 1. Executive Summary & Objective

Phase 3 implements the **Engineering Bill of Materials (BOM) & Cost Foundation** for the **IoT Product R&D Control Center**.

The system will enable hardware engineering teams to define multi-level electronic and mechanical assemblies, assign components from the Component Master, specify engineering reference designators (e.g. `R1, R2`, `C101-C104`, `U1`), calculate recursive line and assembly costs, and establish pricing foundations for both R&D Products and Project Solutions.

---

## 2. Core Architecture Rules & Ownership

### Mandatory Ownership Model:
$$\text{R\&D Product} \longrightarrow \text{Product Version} \longrightarrow \text{Hardware Revision} \longrightarrow \text{Engineering BOM} \longrightarrow \text{BOM Items}$$

1. **Hardware Revision is the EXCLUSIVE Owner**:
   - Each Hardware Revision has at most **one active Engineering BOM** in Phase 3.
   - Products and Product Versions do not own BOMs directly.
2. **Product Type Boundary Isolation**:
   - **`TRADING` Product**: Commercial product with purchase/selling price. **STRICTLY PROHIBITED** from owning Product Versions, Hardware Revisions, or Engineering BOMs.
   - **`PROJECT` Product**: System solution package with `project_items`. **STRICTLY PROHIBITED** from owning Engineering BOMs or copying internal R&D BOM components.
   - **`RND` Product**: The **ONLY** product type permitted to have Engineering BOMs (through its Hardware Revisions).

---

## 3. Database Schema Design

### 3.1. Extend `components` Table
- `estimated_unit_cost DECIMAL(12,4) DEFAULT 0.0000`
- `currency VARCHAR(8) DEFAULT 'USD'`
- `cost_source VARCHAR(64) DEFAULT 'MANUAL_ESTIMATE'`

### 3.2. New Table: `engineering_boms`
- `id VARCHAR(64) PRIMARY KEY`
- `hardware_revision_id VARCHAR(64) UNIQUE NOT NULL` (Indexed, Foreign Key)
- `bom_code VARCHAR(64) UNIQUE NOT NULL`
- `name VARCHAR(255) NOT NULL`
- `description TEXT`
- `status VARCHAR(32) DEFAULT 'Active'` (Draft, Active, Obsolete, Archived)
- `total_cost DECIMAL(12,4) DEFAULT 0.0000`
- `currency VARCHAR(8) DEFAULT 'USD'`
- `target_margin DECIMAL(5,2) DEFAULT 35.00`
- `estimated_selling_price DECIMAL(12,4) DEFAULT 0.0000`
- `notes TEXT`
- `created_by VARCHAR(255)`
- `created_at DATETIME(3)`, `updated_at DATETIME(3)`, `deleted_at DATETIME(3)`

### 3.3. New Table: `bom_items`
- `id VARCHAR(64) PRIMARY KEY`
- `engineering_bom_id VARCHAR(64) NOT NULL` (Indexed)
- `parent_item_id VARCHAR(64)` (Indexed, Nullable for Root Items)
- `item_type VARCHAR(32) NOT NULL` (`COMPONENT`, `SUB_ASSEMBLY`)
- `name VARCHAR(255) NOT NULL`
- `component_id VARCHAR(64)` (Indexed, Nullable for Sub-Assemblies)
- `component_part_number VARCHAR(128)`
- `component_category VARCHAR(100)`
- `quantity DECIMAL(10,4) DEFAULT 1.0000`
- `unit_id VARCHAR(64)`
- `unit_code VARCHAR(32) DEFAULT 'PCS'`
- `unit_cost DECIMAL(12,4) DEFAULT 0.0000`
- `line_cost DECIMAL(12,4) DEFAULT 0.0000`
- `reference_designators TEXT` (e.g. `R1, R2, R3`, `C1, C2`)
- `position INT DEFAULT 0`
- `notes TEXT`
- `created_at DATETIME(3)`, `updated_at DATETIME(3)`, `deleted_at DATETIME(3)`

---

## 4. Cost Calculation & Pricing Mathematical Engine

1. **Leaf Component Line Cost**:
   $$\text{LineCost}_{\text{leaf}} = \text{Quantity} \times \text{UnitCost}$$
2. **Sub-Assembly Line Cost**:
   $$\text{LineCost}_{\text{assembly}} = \text{Quantity} \times \sum_{c \in \text{Children}} \text{LineCost}_c$$
3. **Total BOM Cost**:
   $$\text{TotalCost}_{\text{BOM}} = \sum_{i \in \text{RootItems}} \text{LineCost}_i$$
   *(Zero double counting: Sub-assemblies roll up children once).*
4. **Estimated Selling Price**:
   $$\text{EstimatedSellingPrice} = \frac{\text{TotalCost}_{\text{BOM}}}{1 - (\text{TargetMargin} / 100)}$$
5. **Project Cost Aggregation**:
   $$\text{TotalCost}_{\text{Project}} = \sum \left( \text{TradingPurchasePrice} + \text{RNDActiveBOMCost} + \text{ComponentLineCost} \right)$$

---

## 5. Backend Implementation Plan

1. **Models (`backend/internal/model/bom.go`, `component.go`)**: GORM mappings and associations.
2. **Repositories (`backend/internal/repository/bom_repository.go`)**: CRUD, tree extraction, sorting.
3. **Services (`backend/internal/service/bom_service.go`)**:
   - Hierarchy validation & cycle prevention.
   - Cost calculation engine.
   - Cost rollups for R&D Products and Project Solutions.
4. **Handlers & Router (`backend/internal/handler/bom_handler.go`, `router.go`)**:
   - `GET /revisions/:id/bom`
   - `POST /revisions/:id/bom`
   - `PUT /revisions/:id/bom`
   - `POST /boms/:id/items`
   - `PUT /bom-items/:id`
   - `DELETE /bom-items/:id`
   - `POST /boms/:id/items/reorder`
   - `GET /products/:id/cost-summary`
5. **RBAC & Activity Logging**:
   - Permissions: `bom.view`, `bom.create`, `bom.edit`, `bom.delete`.
   - MySQL `activity_logs` records all BOM lifecycle and item changes.

---

## 6. Frontend UI/UX Plan

1. **`EngineeringBOMWorkspaceView.vue` (`/revisions/:id/bom`)**:
   - **Header**: Hardware revision badge, version context, total BOM cost, margin %, estimated selling price.
   - **Left Panel (BOM Tree Explorer)**: Interactive nested hierarchy with collapsible assemblies, item counters, and cost contribution indicators.
   - **Center Panel (Multi-Level BOM Table)**: Level badges, Part Number, MPN, Category, Qty, Unit, Unit Cost, Line Cost, Reference Designators (JetBrains Mono font), Actions.
   - **Right Panel (Selected Item Inspector)**: Complete Component specs, manufacturer, reference designator pills, cost breakdown, quick navigation to Component Master.
   - **Cost Visuals**: Top 5 Cost Drivers, Cost Contribution by Category bar.
   - **Add/Edit Item Modal**: Item type selection (`COMPONENT` or `SUB_ASSEMBLY`), searchable Component Master picker, quantity, unit, reference designators, notes.
2. **`HardwareRevisionDetailView.vue`**:
   - Add **"Engineering BOM & Cost"** tab + **"Open BOM Workspace"** button.
3. **`ProductDetailView.vue`**:
   - Embed **R&D Product Cost Summary Card** (Estimated BOM Cost, Target Margin, Estimated Selling Price, link to active revision BOM).
   - Embed **Project Cost Aggregation Summary Card** (Trading Product Cost, R&D Product Cost, Direct Component Cost, Total Solution Cost).

---

## 7. Phase Boundary Protections (Strict Exclusions)

Phase 3 strictly excludes:
- Procurement, RFQs, Purchase Orders (PO), Purchase Requests (PR), Supplier Quotations.
- Inventory, Warehouse, Stock tracking, Serial number tracking.
- Manufacturing, Production planning, Shop floor management.
- Invoicing, Accounting, Tax calculations.
- ECO / ECR workflows.
- Prototype management, Testing management.
- Live currency conversion / exchange rate feeds.
