# PHASE 5 ENTITY RELATIONSHIP & DATA FLOW DIAGRAM
## IoT Product R&D Control Center
**Module**: Phase 5 — Testing & Validation Control  
**Document Type**: Entity Relationship Specification & Flow Architecture  
**Status**: ARCHITECTURAL SPECIFICATION  
**Date**: August 2026  

---

## 1. End-to-End System Hierarchy & Flow

```
                      ┌─────────────────────────────────┐
                      │        R&D Product Master       │
                      │          (ProductType)          │
                      └────────────────┬────────────────┘
                                       │ 1
                                       │ N
                      ┌────────────────▼────────────────┐
                      │         Product Version         │
                      │      (v1.0.0, v2.0.0 Major)     │
                      └────────────────┬────────────────┘
                                       │ 1
                                       │ N
                      ┌────────────────▼────────────────┐
                      │        Hardware Revision        │
                      │     (REV-A, REV-B Stepping)     │
                      └───────┬───────────────────┬─────┘
                              │ 1                 │ 1
                              │ 1                 │ N
              ┌───────────────▼──┐       ┌────────▼────────────────┐
              │  Engineering BOM │       │     Prototype Build     │
              │  (Cost Baseline) │       │   (Physical Build Run)  │
              └───────────────┬──┘       └───────┬─────────────────┘
                              │                  │
                              │                  │
                [Freeze BOM Snapshot]            │ 1
                              │                  │
                              └─────────►┌───────▼─────────────────┐
                                         │   Frozen BOM Snapshot   │
                                         │  (Immutable Part Qty)   │
                                         └───────┬─────────────────┘
                                                 │
                                                 │ 1
                                                 │ N
┌─────────────────────────────────┐      ┌───────▼─────────────────┐
│        Master Test Plan         │      │  Prototype Test Session │
│    (Reusable Protocol v1.0)     │      │   (Test Run Assignment) │
└────────────────┬────────────────┘      └───────┬─────────────────┘
                 │                               │
                 │ 1                             │
                 │ N               [Freeze Test Plan Snapshot]
                 │                               │
┌────────────────▼────────────────┐              │
│       Master Test Cases         ├──────────────┘
│  (Min/Max Specs & Procedures)   │
└────────────────┬────────────────┘
                 │
                 │ 1
                 │ N
┌────────────────▼────────────────────────────────────────────────┐
│                     Test Execution Runs                         │
│             (RunNumber: 1, 2, 3... State Machine)               │
└────────┬────────────────────────────────┬───────────────────────┘
         │ 1                              │ 1
         │ N                              │ N
┌────────▼────────────────┐      ┌────────▼────────────────┐
│    Test Measurements    │      │     Test Evidences      │
│  (Tolerance Auto-Check) │      │ (Oscilloscope, Thermal) │
└────────┬────────────────┘      └─────────────────────────┘
         │
         │ (If Failed / Anomaly)
         │ 0..1
┌────────▼────────────────────────────────────────────────┐
│                  Engineering Finding                    │
│    (Defect Severity: CRITICAL / HIGH / MED / LOW)       │
└────────┬────────────────────────────────┬───────────────┘
         │                                │
         │                                │ (If Requires PCB Spin)
         │ 1                              ▼
         │ 1                    ┌──────────────────────────────────┐
┌────────▼────────────────┐     │      is_eco_candidate = true     │
│   Validation Decision   │     │  status = 'DEFERRED_TO_ECO'      │
│ (VALIDATED / REJECTED)  │     └─────────────────┬────────────────┘
└─────────────────────────┘                       │
                                                  │ [STRICT PHASE 5 BOUNDARY]
                                                  ▼
                                ════════════════════════════════════
                                 FUTURE PHASE 6: ECO / ECR WORKFLOW
                                ════════════════════════════════════
                                 Engineering Change Request (ECR)
                                                ↓
                                 Engineering Change Order (ECO)
                                                ↓
                                    New Hardware Revision (REV-B)
```

---

## 2. Relational Cardinality & Foreign Key Schema

