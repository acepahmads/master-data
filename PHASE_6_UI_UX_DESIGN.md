# PHASE 6 — ENGINEERING CHANGE MANAGEMENT (ECR / ECO)
## UI/UX & Frontend Architecture Design Specification

**Project:** IoT Product R&D Control Center  
**Target:** Phase 6 — Engineering Change Management (ECR / ECO)  
**Mode:** UI/UX & Architectural Design (Non-Mutating / Design Only)  
**System Version:** v1.6.0-Design  
**Author:** Antigravity Principal R&D Architect  
**Status:** ARCHITECTURE DESIGN LOCKED  

---

## 1. Executive Summary

Phase 6 introduces a rigorous, enterprise-grade **Engineering Change Management** system to the **IoT Product R&D Control Center**. It bridges the gap between defect discovery in prototype testing (Phase 5) and formal hardware revision baselines (Phase 2 & 3), establishing an auditable, governed pipeline for hardware and firmware evolution.

```
+-------------------------------------------------------------------------------------------------------------+
|                                    PHASE 6 CLOSED-LOOP TRACEABILITY CHAIN                                   |
+-------------------------------------------------------------------------------------------------------------+
|  Phase 5 Lab Bench            Phase 6 Governance                     Phase 2 & 3 Baselines     Phase 4 & 5  |
|                                                                                                             |
|  +--------------------+       +--------------------+                 +--------------------+    +----------+ |
|  | Engineering Finding| ----> | Change Candidate   |                 | Target Revision    |    | Prototype| |
|  | (Severity: CRIT/HI)|       | Status: CANDIDATE  |                 | (REV-B0: DRAFT)    | -> | Build #02| |
|  +--------------------+       +---------+----------+                 +---------+----------+    +----+-----+ |
|                                         |                                      |                    |       |
|                                         v                                      v                    v       |
|                               +--------------------+                 +--------------------+    +----------+ |
|                               | Engineering Change |                 | New BOM Baseline   |    |Regression| |
|                               | Request (ECR)      |                 | (BOM-REV-B0)       |    | Test Run | |
|                               +---------+----------+                 +--------------------+    +----+-----+ |
|                                         |                                                           |       |
|                                         v                                                           v       |
|                               +--------------------+                                           +----------+ |
|                               | Multi-Discipline   |                                           |Validation| |
|                               | Impact & Review    |                                           | Decision | |
|                               +---------+----------+                                           +----------+ |
|                                         |                                                                   |
|                                         v                                                                   |
|                               +--------------------+                                                        |
|                               | Governance Approval|                                                        |
|                               | Chain (Sign-off)   |                                                        |
|                               +---------+----------+                                                        |
|                                         |                                                                   |
|                                         v                                                                   |
|                               +--------------------+                                                        |
|                               | Engineering Change | -------------------------------------------------------+ |
|                               | Order (ECO: Order) |                                                        |
|                               +--------------------+                                                        |
+-------------------------------------------------------------------------------------------------------------+
```

### Core Tenet: Absolute Historical Immutability
Phase 6 **NEVER mutates historical engineering baselines in-place**. When an ECO is executed:
- Historical Revision `REV-A0` remains **100% frozen**.
- Historical BOM `BOM-REV-A0` remains **100% frozen**.
- Historical Prototype Builds (`PROTO-2024-001`), Test Executions, and Validation Decisions remain **immutable**.
- The ECO authorizes and spawns a **NEW baseline** (`REV-B0`, `BOM-REV-B0`), which becomes the parent for subsequent prototype builds and regression test sessions.

---

## 2. Distinction: ECR vs. ECO

