/**
 * Centralized Project Progress & Information Architecture Configuration
 * PRODUCT & MASTER DATA CONTROL CENTER
 * 
 * Defines the complete roadmap, platform foundation, product type hierarchy,
 * engineering lifecycle roadmap (Phases 1–6 Complete), and milestone history.
 */

export const projectMeta = {
  name: "Product & Master Data Control Center",
  subtitle: "Centralized Product, Master Data & Engineering Lifecycle Platform",
  code: "PMD-CC",
  version: "1.0.0",
  currentPhaseId: "phase-6",
  currentPhaseName: "Phase 6 — Engineering Change Management (ECR / ECO)",
  targetPhasesCount: 7, // 7 official engineering phases (Phase 1, 2, 2C, 3, 4, 5, 6)
  systemStatus: "Operational",
  apiStatus: "v1.0.0 Online (Gin + GORM + MySQL)",
  databaseEngine: "MySQL 8.0 Engine (InnoDB)",
  auditDecision: "GO — PHASES 1, 2, 2C, 3, 4, 5, 6 FULL-STACK VERIFIED & AUDITED (100%)",
  auditDate: "August 2026"
};

/**
 * PLATFORM FOUNDATION (Layer 1)
 * Centralized master data capabilities supporting all product types and workflows.
 */
export const platformFoundation = [
  {
    id: "PF-01",
    title: "Product Master",
    code: "PLT-PRD",
    scope: "Trading, Project & R&D Products",
    status: "Complete",
    progress: 100,
    color: "blue",
    icon: "product",
    description: "Centralized product portfolio registry managing all 3 product archetypes with unified taxonomy and lifecycle controls.",
    highlights: ["Unified Product Registry", "3 Product Types Classification", "Dynamic Parametric Specs", "Lifecycle State Tracking"]
  },
  {
    id: "PF-02",
    title: "Component Master",
    code: "PLT-CMP",
    scope: "Semiconductors & Discrete Parts",
    status: "Complete",
    progress: 100,
    color: "emerald",
    icon: "component",
    description: "Central electronics & mechanical component database with MPN, dynamic key-value specifications, and datasheets.",
    highlights: ["MPN & Package Standards", "Dynamic Parametric Specs", "Datasheet & Doc Uploads", "RoHS 3 Compliance"]
  },
  {
    id: "PF-03",
    title: "Supplier & Manufacturer Master",
    code: "PLT-SCM",
    scope: "Global Supply Chain Ecosystem",
    status: "Complete",
    progress: 100,
    color: "amber",
    icon: "manufacturer",
    description: "Global semiconductor fabricators directory, tier classifications, and authorized supplier distribution networks.",
    highlights: ["Fabricator Profiles & Tiers", "Authorized Distributors", "Sourcing Terms & Lead Times", "Supplier Part Numbers"]
  },
  {
    id: "PF-04",
    title: "Category & Units of Measure",
    code: "PLT-TAX",
    scope: "3-Tier Taxonomy & Standard UoM",
    status: "Complete",
    progress: 100,
    color: "purple",
    icon: "category",
    description: "3-tier hierarchical category taxonomy tree and standardized ISO/packaging engineering units of measure.",
    highlights: ["3-Tier Category Tree", "Recursive Item Rollup", "ISO Unit Conversions", "Packaging Types"]
  },
  {
    id: "PF-05",
    title: "Users & Granular RBAC Matrix",
    code: "PLT-USR",
    scope: "Engineering Staff Directory & Security",
    status: "Complete",
    progress: 100,
    color: "cyan",
    icon: "users",
    description: "Centralized user directory with granular 32-permission role-based access control across 6 engineering roles.",
    highlights: ["User Directory & Depts", "32 Granular Permissions", "6 Department Roles", "Activity Audit Logs"]
  },
  {
    id: "PF-06",
    title: "Authentication & Session Security",
    code: "PLT-SEC",
    scope: "Cross-Cutting Security Capability",
    status: "Complete",
    progress: 100,
    color: "indigo",
    icon: "lock",
    description: "Enterprise authentication portal with real JWT tokens, Bcrypt password hashing, session guards, and 401 interceptors.",
    highlights: ["Enterprise Login Portal", "Real JWT Bearer Tokens", "Vue Router Session Guard", "Axios 401 Auto-Redirect"]
  },
  {
    id: "PF-07",
    title: "Product Data Import & Migration",
    code: "PLT-DIMP",
    scope: "Cross-Cutting Data Foundation Capability",
    status: "Implementation Complete",
    progress: 95,
    color: "purple",
    icon: "import",
    description: "Multi-file Excel parser (excelize), canonical normalization, duplicate detection, AI classification co-pilot, triage review, and atomic master dispatch.",
    highlights: ["Multi-File 2026 Scope", "Golang excelize Parser", "Rule & AI Classification", "Data Provenance Lineage"]
  }
];