| Parent Table | Child Table | Relationship | Foreign Key | Cascade Policy | Description |
| :--- | :--- | :---: | :--- | :--- | :--- |
| `products` | `test_plans` | $1:N$ | `test_plans.product_id` | `SET NULL` | Master test plans optionally targeted to a specific product family |
| `hardware_revisions` | `test_plans` | $1:N$ | `test_plans.target_revision_id` | `SET NULL` | Master test plans targeted to specific PCB hardware revision |
| `test_plans` | `test_cases` | $1:N$ | `test_cases.test_plan_id` | `CASCADE` | Reusable test cases defining tolerances and procedures |
| `prototype_builds` | `prototype_test_sessions` | $1:N$ | `prototype_test_sessions.prototype_build_id` | `CASCADE` | Active test plan sessions executed on a physical prototype build |
| `test_plans` | `prototype_test_sessions` | $1:N$ | `prototype_test_sessions.test_plan_id` | `RESTRICT` | Master template reference |
| `prototype_test_sessions` | `test_executions` | $1:N$ | `test_executions.test_session_id` | `CASCADE` | Individual test case execution runs |
| `test_cases` | `test_executions` | $1:N$ | `test_executions.test_case_id` | `RESTRICT` | Reference to original test case specification |
| `test_executions` | `test_measurements` | $1:N$ | `test_measurements.test_execution_id` | `CASCADE` | Numeric readings and tolerance compliance records |
| `test_executions` | `test_evidences` | $1:N$ | `test_evidences.test_execution_id` | `CASCADE` | Waveform images, thermal photos, and logs |
| `prototype_builds` | `engineering_findings` | $1:N$ | `engineering_findings.prototype_build_id` | `CASCADE` | Cross-build defect registry |
| `test_executions` | `engineering_findings` | $1:0..1$ | `engineering_findings.test_execution_id` | `SET NULL` | Direct linkage to the failed test execution |
| `prototype_builds` | `prototype_validation_decisions` | $1:1$ | `prototype_validation_decisions.prototype_build_id` | `CASCADE` | Final formal validation approval record |

---

## 3. Data Immutability & Snapshot Flow

```
MASTER PROTOCOL LEVEL                          PROTOTYPE EXECUTION LEVEL
(Mutable by R&D Lead)                         (Permanently Immutable Record)

┌─────────────────────────┐
│     Master Test Plan    │
│  Code: TP-GW400X-EVT    │
│  Version: v1.0          │
│  Items: 12 Test Cases   │
└────────────┬────────────┘
             │
             │ [Assign Protocol to Prototype Run]
             │
             ▼
┌─────────────────────────┐                   ┌───────────────────────────────────┐
│ Prototype Test Session  │──────────────────►│     Frozen Test Plan Snapshot     │
│ ID: SES-PROTO-2024-001  │                   │ - Serialized JSON Payload         │
│ Status: IN_PROGRESS     │                   │ - 12 Test Cases with Min/Max Specs│
└────────────┬────────────┘                   │ - Nominal Targets & Equipment     │
             │                                └───────────────────────────────────┘
             │                                                  │
             ▼                                                  ▼
┌─────────────────────────┐                   ┌───────────────────────────────────┐
│     Test Executions     │                   │   Historical Traceability Lock    │
│ Run #1: FAILED (5.32V)  │                   │ If Master Test Plan is updated to │
│ Run #2: PASSED (5.08V)  │                   │ v1.1, SES-PROTO-2024-001 remains  │
└─────────────────────────┘                   │ locked to original v1.0 baseline! │
                                              └───────────────────────────────────┘
```

---

## 4. Phase 6 Integration Interface Point

```
┌────────────────────────────────────────────────────────────────────────┐
│                        PHASE 5 FINDING RECORD                          │
├────────────────────────────────────────────────────────────────────────┤
│ ID                 : FND-2024-004                                      │
│ Prototype Build ID : PROTO-2024-001 (GW-400X Alpha Build #1)           │
│ Hardware Revision  : REV-A (Initial PCB Layout)                        │
│ Category           : THERMAL_DEFECT                                    │
│ Title              : Step-Down Inductor L1 Exceeds 85°C at Max Load    │
│ Severity           : HIGH                                              │
│ Status             : DEFERRED_TO_ECO                                   │
│ Proposed Fix       : Increase L1 footprint to 10x10mm 4.7uH 6A inductor│
│ is_eco_candidate   : true                                              │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    │ [Interface Query: is_eco_candidate = true]
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                     FUTURE PHASE 6: ECO ENGINE                         │
├────────────────────────────────────────────────────────────────────────┤
│ 1. ECR (Engineering Change Request) imports open findings.             │
│ 2. Impact analysis calculates PCB tooling & scrap costs.               │
│ 3. Multi-department sign-off approves change.                          │
│ 4. System increments Hardware Revision to REV-B.                       │
│ 5. Prototype #003 instantiated for REV-B with updated BOM.             │
│ 6. Regression test session verifies L1 thermal fix.                    │
└────────────────────────────────────────────────────────────────────────┘
```