| Dimension | Engineering Change Request (ECR) | Engineering Change Order (ECO) |
| :--- | :--- | :--- |
| **Core Question** | *"What needs to change, why, and what are the multi-domain impacts?"* | *"Who is authorized to execute the approved change, how, and by when?"* |
| **Purpose** | Investigation, problem definition, impact analysis, technical feasibility, and cross-functional sign-off. | Authorized implementation blueprint, baseline incrementation, BOM transformation, and verification trigger. |
| **Lifecycle State** | `DRAFT` $\to$ `SUBMITTED` $\to$ `UNDER_REVIEW` $\to$ `IMPACT_ANALYSIS` $\to$ `PENDING_APPROVAL` $\to$ `APPROVED` $\to$ `CONVERTED_TO_ECO` | `DRAFT` $\to$ `APPROVED` $\to$ `IN_IMPLEMENTATION` $\to$ `IMPLEMENTED` $\to$ `VERIFICATION_REQUIRED` $\to$ `VERIFIED` $\to$ `CLOSED` |
| **Author / Actor** | Design Engineers, Test Engineers, Quality Leads, Systems Architects. | Engineering Manager, Change Control Board (CCB), Lead Hardware Engineer. |
| **Output** | Formal justification, Change Items matrix, Impact Matrix, multi-role review endorsements. | New Hardware Revision (`REV-B0`), target BOM baseline, scrap/rework dispositions, regression test mandate. |

---

## 3. Core Domain Entities & Relationships

```
+---------------------------------------------------------------------------------------------+
|                                  PHASE 6 LOGICAL DATA MODEL                                 |
+---------------------------------------------------------------------------------------------+
|                                                                                             |
|   +-----------------------+              +-----------------------+                          |
|   |  EngineeringFinding   | 1          0..1| EngineeringChangeReq  |                          |
|   |  (Phase 5 Source)     |------------->|  (ECR)                |                          |
|   +-----------------------+              +-----------+-----------+                          |
|                                                      | 1                                    |
|                                                      |                                      |
|                                    +-----------------+-----------------+                    |
|                                    | 1..N            | 1..N            | 1                  |
|                                    v                 v                 v                    |
|                         +---------------+ +---------------+ +---------------+               |
|                         |   ChangeItem  | |   ECRReview   | |  ECRApproval  |               |
|                         | (Delta State) | | (Tech Domain) | |  (Governance) |               |
|                         +---------------+ +---------------+ +---------------+               |
|                                                      | 1                                    |
|                                                      | (Converts when APPROVED)             |
|                                                      v                                      |
|                                          +-----------------------+                          |
|                                          | EngineeringChangeOrder|                          |
|                                          |  (ECO)                |                          |
|                                          +-----------+-----------+                          |
|                                                      | 1                                    |
|                                                      |                                      |
|                                    +-----------------+-----------------+                    |
|                                    | 1               | 1               | 1                  |
|                                    v                 v                 v                    |
|                         +---------------+ +---------------+ +---------------+               |
|                         | TargetRevision| | BOMPreviewDiff| |RegressionReq  |               |
|                         | (HardwareRev) | | (Delta Tree)  | |  (Phase 5)    |               |
|                         +---------------+ +---------------+ +---------------+               |
|                                                                                             |
+---------------------------------------------------------------------------------------------+
```

### 3.1 Entity Definitions

#### 1. Engineering Change Request (`ECR`)
- `id`: Unique identifier (`ECR-2026-001`)
- `title`: Summary of proposed change
- `description`: Detailed technical motivation
- `origin`: Category source (`TEST_FINDING`, `CUSTOMER_REQUIREMENT`, `ENGINEERING_REVIEW`, `OBSOLESCENCE`, `COST_REDUCTION`, `REGULATORY`, `RELIABILITY`, `SUPPLIER_CHANGE`, `SECURITY`)
- `sourceFindingId`: Optional reference to `EngineeringFinding` (`FND-xxx`)
- `productId`: Parent R&D Product (`PRD-xxx`)
- `productVersionId`: Affected Product Version (`PV-xxx`)
- `hardwareRevisionId`: Affected Hardware Revision (`REV-xxx`)
- `priority`: Urgency classification (`CRITICAL`, `HIGH`, `MEDIUM`, `LOW`)
- `status`: State machine status (`DRAFT`, `SUBMITTED`, `UNDER_REVIEW`, `IMPACT_ANALYSIS`, `PENDING_APPROVAL`, `APPROVED`, `REJECTED`, `CANCELLED`, `CONVERTED_TO_ECO`)
- `requestedBy`: Author user name and ID
- `assignedLead`: Assigned systems engineer
- `estimatedCost`: Estimated implementation expense ($USD)
- `targetReleaseDate`: Targeted engineering delivery milestone

