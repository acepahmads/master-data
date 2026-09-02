# PHASE 5 ANALYSIS & SYSTEM ARCHITECTURE
## IoT Product R&D Control Center
**Module**: Phase 5 — Testing & Validation Control  
**Document Type**: Technical Architecture & Design Specification  
**Status**: ARCHITECTURAL PROPOSAL (STRICT ANALYSIS MODE)  
**Date**: August 2026  

---

## 1. Executive Summary

Phase 5 introduces **Testing & Validation Control** to the **IoT Product R&D Control Center**. Following the completion and forensic verification of Phases 1–4, Phase 5 establishes the formal engineering test protocol framework, parametric measurement engine, evidence repository, failure defect tracking, and validation governance for physical prototype build runs.

The architecture strictly adheres to the core engineering hierarchy:

$$\text{R\&D Product Master} \longrightarrow \text{Product Version} \longrightarrow \text{Hardware Revision} \longrightarrow \text{Engineering BOM} \longrightarrow \text{Prototype Build} \longrightarrow \mathbf{Testing\ \&\ Validation}$$

Phase 5 preserves historical data integrity through **Immutable Test Plan Snapshots**, server-side tolerance evaluation, granular RBAC, and seamless bridge mechanisms for Phase 6 Engineering Change Orders (ECO).

---

## 2. Current System Context & Inheritance

The system currently manages three distinct Product Types:
1. **TRADING Products**: Commercial sourcing items with supplier pricing, MSRP, and live FX conversions. No engineering versions, BOMs, or prototype builds.
2. **PROJECT Products**: Turnkey system compositions that aggregate Trading items, R&D modules, and direct discrete parts into solution cost rollups.
3. **R&D Products**: Fully engineered IoT devices following the formal 3-tier versioning and revision lifecycle.

### Phase 4 Prototype Build Baseline
Each Prototype Build (`model.PrototypeBuild`) represents a distinct physical manufacturing run instantiated against a specific `HardwareRevision`. It encapsulates:
- Compound uniqueness: `(hardware_revision_id, build_number, deleted_at)`.
- Immutable Frozen Engineering BOM Snapshot (`bom_snapshot_json`, `bom_snapshot_cost`).
- Per-component kitting and readiness tracking ($qty\_per\_board \times build\_units$).
- 8-stage physical assembly checklist (PCB Prep $\rightarrow$ Component Kitting $\rightarrow$ SMT $\rightarrow$ Through-Hole $\rightarrow$ Mechanical $\rightarrow$ Wiring $\rightarrow$ Flashing $\rightarrow$ Initial Power-On).
- Engineering observation notes and activity audit history.

Phase 5 seamlessly extends the Prototype Build lifecycle once assembly is finalized (`Status == 'READY_FOR_TESTING'` or `Status == 'POWER_ON'`).

---

## 3. Core Architectural Decisions

### 3.1 Decision 1: Test Plan Ownership & Versioning Model
- **Evaluation of Models**:
  - *Revision-Direct Ownership*: Fails test reusability (e.g. EVT electrical tests must be redefined from scratch for REV-B).
  - *Build-Direct Ownership*: High data duplication, zero standardization across builds.
  - *Global Templates Only*: Fails revision-specific parametric tolerances (e.g. 5V rail vs 3.3V rail).
- **Selected Architecture: 2-Tier Master Protocol Template + Frozen Test Run Snapshot**:
  1. **Master Test Plan Templates** (`test_plans`): Reusable test protocol specifications created by R&D leads, version-controlled (`v1.0`, `v1.1`), and targeted at specific product families or hardware revisions.
  2. **Prototype Test Sessions** (`prototype_test_sessions`): Instantiated when testing a physical build run (`PROTO-2024-001`). Freezes the test protocol into an **Immutable Test Plan Snapshot** (`test_plan_snapshot_json`).
  3. **Immutability Guarantee**: Subsequent edits or version bumps to the master test plan template will **never alter historical test executions** for already tested prototype builds.