/**
 * ENGINEERING LIFECYCLE ROADMAP (Layer 2: Phases 1 to 6)
 * Note: Authentication is a cross-cutting platform capability; no Phase 7 exists.
 */
export const projectPhases = [
  {
    id: "phase-1",
    number: 1,
    name: "Foundation & Master Data",
    shortName: "Foundation & Masters",
    tagline: "Centralized Product and Master Data Management",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-01",
    completedDate: "2026-08-31",
    description: "Centralized Product and Master Data management. Normalizes products, components, parametric specifications, 3-tier category taxonomy, manufacturers, suppliers, UoM, users, and RBAC matrix with full REST API and MySQL persistence.",
    modules: [
      { id: "M01", name: "Product Master Management", code: "MOD-PRD", status: "completed", description: "Lifecycle tracking, lead engineers, categories, and technical profiles." },
      { id: "M02", name: "Central Components Library", code: "MOD-CMP", status: "completed", description: "Semiconductor parts, passives, MPN codes, package standards, and compliance." },
      { id: "M03", name: "Dynamic Parametric Specs", code: "MOD-SPC", status: "completed", description: "Key-value engineering attributes, tolerances, and electrical characteristics." },
      { id: "M04", name: "Component Documents & Datasheets", code: "MOD-DOC", status: "completed", description: "Multipart file uploads, PDF datasheet attachments, and compliance certificates." },
      { id: "M05", name: "3-Tier Taxonomy & Categories", code: "MOD-CAT", status: "completed", description: "Hierarchical categorization tree with recursive counts." },
      { id: "M06", name: "Manufacturers Directory", code: "MOD-MFG", status: "completed", description: "Semiconductor fabricators, official sites, and tier rankings." },
      { id: "M07", name: "Suppliers & Distributors", code: "MOD-SUP", status: "completed", description: "Authorized supply chain partners, lead times, and terms." },
      { id: "M08", name: "Units of Measure (UoM)", code: "MOD-UOM", status: "completed", description: "Standardized packaging, engineering units, and conversion factors." },
      { id: "M09", name: "User Directory & Profiles", code: "MOD-USR", status: "completed", description: "Engineering staff directory, departments, and credentials." },
      { id: "M10", name: "Role-Based Access Control (RBAC)", code: "MOD-RBC", status: "completed", description: "Granular permission matrix with 32 permissions across 6 roles." },
      { id: "M11", name: "JWT Authentication & Security", code: "MOD-AUT", status: "completed", description: "Token verification, password hashing, and session management." },
      { id: "M12", name: "Activity Logs & Audit Trail", code: "MOD-LOG", status: "completed", description: "Live change capture with actor, target entity, and timestamps." },
      { id: "M13", name: "Golang Gin + GORM REST API", code: "MOD-API", status: "completed", description: "High-performance backend integration with MySQL persistence." },
      { id: "M14", name: "Enterprise UI/UX Design System", code: "MOD-UIX", status: "completed", description: "Deep Navy command center design system and micro-interactions." }
    ]
  },
  {
    id: "phase-2",
    number: 2,
    name: "Product Versioning & Hardware Revision",
    shortName: "Product Versioning",
    tagline: "Controls Engineering Product Evolution",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Controls engineering product evolution. Implements formal hardware revision workflows (Product ➔ Product Version ➔ Hardware Revision), revision branching, schematic attachments, and version-to-version diff comparisons.",
    modules: [
      { id: "M15", name: "Product Version Tree & Branching", code: "MOD-VER", status: "completed", description: "Visual revision hierarchy supporting minor and major hardware branches." },
      { id: "M16", name: "Hardware Revision Matrix (A0, B0, C0)", code: "MOD-REV", status: "completed", description: "PCB revision designations, silicon stepping, and assembly codes." },
      { id: "M17", name: "Revision History & Engineering Changelog", code: "MOD-CHG", status: "completed", description: "Historical log of schematic and component modifications." },
      { id: "M18", name: "Version Lifecycle State Machine", code: "MOD-LCS", status: "completed", description: "Formal states: Draft, In Review, Active, Deprecated, Archived." },
      { id: "M19", name: "Version-to-Version Diff Comparison", code: "MOD-DIF", status: "completed", description: "Side-by-side parametric and baseline delta comparison." },
      { id: "M20", name: "Engineering History & Audit Hook", code: "MOD-ECO-H", status: "completed", description: "Activity logging and persistence foundation for audit governance." }
    ]
  },
  {
    id: "phase-2c",
    number: "2C",
    name: "Product Type & Project Architecture",
    shortName: "Product Type Architecture",
    tagline: "Separates Trading, Project and R&D Product Behavior",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Separates Trading, Project and R&D product behavior. Establishes commercial sourcing terms for Trading products, turnkey solution trees for Project products, and multi-currency real-time exchange rates with product photo galleries.",
    modules: [
      { id: "M21", name: "3-Tier Product Type Taxonomy", code: "MOD-TYP", status: "completed", description: "Architectural classification into TRADING, PROJECT, and R&D products." },
      { id: "M22", name: "Trading Commercial Sourcing Terms", code: "MOD-TRD", status: "completed", description: "Purchase prices, MSRP, MPN, supplier lead times, and warranties." },
      { id: "M23", name: "Project Solution Composition Tree", code: "MOD-PRJ", status: "completed", description: "Hierarchical nesting of Trading items, R&D subsystems, and direct parts." },
      { id: "M24", name: "Multi-Currency Pricing & Live FX Engine", code: "MOD-CUR", status: "completed", description: "Real-time exchange rates (open.er-api.com) + custom manual rate input." },
      { id: "M25", name: "Multi-Image Product Gallery (Max 5)", code: "MOD-GAL", status: "completed", description: "Multi-image upload, drag-and-drop gallery, cover photo selection, and previewer." }
    ]
  },
  {
    id: "phase-3",
    number: 3,
    name: "Engineering BOM & Cost Foundation",
    shortName: "Engineering BOM & Cost",
    tagline: "Defines Multi-Level Engineering BOM and Cost Structure",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Defines multi-level engineering BOM and cost structure. Supports hierarchical sub-assemblies, reference designator mapping, bottom-up recursive cost aggregation, target profit margin markup, and project solution cost rollups.",
    modules: [
      { id: "M26", name: "Multi-Level Engineering BOM Workspace", code: "MOD-BOM-W", status: "completed", description: "Tree-structured BOM items with sub-assemblies, PCBA, and parts breakdown." },
      { id: "M27", name: "Reference Designator Mapper", code: "MOD-REF-D", status: "completed", description: "Designator assignments, duplicate validation, and component count checks." },
      { id: "M28", name: "Recursive Cost Calculation Engine", code: "MOD-CST-R", status: "completed", description: "Bottom-up cost aggregation from component unit costs up to top-level revision." },
      { id: "M29", name: "Profit Margin & Selling Price Derivation", code: "MOD-MRG", status: "completed", description: "Margin percentage markup, target selling price, and dual-currency preview." },
      { id: "M30", name: "Project Solution Cost Aggregator", code: "MOD-PRJ-C", status: "completed", description: "Aggregated rollups of trading hardware, R&D modules, and direct parts." },
      { id: "M31", name: "BOM RBAC & Audit Governance", code: "MOD-BOM-A", status: "completed", description: "Role permissions and live audit activity logging for all BOM modifications." }
    ]
  },
  {
    id: "phase-4",
    number: 4,
    name: "Prototype & Engineering Build Management",
    shortName: "Prototype Builds",
    tagline: "Controls Physical Prototype Builds from Engineering Baselines",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Controls physical prototype builds from engineering baselines. Manages physical build units (#001, #002) tied to Hardware Revisions, frozen immutable BOM snapshots, component preparation tracking, and 8-stage assembly workflow tracking.",
    modules: [
      { id: "M32", name: "Prototype Build Command Center", code: "MOD-PRT-D", status: "completed", description: "KPI metrics, active builds workspace, status distribution, and build stream." },
      { id: "M33", name: "4-Step Prototype Creation Wizard", code: "MOD-PRT-W", status: "completed", description: "Revision selection, build metadata, frozen BOM snapshot, and objectives." },
      { id: "M34", name: "Frozen BOM Snapshot Engine", code: "MOD-BOM-S", status: "completed", description: "Read-only frozen component and cost baseline preserved per prototype build." },
      { id: "M35", name: "Component Preparation & Readiness", code: "MOD-CMP-R", status: "completed", description: "Readiness indicators: Available, Partial, Missing, and Substitute Required." },
      { id: "M36", name: "Multi-Stage Assembly Checklist", code: "MOD-ASM-C", status: "completed", description: "8-stage assembly workflow from PCB prep to power-on and flashing." }
    ]
  },
  {
    id: "phase-5",
    number: 5,
    name: "Testing & Validation Control",
    shortName: "Testing & Validation",
    tagline: "Controls Formal Engineering Testing and Validation",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Controls formal engineering testing and validation. Implements reusable Master Test Plan protocol versioning (DRAFT/RELEASED), frozen session snapshots, multi-run parametric bench execution, engineering defect findings hub, and formal validation sign-off governance.",
    modules: [
      { id: "M37", name: "Master Test Plan & Protocol Versioning", code: "MOD-TST-P", status: "completed", description: "Reusable test protocols, version cloning, and immutable release locking." },
      { id: "M38", name: "Frozen Immutable Test Snapshot Engine", code: "MOD-TST-S", status: "completed", description: "JSON snapshot frozen upon session inception, isolating executions from template edits." },
      { id: "M39", name: "Multi-Run Parametric Execution & Bench", code: "MOD-TST-R", status: "completed", description: "Append-only iterations (Run 1, 2...), server-side tolerance evaluation, and evidence attachments." },
      { id: "M40", name: "Engineering Findings & Defect Hub", code: "MOD-FND-D", status: "completed", description: "Defect logging, severity tiers, root causes, dispositions, and ECO change candidate tagging." },
      { id: "M41", name: "Validation Sign-Off Governance Engine", code: "MOD-VAL-G", status: "completed", description: "Formal validation verdicts with unresolved critical/high finding blocking rules." }
    ]
  },
  {
    id: "phase-6",
    number: 6,
    name: "Engineering Change Management (ECR / ECO)",
    shortName: "Engineering Change",
    tagline: "Controls ECR/ECO and Creation of New Engineering Baselines",
    status: "completed",
    badge: "COMPLETE",
    progress: 100,
    startDate: "2026-08-31",
    completedDate: "2026-08-31",
    description: "Controls ECR/ECO and creation of new engineering baselines. Governs 9 change origins, 5-domain impact analysis, CCB multi-tier sign-offs, atomic target revision spawning (REV-B0), BOM clone & delta transformations, and closed-loop traceability DAG.",
    modules: [
      { id: "M42", name: "ECR / ECO Workflow Engine & Guardrails", code: "MOD-ECO-E", status: "completed", description: "State-driven change lifecycle from proposal to release with strict PLM guardrails." },
      { id: "M43", name: "5-Domain Impact & Cost Delta Matrix", code: "MOD-ECO-I", status: "completed", description: "Electrical, thermal, mechanical, firmware, and scrap cost impact evaluations." },
      { id: "M44", name: "CCB Multi-Tier Approval Chain", code: "MOD-SIG", status: "completed", description: "Tier 1 Tech Lead, Tier 2 Eng Manager, and Tier 3 Quality Director approvals." },
      { id: "M45", name: "Atomic Baseline Spawning & BOM Diff Engine", code: "MOD-ECO-S", status: "completed", description: "Atomic target revision generation, BOM cloning, delta application, and diff engine." },
      { id: "M46", name: "Closed-Loop Traceability Directed Acyclic Graph", code: "MOD-TRC-G", status: "completed", description: "Full lineage linking Defect ➔ ECR ➔ CCB ➔ ECO ➔ Target Rev ➔ Target BOM ➔ Verification." }
    ]
  }
];

/**
 * PRODUCT ARCHITECTURE MODEL
 * Detailed definition of the 3 Product Types and their distinct lifecycle characteristics.
 */
export const systemProductArchitecture = {
  title: "PRODUCT MASTER ARCHITECTURE",
  subtitle: "Unified 3-Tier Classification & Engineering Lifecycle Hierarchy",
  types: [
    {
      id: "trading",
      name: "TRADING PRODUCT",
      code: "TRADING",
      badgeColor: "purple",
      icon: "shopping-bag",
      purpose: "Commercial / Sourced Hardware",
      description: "Ready-made commercial products sourced from external suppliers, distributors, or manufacturers for resale or project integration.",
      characteristics: [
        "Ready-made commercial hardware from external vendors",
        "Commercial sourcing terms (Purchase price, Supplier, MPN)",
        "Derived selling price (Target margin %) & Live multi-currency FX",
        "Supplier warranty, lead time, and multi-image photo gallery",
        "DOES NOT have internal engineering versions, revisions, BOM, prototypes, testing, or ECR/ECO"
      ],
      tree: [
        { label: "Commercial Sourcing Parameters", sub: "Supplier MPN, Purchase Price, Currency, Lead Time, Warranty" },
        { label: "Derived Commercial Selling Price", sub: "Profit Margin % calculation with live multi-currency exchange rates" },
        { label: "Visual Asset Gallery", sub: "Multi-image product photo gallery with cover photo selection" }
      ]
    },
    {
      id: "project",
      name: "PROJECT PRODUCT",
      code: "PROJECT",
      badgeColor: "cyan",
      icon: "layers",
      purpose: "Turnkey / Integrated Solution",
      description: "Turnkey customer solution compositions combining commercial trading products, internal R&D modules, and direct parts.",
      characteristics: [
        "Hierarchical solution composition tree",
        "References Commercial Trading Products with sourcing price rollups",
        "References Internal R&D Subsystems (without duplicating internal engineering data)",
        "References Direct Sourced Components (installation materials, cables, accessories)",
        "Solution-level aggregated cost rollup and target solution selling price"
      ],
      tree: [
        { label: "Commercial Trading Product Line Items", sub: "External commercial items with sourcing price rollup" },
        { label: "Internal R&D Sub-Assembly Modules", sub: "Engineered modules referencing R&D BOM cost aggregations" },
        { label: "Direct Sourced Installation Components", sub: "Direct accessories, connectors, cables, mounting hardware" },
        { label: "Project Cost & Margin Rollup", sub: "Aggregated project solution cost and client quote derivation" }
      ]
    },
    {
      id: "rnd",
      name: "R&D PRODUCT",
      code: "RND",
      badgeColor: "emerald",
      icon: "cpu",
      purpose: "In-House Engineered Hardware",
      description: "Deepest engineering lifecycle for in-house designed hardware, from product version branching down to BOMs, prototypes, tests, and change orders.",
      characteristics: [
        "Formal 3-Tier Hardware Engineering Lifecycle (Product ➔ Product Version ➔ Hardware Revision)",
        "Multi-Level Engineering BOM Structure with Reference Designators & SMT Mapping",
        "Bottom-up recursive BOM cost rollup calculation & target profit margin markup",
        "Physical Prototype Build Runs (#001, #002) with frozen immutable BOM snapshots",
        "Master Test Protocols (DRAFT/RELEASED), multi-run parametric execution bench, and findings hub",
        "Formal Engineering Change Management (ECR / CCB Approvals / ECO / Atomic Target Revision Generation)"
      ],
      tree: [
        { label: "Product Version (e.g. v1.0, v2.0)", sub: "Major architecture releases, functional revisions, and lifecycle states" },
        { label: "Hardware Revision (e.g. REV-A, REV-B)", sub: "Physical PCB stepping, schematic attachments, and silicon stepping" },
        { label: "Engineering BOM Workspace", sub: "Multi-level sub-assemblies, line items, ref designators & recursive cost rollup" },
        { label: "Prototype Build Runs (#001, #002)", sub: "Physical build units with frozen immutable BOM snapshots & assembly tracking" },
        { label: "Testing & Validation Bench", sub: "Frozen test snapshots, multi-run parametric execution & defect findings" },
        { label: "Engineering Change (ECR / ECO)", sub: "CCB multi-tier sign-off, atomic target revision spawning & regression mandate" }
      ]
    }
  ]
};

export const systemArchitecturePipeline = [
  {
    layer: "Frontend Layer",
    tech: "Vue.js 2.7 + Tailwind CSS + Pinia",
    description: "High-density engineering workstation UI, reactive store state, modal system, and responsive layout.",
    status: "Operational",
    statusType: "healthy",
    icon: "terminal",
    metrics: "181 modules • <10ms render"
  },
  {
    layer: "API Gateway Layer",
    tech: "RESTful JSON API (Golang Gin)",
    description: "High-throughput REST routes with CORS, JSON binding, JWT middleware, and role enforcement.",
    status: "Operational",
    statusType: "healthy",
    icon: "server",
    metrics: ":8080 • ~1.2ms avg latency"
  },
  {
    layer: "Data Access Layer (ORM)",
    tech: "GORM v1.25 Engine",
    description: "Type-safe database abstraction, model migrations, preloading relations, and transaction safety.",
    status: "Operational",
    statusType: "healthy",
    icon: "layers",
    metrics: "Auto-Migrate • Prepared Stmts"
  },
  {
    layer: "Relational Persistence",
    tech: "MySQL 8.0 Engine (InnoDB)",
    description: "ACID-compliant storage, foreign key constraints, UTF8MB4 charset, and indexed master lookup.",
    status: "Operational",
    statusType: "healthy",
    icon: "database",
    metrics: "127.0.0.1:3306 • 19 Tables Active"
  }
];

export const developmentMilestones = [
  {
    id: "MLS-01",
    title: "Phase 1 Kickoff & Foundation Architecture",
    category: "Architecture",
    status: "completed",
    date: "Aug 2026",
    description: "Initialized project structure, Vue 2 + Vite build pipeline, Tailwind design system tokens, and Pinia master state store."
  },
  {
    id: "MLS-02",
    title: "Product Master & Engineering Profile Module",
    category: "Master Data",
    status: "completed",
    date: "Aug 2026",
    description: "Implemented full CRUD, dynamic category assignment, lifecycle status management, and thumbnail uploads for IoT products."
  },
  {
    id: "MLS-03",
    title: "Central Components Library & Parametric Specs",
    category: "Master Data",
    status: "completed",
    date: "Aug 2026",
    description: "Built component database with MPN, dynamic key-value specifications, package types, and RoHS 3 compliance."
  },
  {
    id: "MLS-04",
    title: "Golang Gin + GORM Backend Integration (Phase 1B)",
    category: "Backend",
    status: "completed",
    date: "Aug 2026",
    description: "Engineered high-performance REST API in Golang, MySQL database migration, JWT auth middleware, and multipart file upload handling."
  },
  {
    id: "MLS-05",
    title: "Phase 2 — Product Versioning & Hardware Revision UI/UX",
    category: "Product Lifecycle",
    status: "completed",
    date: "Aug 2026",
    description: "Designed complete Phase 2 UI architecture: Product Version Workspace, Version Hero, Hardware Revision Tracks, Version Diff Matrix, and Lifecycle Tree."
  },
  {
    id: "MLS-06",
    title: "Phase 2C — Product Type & Project Architecture",
    category: "Architecture",
    status: "completed",
    date: "Aug 2026",
    description: "Completed architectural classification for Trading, Project, and R&D Products, live multi-currency pricing, and multi-image photo gallery."
  },
  {
    id: "MLS-07",
    title: "Phase 3 — Engineering BOM & Cost Foundation",
    category: "Engineering BOM",
    status: "completed",
    date: "Aug 2026",
    description: "Implemented multi-level hierarchical BOM editor, component quantities, recursive cost rollup, target profit margin markup, and project cost aggregator."
  },
  {
    id: "MLS-08",
    title: "Phase 4 — Prototype & Engineering Build Full-Stack Integration",
    category: "Prototypes",
    status: "completed",
    date: "Aug 2026",
    description: "Full-stack integration for Prototype builds, frozen BOM snapshot engine, component preparation tracking, and 8-stage assembly workflow with MySQL persistence."
  },
  {
    id: "MLS-09",
    title: "Phase 5 — Testing & Validation Control Full-Stack Architecture",
    category: "Testing & Validation",
    status: "completed",
    date: "Aug 2026",
    description: "Master Test Plans (DRAFT/RELEASED), frozen test snapshots, multi-run parametric execution bench, defect findings hub, and formal validation sign-off governance."
  },
  {
    id: "MLS-10",
    title: "Phase 6 — Engineering Change Management (ECR / ECO)",
    category: "Engineering Change",
    status: "completed",
    date: "Aug 2026",
    description: "Enterprise ECR/ECO workflows, 5-domain impact analysis, CCB multi-tier sign-offs, atomic baseline spawning (REV-B0), BOM clone & delta application, and closed-loop traceability DAG."
  }
];

export const moduleStatusMatrix = [
  { domain: "Foundation Shell & Navigation", progress: 100, status: "Complete", count: "4/4 Components", color: "blue" },
  { domain: "Product Master Management", progress: 100, status: "Complete", count: "3/3 Views", color: "blue" },
  { domain: "Components & Parametric Specs", progress: 100, status: "Complete", count: "3/3 Views", color: "emerald" },
  { domain: "Categories & Taxonomy Tree", progress: 100, status: "Complete", count: "2/2 Components", color: "purple" },
  { domain: "Manufacturers & Supply Chain", progress: 100, status: "Complete", count: "2/2 Views", color: "amber" },
  { domain: "Units of Measure (UoM)", progress: 100, status: "Complete", count: "1/1 View", color: "cyan" },
  { domain: "User Directory & RBAC Matrix", progress: 100, status: "Complete", count: "2/2 Views", color: "blue" },
  { domain: "JWT Authentication & Security", progress: 100, status: "Complete", count: "Full Suite", color: "emerald" },
  { domain: "Golang REST API & MySQL Engine", progress: 100, status: "Complete", count: "19 Endpoints", color: "cyan" },
  { domain: "Enterprise UI/UX Design System", progress: 100, status: "Complete", count: "Active", color: "blue" },
  { domain: "Product Versioning & Revisions", progress: 100, status: "Complete", count: "6 Views", color: "emerald" },
  { domain: "Product Type & Sourcing (Phase 2C)", progress: 100, status: "Complete", count: "5 Modules", color: "purple" },
  { domain: "Engineering BOM & Cost (Phase 3)", progress: 100, status: "Complete", count: "6 Modules", color: "emerald" },
  { domain: "Prototype Management (Phase 4)", progress: 100, status: "Complete", count: "5 Full-Stack Modules", color: "indigo" },
  { domain: "Testing & Validation Hub (Phase 5)", progress: 100, status: "Complete", count: "5 Full-Stack Modules", color: "emerald" },
  { domain: "Engineering Change Management (Phase 6)", progress: 100, status: "Complete", count: "7 Full-Stack Views", color: "cyan" }
];

/**
 * Dynamic Calculation Helpers
 */
export function calculateProgressStats() {
  const totalPhases = projectPhases.length;
  const completedPhases = projectPhases.filter(p => p.status === 'completed').length;
  const inProgressPhases = projectPhases.filter(p => p.status === 'in_progress').length;
  
  let totalModules = 0;
  let completedModules = 0;
  let plannedModules = 0;

  projectPhases.forEach(p => {
    p.modules.forEach(m => {
      totalModules++;
      if (m.status === 'completed') {
        completedModules++;
      } else {
        plannedModules++;
      }
    });
  });

  const totalProgressSum = projectPhases.reduce((acc, p) => acc + (p.progress || 0), 0);
  const overallRoadmapProgress = Math.round(totalProgressSum / totalPhases);
  const moduleCompletionProgress = Math.round((completedModules / totalModules) * 100);

  const currentPhase = projectPhases.find(p => p.id === projectMeta.currentPhaseId) || projectPhases[projectPhases.length - 1];
  const nextMilestonePhase = projectPhases.find(p => p.status === 'planned') || projectPhases[projectPhases.length - 1];

  return {
    totalPhases,
    completedPhases,
    inProgressPhases,
    totalModules,
    completedModules,
    plannedModules,
    overallRoadmapProgress,
    moduleCompletionProgress,
    currentPhase,
    nextMilestonePhase
  };
}