#### 2. Change Item (`ChangeItem`)
Discrete technical delta items within an ECR:
- `id`: Item identifier (`CHG-001`)
- `ecrId`: Parent ECR ID
- `sequence`: Ordering index (1, 2, 3...)
- `changeType`: Engineering domain (`SCHEMATIC_CIRCUIT`, `PCB_LAYOUT`, `BOM_COMPONENT`, `COMPONENT_VALUE`, `MECHANICAL_ENCLOSURE`, `FIRMWARE_LOGIC`, `DOCUMENTATION`, `MANUFACTURING_NOTE`)
- `title`: Short summary (e.g. *"LDO Heat Dissipation Rework"*)
- `currentState`: Description/part of existing design (e.g. *"TPS7A4700 in 3x3 QFN"*)
- `proposedState`: Description/part of proposed design (e.g. *"TI LP5907 in SOT-23 with exposed thermal pad"*)
- `reason`: Specific technical justification
- `affectedComponentId`: Optional reference to master `Component` (`CMP-xxx`)
- `affectedReferenceDesignator`: PCB designator (e.g. `U3`, `R12`, `C44`)
- `riskLevel`: Risk tier (`LOW`, `MEDIUM`, `HIGH`, `CRITICAL`)

#### 3. ECR Impact Analysis (`ECRImpactMatrix`)
Multi-dimensional evaluation across technical and operational criteria:
- **Technical Dimensions:**
  - `Electrical`: Signal integrity, power distribution, voltage rail tolerances.
  - `Thermal`: Power dissipation, hotspot distribution, junction temperatures.
  - `Mechanical`: Keep-out zones, height clearance, mounting screw holes, connectors.
  - `Firmware`: Register map changes, pin muxing, peripheral driver initialization.
  - `Software / Cloud`: Payload structures, telemetry schemas, API protocols.
  - `Safety & Regulatory`: FCC/CE EMC compliance, RoHS/REACH compliance, creepage/clearance.
- **Operational Dimensions:**
  - `BOM Cost Delta`: Net unit BOM cost change ($USD).
  - `Schedule Impact`: Estimated schedule slip in engineering weeks.
  - `Scrap / Rework Risk`: Existing prototype inventory disposition (`SCRAP_ALL`, `REWORK_EXISTING`, `USE_AS_IS`).

#### 4. ECR Review (`ECRReview`)
Domain-specific peer endorsements:
- `discipline`: Review specialty (`HARDWARE_EE`, `FIRMWARE_EMBEDDED`, `MECHANICAL_CAD`, `VALIDATION_QA`, `SYSTEMS_ARCH`)
- `reviewerId` / `reviewerName`: Assigned expert
- `status`: Review verdict (`PENDING`, `APPROVED`, `REQUEST_CHANGES`, `REJECTED`)
- `comments`: Technical feedback
- `reviewedAt`: Timestamp

#### 5. ECR Approval Governance (`ECRApproval`)
Formal Change Control Board (CCB) sign-offs:
- `stage`: Approval tier (`TIER_1_TECH_LEAD`, `TIER_2_ENG_MANAGER`, `TIER_3_QUALITY_DIRECTOR`)
- `approverId` / `approverName`: Designated authority
- `decision`: Sign-off status (`PENDING`, `APPROVED`, `REJECTED`)
- `justification`: Governance statement
- `decidedAt`: Timestamp

#### 6. Engineering Change Order (`ECO`)
Authorized implementation order spawned upon ECR approval:
- `id`: Unique ECO identifier (`ECO-2026-001`)
- `ecrId`: Source ECR ID
- `title`: Implementation order title
- `status`: Order status (`DRAFT`, `APPROVED`, `IN_IMPLEMENTATION`, `IMPLEMENTED`, `VERIFICATION_REQUIRED`, `VERIFIED`, `CLOSED`, `CANCELLED`)
- `sourceRevisionId`: Base baseline (`REV-A0`)
- `targetRevisionId`: Spawned increment baseline (`REV-B0`)
- `implementationOwner`: Assigned lead engineer
- `effectiveDate`: Target execution date
- `scrapDisposition`: Existing inventory policy (`SCRAP_EXISTING`, `BENCH_REWORK`, `RUN_OUT`)
- `requiresRegressionTesting`: Boolean flag requiring Phase 5 test session trigger
- `implementedAt`: Execution timestamp
- `verifiedAt`: Final sign-off timestamp