### 3.2 Decision 2: Multi-Type Test Case Evaluation Engine
Test Cases (`test_cases`) support heterogeneous engineering verification paradigms:
1. **NUMERIC_RANGE**: Evaluates continuous values against upper and lower bounds ($Min \le Measured \le Max \implies PASSED$).
2. **NUMERIC_MIN**: Lower-bound threshold testing ($Measured \ge Min \implies PASSED$, e.g., dielectric breakdown voltage $\ge 1500\text{V}$).
3. **NUMERIC_MAX**: Upper-bound limit testing ($Measured \le Max \implies PASSED$, e.g., deep-sleep current $\le 15\mu\text{A}$, thermal temperature $\le 70^\circ\text{C}$).
4. **BOOLEAN_PASS_FAIL**: Binary functional state verification (True/False).
5. **VISUAL_INSPECTION**: Quality / workmanship criteria (IPC-A-610 Class 2 compliance, solder bridging, fiducial alignment).
6. **TEXT_LOG_VERIFICATION**: Serial console log inspection, bootloader banner matching, MAC address allocation check.

### 3.3 Decision 3: Standard Engineering Test Categories
Default system-managed categories:
- `FUNCTIONAL`: Core MCU boot, peripheral activation, UI LEDs, button triggers.
- `ELECTRICAL`: Power rail regulation (3.3V, 5.0V, 1.8V), DC-DC efficiency, switching noise, inrush current.
- `POWER_CONSUMPTION`: Active run current, transmit current, low-power sleep current, battery life verification.
- `RF_WIRELESS`: Tx power output, receiver sensitivity, harmonic suppression, antenna return loss (VSWR).
- `SENSOR_CALIBRATION`: Analog front-end linearity, ADC accuracy, temperature/humidity sensor offset.
- `COMMUNICATIONS`: RS-485 transceiver timing, CAN bus arbitration, Ethernet PHY link speed, LoRa packet error rate (PER).
- `FIRMWARE_INTEGRATION`: JTAG/SWD flashing, factory bootloader integrity, OTA update cycle, watchdog recovery.
- `THERMAL`: Thermal camera hot-spot mapping, MCU junction temperature under 100% compute load.
- `MECHANICAL`: Enclosure fit, IP67 gasket compression, vibration tolerance, connector mating retention force.
- `ENVIRONMENTAL_RELIABILITY`: Operating temperature chamber stress ($-40^\circ\text{C}$ to $+85^\circ\text{C}$), humidity soak.
- `SAFETY_COMPLIANCE`: Isolation breakdown, ESD immunity (IEC 61000-4-2), reverse polarity protection.
- `REGRESSION`: Targeted verification after physical board rework to prevent secondary defects.

### 3.4 Decision 4: Test Execution State Machine & Multi-Run Iterations
- Execution States:
  $$\text{NOT\_STARTED} \longrightarrow \text{IN\_PROGRESS} \longrightarrow \begin{cases} \text{PASSED} \\ \text{FAILED} \\ \text{BLOCKED} \\ \text{SKIPPED} \end{cases}$$
- **Multi-Run History (`RunNumber`)**:
  - In hardware R&D, a prototype that fails Run #1 is often reworked on the lab bench (e.g. swapping a pull-up resistor) and re-tested in Run #2.
  - The system records all execution iterations (`RunNumber: 1, 2, 3...`) with individual measurements, timestamps, and technician signatures.
  - The latest execution run defines the active state of the test case, while all previous runs remain permanently preserved in the historical audit trail.

### 3.5 Decision 5: Test Evidence Repository
- Supports multi-format evidence attachments bound to `TestExecution`:
  - Waveform captures (`.png`, `.jpg` oscilloscope traces).
  - Thermal infrared heatmaps (`.jpg`, `.png`).
  - Datalog exports (`.csv`, `.xlsx` power analyzer logs).
  - Terminal logs (`.txt`, `.log` UART console captures).
  - Inspection reports (`.pdf` external lab compliance certificates).
- **Evidence Governance**: Evidence files are write-once and become permanently locked upon formal prototype validation sign-off.

### 3.6 Decision 6: Engineering Findings (Defects) & Severity Matrix
When a test execution fails or reveals an unexpected anomaly, an **Engineering Finding** (`engineering_findings`) is created:
- **Severity Levels**:
  - `CRITICAL`: Safety hazard, component smoke/damage, bricked board. **Strictly blocks Validation Decision sign-off**.
  - `HIGH`: Major functional failure, primary interface non-operational. Requires formal engineering disposition before validation.
  - `MEDIUM`: Parametric value out of tolerance, non-fatal performance degradation.
  - `LOW`: Minor cosmetic or bench anomaly with negligible impact.
  - `INFO`: Design observation or bench note for documentation.