---

## 4. Screen Architecture & Navigation Flows

```
+---------------------------------------------------------------------------------------------+
|                                    NAVIGATION SITE MAP                                      |
+---------------------------------------------------------------------------------------------+
|                                                                                             |
|   /changes                                                                                  |
|   ├── /ecr                   (ECR Registry Catalog)                                         |
|   │   ├── /ecr/create        (ECR Multi-Step Authoring Wizard)                              |
|   │   └── /ecr/:id           (ECR 7-Tab Detailed Engineering Workspace)                    |
|   ├── /eco                   (ECO Registry Catalog)                                         |
|   │   ├── /eco/create        (ECO Conversion & Execution Setup)                             |
|   │   └── /eco/:id           (ECO 6-Tab Implementation & Baseline Workspace)                |
|   └── /changes/traceability  (End-to-End Visual Lineage Matrix)                             |
|                                                                                             |
+---------------------------------------------------------------------------------------------+
```

---

## 5. UI/UX Wireframe & Layout Specifications

### 5.1 Screen 1: Change Control Center (`/changes`)
High-altitude executive dashboard displaying change velocity, aging, and review bottlenecks.

```
+---------------------------------------------------------------------------------------------------------------------+
| [⚡] ENGINEERING CHANGE CONTROL CENTER (ECR / ECO)                      [+ New Change Request] [Traceability Matrix]|
+---------------------------------------------------------------------------------------------------------------------+
|                                                                                                                     |
|  +-------------------+  +-------------------+  +-------------------+  +-------------------+  +-------------------+  |
|  | OPEN ECRs         |  | PENDING REVIEWS   |  | PENDING APPROVALS |  | ACTIVE ECOs       |  | REGRESSION QUEUE  |  |
|  |   12              |  |   5               |  |   3 (1 Critical)  |  |   4 (In Exec)     |  |   2 Revisions     |  |
|  | ▲ +2 this week    |  | ⏱ Avg Age: 3.2d   |  | 🔒 CCB Tier 2     |  | 🎯 1 Effective Soon|  | 🧪 Phase 5 Link   |  |
|  +-------------------+  +-------------------+  +-------------------+  +-------------------+  +-------------------+  |
|                                                                                                                     |
|  +-------------------------------------------------------------+  +-----------------------------------------------+ |
|  | ACTIVE CHANGE PIPELINE STAGES                               |  | CHANGE CANDIDATES FROM TESTING (PHASE 5)      | |
|  |                                                             |  | [From Defect Findings with CANDIDATE status]  | |
|  | ECR Inception    Review & Analysis   CCB Approval   ECO Exec|  |                                               | |
|  | [ 4 Drafts ] -> [ 5 In Review ] -> [ 3 Appr ] -> [ 4 ECO ] |  | • FND-2024-001: LDO Thermal Runaway (GW-400X) | |
|  |                                                             |  |   Severity: CRITICAL | [Create ECR ➔]         | |
|  | (Clicking any stage filters the global registry below)      |  | • FND-2024-002: RF 2.4GHz Antenna Return Loss | |
|  +-------------------------------------------------------------+  |   Severity: HIGH     | [Create ECR ➔]         | |
|                                                                   +-----------------------------------------------+ |
|  +---------------------------------------------------------------------------------------------------------------+ |
|  | RECENT ENGINEERING CHANGE REQUESTS & ORDERS                                                    [View Full ECRs]| |
|  | Code         Title                          Product       Target Rev  Priority    Status             Owner    | |
|  | ECR-2026-001 LDO Thermal Decoupling Rework  GW-400X Edge  REV-B0      CRITICAL    PENDING_APPROVAL   A. Chen  | |
|  | ECR-2026-002 Flash Memory Alternative Sourcing ESP32-Node REV-A1      MEDIUM      UNDER_REVIEW       M. Vance | |
|  | ECO-2026-001 LoRa PA Output Filter Tuning   GW-400X Edge  REV-B0      HIGH        IN_IMPLEMENTATION  S. Jenkins| |
|  +---------------------------------------------------------------------------------------------------------------+ |
+---------------------------------------------------------------------------------------------------------------------+
```

---

### 5.2 Screen 2: ECR Workspace (`/ecr/:id`)
Comprehensive 7-tab workbench for investigation, impact analysis, and review endorsements.

```
+---------------------------------------------------------------------------------------------------------------------+
| [← Back to ECRs]  ECR-2026-001: LDO Thermal Decoupling & Ground Plane Rework                                        |
| Status: [ PENDING_APPROVAL ]  Priority: [ CRITICAL ]  Product: Smart Gateway (GW-400X)  Revision: REV-A0 -> [REV-B0] |
+---------------------------------------------------------------------------------------------------------------------+
| [Tab 1: Overview] [Tab 2: Source & Reason] [Tab 3: Impact Analysis] [Tab 4: Proposed Changes] [Tab 5: Reviews]      |
| [Tab 6: Approval Chain] [Tab 7: Audit Timeline]                                                                     |
+---------------------------------------------------------------------------------------------------------------------+
|                                                                                                                     |
| TAB 4: PROPOSED CHANGE ITEMS (DELTA SPECIFICATION)                                                                  |
| +-----------------------------------------------------------------------------------------------------------------+ |
| | Item #1: Power Rail Regulation (BOM Component Replacement)                                    Risk: [ HIGH ]    | |
| | Area: Power Management (PMIC / LDO) | Reference Designators: U3, C18, C19                                           | |
| |                                                                                                                 | |
| | CURRENT BASELINE (REV-A0)                         PROPOSED BASELINE (TARGET: REV-B0)                            | |
| | Component: CMP-003 (TLV73333PDBVR)                Component: CMP-088 (LP5907MFX-3.3)                            | |
| | Package: SOT-23-5                                 Package: SOT-23-5 with Enhanced Exposed Thermal Ground Pad    | |
| | Max Junction Temp: 125°C                          Max Junction Temp: 150°C (Automotive Qualified)               | |
| | Unit Cost: $0.18                                  Unit Cost: $0.24 (+$0.06 delta)                               | |
| |                                                                                                                 | |
| | Technical Justification: Prototype thermal chamber tests (PROTO-2024-001) registered 128°C under max compute    | |
| | load. Upgrading to high-temperature LDO with thermal via stitching under pad.                                  | |
| +-----------------------------------------------------------------------------------------------------------------+ |
| | Item #2: PCB Ground Plane Thermal Relief (PCB Layout Change)                                  Risk: [ MEDIUM ]  | |
| | Area: Layer 2 Internal Ground Plane | Layer: L2 Ground Inner Layer                                              | |
| | CURRENT: Standard 4-spoke thermal relief spoke width 0.25mm                                                     | |
| | PROPOSED: Solid copper fill connect on pin 2 GND with 6x 0.3mm thermal stitching vias to bottom heatsink plane  | |
| +-----------------------------------------------------------------------------------------------------------------+ |
|                                                                                                                     |
| TAB 3: IMPACT ANALYSIS MATRIX                                                                                       |
| +-----------------------------------------------------------------------------------------------------------------+ |
| | Domain               Impact Level    Assessment & Containment Rationale                                         | |
| | Electrical           [ MEDIUM ]      Output ripple improves from 45mVpp to 18mVpp. PSRR increases +12dB.        | |
| | Thermal              [ POSITIVE ]    Junction temp reduced from 128°C to 84°C under 100% duty cycle.            | |
| | Mechanical           [ NONE ]        Component footprint footprint-compatible; zero enclosure clash.           | |
| | Firmware             [ NONE ]        Zero register/software changes required for analog LDO replacement.        | |
| | BOM Cost Impact      [ +$0.06/unit ] Net unit cost increase +$0.06. Covered within product margin tolerance.    | |
| | Existing Prototypes  [ SCRAP_LATER ] Existing PROTO-2024-001 usable for FW bench development until REV-B arrives| |
| +-----------------------------------------------------------------------------------------------------------------+ |
+---------------------------------------------------------------------------------------------------------------------+
```

---