- **Finding Lifecycle**:
  $$\text{OPEN} \longrightarrow \text{IN\_INVESTIGATION} \longrightarrow \begin{cases} \text{FIXED\_IN\_REWORK} \\ \text{DEFERRED\_TO\_ECO} \\ \text{ACCEPTED\_RISK} \\ \text{CLOSED} \end{cases}$$

### 3.7 Decision 7: Formal Validation Decision Sign-Off
Formal engineering validation governance for a Prototype Build:
- **Validation Decisions**:
  - `PENDING_REVIEW`: Testing in progress or awaiting formal review.
  - `VALIDATED`: All test cases passed, 0 unresolved `CRITICAL` or `HIGH` findings.
  - `CONDITIONALLY_VALIDATED`: Approved with documented constraints (e.g. "Approved for bench RF tuning only; full environmental chamber testing deferred to REV-B"). Mandatory disposition justification required.
  - `REJECTED`: Prototype failed critical requirements; revision redesign or board spin required.
- **Separation of Duties**: Requires role with `testing.validate` permission (`ROLE-RNDM`, `ROLE-QA`). Testers cannot validate their own builds without lead approval.

---

## 4. UI/UX Architecture

```
                    ┌────────────────────────────────────────────────────────┐
                    │       IoT Product R&D Control Center - Phase 5         │
                    └───────────────────────────┬────────────────────────────┘
                                                │
         ┌───────────────────────┬──────────────┴────────┬──────────────────────┐
         ▼                       ▼                       ▼                      ▼
┌──────────────────┐   ┌──────────────────┐    ┌──────────────────┐   ┌──────────────────┐
│ Testing Dashboard│   │Test Plan Registry│    │Prototype Details │   │ Findings Registry│
│  (/testing)      │   │ (/test-plans)    │    │(/prototypes/:id) │   │  (/findings)     │
├──────────────────┤   ├──────────────────┤    ├──────────────────┤   ├──────────────────┤
│• Total Builds    │   │• Master Templates│    │• Overview        │   │• Open Defects    │
│• Tests Pass Rate │   │• Category Matrix │    │• BOM Snapshot    │   │• Severity Matrix │
│• Active Sessions │   │• Versioning      │    │• Component Prep  │   │• Root Causes     │
│• Open Findings   │   │• Parametric Specs│    │• Assembly Stages │   │• ECO Candidates  │
│• Validation KPIs │   │• Clone / Export  │    │• Testing & Valid*│   │• Dispositions    │
└──────────────────┘   └──────────────────┘    │• Eng Notes       │   └──────────────────┘
                                               │• Activity History│
                                               └──────────────────┘
```

### 4.1 Testing Command Center (`/testing`)
- **Top KPI Strip**: Builds Under Test, Overall Test Pass Rate %, Total Executions, Active Findings (Critical/High), Pending Validations.
- **Prototype Testing Matrix Table**: Live status, pass/fail progress bar, active findings count, and validation status per prototype build.
- **Severity Breakdown & Category Distribution**: Visual breakdown of findings and test execution results.
- **Live Test Event Stream**: Real-time log of test completions, measurements, and finding submissions.

### 4.2 Test Plan & Protocol Matrix (`/test-plans`, `/test-plans/:id`)
- Master protocol catalog with version branching (`v1.0`, `v1.1`), product family targeting, and reusable test case library.
- Parameter editor for specifying nominal targets, tolerances, equipment, and firmware requirements.

### 4.3 Prototype Detail Integration (`/prototypes/:id` -> Tab 7: `Testing & Validation`)
- Add **Tab 7: Testing & Validation** to `PrototypeDetailView.vue`.
- Displays:
  - Header: Formal Validation Decision Status Banner (`VALIDATED`, `PENDING_REVIEW`, `CONDITIONALLY_VALIDATED`, `REJECTED`) with sign-off details.
  - Test Session Selector: Associated test protocols and frozen snapshot metadata.
  - Interactive Test Execution Matrix: Grouped by category with status pills, measured values, and Pass/Fail buttons.
  - Interactive Measurement Modal: Step-by-step entry with real-time tolerance compliance indicator.
  - Defect Drawer: List of linked engineering findings with severity badges.
  - Attached Evidence Viewer: Oscilloscope screenshots, thermal captures, and serial logs.

---

## 5. Backend Architecture & Relational Schema Proposal

### 5.1 Relational Schema (`MySQL 8.0`)

```
┌─────────────────────────────────┐           ┌─────────────────────────────────┐
│           test_plans            │ 1       N │           test_cases            │
│─────────────────────────────────│───────────│─────────────────────────────────│
│ id (PK)                         │           │ id (PK)                         │
│ code (UNIQUE)                   │           │ test_plan_id (FK)               │
│ name                            │           │ code                            │
│ version (e.g. '1.0')            │           │ title                           │
│ product_id (FK, nullable)       │           │ category                        │
│ target_revision_id (FK, null)   │           │ evaluation_type                 │
│ description                     │           │ min_value, max_value, target    │
│ status ('ACTIVE', 'DRAFT')      │           │ unit_of_measure                 │
│ created_by, created_at, etc.    │           │ procedure, acceptance_criteria  │
└────────────────┬────────────────┘           │ sequence                        │
                 │ 1                          └────────────────┬────────────────┘
                 │                                             │ 1
                 │ N                                           │ N
┌────────────────▼────────────────┐           ┌────────────────▼────────────────┐
│    prototype_test_sessions      │ 1       N │         test_executions         │
│─────────────────────────────────│───────────│─────────────────────────────────│
│ id (PK)                         │           │ id (PK)                         │
│ prototype_build_id (FK)         │           │ test_session_id (FK)            │
│ test_plan_id (FK)               │           │ test_case_id (FK)               │
│ test_plan_snapshot_json (LONG)  │           │ run_number (INT, e.g. 1, 2, 3)  │
│ total_cases, passed, failed     │           │ status ('PASSED', 'FAILED',...) │
│ session_status ('IN_PROGRESS'..)│           │ executed_by, executed_at        │
│ started_at, completed_at        │           │ notes                           │
└────────────────┬────────────────┘           └────────┬───────────────┬────────┘
                 │ 1                                   │ 1             │ 1
                 │                                     │               │
                 │ 1                                   │ N             │ N
┌────────────────▼────────────────┐           ┌────────▼────────┐ ┌────▼──────────────┐
│ prototype_validation_decisions  │           │test_measurements│ │  test_evidences   │
│─────────────────────────────────│           │─────────────────│ │───────────────────│
│ id (PK)                         │           │ id (PK)         │ │ id (PK)           │
│ prototype_build_id (FK, UNIQUE) │           │ test_exec_id(FK)│ │ test_exec_id (FK) │
│ test_session_id (FK)            │           │ measured_value  │ │ file_name, path   │
│ decision ('VALIDATED', etc.)    │           │ unit            │ │ file_type, size   │
│ conditions, disposition_notes   │           │ is_in_tolerance │ │ caption           │
│ validated_by, validated_at      │           │ sample_index    │ │ uploaded_by, at   │
└─────────────────────────────────┘           └─────────────────┘ └───────────────────┘
                                                       │ 1
                                                       │ 
                                                       │ 0..1
                                              ┌────────▼────────────────┐
                                              │  engineering_findings   │
                                              │─────────────────────────│
                                              │ id (PK)                 │
                                              │ prototype_build_id (FK) │
                                              │ test_execution_id (FK)  │
                                              │ code (UNIQUE)           │
                                              │ title, category         │
                                              │ severity ('CRITICAL'..) │
                                              │ status ('OPEN', 'CLOSED')
                                              │ root_cause, proposed_fix│
                                              │ is_eco_candidate (BOOL) │
                                              │ created_by, created_at  │
                                              └─────────────────────────┘
```

---

## 6. Granular RBAC Permissions

| Permission Code | Module | Action | Description | Granted Roles |
| :--- | :--- | :--- | :--- | :--- |
| `testing.view` | `testing` | `view` | View test plans, prototype test sessions, executions, and dashboards | `ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-HWENG`, `ROLE-FWENG`, `ROLE-QA`, `ROLE-VIEW` |
| `testplans.manage` | `testing` | `managePlans` | Create, edit, clone, and version master test plan templates and test cases | `ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-HWENG`, `ROLE-QA` |
| `testing.execute` | `testing` | `execute` | Execute test cases, record parametric measurements, and upload evidence | `ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-HWENG`, `ROLE-FWENG`, `ROLE-QA` |
| `findings.manage` | `testing` | `manageFindings` | Create, update, assign, and disposition engineering findings | `ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-HWENG`, `ROLE-FWENG`, `ROLE-QA` |
| `testing.validate` | `testing` | `validate` | Formally approve, conditionally approve, or reject Prototype Validation Decisions | `ROLE-ADMIN`, `ROLE-RNDM`, `ROLE-QA` |