### 5.3 Screen 3: ECO Implementation & BOM Diff Workspace (`/eco/:id`)
Operational execution workspace visualizing the before-and-after baseline transformation.

```
+---------------------------------------------------------------------------------------------------------------------+
| [⚡] ENGINEERING CHANGE ORDER: ECO-2026-001                                            [Print ECO Order] [Close ECO] |
| Source: [ECR-2026-001]  Base Baseline: [REV-A0] ➔ Target Baseline: [REV-B0] (DRAFT)  Status: [ IN_IMPLEMENTATION ]  |
+---------------------------------------------------------------------------------------------------------------------+
| [Tab 1: Order Details] [Tab 2: Baseline Transformation] [Tab 3: BOM Change Diff] [Tab 4: Verification Requirements] |
+---------------------------------------------------------------------------------------------------------------------+
|                                                                                                                     |
| TAB 3: BOM BASELINE CHANGE PREVIEW (REV-A0 vs. REV-B0)                                                              |
|                                                                                                                     |
|  Line  RefDes   Component MPN            Current (REV-A0)      Target (REV-B0)       Cost Delta    Change Status    |
|  -----------------------------------------------------------------------------------------------------------------  |
|  01    U1       STM32H753BIT6            1 pcs ($14.20)        1 pcs ($14.20)        $0.00         [ UNCHANGED ]    |
|  02    U3       LP5907MFX-3.3            TLV73333 ($0.18)      LP5907 ($0.24)        +$0.06        [ MODIFIED ] 🔄  |
|  03    C18      10uF 0603 X5R 10V        1 pcs ($0.03)         2 pcs ($0.06)         +$0.03        [ MODIFIED ] 🔄  |
|  04    R19      1.2k 0402 1%             ---                   1 pcs ($0.01)         +$0.01        [ ADDED ] ➕     |
|  05    D4       BAT54WS Schottky         1 pcs ($0.04)         ---                   -$0.04        [ REMOVED ] ❌   |
|  -----------------------------------------------------------------------------------------------------------------  |
|  TOTAL BOM COST:                         $44.90 / unit         $44.96 / unit         +$0.06/unit   (+0.13%)         |
|                                                                                                                     |
| TAB 4: POST-IMPLEMENTATION VERIFICATION MANDATE                                                                     |
| +-----------------------------------------------------------------------------------------------------------------+ |
| | [⚠️ REGRESSION TESTING REQUIRED UPON REV-B0 FABRICATION]                                                         | |
| |                                                                                                                 | |
| | Mandatory Test Protocols (Phase 5 Linkage):                                                                     | |
| | 1. Master Protocol: TP-GW400X-EVT (Thermal Chamber Heatmap @ 65°C Ambient)                                      | |
| | 2. Master Protocol: TP-GW400X-EVT (3.3V Power Rail Ripple & Noise Regulation)                                    | |
| |                                                                                                                 | |
| | Execution Trigger:                                                                                              | |
| | When Hardware Revision REV-B0 is submitted and prototype build is assembled in Phase 4, the testing team        | |
| | must initiate a linked validation session referencing ECO-2026-001.                                             | |
| +-----------------------------------------------------------------------------------------------------------------+ |
+---------------------------------------------------------------------------------------------------------------------+
```

---

### 5.4 Screen 4: Visual Traceability Matrix (`/changes/traceability`)
Visual graph demonstrating the end-to-end lineage from lab failure to new validated hardware revision.

```
+---------------------------------------------------------------------------------------------------------------------+
| [🌐] CLOSED-LOOP ENGINEERING CHANGE TRACEABILITY GRAPH                                    [Filter by Product: GW-400X]|
+---------------------------------------------------------------------------------------------------------------------+
|                                                                                                                     |
|  +--------------------+      +--------------------+      +--------------------+      +--------------------+         |
|  | Phase 5 Finding    |      | Phase 6 ECR        |      | Phase 6 CCB Appr   |      | Phase 6 ECO        |         |
|  | FND-2024-001       | ===> | ECR-2026-001       | ===> | CCB Tier 1, 2 & 3  | ===> | ECO-2026-001       |         |
|  | LDO Overheating    |      | Thermal Redesign   |      | Approved Unanimous |      | Baseline Increment |         |
|  +--------------------+      +--------------------+      +--------------------+      +---------+----------+         |
|                                                                                                |                    |
|                                                                                                v                    |
|  +--------------------+      +--------------------+      +--------------------+      +---------+----------+         |
|  | Phase 5 Validation |      | Phase 5 Regress    |      | Phase 4 Prototype  |      | Phase 2 & 3 Baseline|        |
|  | VAL-2026-002       | <=== | SESS-2026-002      | <=== | PROTO-2026-002     | <=== | Hardware Rev: REV-B0|        |
|  | Status: VALIDATED  |      | Pass Rate: 100.0%  |      | Build Qty: 5 Units |      | BOM: BOM-REV-B0    |         |
|  +--------------------+      +--------------------+      +--------------------+      +--------------------+         |
|                                                                                                                     |
+---------------------------------------------------------------------------------------------------------------------+
```

---

## 6. Frontend State Architecture & Pinia Store Design

To maintain strict isolation during this design phase, structured mock stores are proposed:

### 6.1 `src/stores/changeStore.js` (Proposed Pinia Store Structure)

```javascript
import { defineStore } from 'pinia';

export const useChangeStore = defineStore('change', {
  state: () => ({
    // Dashboard & Metrics
    dashboardMetrics: {
      totalEcrs: 18,
      openEcrs: 12,
      pendingReviews: 5,
      pendingApprovals: 3,
      activeEcos: 4,
      regressionQueueCount: 2,
      closedChanges: 34,
      averageCycleTimeDays: 4.8
    },

    // ECR Registry & Active Workspace
    ecrs: [],
    currentEcr: null,
    ecrFilters: {
      status: 'ALL',
      priority: 'ALL',
      productId: 'ALL',
      origin: 'ALL',
      search: ''
    },

    // ECO Registry & Active Order
    ecos: [],
    currentEco: null,
    ecoFilters: {
      status: 'ALL',
      productId: 'ALL',
      search: ''
    },

    // Lineage & Traceability Map
    traceabilityChains: [],

    loading: false,
    error: null
  }),

  getters: {
    criticalEcrs: (state) => state.ecrs.filter(e => e.priority === 'CRITICAL'),
    pendingReviewEcrs: (state) => state.ecrs.filter(e => e.status === 'UNDER_REVIEW' || e.status === 'IMPACT_ANALYSIS'),
    pendingApprovalEcrs: (state) => state.ecrs.filter(e => e.status === 'PENDING_APPROVAL'),
    activeEcosList: (state) => state.ecos.filter(e => e.status === 'IN_IMPLEMENTATION' || e.status === 'APPROVED'),
    getEcrById: (state) => (id) => state.ecrs.find(e => e.id === id || e.code === id),
    getEcoById: (state) => (id) => state.ecos.find(e => e.id === id || e.code === id)
  },

  actions: {
    // ECR Lifecycle Actions (Mock Design Interfaces)
    async fetchDashboardStats() { /* Mock load */ },
    async fetchEcrs(filters) { /* Mock load */ },
    async fetchEcrById(id) { /* Mock load */ },
    async createEcr(payload) { /* Design Mock */ },
    async submitEcrForReview(id) { /* Design Mock */ },
    async submitReviewFeedback(ecrId, reviewPayload) { /* Design Mock */ },
    async submitApprovalDecision(ecrId, approvalPayload) { /* Design Mock */ },
    
    // ECO Lifecycle Actions (Mock Design Interfaces)
    async convertEcrToEco(ecrId, ecoPayload) { /* Design Mock */ },
    async fetchEcos(filters) { /* Mock load */ },
    async fetchEcoById(id) { /* Mock load */ },
    async executeEcoBaseline(ecoId) { /* Design Mock */ },
    async markVerificationComplete(ecoId, testSessionId) { /* Design Mock */ }
  }
});
```

---

## 7. Reusable Vue Component Inventory