---

## 7. Activity Logging Event Taxonomy

Standardized event taxonomy recorded in `activity_logs`:

| Action Constant | Entity Type | Description Format | Badge Color |
| :--- | :--- | :--- | :--- |
| `TEST_PLAN_CREATED` | `TestPlan` | Created master test plan '%s' (v%s) with %d test cases | `blue` |
| `TEST_SESSION_INITIATED` | `PrototypeTestSession` | Initiated test session for prototype %s with frozen test snapshot (%d cases) | `indigo` |
| `TEST_CASE_PASSED` | `TestExecution` | Test case %s (%s) PASSED on %s (Measured: %s %s) | `emerald` |
| `TEST_CASE_FAILED` | `TestExecution` | Test case %s (%s) FAILED on %s (Measured: %s %s, expected %s to %s) | `red` |
| `MEASUREMENT_RECORDED` | `TestMeasurement` | Recorded measurement for %s: %s %s (Tolerance: %s) | `cyan` |
| `EVIDENCE_ATTACHED` | `TestEvidence` | Attached test evidence '%s' (%s) to test execution %s | `slate` |
| `ENGINEERING_FINDING_CREATED` | `EngineeringFinding` | Logged %s finding '%s' for prototype %s (Severity: %s) | `amber` |
| `ENGINEERING_FINDING_DISPOSITIONED` | `EngineeringFinding` | Finding %s dispositioned to %s: '%s' | `purple` |
| `PROTOTYPE_VALIDATED` | `ValidationDecision` | Prototype %s formally VALIDATED by %s | `emerald` |
| `PROTOTYPE_CONDITIONALLY_VALIDATED` | `ValidationDecision` | Prototype %s CONDITIONALLY VALIDATED with %d constraints | `amber` |
| `PROTOTYPE_VALIDATION_REJECTED` | `ValidationDecision` | Prototype %s validation REJECTED: %s | `red` |

---

## 8. Phase 6 ECO / ECR Integration Boundary

### Strict Boundary Definition
- **Phase 5 Scope**:
  - When an engineering finding is confirmed to require physical PCB schematic, layout, or component changes, the engineer sets:
    - `is_eco_candidate = true`
    - `status = 'DEFERRED_TO_ECO'`
    - `proposed_fix` = description of required hardware changes.
  - Phase 5 stops **strictly** at marking the finding as an ECO candidate.
- **Future Phase 6 Scope**:
  - In Phase 6, the Engineering Change Order engine will query all open findings where `is_eco_candidate == true` and link them as justification for initiating an Engineering Change Request (ECR) $\rightarrow$ Engineering Change Order (ECO).
  - Approved ECOs will trigger the creation of a new `HardwareRevision` (e.g. `REV-B`) or new `ProductVersion` (e.g. `v2.0.0`), completing the closed-loop R&D engineering cycle.

---

## 9. Architectural Risk Assessment & Mitigations

| Risk | Impact | Mitigation Strategy |
| :--- | :--- | :--- |
| **Historical Test Drift** | High | Frozen JSON test plan snapshots preserve exact test parameters at time of execution, preventing test drift when master templates are modified. |
| **Measurement Type Rigidity** | Medium | Multi-type evaluation engine (`NUMERIC_RANGE`, `NUMERIC_MIN`, `NUMERIC_MAX`, `BOOLEAN`, `VISUAL`, `TEXT`) supports all hardware engineering disciplines. |
| **Premature Validation** | High | Server-side validation rule blocks `VALIDATED` status if unresolved `CRITICAL` or `HIGH` findings exist. |
| **Evidence Loss** | Medium | Multipart storage with file checksums, locked permanently upon validation decision approval. |
| **Performance with Large Datalogs** | Low | High-frequency time-series datalogs stored as attached artifacts/files rather than individual relational rows. |

---

## 10. Final Architecture Recommendation

$$\mathbf{ARCHITECTURE\ VERDICT:\ GO}$$

The proposed Phase 5 Testing & Validation Control architecture is robust, fully normalized, preserves complete historical test immutability, enforces strict product hierarchy guards, and cleanly prepares for Phase 6 ECO governance without any scope creep.