| Component Name | File Path Proposal | Functional Purpose |
| :--- | :--- | :--- |
| `ChangeStatusBadge.vue` | `src/components/changes/ChangeStatusBadge.vue` | Color-coded status chip with pulsed indicators for `IN_IMPLEMENTATION` and `PENDING_APPROVAL`. |
| `ChangePriorityBadge.vue` | `src/components/changes/ChangePriorityBadge.vue` | Visual urgency indicator (`CRITICAL`: Red, `HIGH`: Amber, `MEDIUM`: Blue, `LOW`: Slate). |
| `ChangeImpactMatrix.vue` | `src/components/changes/ChangeImpactMatrix.vue` | Interactive 6-domain technical & financial impact grid with risk ratings. |
| `ChangeItemCard.vue` | `src/components/changes/ChangeItemCard.vue` | Before-and-after delta card showing Current vs. Proposed states, schematics, and designators. |
| `BOMChangeDiffTable.vue` | `src/components/changes/BOMChangeDiffTable.vue` | Side-by-side BOM line item diff showing `ADDED`, `REMOVED`, `MODIFIED`, and cost variance. |
| `ApprovalFlowTimeline.vue` | `src/components/changes/ApprovalFlowTimeline.vue` | Multi-tier CCB approval chain stepper with sign-off signatures and justification popovers. |
| `TraceabilityGraph.vue` | `src/components/changes/TraceabilityGraph.vue` | Visual DAG lineage connecting Finding $\to$ ECR $\to$ ECO $\to$ Revision $\to$ Prototype $\to$ Test. |
| `RegressionBanner.vue` | `src/components/changes/RegressionBanner.vue` | Callout alerting engineering teams to mandatory Phase 5 re-testing requirements. |

---

## 8. Integration Points with Existing Phases

### 8.1 Phase 5 Defect Ingestion
- In `EngineeringFindingsView.vue`, when a finding has `change_candidate_status == 'CANDIDATE'`, an action button `[Create ECR]` is displayed.
- Clicking `[Create ECR]` opens the ECR Authoring Wizard pre-populated with finding title, description, category, affected prototype build, and hardware revision ID.

### 8.2 Phase 2 Hardware Revision Integration
- When an ECO is executed, the system initiates the creation of `HardwareRevision` (e.g. `REV-B0`) linked to parent `ProductVersion`.
- The revision detail view displays an **Engineering Changes** tab linking to source `ECO-2026-001`.

### 8.3 Phase 3 Engineering BOM Integration
- The target revision's BOM (`BOM-REV-B0`) is cloned from the source revision's BOM (`BOM-REV-A0`), and the Change Items are applied to the new draft BOM.
- Historical `BOM-REV-A0` remains permanently locked.

### 8.4 Phase 4 Prototype Build Integration
- New prototype builds (`PROTO-2026-002`) are created referencing the new `REV-B0`.
- The prototype build detail displays an **ECO Provenance** tag.

---

## 9. Non-Scope Boundaries (Strict Exclusions)

To maintain architectural focus on R&D product engineering control, the following domains are strictly **OUT OF SCOPE** for Phase 6:
1. **No ERP / Procurement Management:** Zero purchase orders, vendor bidding, RFQs, or invoice reconciliation.
2. **No Warehouse / Physical Inventory:** Zero bin locations, stock transfers, or physical inventory valuations.
3. **No Automated CAD / EDA Layout Editing:** The system manages engineering metadata, change orders, and BOMs; it does not parse or rewrite binary Altium / KiCad PCB layout files.
4. **No Factory Floor Execution / MES:** Zero SMT feeder automation or pick-and-place machine programming.

---

## 10. Architecture Verdict & Next Steps

$$\mathbf{PHASE\ 6\ UI/UX\ \&\ FRONTEND\ ARCHITECTURE\ VERDICT:}\quad \mathbf{GO}$$

### Key Deliverables Completed:
- [x] Clear architectural separation between ECR and ECO.
- [x] Absolute historical baseline immutability strategy defined.
- [x] Decoupled ingestion from Phase 5 engineering findings (`CANDIDATE`).
- [x] Reusable change item delta schema (Current vs Proposed).
- [x] Multi-domain impact matrix and risk analysis framework.
- [x] Multi-tier CCB review and approval governance.
- [x] BOM change preview diff engine.
- [x] Closed-loop traceability chain connecting all 6 lifecycle phases.
- [x] Zero production code modified during this design session.

Phase 6 is fully designed and ready for the formal **Implementation Plan & Review**.
