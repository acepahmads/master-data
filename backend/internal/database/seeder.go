package database

import (
	"encoding/json"
	"log"
	"time"

	"iot-rd-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) {
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		log.Println("[INFO] Seeding initial IoT Master Data permissions, roles, and users...")

		// 1. SEED PERMISSIONS
	permissions := []model.Permission{
		{ID: "PERM-PRD-VIEW", Module: "products", Action: "view", Code: "products.view", Description: "View products"},
		{ID: "PERM-PRD-CREATE", Module: "products", Action: "create", Code: "products.create", Description: "Create products"},
		{ID: "PERM-PRD-EDIT", Module: "products", Action: "edit", Code: "products.edit", Description: "Edit products"},
		{ID: "PERM-PRD-DELETE", Module: "products", Action: "delete", Code: "products.delete", Description: "Delete products"},
		{ID: "PERM-PRD-EXPORT", Module: "products", Action: "export", Code: "products.export", Description: "Export products"},
		{ID: "PERM-PRD-APPROVE", Module: "products", Action: "approve", Code: "products.approve", Description: "Approve products"},

		{ID: "PERM-CMP-VIEW", Module: "components", Action: "view", Code: "components.view", Description: "View components"},
		{ID: "PERM-CMP-CREATE", Module: "components", Action: "create", Code: "components.create", Description: "Create components"},
		{ID: "PERM-CMP-EDIT", Module: "components", Action: "edit", Code: "components.edit", Description: "Edit components"},
		{ID: "PERM-CMP-DELETE", Module: "components", Action: "delete", Code: "components.delete", Description: "Delete components"},
		{ID: "PERM-CMP-EXPORT", Module: "components", Action: "export", Code: "components.export", Description: "Export components"},
		{ID: "PERM-CMP-APPROVE", Module: "components", Action: "approve", Code: "components.approve", Description: "Approve components"},

		{ID: "PERM-CAT-VIEW", Module: "categories", Action: "view", Code: "categories.view", Description: "View categories"},
		{ID: "PERM-CAT-CREATE", Module: "categories", Action: "create", Code: "categories.create", Description: "Create categories"},
		{ID: "PERM-CAT-EDIT", Module: "categories", Action: "edit", Code: "categories.edit", Description: "Edit categories"},
		{ID: "PERM-CAT-DELETE", Module: "categories", Action: "delete", Code: "categories.delete", Description: "Delete categories"},

		{ID: "PERM-MFG-VIEW", Module: "manufacturers", Action: "view", Code: "manufacturers.view", Description: "View manufacturers"},
		{ID: "PERM-MFG-CREATE", Module: "manufacturers", Action: "create", Code: "manufacturers.create", Description: "Create manufacturers"},
		{ID: "PERM-MFG-EDIT", Module: "manufacturers", Action: "edit", Code: "manufacturers.edit", Description: "Edit manufacturers"},
		{ID: "PERM-MFG-DELETE", Module: "manufacturers", Action: "delete", Code: "manufacturers.delete", Description: "Delete manufacturers"},

		{ID: "PERM-SUP-VIEW", Module: "suppliers", Action: "view", Code: "suppliers.view", Description: "View suppliers"},
		{ID: "PERM-SUP-CREATE", Module: "suppliers", Action: "create", Code: "suppliers.create", Description: "Create suppliers"},
		{ID: "PERM-SUP-EDIT", Module: "suppliers", Action: "edit", Code: "suppliers.edit", Description: "Edit suppliers"},
		{ID: "PERM-SUP-DELETE", Module: "suppliers", Action: "delete", Code: "suppliers.delete", Description: "Delete suppliers"},

		{ID: "PERM-UNT-VIEW", Module: "units", Action: "view", Code: "units.view", Description: "View units"},
		{ID: "PERM-UNT-CREATE", Module: "units", Action: "create", Code: "units.create", Description: "Create units"},
		{ID: "PERM-UNT-EDIT", Module: "units", Action: "edit", Code: "units.edit", Description: "Edit units"},
		{ID: "PERM-UNT-DELETE", Module: "units", Action: "delete", Code: "units.delete", Description: "Delete units"},

		{ID: "PERM-ADM-USERS", Module: "administration", Action: "manageUsers", Code: "administration.manageUsers", Description: "Manage user accounts"},
		{ID: "PERM-ADM-ROLES", Module: "administration", Action: "manageRoles", Code: "administration.manageRoles", Description: "Manage roles & permissions"},
		{ID: "PERM-ADM-LOGS", Module: "administration", Action: "auditLogs", Code: "administration.auditLogs", Description: "View audit logs"},
		{ID: "PERM-ADM-SETTINGS", Module: "administration", Action: "systemSettings", Code: "administration.systemSettings", Description: "Configure system settings"},
	}
	db.Create(&permissions)

	// 2. SEED ROLES
	roles := []model.Role{
		{
			ID:          "ROLE-ADMIN",
			Name:        "Super Admin",
			Code:        "SUPER_ADMIN",
			Description: "Complete unrestricted administrative control across all master data modules, users, security policies, and system settings.",
			IsSystem:    true,
			Status:      "Active",
		},
		{
			ID:          "ROLE-RNDM",
			Name:        "R&D Manager",
			Code:        "RND_MANAGER",
			Description: "Engineering department management with full access to approve components, release product masters, and assign engineering leads.",
			IsSystem:    false,
			Status:      "Active",
		},
		{
			ID:          "ROLE-HWENG",
			Name:        "Hardware Engineer",
			Code:        "HW_ENGINEER",
			Description: "Engineering contributor capable of creating components, maintaining technical specs, and attaching engineering documentation.",
			IsSystem:    false,
			Status:      "Active",
		},
		{
			ID:          "ROLE-FWENG",
			Name:        "Firmware Engineer",
			Code:        "FW_ENGINEER",
			Description: "Embedded software engineer focused on MCU memory maps, pinouts, interfaces, and product master references.",
			IsSystem:    false,
			Status:      "Active",
		},
		{
			ID:          "ROLE-PURCH",
			Name:        "Purchasing Specialist",
			Code:        "PURCHASING",
			Description: "Supply chain management focused on suppliers, manufacturers, pricing terms, lead times, and unit conversions.",
			IsSystem:    false,
			Status:      "Active",
		},
		{
			ID:          "ROLE-VIEW",
			Name:        "Viewer",
			Code:        "VIEWER",
			Description: "Read-only access across all master data screens for external auditors and cross-functional observers.",
			IsSystem:    false,
			Status:      "Active",
		},
	}
	db.Create(&roles)

	// Assign all permissions to Super Admin
	for _, p := range permissions {
		db.Create(&model.RolePermission{RoleID: "ROLE-ADMIN", PermissionID: p.ID})
		if p.Action == "view" {
			db.Create(&model.RolePermission{RoleID: "ROLE-VIEW", PermissionID: p.ID})
		}
	}

	// 3. SEED USERS
	pwdHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	now := time.Now()
	users := []model.User{
		{
			ID:           "USR-001",
			Name:         "Alex Chen",
			Username:     "achen",
			Email:        "alex.chen@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-ADMIN",
			Department:   "Hardware R&D",
			Status:       "Active",
			TwoFactor:    true,
			LastLoginAt:  &now,
			Avatar:       "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:           "USR-002",
			Name:         "Dr. Sarah Jenkins",
			Username:     "sjenkins",
			Email:        "sarah.jenkins@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-RNDM",
			Department:   "Engineering Leadership",
			Status:       "Active",
			TwoFactor:    true,
			Avatar:       "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:           "USR-003",
			Name:         "Marcus Vance",
			Username:     "mvance",
			Email:        "marcus.vance@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-HWENG",
			Department:   "Electronics & RF",
			Status:       "Active",
			TwoFactor:    true,
			Avatar:       "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:           "USR-004",
			Name:         "Elena Rostova",
			Username:     "erostova",
			Email:        "elena.rostova@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-FWENG",
			Department:   "Embedded Systems",
			Status:       "Active",
			TwoFactor:    true,
			Avatar:       "https://images.unsplash.com/photo-1580489944761-15a19d654956?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:           "USR-005",
			Name:         "David Tanaka",
			Username:     "dtanaka",
			Email:        "david.tanaka@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-PURCH",
			Department:   "Supply Chain & Sourcing",
			Status:       "Active",
			TwoFactor:    false,
			Avatar:       "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:           "USR-006",
			Name:         "Maya Lin",
			Username:     "mlin",
			Email:        "maya.lin@iotcontrol.io",
			PasswordHash: string(pwdHash),
			RoleID:       "ROLE-VIEW",
			Department:   "Quality & Compliance",
			Status:       "Active",
			TwoFactor:    false,
			Avatar:       "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100&auto=format&fit=crop&q=80",
		},
	}
	db.Create(&users)
	}

	// 4. SEED UNITS (UoM)
	var unitCount int64
	db.Model(&model.Unit{}).Count(&unitCount)
	if unitCount == 0 {
		units := []model.Unit{
			{ID: "UOM-001", Code: "PCS", Name: "Piece", Symbol: "pcs", Category: "Count", Description: "Discrete individual component unit", Status: "Active", IsDefault: true},
			{ID: "UOM-002", Code: "UNIT", Name: "Unit / Device", Symbol: "ea", Category: "Count", Description: "Finished assembly or modular device", Status: "Active", IsDefault: false},
			{ID: "UOM-003", Code: "REEL", Name: "SMD Tape Reel", Symbol: "rl", Category: "Packaging", Description: "Standard carrier tape reel (e.g. 3000/5000 pcs)", Status: "Active", IsDefault: false},
			{ID: "UOM-004", Code: "TUBE", Name: "IC Rail / Tube", Symbol: "tb", Category: "Packaging", Description: "Anti-static extruded stick/tube container", Status: "Active", IsDefault: false},
			{ID: "UOM-005", Code: "TRAY", Name: "JEDEC Matrix Tray", Symbol: "try", Category: "Packaging", Description: "High-temperature plastic IC shipping matrix tray", Status: "Active", IsDefault: false},
			{ID: "UOM-006", Code: "SET", Name: "Kit / Assembly Set", Symbol: "set", Category: "Packaging", Description: "Bundled kit of related fasteners or accessories", Status: "Active", IsDefault: false},
			{ID: "UOM-007", Code: "METER", Name: "Linear Meter", Symbol: "m", Category: "Length", Description: "Wire, coaxial cable, or heat shrink tubing", Status: "Active", IsDefault: false},
			{ID: "UOM-008", Code: "CM", Name: "Centimeter", Symbol: "cm", Category: "Length", Description: "Flexible flat cables (FFC) or thermal pads", Status: "Active", IsDefault: false},
			{ID: "UOM-009", Code: "GRAM", Name: "Gram", Symbol: "g", Category: "Weight", Description: "Solder paste, potting epoxy, thermal grease", Status: "Active", IsDefault: false},
			{ID: "UOM-010", Code: "KG", Name: "Kilogram", Symbol: "kg", Category: "Weight", Description: "Bulk raw plastics, potting silicone drums", Status: "Active", IsDefault: false},
		}
		db.Create(&units)
	}

	// 5. SEED MANUFACTURERS
	var mfgCount int64
	db.Model(&model.Manufacturer{}).Count(&mfgCount)
	if mfgCount == 0 {
		manufacturers := []model.Manufacturer{
			{ID: "MFG-STM", Name: "STMicroelectronics", Code: "STM", Country: "Switzerland", CountryCode: "CH", Website: "https://www.st.com", TotalComponents: 78, Status: "Approved", Rating: "Tier 1 Preferred", ContactPerson: "Jean-Luc Picard", ContactEmail: "technical.sales@st.com", Description: "Global semiconductor leader delivering intelligent and energy-efficient solutions for automotive, industrial, and IoT applications."},
			{ID: "MFG-ESPR", Name: "Espressif Systems", Code: "ESPR", Country: "China", CountryCode: "CN", Website: "https://www.espressif.com", TotalComponents: 42, Status: "Approved", Rating: "Tier 1 Preferred", ContactPerson: "Wang Qiang", ContactEmail: "sales@espressif.com", Description: "Fabless semiconductor company specializing in low-power Wi-Fi and Bluetooth wireless SoCs for AIoT hardware."},
			{ID: "MFG-BSCH", Name: "Bosch Sensortec", Code: "BSCH", Country: "Germany", CountryCode: "DE", Website: "https://www.bosch-sensortec.com", TotalComponents: 29, Status: "Approved", Rating: "Tier 1 Preferred", ContactPerson: "Klaus Schmidt", ContactEmail: "contact@bosch-sensortec.com", Description: "Pioneer and world leader in MEMS environmental and inertial sensors for consumer and industrial electronics."},
			{ID: "MFG-TI", Name: "Texas Instruments", Code: "TXN", Country: "United States", CountryCode: "US", Website: "https://www.ti.com", TotalComponents: 95, Status: "Approved", Rating: "Tier 1 Preferred", ContactPerson: "Emily Rogers", ContactEmail: "support@ti.com", Description: "Global semiconductor company that designs, manufactures, tests and sells analog and embedded processing chips."},
			{ID: "MFG-QCTL", Name: "Quectel Wireless", Code: "QCTL", Country: "China", CountryCode: "CN", Website: "https://www.quectel.com", TotalComponents: 34, Status: "Approved", Rating: "Approved", ContactPerson: "Li Ming", ContactEmail: "global-sales@quectel.com", Description: "Leading global supplier of cellular IoT, GNSS, satellite, Wi-Fi, and Bluetooth modules and smart antennas."},
			{ID: "MFG-AMPH", Name: "Amphenol Industrial", Code: "AMPH", Country: "United States", CountryCode: "US", Website: "https://www.amphenol-industrial.com", TotalComponents: 48, Status: "Approved", Rating: "Approved", ContactPerson: "Robert Miller", ContactEmail: "custservice@amphenol-aio.com", Description: "World leader in harsh environment industrial interconnect systems, circular M-series connectors, and custom cabling."},
			{ID: "MFG-HMM", Name: "Hammond Manufacturing", Code: "HMM", Country: "Canada", CountryCode: "CA", Website: "https://www.hammfg.com", TotalComponents: 26, Status: "Approved", Rating: "Approved", ContactPerson: "Dave MacLeod", ContactEmail: "sales@hammfg.com", Description: "Leading manufacturer of standard and custom electronic enclosures, die-cast aluminum boxes, and racks."},
			{ID: "MFG-MUR", Name: "Murata Electronics", Code: "MUR", Country: "Japan", CountryCode: "JP", Website: "https://www.murata.com", TotalComponents: 110, Status: "Approved", Rating: "Tier 1 Preferred", ContactPerson: "Kenji Sato", ContactEmail: "inquiry@murata.com", Description: "Worldwide leader in advanced ceramic electronic passive components, power modules, and RF communication modules."},
		}
		db.Create(&manufacturers)
	}

	// 6. SEED SUPPLIERS
	var supCount int64
	db.Model(&model.Supplier{}).Count(&supCount)
	if supCount == 0 {
		suppliers := []model.Supplier{
			{ID: "SUP-DIGI", Name: "DigiKey Electronics", Code: "DIGIKEY-US", Country: "United States", CountryCode: "US", ContactPerson: "Dave Larson", Email: "sales@digikey.com", Phone: "+1 800-344-4539", Tier: "Tier 1 Authorized Distributor", Status: "Active", TotalComponents: 185, Terms: "Net 30 Days", LeadTimeAvg: "2-4 Days", Website: "https://www.digikey.com", Description: "Primary authorized franchised distributor with worldwide same-day dispatch and EDI integration."},
			{ID: "SUP-MOUS", Name: "Mouser Electronics", Code: "MOUSER-US", Country: "United States", CountryCode: "US", ContactPerson: "Sarah Jenkins", Email: "orders@mouser.com", Phone: "+1 800-346-6873", Tier: "Tier 1 Authorized Distributor", Status: "Active", TotalComponents: 140, Terms: "Net 30 Days", LeadTimeAvg: "2-5 Days", Website: "https://www.mouser.com", Description: "Leading authorized distributor for rapid new product introduction (NPI) prototypes and engineering sample batches."},
			{ID: "SUP-ARRW", Name: "Arrow Electronics", Code: "ARROW-GLB", Country: "United States", CountryCode: "US", ContactPerson: "Michael Vance", Email: "americas@arrow.com", Phone: "+1 855-326-4757", Tier: "Global Volume Partner", Status: "Active", TotalComponents: 92, Terms: "Net 60 Days", LeadTimeAvg: "5-10 Days", Website: "https://www.arrow.com", Description: "Global supply-chain aggregation partner for high volume production runs and buffer stock agreements."},
			{ID: "SUP-LCSC", Name: "LCSC Electronics", Code: "LCSC-CN", Country: "China", CountryCode: "CN", ContactPerson: "Chloe Zhang", Email: "service@lcsc.com", Phone: "+86 755-8395-5800", Tier: "Fast-Turn Regional Partner", Status: "Active", TotalComponents: 76, Terms: "Prepaid / Credit Card", LeadTimeAvg: "5-8 Days", Website: "https://www.lcsc.com", Description: "Specialized Asian marketplace for cost-effective passives, Asian semiconductor alternates, and PCB-ready SMD reels."},
		}
		db.Create(&suppliers)
	}

	// 7. SEED CATEGORIES (Hierarchical Tree)
	var catCount int64
	db.Model(&model.ComponentCategory{}).Count(&catCount)
	if catCount == 0 {
		categories := []model.ComponentCategory{
			{ID: "CAT-ELEC", Name: "Electronic Components", Code: "ELEC", Description: "All active, passive, and integrated semiconductor electronic hardware components.", ParentID: nil, Status: "Active", TotalComponents: 342},
			{ID: "CAT-MCU", Name: "Microcontrollers & Processors", Code: "ELEC-MCU", Description: "Core microcontrollers, DSPs, application processors, and SoCs.", ParentID: strPtr("CAT-ELEC"), Status: "Active", TotalComponents: 74},
			{ID: "CAT-SENS", Name: "Environmental & Inertial Sensors", Code: "ELEC-SENS", Description: "MEMS sensors, gas, temperature, pressure, accelerometer, and gyro sensors.", ParentID: strPtr("CAT-ELEC"), Status: "Active", TotalComponents: 56},
			{ID: "CAT-RF", Name: "Wireless & Cellular Modules", Code: "ELEC-RF", Description: "Wi-Fi, Bluetooth BLE, LoRa, Cellular 4G/5G/NB-IoT, and GNSS modules.", ParentID: strPtr("CAT-ELEC"), Status: "Active", TotalComponents: 48},
			{ID: "CAT-PWR", Name: "Power Management & Conversion", Code: "ELEC-PWR", Description: "Buck/boost converters, PMICs, LDOs, battery chargers, and protection circuits.", ParentID: strPtr("CAT-ELEC"), Status: "Active", TotalComponents: 82},
			{ID: "CAT-PASS", Name: "Passive Components", Code: "ELEC-PASS", Description: "SMD MLCC capacitors, precision resistors, power inductors, and crystals.", ParentID: strPtr("CAT-ELEC"), Status: "Active", TotalComponents: 82},

			{ID: "CAT-MECH", Name: "Mechanical & Structural", Code: "MECH", Description: "Physical enclosures, thermal management, brackets, and structural hardware.", ParentID: nil, Status: "Active", TotalComponents: 124},
			{ID: "CAT-ENC", Name: "IP67/IP68 Enclosures", Code: "MECH-ENC", Description: "Die-cast aluminum, polycarbonate, and extruded weatherproof enclosures.", ParentID: strPtr("CAT-MECH"), Status: "Active", TotalComponents: 38},
			{ID: "CAT-CONN", Name: "Connectors & Cabling", Code: "MECH-CONN", Description: "M12 industrial connectors, terminal blocks, SMA RF jacks, and USB-C.", ParentID: strPtr("CAT-MECH"), Status: "Active", TotalComponents: 54},
			{ID: "CAT-FAST", Name: "Fasteners & Standoffs", Code: "MECH-FAST", Description: "Stainless steel M3/M4 screws, brass PCB standoffs, and sealing gaskets.", ParentID: strPtr("CAT-MECH"), Status: "Active", TotalComponents: 32},

			{ID: "CAT-OPT", Name: "Optoelectronics & Display", Code: "OPT", Description: "Displays, indicators, light guides, and optical filters.", ParentID: nil, Status: "Active", TotalComponents: 45},
			{ID: "CAT-DISP", Name: "OLED & E-Paper Panels", Code: "OPT-DISP", Description: "Monochrome and color graphical SPI/I2C display modules.", ParentID: strPtr("CAT-OPT"), Status: "Active", TotalComponents: 24},
			{ID: "CAT-LED", Name: "LED Indicators & Light Pipes", Code: "OPT-LED", Description: "SMD multi-color LEDs and polycarbonate optical light guides.", ParentID: strPtr("CAT-OPT"), Status: "Active", TotalComponents: 21},
		}
		db.Create(&categories)
	}

	// 8. SEED PRODUCTS
	var prodCount int64
	db.Model(&model.Product{}).Count(&prodCount)
	if prodCount == 0 {
		products := []model.Product{
		{
			ID:                "PRD-2024-001",
			Code:              "PRD-2024-001",
			Name:              "Smart Edge Industrial Gateway (GW-400X)",
			Category:          "Edge Gateways",
			CategoryID:        strPtr("CAT-GW"),
			Status:            "Active",
			Description:       "Enterprise quad-core ARM Cortex gateway with dual Gigabit Ethernet, RS-485 Modbus interface, and integrated 4G LTE fallback. Designed for smart factory IoT deployments.",
			LongDescription:   "The GW-400X is an enterprise-grade IoT edge computer capable of running containerized telemetry processing workloads at the edge. Features robust surge protection on all I/O, dual power inputs (9-36V DC), and an IP50 metal enclosure suitable for DIN rail or wall mounting.",
			Image:             "https://images.unsplash.com/photo-1518770660439-4636190af475?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Dr. Sarah Jenkins",
			TargetMarket:      "Industrial Automation & Smart Factories",
			PowerRating:       "12-24V DC / 15W Max",
			IngressProtection: "IP50 (DIN-rail mountable)",
			OperatingTemp:     "-40°C to +85°C",
			Certifications:    `["CE", "FCC Part 15B", "RoHS 3", "UL 61010-1"]`,
			ComponentsCount:   142,
		},
		{
			ID:                "PRD-2024-002",
			Code:              "PRD-2024-002",
			Name:              "Multi-Sensor Environmental Node (EN-100)",
			Category:          "Sensors & Nodes",
			CategoryID:        strPtr("CAT-SENS"),
			Status:            "Active",
			Description:       "Solar-assisted environmental telemetry node monitoring PM2.5, VOC, Temperature, Humidity, and Barometric pressure over LoRaWAN.",
			LongDescription:   "Compact autonomous node engineered for outdoor urban air quality monitoring and greenhouse automation. Boasts 5-year battery life when paired with integrated solar harvester.",
			Image:             "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Marcus Vance",
			TargetMarket:      "Smart Cities & Smart Agriculture",
			PowerRating:       "3.6V Li-SOCl2 / Solar Cell",
			IngressProtection: "IP67 Weatherproof",
			OperatingTemp:     "-20°C to +70°C",
			Certifications:    `["CE", "LoRaWAN Regional 915/868", "RoHS"]`,
			ComponentsCount:   68,
		},
		{
			ID:                "PRD-2024-003",
			Code:              "PRD-2024-003",
			Name:              "Solar Asset Tracker & Telemetry Hub (AT-550)",
			Category:          "Asset Tracking",
			CategoryID:        strPtr("CAT-TRACK"),
			Status:            "Active",
			Description:       "Ruggedized GNSS + NB-IoT / LTE-M asset tracker with Bluetooth 5.2 beacon receiver and 3-axis crash detection accelerometer.",
			LongDescription:   "Designed for intermodal container tracking and heavy equipment fleet monitoring. Features magnetic mounting, tamper detection, and rechargeable solar backup.",
			Image:             "https://images.unsplash.com/photo-1558346490-a72e53ae2d4f?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Elena Rostova",
			TargetMarket:      "Cold Chain & Logistics",
			PowerRating:       "3.7V 5000mAh LiPo + High Efficiency Solar",
			IngressProtection: "IP68 & IK08 Shockproof",
			OperatingTemp:     "-30°C to +75°C",
			Certifications:    `["FCC", "PTCRB", "AT&T / Verizon Certified"]`,
			ComponentsCount:   84,
		},
		{
			ID:                "PRD-2024-004",
			Code:              "PRD-2024-004",
			Name:              "Sub-GHz Irrigation Valve Controller (IC-900)",
			Category:          "Sensors & Nodes",
			CategoryID:        strPtr("CAT-SENS"),
			Status:            "Draft",
			Description:       "Battery powered pulse valve actuator with 4 dry-contact pulse meter inputs and long-range proprietary 915MHz RF link.",
			LongDescription:   "Next generation precision agriculture valve controller with ultra-low sleep current (<5uA) for remote commercial vineyard installations.",
			Image:             "https://images.unsplash.com/photo-1581092335397-9583fe92d232?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Marcus Vance",
			TargetMarket:      "AgTech & Water Utilities",
			PowerRating:       "4x D-Cell Alkaline (3-Year Lifespan)",
			IngressProtection: "IP68 Submersible 1m",
			OperatingTemp:     "-10°C to +60°C",
			Certifications:    `["FCC Part 15C", "RoHS"]`,
			ComponentsCount:   45,
		},
		{
			ID:                "PRD-2024-005",
			Code:              "PRD-2024-005",
			Name:              "Smart Energy Sub-Metering Gateway (EMG-200)",
			Category:          "Edge Gateways",
			CategoryID:        strPtr("CAT-GW"),
			Status:            "In Review",
			Description:       "Multi-channel 3-phase power quality and energy meter with Wi-Fi, Ethernet, and RS485 interfaces.",
			LongDescription:   "High precision ANSI C12.20 Class 0.2 energy sub-metering appliance for commercial tenant billing and building energy management systems.",
			Image:             "https://images.unsplash.com/photo-1517420704952-d9f39e95b43e?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Alex Chen",
			TargetMarket:      "Commercial Real Estate & Microgrids",
			PowerRating:       "85-265V AC / 5W",
			IngressProtection: "IP40",
			OperatingTemp:     "-25°C to +70°C",
			Certifications:    `["UL 61010-2-030", "CE", "MID Approved"]`,
			ComponentsCount:   110,
		},
		{
			ID:                "PRD-2024-006",
			Code:              "PRD-2024-006",
			Name:              "High-Precision RTK Fleet Beacon (FB-800)",
			Category:          "Asset Tracking",
			CategoryID:        strPtr("CAT-TRACK"),
			Status:            "Active",
			Description:       "Centimeter-level GNSS RTK positioning beacon with dual frequency L1/L5 antenna and CAN-bus integration.",
			LongDescription:   "Precision navigation module for autonomous mobile robots (AMRs), automated guided vehicles (AGVs), and smart construction equipment.",
			Image:             "https://images.unsplash.com/photo-1581092580497-e0d23cbdf1dc?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Dr. Sarah Jenkins",
			TargetMarket:      "Robotics & Autonomous Heavy Equipment",
			PowerRating:       "9-36V DC via M12 Connector",
			IngressProtection: "IP67",
			OperatingTemp:     "-40°C to +85°C",
			Certifications:    `["ISO 16750", "ECE-R10", "CE"]`,
			ComponentsCount:   95,
		},
	}
	db.Create(&products)
	}

	// 9. SEED COMPONENTS WITH DYNAMIC SPECS AND DOCUMENTS
	var compCount int64
	db.Model(&model.Component{}).Count(&compCount)
	if compCount == 0 {
		c1 := model.Component{
		ID:                  "CMP-001",
		PartNumber:          "STM32H753BIT6",
		InternalCode:        "IPN-MCU-001",
		Name:                "Arm Cortex-M7 480MHz High-Performance MCU",
		CategoryID:          "CAT-MCU",
		ManufacturerID:      "MFG-STM",
		SupplierID:          "SUP-DIGI",
		UnitID:              "UOM-001",
		PackageType:         "LQFP-208",
		MountingType:        "Surface Mount (SMD)",
		LeadTime:            "12 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "High-performance 32-bit RISC core with 2MB Flash, 1MB RAM, Cryptographic engine, and Gigabit Ethernet MAC.",
		DetailedDescription: "The STM32H753xI devices are based on the high-performance Arm Cortex-M7 32-bit RISC core operating at up to 480 MHz. Features Hardware floating point unit (FPU), L1 cache (16 KB I + 16 KB D), high-speed crypto accelerator, and advanced analog peripherals.",
		Image:               "https://images.unsplash.com/photo-1518770660439-4636190af475?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
		Specs: []model.ComponentSpecification{
			{ID: "SPEC-101", SpecificationName: "Core Architecture", SpecificationValue: "Arm Cortex-M7", SpecificationUnit: "", DisplayOrder: 1},
			{ID: "SPEC-102", SpecificationName: "Clock Speed", SpecificationValue: "480", SpecificationUnit: "MHz", DisplayOrder: 2},
			{ID: "SPEC-103", SpecificationName: "Flash Memory", SpecificationValue: "2048", SpecificationUnit: "KB", DisplayOrder: 3},
			{ID: "SPEC-104", SpecificationName: "RAM Size", SpecificationValue: "1024", SpecificationUnit: "KB", DisplayOrder: 4},
			{ID: "SPEC-105", SpecificationName: "Supply Voltage (Min/Max)", SpecificationValue: "1.62 - 3.6", SpecificationUnit: "V", DisplayOrder: 5},
			{ID: "SPEC-106", SpecificationName: "Operating Temperature", SpecificationValue: "-40 to +85", SpecificationUnit: "°C", DisplayOrder: 6},
		},
		Documents: []model.ComponentDocument{
			{ID: "DOC-101", DocumentType: "Datasheet", FileName: "STM32H753BIT6_Datasheet_Rev6.pdf", FilePath: "/uploads/STM32H753BIT6_Datasheet_Rev6.pdf", FileSize: "3.4 MB", MimeType: "application/pdf", UploadedBy: "Alex Chen"},
		},
	}
	db.Create(&c1)

	c2 := model.Component{
		ID:                  "CMP-002",
		PartNumber:          "ESP32-S3-WROOM-1-N8R8",
		InternalCode:        "IPN-WIFI-002",
		Name:                "Wi-Fi 4 (802.11 b/g/n) & Bluetooth 5 (LE) Module",
		CategoryID:          "CAT-RF",
		ManufacturerID:      "MFG-ESPR",
		SupplierID:          "SUP-MOUS",
		UnitID:              "UOM-001",
		PackageType:         "SMD Module (41-pin)",
		MountingType:        "Surface Mount",
		LeadTime:            "6 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Xtensa 32-bit LX7 dual-core module with 8MB Octal SPI Flash and 8MB Octal PSRAM + PCB antenna.",
		DetailedDescription: "ESP32-S3-WROOM-1 is a powerful, generic Wi-Fi + Bluetooth LE MCU module with a rich set of peripherals, vector instructions for AI acceleration, and reliable security mechanisms.",
		Image:               "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
		Specs: []model.ComponentSpecification{
			{ID: "SPEC-201", SpecificationName: "Core Processor", SpecificationValue: "Dual-core Xtensa 32-bit LX7", SpecificationUnit: "", DisplayOrder: 1},
			{ID: "SPEC-202", SpecificationName: "Frequency", SpecificationValue: "240", SpecificationUnit: "MHz", DisplayOrder: 2},
			{ID: "SPEC-203", SpecificationName: "Embedded Flash", SpecificationValue: "8", SpecificationUnit: "MB", DisplayOrder: 3},
		},
	}
	db.Create(&c2)

	c3 := model.Component{
		ID:                  "CMP-003",
		PartNumber:          "BME680",
		InternalCode:        "IPN-SENS-003",
		Name:                "Gas, Humidity, Pressure & Temperature 4-in-1 Sensor",
		CategoryID:          "CAT-SENS",
		ManufacturerID:      "MFG-BSCH",
		SupplierID:          "SUP-DIGI",
		UnitID:              "UOM-001",
		PackageType:         "LGA-8 (3.0 x 3.0 x 0.93 mm)",
		MountingType:        "Surface Mount",
		LeadTime:            "8 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Integrated environmental sensor combining VOC gas sensor with barometric pressure, humidity, and ambient temperature.",
		DetailedDescription: "The BME680 is an integrated environmental sensor developed specifically for mobile applications and wearables where size and low power consumption are key requirements.",
		Image:               "https://images.unsplash.com/photo-1558346490-a72e53ae2d4f?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
		Specs: []model.ComponentSpecification{
			{ID: "SPEC-301", SpecificationName: "Supply Voltage VDD", SpecificationValue: "1.71 - 3.6", SpecificationUnit: "V", DisplayOrder: 1},
			{ID: "SPEC-302", SpecificationName: "Humidity Accuracy", SpecificationValue: "±3", SpecificationUnit: "% RH", DisplayOrder: 2},
			{ID: "SPEC-303", SpecificationName: "Pressure Accuracy", SpecificationValue: "±0.6", SpecificationUnit: "hPa", DisplayOrder: 3},
		},
	}
	db.Create(&c3)

	c4 := model.Component{
		ID:                  "CMP-004",
		PartNumber:          "BG95-M3",
		InternalCode:        "IPN-RF-004",
		Name:                "Multi-Mode LTE Cat M1/NB2/EGPRS & GNSS Module",
		CategoryID:          "CAT-RF",
		ManufacturerID:      "MFG-QCTL",
		SupplierID:          "SUP-ARRW",
		UnitID:              "UOM-001",
		PackageType:         "LGA-102 (23.6 x 19.9 x 2.2 mm)",
		MountingType:        "Surface Mount",
		LeadTime:            "10 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Ultra-compact multi-mode LPWA module supporting LTE Cat M1, Cat NB2, EGPRS, and integrated GNSS.",
		DetailedDescription: "BG95-M3 achieves maximum data rates of 588 kbps downlink and 1119 kbps uplink under LTE Cat M1. Features ultra-low power consumption via Power Saving Mode (PSM) and eDRX.",
		Image:               "https://images.unsplash.com/photo-1518770660439-4636190af475?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
	}
	db.Create(&c4)

	c5 := model.Component{
		ID:                  "CMP-005",
		PartNumber:          "TPS54302DDCR",
		InternalCode:        "IPN-PWR-005",
		Name:                "4.5V to 28V Input, 3A Synchronous Step-Down Converter",
		CategoryID:          "CAT-PWR",
		ManufacturerID:      "MFG-TI",
		SupplierID:          "SUP-DIGI",
		UnitID:              "UOM-003",
		PackageType:         "SOT-23-6 (THIN)",
		MountingType:        "Surface Mount",
		LeadTime:            "4 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "High-efficiency 3A synchronous buck converter with integrated 85mΩ and 40mΩ MOSFETs and 400kHz fixed switching frequency.",
		DetailedDescription: "The TPS54302 is a 28V, 3A step-down converter with integrated high and low-side FETs. Fixed-frequency peak current-mode control delivers fast transient response.",
		Image:               "https://images.unsplash.com/photo-1581092335397-9583fe92d232?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
	}
	db.Create(&c5)

	c6 := model.Component{
		ID:                  "CMP-006",
		PartNumber:          "ENC-IP67-ALU-120",
		InternalCode:        "IPN-MECH-006",
		Name:                "Die-Cast Aluminum IP67 Weatherproof Enclosure 120x80x40mm",
		CategoryID:          "CAT-ENC",
		ManufacturerID:      "MFG-HMM",
		SupplierID:          "SUP-MOUS",
		UnitID:              "UOM-001",
		PackageType:         "Box Enclosure",
		MountingType:        "Wall / Flange Mount",
		LeadTime:            "2 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Heavy-duty die-cast aluminum alloy IP67 enclosure with silicone rubber gasket and stainless steel screws.",
		DetailedDescription: "Engineered for outdoor telemetry gateways requiring EMI shielding and rugged impact protection. Sandblast anodized dark slate finish.",
		Image:               "https://images.unsplash.com/photo-1581092580497-e0d23cbdf1dc?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
	}
	db.Create(&c6)

	c7 := model.Component{
		ID:                  "CMP-007",
		PartNumber:          "M12-8P-FEM-PNL",
		InternalCode:        "IPN-CONN-007",
		Name:                "M12 8-Pin Female A-Coded Panel Mount Connector",
		CategoryID:          "CAT-CONN",
		ManufacturerID:      "MFG-AMPH",
		SupplierID:          "SUP-DIGI",
		UnitID:              "UOM-001",
		PackageType:         "Panel Mount (PG9 / M16)",
		MountingType:        "Front Panel Screw-in",
		LeadTime:            "3 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Waterproof M12 A-Coded 8-pole female panel receptacle with 200mm pre-stripped wires.",
		DetailedDescription: "Ruggedized nickel-plated brass industrial sensor interface connector rated IP68 when mated. Ideal for RS-485 and power supply inputs.",
		Image:               "https://images.unsplash.com/photo-1517420704952-d9f39e95b43e?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
	}
	db.Create(&c7)

	c8 := model.Component{
		ID:                  "CMP-008",
		PartNumber:          "GRM188R61A106KE69D",
		InternalCode:        "IPN-PASS-008",
		Name:                "Capacitor Ceramic 10µF 10V X5R 0603 (1608 Metric)",
		CategoryID:          "CAT-PASS",
		ManufacturerID:      "MFG-MUR",
		SupplierID:          "SUP-LCSC",
		UnitID:              "UOM-003",
		PackageType:         "0603 SMD",
		MountingType:        "Surface Mount",
		LeadTime:            "2 Weeks",
		RoHSCompliant:       true,
		REACHCompliant:      true,
		ShortDescription:    "Multi-layer ceramic chip capacitor 10µF ±10% 10V X5R dielectric in standard 0603 footprint.",
		DetailedDescription: "High volumetric capacitance MLCC for power decoupling and DC-DC filter stages in compact IoT circuits.",
		Image:               "https://images.unsplash.com/photo-1518770660439-4636190af475?w=300&auto=format&fit=crop&q=80",
		Status:              "Active",
	}
	db.Create(&c8)
	}

	// 10. SEED INITIAL ACTIVITIES
	var actCount int64
	db.Model(&model.ActivityLog{}).Count(&actCount)
	if actCount == 0 {
		activities := []model.ActivityLog{
		{
			ID:          "ACT-001",
			UserID:      "USR-001",
			UserName:    "Alex Chen",
			UserAvatar:  "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
			Module:      "components",
			EntityType:  "Component",
			EntityID:    "CMP-001",
			EntityName:  "STM32H753BIT6",
			Action:      "Updated Component Specs",
			Description: "Updated parametric technical specifications for Arm Cortex-M7 MCU.",
			BadgeColor:  "blue",
		},
		{
			ID:          "ACT-002",
			UserID:      "USR-002",
			UserName:    "Dr. Sarah Jenkins",
			UserAvatar:  "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100&auto=format&fit=crop&q=80",
			Module:      "products",
			EntityType:  "Product",
			EntityID:    "PRD-2024-001",
			EntityName:  "Smart Edge Industrial Gateway (GW-400X)",
			Action:      "Created Product Master",
			Description: "Created new industrial edge gateway hardware product master.",
			BadgeColor:  "emerald",
		},
		{
			ID:          "ACT-003",
			UserID:      "USR-005",
			UserName:    "David Tanaka",
			UserAvatar:  "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&auto=format&fit=crop&q=80",
			Module:      "suppliers",
			EntityType:  "Supplier",
			EntityID:    "SUP-DIGI",
			EntityName:  "DigiKey Electronics",
			Action:      "Added Supplier Terms",
			Description: "Updated procurement terms and lead times for DigiKey.",
			BadgeColor:  "indigo",
		},
	}
	db.Create(&activities)
	}

	log.Println("[INFO] Master Data seeding completed successfully.")
	SeedPhase2(db)
}

func SeedPhase2(db *gorm.DB) {
	// Seed Permissions for Versions & Revisions if not present
	var permCount int64
	db.Model(&model.Permission{}).Where("module IN ('versions', 'revisions')").Count(&permCount)
	if permCount == 0 {
		phase2Perms := []model.Permission{
			{ID: "PERM-VER-VIEW", Module: "versions", Action: "view", Code: "versions.view", Description: "View product versions"},
			{ID: "PERM-VER-CREATE", Module: "versions", Action: "create", Code: "versions.create", Description: "Create product versions"},
			{ID: "PERM-VER-EDIT", Module: "versions", Action: "edit", Code: "versions.edit", Description: "Edit product versions"},
			{ID: "PERM-VER-DELETE", Module: "versions", Action: "delete", Code: "versions.delete", Description: "Delete product versions"},

			{ID: "PERM-REV-VIEW", Module: "revisions", Action: "view", Code: "revisions.view", Description: "View hardware revisions"},
			{ID: "PERM-REV-CREATE", Module: "revisions", Action: "create", Code: "revisions.create", Description: "Create hardware revisions"},
			{ID: "PERM-REV-EDIT", Module: "revisions", Action: "edit", Code: "revisions.edit", Description: "Edit hardware revisions"},
			{ID: "PERM-REV-DELETE", Module: "revisions", Action: "delete", Code: "revisions.delete", Description: "Delete hardware revisions"},
		}
		db.Create(&phase2Perms)

		// Grant to Super Admin & R&D Manager
		for _, p := range phase2Perms {
			db.Create(&model.RolePermission{RoleID: "ROLE-ADMIN", PermissionID: p.ID})
			db.Create(&model.RolePermission{RoleID: "ROLE-RNDM", PermissionID: p.ID})
			if p.Action != "delete" {
				db.Create(&model.RolePermission{RoleID: "ROLE-HWENG", PermissionID: p.ID})
			}
			if p.Action == "view" {
				db.Create(&model.RolePermission{RoleID: "ROLE-FWENG", PermissionID: p.ID})
				db.Create(&model.RolePermission{RoleID: "ROLE-QAENG", PermissionID: p.ID})
				db.Create(&model.RolePermission{RoleID: "ROLE-PROCM", PermissionID: p.ID})
			}
		}
	}

	// Ensure PRD-2024-001 exists and is active
	var p1 model.Product
	if err := db.Unscoped().Where("id = ? OR code = ?", "PRD-2024-001", "PRD-2024-001").First(&p1).Error; err != nil {
		p1 = model.Product{
			ID:                "PRD-2024-001",
			Code:              "PRD-2024-001",
			Name:              "Smart Edge Industrial Gateway (GW-400X)",
			Category:          "Edge Gateways",
			CategoryID:        strPtr("CAT-GW"),
			Status:            "Active",
			Description:       "Enterprise quad-core ARM Cortex gateway with dual Gigabit Ethernet, RS-485 Modbus interface, and integrated 4G LTE fallback. Designed for smart factory IoT deployments.",
			LongDescription:   "The GW-400X is an enterprise-grade IoT edge computer capable of running containerized telemetry processing workloads at the edge. Features robust surge protection on all I/O, dual power inputs (9-36V DC), and an IP50 metal enclosure suitable for DIN rail or wall mounting.",
			Image:             "https://images.unsplash.com/photo-1518770660439-4636190af475?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Dr. Sarah Jenkins",
			TargetMarket:      "Industrial Automation & Smart Factories",
			PowerRating:       "12-24V DC / 15W Max",
			IngressProtection: "IP50 (DIN-rail mountable)",
			OperatingTemp:     "-40°C to +85°C",
			Certifications:    `["CE", "FCC Part 15B", "RoHS 3", "UL 61010-1"]`,
			ComponentsCount:   142,
		}
		db.Create(&p1)
	} else if p1.DeletedAt.Valid {
		db.Unscoped().Model(&p1).Update("deleted_at", nil)
	}

	// Seed Versions if empty
	var verCount int64
	db.Model(&model.ProductVersion{}).Count(&verCount)
	if verCount > 0 {
		return
	}

	log.Println("[INFO] Seeding initial Phase 2 Product Versions & Hardware Revisions...")

	versions := []model.ProductVersion{
		{
			ID:              "VER-001",
			ProductID:       "PRD-2024-001",
			ProductCode:     "PRD-2024-001",
			ProductName:     "IoT Environmental Gateway GW-400X",
			VersionNumber:   "1.0.0",
			VersionName:     "Initial EVT Engineering Release",
			Status:          "Archived",
			CurrentRevision: "REV-A",
			RevisionsCount:  1,
			Owner:           "Alex Chen",
			OwnerAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2025-09-10"),
			Description:     "First production prototype of the environmental edge gateway with basic LoRaWAN and Modbus RTU telemetry.",
			ReleaseNotes:    "1. Initial 4-layer PCB design with STM32F4 core.\n2. Integrated SX1262 LoRa transceiver.\n3. Basic RS-485 isolation barrier.",
			Documents:       `[{"name":"SCH-GW400X-v1.0.0.pdf","size":"2.4 MB","type":"PDF"},{"name":"GERBER-GW400X-v1.0.0.zip","size":"8.1 MB","type":"ZIP"}]`,
			CreatedAt:       time.Now().AddDate(0, -11, 0),
			UpdatedAt:       time.Now().AddDate(0, -11, 0),
		},
		{
			ID:              "VER-002",
			ProductID:       "PRD-2024-001",
			ProductCode:     "PRD-2024-001",
			ProductName:     "IoT Environmental Gateway GW-400X",
			VersionNumber:   "1.2.0",
			VersionName:     "DVT Industrial Field Release",
			Status:          "Deprecated",
			CurrentRevision: "REV-B",
			RevisionsCount:  1,
			Owner:           "Alex Chen",
			OwnerAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-02-14"),
			Description:     "Field hardened design with wide DC input power regulation (9-36V) and dual SIM cellular modem slot.",
			ReleaseNotes:    "1. Replaced linear LDO with high-efficiency buck regulator.\n2. Added hardware watchdog IC.\n3. Enhanced ESD protection to IEC 61000-4-2 Level 4.",
			BaseVersionID:   strPtr("VER-001"),
			Documents:       `[{"name":"SCH-GW400X-v1.2.0.pdf","size":"2.9 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, -6, 0),
			UpdatedAt:       time.Now().AddDate(0, -6, 0),
		},
		{
			ID:              "VER-003",
			ProductID:       "PRD-2024-001",
			ProductCode:     "PRD-2024-001",
			ProductName:     "IoT Environmental Gateway GW-400X",
			VersionNumber:   "2.0.0",
			VersionName:     "Major Edge Compute Architecture",
			Status:          "Active",
			CurrentRevision: "REV-A",
			RevisionsCount:  1,
			Owner:           "Alex Chen",
			OwnerAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-05-30"),
			Description:     "Major revision transitioning to STM32MP157 dual-core Linux MPU with on-board eMMC and Gigabit Ethernet.",
			ReleaseNotes:    "1. Upgraded CPU architecture to STM32MP157.\n2. Added 512MB DDR3L RAM and 4GB eMMC storage.\n3. Integrated Gigabit Ethernet PHY.",
			BaseVersionID:   strPtr("VER-002"),
			Documents:       `[{"name":"SCH-GW400X-v2.0.0.pdf","size":"3.8 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, -3, 0),
			UpdatedAt:       time.Now().AddDate(0, -3, 0),
		},
		{
			ID:              "VER-004",
			ProductID:       "PRD-2024-001",
			ProductCode:     "PRD-2024-001",
			ProductName:     "IoT Environmental Gateway GW-400X",
			VersionNumber:   "2.1.0",
			VersionName:     "LTE-M / NB-IoT Cellular Edge Hub",
			Status:          "Active",
			CurrentRevision: "REV-B",
			RevisionsCount:  2,
			Owner:           "Alex Chen",
			OwnerAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-08-15"),
			Description:     "Current production version with certified Quectel BG95-M3 cellular modem, ultra-low power sleep states, and IP67 enclosure mounting.",
			ReleaseNotes:    "1. Integrated Quectel BG95-M3 LTE-M / NB-IoT.\n2. Power consumption reduced by 22%.\n3. Certified for CE and FCC Part 15 Class B.",
			BaseVersionID:   strPtr("VER-003"),
			Documents:       `[{"name":"SCH-GW400X-v2.1.0.pdf","size":"4.1 MB","type":"PDF"},{"name":"GERBER-GW400X-v2.1.0-REVB.zip","size":"9.4 MB","type":"ZIP"},{"name":"FCC-CERTIFICATE-v2.1.0.pdf","size":"1.8 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, -1, 0),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              "VER-005",
			ProductID:       "PRD-2024-002",
			ProductCode:     "PRD-2024-002",
			ProductName:     "Smart Air Quality Monitor AQM-200",
			VersionNumber:   "1.0.0",
			VersionName:     "Commercial Indoor Monitor Release",
			Status:          "Active",
			CurrentRevision: "REV-C",
			RevisionsCount:  3,
			Owner:           "Sarah Connor",
			OwnerAvatar:     "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-03-20"),
			Description:     "Precision indoor air monitor with Sensirion SCD40 CO2, SHT40 Temperature/Humidity, and PM2.5 laser scattering.",
			ReleaseNotes:    "1. High precision NDIR optical CO2 sensor.\n2. 2.13 inch E-ink display.\n3. USB Type-C and internal LiPo battery charging.",
			Documents:       `[{"name":"SCH-AQM200-v1.0.0.pdf","size":"2.1 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, -5, 0),
			UpdatedAt:       time.Now().AddDate(0, -5, 0),
		},
		{
			ID:              "VER-006",
			ProductID:       "PRD-2024-002",
			ProductCode:     "PRD-2024-002",
			ProductName:     "Smart Air Quality Monitor AQM-200",
			VersionNumber:   "1.1.0",
			VersionName:     "VOC & Formaldehyde Sensor Expansion",
			Status:          "Draft",
			CurrentRevision: "REV-A",
			RevisionsCount:  1,
			Owner:           "Sarah Connor",
			OwnerAvatar:     "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-09-30"),
			Description:     "Next minor iteration adding SGP40 multi-pixel gas sensor for Total Volatile Organic Compounds (TVOC).",
			ReleaseNotes:    "1. Integrated SGP40 MOx gas sensor.\n2. Calibration algorithm for TVOC index.",
			BaseVersionID:   strPtr("VER-005"),
			Documents:       `[{"name":"SCH-AQM200-v1.1.0-DRAFT.pdf","size":"2.5 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, 0, -10),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              "VER-007",
			ProductID:       "PRD-2024-003",
			ProductCode:     "PRD-2024-003",
			ProductName:     "Industrial Vibration & Temperature Sensor ISN-100",
			VersionNumber:   "1.0.0",
			VersionName:     "Explosion-Proof ATEX Sensor Node",
			Status:          "Active",
			CurrentRevision: "REV-B",
			RevisionsCount:  2,
			Owner:           "David Kim",
			OwnerAvatar:     "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80",
			ReleaseDate:     strPtr("2026-04-18"),
			Description:     "ATEX Zone 1 certified 3-axis piezoresistive vibration monitoring node with BLE 5.2 and 10-year battery life.",
			ReleaseNotes:    "1. 3-axis 10kHz vibration sensing.\n2. Potting epoxy sealed CNC aluminum enclosure.\n3. Magnetic contactless activation.",
			Documents:       `[{"name":"SCH-ISN100-v1.0.0.pdf","size":"1.9 MB","type":"PDF"},{"name":"ATEX-ZONE1-CERTIFICATE.pdf","size":"3.1 MB","type":"PDF"}]`,
			CreatedAt:       time.Now().AddDate(0, -4, 0),
			UpdatedAt:       time.Now().AddDate(0, -4, 0),
		},
	}
	db.Create(&versions)

	revisions := []model.HardwareRevision{
		{
			ID:                "REV-001",
			ProductVersionID:  "VER-004",
			VersionNumber:     "2.1.0",
			ProductCode:       "PRD-2024-001",
			ProductName:       "IoT Environmental Gateway GW-400X",
			Code:              "REV-A",
			Name:              "Initial PCB Prototype Layout",
			Status:            "Superseded",
			PCBRevision:       "PCB-v2.1-A",
			SchematicRevision: "SCH-v2.1-01",
			Engineer:          "Marcus Vance",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-07-15"),
			ChangeSummary:     "Initial prototype schematic and 6-layer high-speed layout for STM32MP157 and BG95 cellular module.",
			ChangesList:       `["Initial 6-layer stackup design with controlled 50Ω impedance.","Placed STM32MP157 MPU, DDR3L RAM, and eMMC in tight cluster.","Integrated Quectel BG95-M3 LGA footprint."]`,
			Documents:         `[{"name":"SCH-GW400X-REV-A.pdf","size":"3.9 MB","type":"PDF"},{"name":"PCB-GERBER-GW400X-REV-A.zip","size":"8.8 MB","type":"ZIP"}]`,
			CreatedAt:         time.Now().AddDate(0, -2, 0),
			UpdatedAt:         time.Now().AddDate(0, -2, 0),
		},
		{
			ID:                "REV-002",
			ProductVersionID:  "VER-004",
			VersionNumber:     "2.1.0",
			ProductCode:       "PRD-2024-001",
			ProductName:       "IoT Environmental Gateway GW-400X",
			Code:              "REV-B",
			Name:              "RF Power Filtering & RS485 Isolation Improvement",
			Status:            "Released",
			PCBRevision:       "PCB-v2.1-B",
			SchematicRevision: "SCH-v2.1-02",
			Engineer:          "Marcus Vance",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-08-15"),
			ChangeSummary:     "Optimized cellular RF antenna matching network, increased power plane copper pour around DC-DC converters, and added ESD suppression arrays.",
			ChangesList:       `["Added pi-filter matching components (L12, C45, C46) to optimize cellular LTE-M antenna VSWR < 1.8.","Separated analog ground (AGND) and digital ground (DGND).","Upgraded RS-485 isolation barrier with TI ISO1410 galvanic isolator."]`,
			Documents:         `[{"name":"SCH-GW400X-REV-B.pdf","size":"4.1 MB","type":"PDF"},{"name":"PCB-GERBER-GW400X-REV-B.zip","size":"9.4 MB","type":"ZIP"}]`,
			CreatedAt:         time.Now().AddDate(0, -1, 0),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                "REV-003",
			ProductVersionID:  "VER-005",
			VersionNumber:     "1.0.0",
			ProductCode:       "PRD-2024-002",
			ProductName:       "Smart Air Quality Monitor AQM-200",
			Code:              "REV-A",
			Name:              "Air Quality Optical Chamber Layout",
			Status:            "Superseded",
			PCBRevision:       "PCB-AQM-A0",
			SchematicRevision: "SCH-AQM-01",
			Engineer:          "Sarah Connor",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-01-28"),
			ChangeSummary:     "First engineering prototype with optical laser dust sensor and discrete NDIR CO2 sensor header.",
			ChangesList:       `["Dedicated 3.3V ultra-low-noise analog power rail.","E-ink display SPI interface pinout configuration."]`,
			Documents:         `[{"name":"SCH-AQM200-REV-A.pdf","size":"1.8 MB","type":"PDF"}]`,
			CreatedAt:         time.Now().AddDate(0, -7, 0),
			UpdatedAt:         time.Now().AddDate(0, -7, 0),
		},
		{
			ID:                "REV-004",
			ProductVersionID:  "VER-005",
			VersionNumber:     "1.0.0",
			ProductCode:       "PRD-2024-002",
			ProductName:       "Smart Air Quality Monitor AQM-200",
			Code:              "REV-B",
			Name:              "SCD40 Photoacoustic Integration",
			Status:            "Superseded",
			PCBRevision:       "PCB-AQM-B0",
			SchematicRevision: "SCH-AQM-02",
			Engineer:          "Sarah Connor",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-02-25"),
			ChangeSummary:     "Switched CO2 sensor from bulky NDIR module to miniature photoacoustic Sensirion SCD40 SMD.",
			ChangesList:       `["Reduced PCB area by 34% by moving to SCD40 SMD footprint.","Added thermal relief barrier to isolate ambient MCU heat."]`,
			Documents:         `[{"name":"SCH-AQM200-REV-B.pdf","size":"2.0 MB","type":"PDF"}]`,
			CreatedAt:         time.Now().AddDate(0, -6, 0),
			UpdatedAt:         time.Now().AddDate(0, -6, 0),
		},
		{
			ID:                "REV-005",
			ProductVersionID:  "VER-005",
			VersionNumber:     "1.0.0",
			ProductCode:       "PRD-2024-002",
			ProductName:       "Smart Air Quality Monitor AQM-200",
			Code:              "REV-C",
			Name:              "Mass Production Release & Enclosure Fit",
			Status:            "Released",
			PCBRevision:       "PCB-AQM-C0",
			SchematicRevision: "SCH-AQM-03",
			Engineer:          "Sarah Connor",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-03-20"),
			ChangeSummary:     "Final mass production layout with rounded board corners, snap-fit enclosure mounting holes, and test points for ICT test bed.",
			ChangesList:       `["Added 28 gold-plated test pads for automated factory testing.","Adjusted board corner radius to R=3.5mm for enclosure fit."]`,
			Documents:         `[{"name":"SCH-AQM200-REV-C.pdf","size":"2.1 MB","type":"PDF"}]`,
			CreatedAt:         time.Now().AddDate(0, -5, 0),
			UpdatedAt:         time.Now().AddDate(0, -5, 0),
		},
		{
			ID:                "REV-006",
			ProductVersionID:  "VER-007",
			VersionNumber:     "1.0.0",
			ProductCode:       "PRD-2024-003",
			ProductName:       "Industrial Vibration Sensor ISN-100",
			Code:              "REV-A",
			Name:              "Piezoresistive Analog Front-End",
			Status:            "Superseded",
			PCBRevision:       "PCB-ISN-A",
			SchematicRevision: "SCH-ISN-01",
			Engineer:          "David Kim",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-03-10"),
			ChangeSummary:     "Analog instrumentation amplifier front-end for 10kHz piezoresistive vibration transducer.",
			ChangesList:       `["High CMRR differential amplifier stage with low-noise op-amp.","BLE 5.2 ceramic chip antenna layout."]`,
			Documents:         `[{"name":"SCH-ISN100-REV-A.pdf","size":"1.7 MB","type":"PDF"}]`,
			CreatedAt:         time.Now().AddDate(0, -5, 0),
			UpdatedAt:         time.Now().AddDate(0, -5, 0),
		},
		{
			ID:                "REV-007",
			ProductVersionID:  "VER-007",
			VersionNumber:     "1.0.0",
			ProductCode:       "PRD-2024-003",
			ProductName:       "Industrial Vibration Sensor ISN-100",
			Code:              "REV-B",
			Name:              "ATEX Intrinsic Safety Barrier & Potting Alignment",
			Status:            "Released",
			PCBRevision:       "PCB-ISN-B",
			SchematicRevision: "SCH-ISN-02",
			Engineer:          "David Kim",
			ApprovedBy:        strPtr("Alex Chen"),
			ReleaseDate:       strPtr("2026-04-18"),
			ChangeSummary:     "Added current-limiting resistors and triple redundant Zener barrier diodes for ATEX Zone 1 explosion safety compliance.",
			ChangesList:       `["Added intrinsic safety current-limiting resistor array.","Enforced minimum 2.0mm clearance between traces."]`,
			Documents:         `[{"name":"SCH-ISN100-REV-B.pdf","size":"1.9 MB","type":"PDF"},{"name":"ATEX-SAFETY-CALCULATIONS.pdf","size":"2.8 MB","type":"PDF"}]`,
			CreatedAt:         time.Now().AddDate(0, -4, 0),
			UpdatedAt:         time.Now().AddDate(0, -4, 0),
		},
	}
	db.Create(&revisions)

	log.Println("[INFO] Phase 2 Product Versions & Hardware Revisions seeded successfully.")
}

func SeedPhase2C(db *gorm.DB) {
	// 1. Ensure existing products have ProductType = 'MANUFACTURE'
	db.Model(&model.Product{}).Where("product_type IS NULL OR product_type = ''").Update("product_type", model.ProductTypeManufacture)

	// 2. Seed Trading Products if not present
	var trdCount int64
	db.Model(&model.Product{}).Where("code = 'PRD-TRD-001'").Count(&trdCount)
	if trdCount == 0 {
		trdProducts := []model.Product{
			{
				ID:                "PRD-TRD-001",
				Code:              "PRD-TRD-001",
				Name:              "Commercial Multi-Parameter Weather Station (WXT536)",
				ProductType:       model.ProductTypeTrading,
				Category:          "Sensors & Nodes",
				Status:            "Active",
				Description:       "Externally sourced high-precision all-in-one ultrasonic meteorological sensor.",
				LongDescription:   "Ultrasonic solid-state weather station measuring wind speed, wind direction, liquid precipitation, barometric pressure, temperature, and relative humidity with zero moving parts. Certified for industrial meteorological networks and harsh coastal deployments.",
				Image:             "https://images.unsplash.com/photo-1590055531615-f16d36ffe8ec?w=400&auto=format&fit=crop&q=80",
				ProjectLead:       "David Tanaka",
				TargetMarket:      "Meteorology & Smart Agriculture",
				PowerRating:       "9-32 VDC, 0.5W",
				IngressProtection: "IP66 / NEMA 4X",
				OperatingTemp:     "-52°C to +60°C",
				Certifications:    `["CE","FCC","WMO-No. 8","RoHS"]`,
				ComponentsCount:   0,
				CreatedAt:         time.Now().AddDate(0, -3, 0),
				UpdatedAt:         time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                "PRD-TRD-002",
				Code:              "PRD-TRD-002",
				Name:              "Monocrystalline Solar Power Supply Unit (SPU-60W)",
				ProductType:       model.ProductTypeTrading,
				Category:          "Power Systems",
				Status:            "Active",
				Description:       "Ready-made commercial solar power package with MPPT charger and LiFePO4 battery.",
				LongDescription:   "Industrial turnkey solar power package including 60W monocrystalline PV panel, intelligent MPPT charge controller, 24V 20Ah LiFePO4 battery pack, and IP67 weather-resistant enclosure with mast mounting hardware.",
				Image:             "https://images.unsplash.com/photo-1509391365360-2e959784a276?w=400&auto=format&fit=crop&q=80",
				ProjectLead:       "David Tanaka",
				TargetMarket:      "Off-Grid Remote Infrastructure",
				PowerRating:       "24V DC / 60W Peak",
				IngressProtection: "IP67",
				OperatingTemp:     "-20°C to +60°C",
				Certifications:    `["CE","UL1703","UN38.3","IEC61215"]`,
				ComponentsCount:   0,
				CreatedAt:         time.Now().AddDate(0, -2, 0),
				UpdatedAt:         time.Now().AddDate(0, -2, 0),
			},
		}
		db.Create(&trdProducts)

		// Seed Trading Details
		trdDetails := []model.TradingProductDetail{
			{
				ID:                     "TRD-DET-001",
				ProductID:              "PRD-TRD-001",
				ManufacturerID:         strPtr("MFG-BOSCH"),
				ManufacturerName:       "Vaisala Instruments Oyj",
				SupplierID:             strPtr("SUP-DIGI"),
				SupplierName:           "DigiKey Electronics",
				ManufacturerPartNumber: "WXT536-6A1B1A",
				SupplierPartNumber:     "DK-WXT536-ND",
				PurchasePrice:          1250.00,
				SellingPrice:           1890.00,
				Currency:               "USD",
				LeadTimeDays:           14,
				WarrantyInformation:    "24 Months Factory Limited Warranty",
				Notes:                  "Includes RS-485 Modbus ASCII/RTU & SDI-12 digital outputs. Pre-calibrated calibration certificate included in shipment.",
				CreatedAt:              time.Now().AddDate(0, -3, 0),
				UpdatedAt:              time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                     "TRD-DET-002",
				ProductID:              "PRD-TRD-002",
				ManufacturerID:         strPtr("MFG-TI"),
				ManufacturerName:       "Mean Well Enterprises Co., Ltd.",
				SupplierID:             strPtr("SUP-MOUS"),
				SupplierName:           "Mouser Electronics",
				ManufacturerPartNumber: "SPU-60W-24V-IP67",
				SupplierPartNumber:     "MSR-SPU60W24",
				PurchasePrice:          280.00,
				SellingPrice:           420.00,
				Currency:               "USD",
				LeadTimeDays:           7,
				WarrantyInformation:    "36 Months Standard Manufacturer Warranty",
				Notes:                  "High-efficiency MPPT charge controller with built-in battery temperature compensation and low-temperature cutoff protection.",
				CreatedAt:              time.Now().AddDate(0, -2, 0),
				UpdatedAt:              time.Now().AddDate(0, -2, 0),
			},
		}
		db.Create(&trdDetails)
	}

	// 3. Seed Project Product if not present
	var prjCount int64
	db.Model(&model.Product{}).Where("code = 'PRD-PRJ-001'").Count(&prjCount)
	if prjCount == 0 {
		prjProduct := model.Product{
			ID:                "PRD-PRJ-001",
			Code:              "PRD-PRJ-001",
			Name:              "Environmental Monitoring Station Package (EMS-2000)",
			ProductType:       model.ProductTypeProject,
			Category:          "Integrated Solutions",
			Status:            "Active",
			Description:       "Turnkey environmental station package integrating R&D telemetry, commercial weather sensors, and solar power.",
			LongDescription:   "Complete turnkey field-deployable solution designed for smart municipal air monitoring, industrial emission monitoring, and meteorological observation networks. Houses core R&D edge gateway and sensor modules, paired with commercial ultrasonic weather sensors and solar power supplies.",
			Image:             "https://images.unsplash.com/photo-1542601906990-b4d3fb778b09?w=400&auto=format&fit=crop&q=80",
			ProjectLead:       "Marcus Vance",
			TargetMarket:      "Smart City & Environmental Compliance",
			PowerRating:       "Solar / Battery (24V DC)",
			IngressProtection: "IP67 System Level",
			OperatingTemp:     "-30°C to +60°C",
			Certifications:    `["CE","FCC","NEMA 4X","RoHS"]`,
			ComponentsCount:   0,
			CreatedAt:         time.Now().AddDate(0, -2, 0),
			UpdatedAt:         time.Now().AddDate(0, -2, 0),
		}
		db.Create(&prjProduct)

		// Seed Project Items Hierarchy
		prjItems := []model.ProjectItem{
			// Sub Assembly 1: Core Telemetry & Acquisition Unit
			{
				ID:               "PRJ-ITM-001",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     nil,
				ItemType:         model.ProjectItemTypeSubAssembly,
				Name:             "Core Telemetry & Acquisition Unit",
				SortOrder:        1,
				Notes:            "Primary ingress-protected enclosure housing internal electronics and transceivers.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-002",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-001"),
				ItemType:         model.ProjectItemTypeProduct,
				Name:             "Smart Edge Gateway GW-400X",
				ProductID:        strPtr("PRD-2024-001"),
				ProductCode:      "PRD-2024-001",
				ProductType:      model.ProductTypeManufacture,
				Quantity:         1.0,
				UnitName:         "pcs",
				SortOrder:        1,
				Notes:            "References internal Manufactured Product Master PRD-2024-001. Handles LTE Cat-M1 and Modbus polling.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-003",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-001"),
				ItemType:         model.ProjectItemTypeProduct,
				Name:             "Smart Air Quality Monitor AQM-200",
				ProductID:        strPtr("PRD-2024-002"),
				ProductCode:      "PRD-2024-002",
				ProductType:      model.ProductTypeManufacture,
				Quantity:         1.0,
				UnitName:         "pcs",
				SortOrder:        2,
				Notes:            "References internal R&D Product Master PRD-2024-002 for particulate and CO2 sensing.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-004",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-001"),
				ItemType:         model.ProjectItemTypeComponent,
				Name:             "Industrial RS-485 Shielded Communication Harness (3m)",
				ComponentID:      strPtr("CMP-004"),
				ComponentMPN:     "SN65HVD72DR",
				Quantity:         2.0,
				UnitName:         "pcs",
				SortOrder:        3,
				Notes:            "Custom outdoor grade cable harness linking gateway to sensor heads.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},

			// Sub Assembly 2: Meteorological & Environmental Sourcing Bundle
			{
				ID:               "PRJ-ITM-005",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     nil,
				ItemType:         model.ProjectItemTypeSubAssembly,
				Name:             "Meteorological Sourcing Package",
				SortOrder:        2,
				Notes:            "Commercial external sensor assembly.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-006",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-005"),
				ItemType:         model.ProjectItemTypeProduct,
				Name:             "Commercial Weather Station (WXT536)",
				ProductID:        strPtr("PRD-TRD-001"),
				ProductCode:      "PRD-TRD-001",
				ProductType:      model.ProductTypeTrading,
				Quantity:         1.0,
				UnitName:         "pcs",
				SortOrder:        1,
				Notes:            "References Trading Product PRD-TRD-001. Connected via RS-485 Modbus RTU.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},

			// Sub Assembly 3: Autonomous Power Infrastructure
			{
				ID:               "PRJ-ITM-007",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     nil,
				ItemType:         model.ProjectItemTypeSubAssembly,
				Name:             "Autonomous Solar Power Infrastructure",
				SortOrder:        3,
				Notes:            "External solar power generation and battery storage system.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-008",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-007"),
				ItemType:         model.ProjectItemTypeProduct,
				Name:             "Monocrystalline Solar Power Supply Unit (SPU-60W)",
				ProductID:        strPtr("PRD-TRD-002"),
				ProductCode:      "PRD-TRD-002",
				ProductType:      model.ProductTypeTrading,
				Quantity:         1.0,
				UnitName:         "pcs",
				SortOrder:        1,
				Notes:            "References Trading Product PRD-TRD-002. Delivers 24V regulated DC power.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
			{
				ID:               "PRJ-ITM-009",
				ProjectProductID: "PRD-PRJ-001",
				ParentItemID:     strPtr("PRJ-ITM-007"),
				ItemType:         model.ProjectItemTypeComponent,
				Name:             "Galvanized Steel Mast Mounting Brackets",
				ComponentID:      nil,
				ComponentMPN:     "HD-MBKT-60",
				Quantity:         4.0,
				UnitName:         "pcs",
				SortOrder:        2,
				Notes:            "Heavy-duty pole mounting brackets with stainless steel U-bolts.",
				Status:           "Active",
				CreatedAt:        time.Now().AddDate(0, -2, 0),
				UpdatedAt:        time.Now().AddDate(0, -2, 0),
			},
		}
		db.Create(&prjItems)
	}

	log.Println("[INFO] Phase 2C Trading Products, Project Products & Hierarchies seeded successfully.")
}

func SeedPhase3(db *gorm.DB) {
	// 1. Backfill Component Estimated Unit Costs
	componentCosts := map[string]struct {
		Cost   float64
		Source string
	}{
		"CMP-001": {Cost: 8.50, Source: "MANUAL_ESTIMATE"},     // STM32F407VGT6
		"CMP-002": {Cost: 3.80, Source: "MANUAL_ESTIMATE"},     // ESP32-WROOM-32D
		"CMP-003": {Cost: 4.85, Source: "SUPPLIER_REFERENCE"},  // BME680
		"CMP-004": {Cost: 1.25, Source: "SUPPLIER_REFERENCE"},  // SN65HVD72DR
		"CMP-005": {Cost: 1.90, Source: "SUPPLIER_REFERENCE"},  // LM2596S-5.0
		"CMP-006": {Cost: 0.15, Source: "MANUAL_ESTIMATE"},     // 100uF 25V Capacitor
	}

	for id, data := range componentCosts {
		db.Model(&model.Component{}).Where("id = ?", id).Updates(map[string]interface{}{
			"estimated_unit_cost": data.Cost,
			"currency":            "USD",
			"cost_source":         data.Source,
		})
	}

	// 2. Seed Permissions for Engineering BOM if not present
	var permCount int64
	db.Model(&model.Permission{}).Where("module = 'bom'").Count(&permCount)
	if permCount == 0 {
		bomPerms := []model.Permission{
			{ID: "PERM-BOM-VIEW", Module: "bom", Action: "view", Code: "bom.view", Description: "View engineering BOMs and cost rollups"},
			{ID: "PERM-BOM-CREATE", Module: "bom", Action: "create", Code: "bom.create", Description: "Create engineering BOMs and add items"},
			{ID: "PERM-BOM-EDIT", Module: "bom", Action: "edit", Code: "bom.edit", Description: "Edit BOM items, quantities, and reference designators"},
			{ID: "PERM-BOM-DELETE", Module: "bom", Action: "delete", Code: "bom.delete", Description: "Delete BOM items and structures"},
		}
		db.Create(&bomPerms)

		// Grant to Super Admin, Hardware Engineer, and R&D Manager
		for _, p := range bomPerms {
			db.Create(&model.RolePermission{RoleID: "ROLE-ADMIN", PermissionID: p.ID})
			db.Create(&model.RolePermission{RoleID: "ROLE-HWENG", PermissionID: p.ID})
			db.Create(&model.RolePermission{RoleID: "ROLE-RNDMGR", PermissionID: p.ID})
		}
	}

	// 3. Seed Engineering BOM for REV-001 (PRD-2024-001 v2.1.0 REV-A)
	var bomCount int64
	db.Model(&model.EngineeringBOM{}).Where("hardware_revision_id = 'REV-001'").Count(&bomCount)
	if bomCount == 0 {
		bom := model.EngineeringBOM{
			ID:                    "BOM-REV-001",
			HardwareRevisionID:    "REV-001",
			BOMCode:               "BOM-GW400X-2.1.0-A",
			Name:                  "Smart Edge Industrial Gateway Main Board Engineering BOM",
			Description:           "Production engineering Bill of Materials for REV-A gateway hardware architecture with multi-level subsystems.",
			Status:                model.BOMStatusActive,
			TotalCost:             22.45,
			Currency:              "USD",
			TargetMargin:          35.00,
			EstimatedSellingPrice: 34.54,
			Notes:                 "Initial released baseline including microcontroller, power conditioning, and RS485 isolated transceiver stages.",
			CreatedBy:             "Alex Chen",
			CreatedAt:             time.Now().AddDate(0, -3, 0),
			UpdatedAt:             time.Now().AddDate(0, -3, 0),
		}
		db.Create(&bom)

		// Seed Multi-Level BOM Items
		bomItems := []model.BOMItem{
			// SUB-ASSEMBLY 1: Core Processing Subsystem
			{
				ID:                   "BOM-ITM-001",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         nil,
				ItemType:             model.BOMItemTypeSubAssembly,
				Name:                 "Core Processing & Microcontroller Subsystem",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             9.10,
				LineCost:             9.10,
				Position:             1,
				Notes:                "High-performance STM32 ARM Cortex-M4 MCU with high-stability decoupling network.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-002",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-001"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "STM32F407VGT6 ARM Cortex-M4 MCU (168MHz, 1MB Flash)",
				ComponentID:          strPtr("CMP-001"),
				ComponentPartNumber:  "STM32F407VGT6",
				ComponentCategory:    "Microcontrollers",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             8.50,
				LineCost:             8.50,
				ReferenceDesignators: "U1",
				Position:             1,
				Notes:                "Main system processor running FreeRTOS telemetry firmware.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-003",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-001"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "100uF 25V SMD Decoupling Capacitors",
				ComponentID:          strPtr("CMP-006"),
				ComponentPartNumber:  "EEE-FK1E101P",
				ComponentCategory:    "Passives",
				Quantity:             4.0,
				UnitCode:             "PCS",
				UnitCost:             0.15,
				LineCost:             0.60,
				ReferenceDesignators: "C1, C2, C3, C4",
				Position:             2,
				Notes:                "Low-ESR supply decoupling capacitors placed adjacent to MCU VDD pins.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},

			// SUB-ASSEMBLY 2: Power Stage
			{
				ID:                   "BOM-ITM-004",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         nil,
				ItemType:             model.BOMItemTypeSubAssembly,
				Name:                 "Wide-Range DC-DC Power Stage",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             2.20,
				LineCost:             2.20,
				Position:             2,
				Notes:                "Regulates 9-36V industrial DC input down to 5.0V system rail.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-005",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-004"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "LM2596S-5.0 Step-Down Buck Regulator (3A Output)",
				ComponentID:          strPtr("CMP-005"),
				ComponentPartNumber:  "LM2596S-5.0/NOPB",
				ComponentCategory:    "Power Management",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             1.90,
				LineCost:             1.90,
				ReferenceDesignators: "U2",
				Position:             1,
				Notes:                "Primary step-down converter with thermal shutdown protection.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-006",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-004"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "Bulk Storage Capacitors (100uF 25V)",
				ComponentID:          strPtr("CMP-006"),
				ComponentPartNumber:  "EEE-FK1E101P",
				ComponentCategory:    "Passives",
				Quantity:             2.0,
				UnitCode:             "PCS",
				UnitCost:             0.15,
				LineCost:             0.30,
				ReferenceDesignators: "C10, C11",
				Position:             2,
				Notes:                "Buck converter input/output smoothing capacitors.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},

			// SUB-ASSEMBLY 3: Industrial Communications & Sensor Interface
			{
				ID:                   "BOM-ITM-007",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         nil,
				ItemType:             model.BOMItemTypeSubAssembly,
				Name:                 "Industrial Communications & Sensor Interface",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             11.15,
				LineCost:             11.15,
				Position:             3,
				Notes:                "RS-485 Modbus interface, Wi-Fi/BLE transceiver module, and onboard BME680.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-008",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-007"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "SN65HVD72DR RS-485 Isolated Transceiver",
				ComponentID:          strPtr("CMP-004"),
				ComponentPartNumber:  "SN65HVD72DR",
				ComponentCategory:    "Interface ICs",
				Quantity:             2.0,
				UnitCode:             "PCS",
				UnitCost:             1.25,
				LineCost:             2.50,
				ReferenceDesignators: "U3, U4",
				Position:             1,
				Notes:                "Dual isolated RS-485 transceiver channels with IEC ESD protection.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-009",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-007"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "ESP32-WROOM-32D Wi-Fi / BLE 4.2 Module",
				ComponentID:          strPtr("CMP-002"),
				ComponentPartNumber:  "ESP32-WROOM-32D",
				ComponentCategory:    "Wireless Modules",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             3.80,
				LineCost:             3.80,
				ReferenceDesignators: "MOD1",
				Position:             2,
				Notes:                "Wireless provisioning and local dashboard BLE beacon interface.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
			{
				ID:                   "BOM-ITM-010",
				EngineeringBOMID:     "BOM-REV-001",
				ParentItemID:         strPtr("BOM-ITM-007"),
				ItemType:             model.BOMItemTypeComponent,
				Name:                 "BME680 Environmental Gas & Temp Sensor",
				ComponentID:          strPtr("CMP-003"),
				ComponentPartNumber:  "BME680",
				ComponentCategory:    "Sensors",
				Quantity:             1.0,
				UnitCode:             "PCS",
				UnitCost:             4.85,
				LineCost:             4.85,
				ReferenceDesignators: "U5",
				Position:             3,
				Notes:                "Internal enclosure humidity, temperature, and VOC air quality monitor.",
				CreatedAt:            time.Now().AddDate(0, -3, 0),
				UpdatedAt:            time.Now().AddDate(0, -3, 0),
			},
		}
		db.Create(&bomItems)
	}

	log.Println("[INFO] Phase 3 Engineering BOM & Cost Foundation seeded successfully.")
}

// SeedMultiCurrency seeds exchange rates and enriches trading products with multi-currency fields.
func SeedMultiCurrency(db *gorm.DB) {
	now := time.Now()
	effDate := now.Format("2006-01-02")

	rates := []model.ExchangeRate{
		{
			ID:            "EXR-USD-IDR",
			BaseCurrency:  "USD",
			QuoteCurrency: "IDR",
			Rate:          17749.00,
			Provider:      "ExchangeRate-API (Live)",
			Source:        "Live Market FX",
			FetchedAt:     now,
			EffectiveDate: effDate,
		},
		{
			ID:            "EXR-EUR-IDR",
			BaseCurrency:  "EUR",
			QuoteCurrency: "IDR",
			Rate:          20520.00,
			Provider:      "ExchangeRate-API (Live)",
			Source:        "Live Market FX",
			FetchedAt:     now,
			EffectiveDate: effDate,
		},
		{
			ID:            "EXR-SGD-IDR",
			BaseCurrency:  "SGD",
			QuoteCurrency: "IDR",
			Rate:          13920.00,
			Provider:      "ExchangeRate-API (Live)",
			Source:        "Live Market FX",
			FetchedAt:     now,
			EffectiveDate: effDate,
		},
		{
			ID:            "EXR-USD-EUR",
			BaseCurrency:  "USD",
			QuoteCurrency: "EUR",
			Rate:          0.86,
			Provider:      "ExchangeRate-API (Live)",
			Source:        "Live Market FX",
			FetchedAt:     now,
			EffectiveDate: effDate,
		},
	}

	for _, r := range rates {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"rate", "provider", "source", "fetched_at", "effective_date"}),
		}).Create(&r)
	}

	// Update existing Trading Products with multi-currency values
	var tradingDetails []model.TradingProductDetail
	db.Find(&tradingDetails)
	for _, d := range tradingDetails {
		if d.PurchaseCurrency == "" {
			d.PurchaseCurrency = "USD"
		}
		if d.SellingCurrency == "" {
			d.SellingCurrency = "USD"
		}
		if d.ProfitMarginPercentage <= 0 {
			if d.PurchasePrice > 0 && d.SellingPrice > 0 {
				d.ProfitMarginPercentage = ((d.SellingPrice - d.PurchasePrice) / d.PurchasePrice) * 100
			} else {
				d.ProfitMarginPercentage = 25.0
			}
		}
		if d.ExchangeRateMode == "" {
			d.ExchangeRateMode = "AUTO"
		}
		d.LiveExchangeRate = 17749.00
		if d.ManualExchangeRate <= 0 {
			d.ManualExchangeRate = 17749.00
		}
		d.EffectiveExchangeRate = 17749.00
		d.ExchangeRateProvider = "ExchangeRate-API (Live)"
		d.ExchangeRateUpdatedAt = &now
		d.PurchasePriceInIDR = d.PurchasePrice * 17749.00
		d.SellingPriceInIDR = d.SellingPrice * 17749.00
		db.Save(&d)
	}

	log.Println("[INFO] Multi-Currency Pricing & Profit Engine seeded successfully.")
}

// SeedPhase4 seeds Prototype & Engineering Build Management permissions and initial prototype runs.
func SeedPhase4(db *gorm.DB) {
	// 1. Seed Phase 4 Permissions
	permissions := []model.Permission{
		{ID: "PERM-PRT-VIEW", Module: "prototypes", Action: "view", Code: "prototypes.view", Description: "View prototype builds, BOM snapshots, and assembly progress"},
		{ID: "PERM-PRT-CREATE", Module: "prototypes", Action: "create", Code: "prototypes.create", Description: "Create prototype builds with frozen BOM snapshot"},
		{ID: "PERM-PRT-EDIT", Module: "prototypes", Action: "edit", Code: "prototypes.edit", Description: "Update prototype build metadata and status"},
		{ID: "PERM-PRT-DELETE", Module: "prototypes", Action: "delete", Code: "prototypes.delete", Description: "Delete prototype builds"},
		{ID: "PERM-PRT-EXEC", Module: "prototypes", Action: "execute", Code: "prototype.build.execute", Description: "Advance assembly stages and component readiness"},
		{ID: "PERM-PRT-NOTE", Module: "prototypes", Action: "createNote", Code: "prototype.notes.create", Description: "Create prototype engineering observation notes"},
	}

	for _, p := range permissions {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
	}

	// 2. Assign Permissions to Roles
	rolesToGrant := []string{"ROLE-ADMIN", "ROLE-RNDM", "ROLE-HWENG", "ROLE-FWSPEC", "ROLE-QA"}
	for _, roleID := range rolesToGrant {
		for _, p := range permissions {
			rp := model.RolePermission{RoleID: roleID, PermissionID: p.ID}
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rp)
		}
	}

	// 3. Seed Initial Prototype Builds if table is empty
	var count int64
	db.Model(&model.PrototypeBuild{}).Count(&count)
	if count == 0 {
		now := time.Now()
		target1 := "2026-09-15"
		target2 := "2026-09-05"
		target3 := "2026-09-30"

		bomSnapshot1JSON := `{"snapshotId":"SNP-REV001-001","frozenAt":"2026-08-25T09:30:00Z","hardwareRevisionId":"REV-001","hardwareRevisionCode":"REV-A","totalItems":5,"totalQuantity":28,"totalEstimatedCost":22.45,"currency":"USD","items":[{"id":"SNP-ITM-01","bomItemId":"BOM-ITM-001","componentId":"CMP-001","componentName":"ESP32-S3-WROOM-1-N16R8","mpn":"ESP32-S3-WROOM-1","category":"Microcontrollers","quantity":1,"unitCode":"PCS","unitCost":3.45,"lineCost":3.45,"refDesignators":"U1","position":1,"status":"FROZEN_SNAPSHOT"},{"id":"SNP-ITM-02","bomItemId":"BOM-ITM-002","componentId":"CMP-002","componentName":"SX1262 LoRa Transceiver IC","mpn":"SX1262IMLTRT","category":"Wireless & RF","quantity":1,"unitCode":"PCS","unitCost":4.80,"lineCost":4.80,"refDesignators":"U2","position":2,"status":"FROZEN_SNAPSHOT"},{"id":"SNP-ITM-03","bomItemId":"BOM-ITM-004","componentId":"CMP-004","componentName":"TPS54302 Synchronous Step-Down Converter","mpn":"TPS54302DDCR","category":"Power Management","quantity":1,"unitCode":"PCS","unitCost":1.15,"lineCost":1.15,"refDesignators":"U3","position":3,"status":"FROZEN_SNAPSHOT"},{"id":"SNP-ITM-04","bomItemId":"BOM-ITM-005","componentId":"CMP-006","componentName":"Multi-layer Ceramic Capacitor 10uF 50V","mpn":"CL21A106KOQNNNE","category":"Passive Components","quantity":10,"unitCode":"PCS","unitCost":0.08,"lineCost":0.80,"refDesignators":"C1-C10","position":4,"status":"FROZEN_SNAPSHOT"},{"id":"SNP-ITM-05","bomItemId":"BOM-ITM-006","componentId":"CMP-007","componentName":"Thick Film Chip Resistor 10k 1%","mpn":"RC0603FR-0710KL","category":"Passive Components","quantity":15,"unitCode":"PCS","unitCost":0.02,"lineCost":0.30,"refDesignators":"R1-R15","position":5,"status":"FROZEN_SNAPSHOT"}]}`

		build1 := model.PrototypeBuild{
			ID:                   "PROTO-2024-001",
			Code:                 "PROTO-2024-001",
			Name:                 "GW-400X Alpha Functional Test Rig",
			HardwareRevisionID:   "REV-001",
			ProductID:            "PRD-2024-001",
			ProductCode:          "PRD-2024-001",
			ProductName:          "Smart Edge Gateway (GW-400X)",
			VersionID:            "VER-001",
			VersionNumber:        "v1.0.0",
			VersionName:          "v1.0.0 Base Commercial Architecture",
			HardwareRevisionCode: "REV-A",
			HardwareRevisionName: "REV-A Prototype Alpha",
			BuildNumber:          1,
			Quantity:             5,
			Purpose:              "Functional Prototype",
			Objective:            "Validate multi-protocol gateway processing, LoRa SX1262 telemetry, and 9-36V DC buck regulator stability under simulated thermal load.",
			Status:               model.PrototypeStatusAssembly,
			EngineeringOwner:     "Marcus Vance",
			TargetCompletionAt:   &target1,
			BOMSnapshotJSON:      bomSnapshot1JSON,
			BOMSnapshotCost:      22.45,
			BOMSnapshotID:        "SNP-REV001-001",
			BOMSnapshotCurrency:  "USD",
			ReadinessPercentage:  100.0,
			CreatedAt:            now.AddDate(0, 0, -7),
			UpdatedAt:            now,
		}

		build2 := model.PrototypeBuild{
			ID:                   "PROTO-2024-002",
			Code:                 "PROTO-2024-002",
			Name:                 "GW-400X PCB Validation Sample",
			HardwareRevisionID:   "REV-001",
			ProductID:            "PRD-2024-001",
			ProductCode:          "PRD-2024-001",
			ProductName:          "Smart Edge Gateway (GW-400X)",
			VersionID:            "VER-001",
			VersionNumber:        "v1.0.0",
			VersionName:          "v1.0.0 Base Commercial Architecture",
			HardwareRevisionCode: "REV-A",
			HardwareRevisionName: "REV-A Prototype Alpha",
			BuildNumber:          2,
			Quantity:             2,
			Purpose:              "PCB Validation",
			Objective:            "Specialized 2-board build for 50-ohm RF impedance coupon measurement and micro-via cross-section analysis.",
			Status:               model.PrototypeStatusCompleted,
			EngineeringOwner:     "Elena Rostova",
			TargetCompletionAt:   &target2,
			BOMSnapshotJSON:      bomSnapshot1JSON,
			BOMSnapshotCost:      22.45,
			BOMSnapshotID:        "SNP-REV001-002",
			BOMSnapshotCurrency:  "USD",
			ReadinessPercentage:  100.0,
			CreatedAt:            now.AddDate(0, 0, -12),
			UpdatedAt:            now,
		}

		build3 := model.PrototypeBuild{
			ID:                   "PROTO-2024-003",
			Code:                 "PROTO-2024-003",
			Name:                 "Power Optimization EVT Gateway Rig",
			HardwareRevisionID:   "REV-002",
			ProductID:            "PRD-2024-001",
			ProductCode:          "PRD-2024-001",
			ProductName:          "Smart Edge Gateway (GW-400X)",
			VersionID:            "VER-001",
			VersionNumber:        "v1.0.0",
			VersionName:          "v1.0.0 Base Commercial Architecture",
			HardwareRevisionCode: "REV-B",
			HardwareRevisionName: "REV-B Power Optimization",
			BuildNumber:          1,
			Quantity:             10,
			Purpose:              "Engineering Validation",
			Objective:            "Evaluate deep sleep current reduction with new ultra-low Iq LDO and sleep clock circuit.",
			Status:               model.PrototypeStatusPreparing,
			EngineeringOwner:     "Alex Chen",
			TargetCompletionAt:   &target3,
			BOMSnapshotJSON:      bomSnapshot1JSON,
			BOMSnapshotCost:      22.45,
			BOMSnapshotID:        "SNP-REV002-001",
			BOMSnapshotCurrency:  "USD",
			ReadinessPercentage:  80.0,
			CreatedAt:            now.AddDate(0, 0, -2),
			UpdatedAt:            now,
		}

		db.Create(&build1)
		db.Create(&build2)
		db.Create(&build3)

		// Seed Component Preparations for Build 1
		compPreps1 := []model.PrototypeComponentPreparation{
			{ID: "PREP-001-01", PrototypeBuildID: "PROTO-2024-001", ComponentID: strPtr("CMP-001"), ComponentName: "ESP32-S3-WROOM-1-N16R8", ComponentPartNumber: "ESP32-S3-WROOM-1", ComponentCategory: "Microcontrollers", RequiredQuantity: 5, AvailableQuantity: 5, MissingQuantity: 0, UnitCode: "PCS", Status: model.CompPrepAvailable, Notes: "Reel inspected and verified in lab bin A-12", CreatedAt: now, UpdatedAt: now},
			{ID: "PREP-001-02", PrototypeBuildID: "PROTO-2024-001", ComponentID: strPtr("CMP-002"), ComponentName: "SX1262 LoRa Transceiver IC", ComponentPartNumber: "SX1262IMLTRT", ComponentCategory: "Wireless & RF", RequiredQuantity: 5, AvailableQuantity: 5, MissingQuantity: 0, UnitCode: "PCS", Status: model.CompPrepAvailable, Notes: "Bake before reflow (MSL 3 compliant)", CreatedAt: now, UpdatedAt: now},
			{ID: "PREP-001-03", PrototypeBuildID: "PROTO-2024-001", ComponentID: strPtr("CMP-004"), ComponentName: "TPS54302 Synchronous Step-Down Converter", ComponentPartNumber: "TPS54302DDCR", ComponentCategory: "Power Management", RequiredQuantity: 5, AvailableQuantity: 5, MissingQuantity: 0, UnitCode: "PCS", Status: model.CompPrepAvailable, Notes: "SOT-23-6 package verified", CreatedAt: now, UpdatedAt: now},
			{ID: "PREP-001-04", PrototypeBuildID: "PROTO-2024-001", ComponentID: strPtr("CMP-006"), ComponentName: "Multi-layer Ceramic Capacitor 10uF 50V", ComponentPartNumber: "CL21A106KOQNNNE", ComponentCategory: "Passive Components", RequiredQuantity: 50, AvailableQuantity: 50, MissingQuantity: 0, UnitCode: "PCS", Status: model.CompPrepAvailable, Notes: "SMD 0805 taped", CreatedAt: now, UpdatedAt: now},
			{ID: "PREP-001-05", PrototypeBuildID: "PROTO-2024-001", ComponentID: strPtr("CMP-007"), ComponentName: "Thick Film Chip Resistor 10k 1%", ComponentPartNumber: "RC0603FR-0710KL", ComponentCategory: "Passive Components", RequiredQuantity: 75, AvailableQuantity: 75, MissingQuantity: 0, UnitCode: "PCS", Status: model.CompPrepAvailable, Notes: "SMD 0603 10k reel available", CreatedAt: now, UpdatedAt: now},
		}
		db.Create(&compPreps1)

		// Seed 8 Assembly Stages for Build 1
		doneDate := now.AddDate(0, 0, -2)
		assemblyStages1 := []model.PrototypeAssemblyStage{
			{ID: "STG-001-01", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStagePCBPreparation, Sequence: 1, Name: "PCB Preparation & Solder Paste Stencil", Description: "Laser stainless stencil alignment and Lead-Free SAC305 paste deposition.", Status: model.StageStatusCompleted, PerformedBy: "Marcus Vance", CompletedAt: &doneDate, Notes: "Stencil thickness 0.12mm verified on optical microscope.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-02", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageComponentPreparation, Sequence: 2, Name: "Component Kitting & Reel Setup", Description: "Feeder loading and component verification for high-speed SMD assembly.", Status: model.StageStatusCompleted, PerformedBy: "Marcus Vance", CompletedAt: &doneDate, Notes: "All reels loaded and optical fiducials calibrated.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-03", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageSMTAssembly, Sequence: 3, Name: "SMT Pick & Place Assembly", Description: "Automated component placement for top and bottom PCBA layers.", Status: model.StageStatusCompleted, PerformedBy: "Marcus Vance", CompletedAt: &doneDate, Notes: "5 boards placed in 14 minutes. Zero component dropouts.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-04", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageThroughHoleAssembly, Sequence: 4, Name: "Reflow Soldering & AOI Inspection", Description: "9-zone nitrogen reflow oven profile and automated optical inspection.", Status: model.StageStatusCompleted, PerformedBy: "Marcus Vance", CompletedAt: &doneDate, Notes: "Peak temperature 242°C. Zero solder bridging observed.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-05", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageMechanicalFitting, Sequence: 5, Name: "Through-Hole Connectors & Hand Assembly", Description: "Terminal blocks, SMA antenna connectors, and RJ45 Ethernet jacks.", Status: model.StageStatusInProgress, PerformedBy: "Marcus Vance", Notes: "Currently hand-soldering IP67 SMA connectors on unit #3.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-06", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageWiring, Sequence: 6, Name: "Mechanical Fitting & Enclosure Assembly", Description: "Extruded aluminum enclosure insertion, thermal pad mounting, and gaskets.", Status: model.StageStatusNotStarted, PerformedBy: "Elena Rostova", Notes: "Pending completion of through-hole soldering.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-07", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageFirmwareFlashing, Sequence: 7, Name: "Firmware Bootloader Flashing", Description: "JTAG/UART flashing of v1.0.0 factory bootloader and MAC allocation.", Status: model.StageStatusNotStarted, PerformedBy: "Sarah Jenkins", Notes: "Flashing test fixture prepared at Lab Bench 4.", CreatedAt: now, UpdatedAt: now},
			{ID: "STG-001-08", PrototypeBuildID: "PROTO-2024-001", Stage: model.AssemblyStageInitialPowerOn, Sequence: 8, Name: "Initial Power-On & Voltage Rail Validation", Description: "Current limiting bench supply test: 3.3V, 5.0V, and RF LDO rails.", Status: model.StageStatusNotStarted, PerformedBy: "Marcus Vance", Notes: "DMM and oscilloscope probes assigned.", CreatedAt: now, UpdatedAt: now},
		}
		db.Create(&assemblyStages1)

		// Seed Engineering Notes for Build 1
		notes1 := []model.PrototypeEngineeringNote{
			{ID: "NOT-001-01", PrototypeBuildID: "PROTO-2024-001", Category: model.NoteCategoryAssemblyIssue, Title: "SMA Connector Thermal Mass", Content: "The ground plane on J1 requires elevated iron temperature (360°C) during hand soldering due to heavy 2oz copper ground fill.", CreatedBy: "Marcus Vance", CreatedAt: now.AddDate(0, 0, -1), UpdatedAt: now.AddDate(0, 0, -1)},
			{ID: "NOT-001-02", PrototypeBuildID: "PROTO-2024-001", Category: model.NoteCategoryDesignObservation, Title: "Thermal Pad Thickness on MCU", Content: "Thermal pad thickness of 1.0mm provides optimal compression against the aluminum casing without putting stress on the QFN solder joints.", CreatedBy: "Elena Rostova", CreatedAt: now.AddDate(0, 0, -3), UpdatedAt: now.AddDate(0, 0, -3)},
		}
		db.Create(&notes1)
	}

	log.Println("[INFO] Phase 4 Prototype & Engineering Build Management seeded successfully.")
}

func SeedPhase5(db *gorm.DB) {
	// 1. Seed Permissions
	permissions := []model.Permission{
		{ID: "PERM-TST-VIEW", Module: "testing", Action: "view", Code: "testing.view", Description: "View test plans, sessions, executions, and findings"},
		{ID: "PERM-TST-PLN-MNG", Module: "testing", Action: "managePlans", Code: "testplans.manage", Description: "Create, edit, clone, and release test plan versions"},
		{ID: "PERM-TST-EXEC", Module: "testing", Action: "execute", Code: "testing.execute", Description: "Execute tests, record measurements, and upload evidence"},
		{ID: "PERM-FND-MNG", Module: "testing", Action: "manageFindings", Code: "findings.manage", Description: "Create, assign, and disposition engineering findings"},
		{ID: "PERM-TST-VAL", Module: "testing", Action: "validate", Code: "testing.validate", Description: "Create formal prototype validation decisions"},
	}

	for _, p := range permissions {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
	}

	// 2. Assign Permissions to Roles
	// ROLE-ADMIN & ROLE-RNDM get everything
	for _, p := range permissions {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-ADMIN", PermissionID: p.ID})
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-RNDM", PermissionID: p.ID})
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-QA", PermissionID: p.ID})
	}
	// ROLE-HWENG & ROLE-FWENG get view, plan manage (HW), execute, findings manage
	hwPerms := []string{"PERM-TST-VIEW", "PERM-TST-PLN-MNG", "PERM-TST-EXEC", "PERM-FND-MNG"}
	for _, pid := range hwPerms {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-HWENG", PermissionID: pid})
	}
	fwPerms := []string{"PERM-TST-VIEW", "PERM-TST-EXEC", "PERM-FND-MNG"}
	for _, pid := range fwPerms {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-FWENG", PermissionID: pid})
	}
	// ROLE-VIEW gets view only
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-VIEW", PermissionID: "PERM-TST-VIEW"})

	// 3. Seed Master Test Plans if empty
	var planCount int64
	db.Model(&model.TestPlan{}).Count(&planCount)
	if planCount == 0 {
		now := time.Now()
		relDate := now.AddDate(0, 0, -10)

		plan1 := model.TestPlan{
			ID:               "TP-001",
			Code:             "TP-GW400X-EVT",
			Name:             "Smart Edge Gateway EVT Qualification Suite",
			Description:      "Comprehensive engineering validation test (EVT) protocol covering power rails, RF transceiver metrics, sensor interfaces, and thermal stress.",
			Category:         "ELECTRICAL",
			ProductID:        strPtr("PRD-2024-001"),
			TargetRevisionID: strPtr("REV-001"),
			Status:           model.TestPlanStatusActive,
			CreatedBy:        "Sarah Jenkins",
			CreatedAt:        now.AddDate(0, 0, -15),
			UpdatedAt:        now,
		}
		db.Create(&plan1)

		ver1 := model.TestPlanVersion{
			ID:            "TPV-001-01",
			TestPlanID:    "TP-001",
			VersionNumber: "v1.0",
			Status:        model.PlanVersionStatusReleased,
			Description:   "Baseline EVT test protocol for REV-A initial PCB prototypes.",
			ReleaseNotes:  "Initial formal release approved by R&D Lead.",
			CreatedBy:     "Sarah Jenkins",
			ReleasedBy:    "Sarah Jenkins",
			ReleasedAt:    &relDate,
			CreatedAt:     now.AddDate(0, 0, -15),
			UpdatedAt:     relDate,
		}
		db.Create(&ver1)

		// 8 Test Cases for Version 1.0
		min3V3 := 3.15
		max3V3 := 3.45
		tgt3V3 := 3.30

		min5V := 4.75
		max5V := 5.25
		tgt5V := 5.00

		maxSleepCurrent := 25.0 // uA
		tgtSleepCurrent := 15.0

		minTxPower := 13.0 // dBm
		maxTxPower := 15.5
		tgtTxPower := 14.0

		maxPcbTemp := 70.0 // degC
		tgtPcbTemp := 45.0

		expTrue := true

		cases := []model.TestCase{
			{
				ID:                "TC-001",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          1,
				Code:              "TC-PWR-01",
				Name:              "3.3V System Rail Regulation Test",
				Category:          "ELECTRICAL",
				Description:       "Verify main 3.3V buck converter output rail under nominal 12V DC input.",
				TestType:          model.TestTypeNumericRange,
				Unit:              "V",
				MinimumValue:      &min3V3,
				MaximumValue:      &max3V3,
				TargetValue:       &tgt3V3,
				Instructions:      "Connect 12.0V DC bench supply to J1 input. Measure TP3 (3.3V) with calibrated 6.5-digit DMM.",
				AcceptanceNotes:   "Voltage must remain within 3.3V ± 4.5% tolerance.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-002",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          2,
				Code:              "TC-PWR-02",
				Name:              "5.0V Auxiliary Power Rail Output",
				Category:          "ELECTRICAL",
				Description:       "Verify 5.0V switching converter rail for RS-485 transceiver and expansion modules.",
				TestType:          model.TestTypeNumericRange,
				Unit:              "V",
				MinimumValue:      &min5V,
				MaximumValue:      &max5V,
				TargetValue:       &tgt5V,
				Instructions:      "Probe TP5 (5.0V Rail) with oscilloscope to check DC level and ripple noise (<50mVpp).",
				AcceptanceNotes:   "5.0V ± 5% regulation required.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-003",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          3,
				Code:              "TC-PWR-03",
				Name:              "Deep Sleep Quiescent Current Consumption",
				Category:          "POWER_CONSUMPTION",
				Description:       "Measure total board current draw when MCU is in deep sleep with RF powered down.",
				TestType:          model.TestTypeNumericMax,
				Unit:              "uA",
				MaximumValue:      &maxSleepCurrent,
				TargetValue:       &tgtSleepCurrent,
				Instructions:      "Place Keysight Current Analyzer in series with 12V input. Issue CLI sleep command. Record 60-second average.",
				AcceptanceNotes:   "Sleep current must not exceed 25uA.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-004",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          4,
				Code:              "TC-RF-01",
				Name:              "LoRa Transmit Power at 915 MHz",
				Category:          "RF_WIRELESS",
				Description:       "Measure calibrated Tx power output from SMA connector into 50-ohm spectrum analyzer.",
				TestType:          model.TestTypeNumericRange,
				Unit:              "dBm",
				MinimumValue:      &minTxPower,
				MaximumValue:      &maxTxPower,
				TargetValue:       &tgtTxPower,
				Instructions:      "Connect SMA J2 via 20dB attenuator to R&S Spectrum Analyzer. Execute continuous wave transmit carrier at 915.0 MHz.",
				AcceptanceNotes:   "Tx power must measure between 13.0 and 15.5 dBm.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-005",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          5,
				Code:              "TC-FNC-01",
				Name:              "Dual-Band Wi-Fi Association & Cloud Telemetry Ping",
				Category:          "FUNCTIONAL",
				Description:       "Verify ESP32-S3 Wi-Fi connection to enterprise AP and MQTT cloud telemetry handshake.",
				TestType:          model.TestTypeBooleanPassFail,
				ExpectedBoolean:   &expTrue,
				Instructions:      "Boot device with test Wi-Fi credentials. Observe serial log for MQTT CONNACK status.",
				AcceptanceNotes:   "Device must connect and publish test payload in under 5.0 seconds.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-006",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          6,
				Code:              "TC-COM-01",
				Name:              "Isolated RS-485 Loopback Transceive",
				Category:          "COMMUNICATIONS",
				Description:       "Verify bidirectional communication across isolated RS-485 interface at 115200 baud.",
				TestType:          model.TestTypeBooleanPassFail,
				ExpectedBoolean:   &expTrue,
				Instructions:      "Attach loopback plug across A/B terminal pins. Execute RS485 loopback packet echo tool.",
				AcceptanceNotes:   "1000/1000 packets must echo without bit errors (0% PER).",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-007",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          7,
				Code:              "TC-THM-01",
				Name:              "Max Compute Load Thermal Chamber Heatmap",
				Category:          "THERMAL",
				Description:       "Record steady-state temperature on MCU junction, LDO, and buck inductor under 100% duty cycle.",
				TestType:          model.TestTypeNumericMax,
				Unit:              "degC",
				MaximumValue:      &maxPcbTemp,
				TargetValue:       &tgtPcbTemp,
				Instructions:      "Run benchmark firmware for 60 minutes in 25°C ambient room. Capture thermal FLIR infrared image.",
				AcceptanceNotes:   "No component surface may exceed 70.0°C.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                "TC-008",
				TestPlanVersionID: "TPV-001-01",
				Sequence:          8,
				Code:              "TC-VIS-01",
				Name:              "IPC-A-610 Class 2 PCBA Solder Quality Inspection",
				Category:          "VISUAL_INSPECTION",
				Description:       "Microscopic inspection of SMT fillets, QFN wetting, and solder mask clearance.",
				TestType:          model.TestTypeVisualInspection,
				Instructions:      "Inspect top and bottom layers under 20x stereo microscope. Verify zero solder balls, tombstoning, or voiding > 15%.",
				AcceptanceNotes:   "100% visual compliance with IPC-A-610 Class 2 standard.",
				Status:            "ACTIVE",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		}
		db.Create(&cases)

		// Draft Version 1.1 for TP-001
		ver2 := model.TestPlanVersion{
			ID:            "TPV-001-02",
			TestPlanID:    "TP-001",
			VersionNumber: "v1.1",
			Status:        model.PlanVersionStatusDraft,
			Description:   "Proposed EVT update adding high-speed CAN-FD arbitration test and tighter 3.3V tolerance.",
			ReleaseNotes:  "Draft in progress by Marcus Vance for REV-B preparation.",
			CreatedBy:     "Marcus Vance",
			CreatedAt:     now.AddDate(0, 0, -2),
			UpdatedAt:     now,
		}
		db.Create(&ver2)

		// Seed initial Test Session on PROTO-2024-001
		snapshotItems := make([]model.FrozenTestCaseItem, 0, len(cases))
		for _, c := range cases {
			snapshotItems = append(snapshotItems, model.FrozenTestCaseItem{
				ID:              c.ID,
				Sequence:        c.Sequence,
				Code:            c.Code,
				Name:            c.Name,
				Category:        c.Category,
				TestType:        c.TestType,
				Unit:            c.Unit,
				MinimumValue:    c.MinimumValue,
				MaximumValue:    c.MaximumValue,
				TargetValue:     c.TargetValue,
				ExpectedBoolean: c.ExpectedBoolean,
				ExpectedText:    c.ExpectedText,
				Instructions:    c.Instructions,
			})
		}
		snapJSON, _ := json.Marshal(model.FrozenTestPlanSnapshot{
			TestPlanID:     plan1.ID,
			TestPlanCode:   plan1.Code,
			TestPlanName:   plan1.Name,
			Category:       plan1.Category,
			VersionID:      ver1.ID,
			VersionNumber:  ver1.VersionNumber,
			TotalTestCases: len(cases),
			FrozenAt:       relDate,
			FrozenBy:       "Sarah Jenkins",
			TestCases:      snapshotItems,
		})

		sessStartTime := now.AddDate(0, 0, -5)
		sessDoneTime := now.AddDate(0, 0, -1)

		session1 := model.PrototypeTestSession{
			ID:                   "SESS-2024-001",
			PrototypeBuildID:     "PROTO-2024-001",
			TestPlanID:           plan1.ID,
			TestPlanVersionID:    ver1.ID,
			SessionCode:          "SESS-PROTO001-EVT01",
			Name:                 "PROTO-2024-001 Initial EVT Verification Run",
			Description:          "Primary laboratory qualification run on Build Unit #1 at Room Ambient (25°C).",
			Status:               model.SessionStatusCompleted,
			TestPlanSnapshotJSON: string(snapJSON),
			TotalCases:           8,
			PassedCases:          6,
			FailedCases:          1,
			BlockedCases:         0,
			SkippedCases:         1,
			PassRatePercentage:   75.0,
			StartedAt:            &sessStartTime,
			CompletedAt:          &sessDoneTime,
			CreatedBy:            "Marcus Vance",
			CreatedAt:            now.AddDate(0, 0, -5),
			UpdatedAt:            sessDoneTime,
		}
		db.Create(&session1)

		// Seed Executions for Session 1
		v33 := 3.28
		v50 := 5.04
		iSleep := 18.4
		pTx := 14.2
		tPcb := 74.8 // Fails max 70.0!

		exec1 := model.TestExecution{
			ID:                     "EXEC-001-01",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-001",
			TestCaseCode:           "TC-PWR-01",
			TestCaseName:           "3.3V System Rail Regulation Test",
			Category:               "ELECTRICAL",
			TestType:               model.TestTypeNumericRange,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredValue:          &v33,
			Unit:                   "V",
			CalculatedResult:       "PASSED",
			Notes:                  "Measured 3.28V DC at TP3. Output clean and stable.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec2 := model.TestExecution{
			ID:                     "EXEC-001-02",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-002",
			TestCaseCode:           "TC-PWR-02",
			TestCaseName:           "5.0V Auxiliary Power Rail Output",
			Category:               "ELECTRICAL",
			TestType:               model.TestTypeNumericRange,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredValue:          &v50,
			Unit:                   "V",
			CalculatedResult:       "PASSED",
			Notes:                  "Measured 5.04V DC at TP5. Ripple < 25mVpp.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec3 := model.TestExecution{
			ID:                     "EXEC-001-03",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-003",
			TestCaseCode:           "TC-PWR-03",
			TestCaseName:           "Deep Sleep Quiescent Current Consumption",
			Category:               "POWER_CONSUMPTION",
			TestType:               model.TestTypeNumericMax,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredValue:          &iSleep,
			Unit:                   "uA",
			CalculatedResult:       "PASSED",
			Notes:                  "Average sleep current 18.4 uA across 60 seconds.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec4 := model.TestExecution{
			ID:                     "EXEC-001-04",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-004",
			TestCaseCode:           "TC-RF-01",
			TestCaseName:           "LoRa Transmit Power at 915 MHz",
			Category:               "RF_WIRELESS",
			TestType:               model.TestTypeNumericRange,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredValue:          &pTx,
			Unit:                   "dBm",
			CalculatedResult:       "PASSED",
			Notes:                  "Spectrum analyzer recorded 14.2 dBm at 915.0 MHz center frequency.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec5 := model.TestExecution{
			ID:                     "EXEC-001-05",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-005",
			TestCaseCode:           "TC-FNC-01",
			TestCaseName:           "Dual-Band Wi-Fi Association & Cloud Telemetry Ping",
			Category:               "FUNCTIONAL",
			TestType:               model.TestTypeBooleanPassFail,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredBoolean:        &expTrue,
			CalculatedResult:       "PASSED",
			Notes:                  "MQTT telemetry confirmed connected to AWS IoT Core in 2.8s.",
			ExecutedBy:             "Sarah Jenkins",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec6 := model.TestExecution{
			ID:                     "EXEC-001-06",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-006",
			TestCaseCode:           "TC-COM-01",
			TestCaseName:           "Isolated RS-485 Loopback Transceive",
			Category:               "COMMUNICATIONS",
			TestType:               model.TestTypeBooleanPassFail,
			RunNumber:              1,
			Status:                 model.ExecStatusPassed,
			MeasuredBoolean:        &expTrue,
			CalculatedResult:       "PASSED",
			Notes:                  "1000/1000 packets echoed without transmission framing errors.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec7 := model.TestExecution{
			ID:                     "EXEC-001-07",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-007",
			TestCaseCode:           "TC-THM-01",
			TestCaseName:           "Max Compute Load Thermal Chamber Heatmap",
			Category:               "THERMAL",
			TestType:               model.TestTypeNumericMax,
			RunNumber:              1,
			Status:                 model.ExecStatusFailed,
			MeasuredValue:          &tPcb,
			Unit:                   "degC",
			CalculatedResult:       "FAILED",
			Notes:                  "Inductor L1 and adjacent LDO reached 74.8°C (limit: 70.0°C) during continuous benchmark.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}
		exec8 := model.TestExecution{
			ID:                     "EXEC-001-08",
			PrototypeTestSessionID: session1.ID,
			TestCaseSnapshotID:     "TC-008",
			TestCaseCode:           "TC-VIS-01",
			TestCaseName:           "IPC-A-610 Class 2 PCBA Solder Quality Inspection",
			Category:               "VISUAL_INSPECTION",
			TestType:               model.TestTypeVisualInspection,
			RunNumber:              1,
			Status:                 model.ExecStatusSkipped,
			Notes:                  "Skipped for initial prototype bench bring-up. Scheduled before final box build.",
			ExecutedBy:             "Marcus Vance",
			StartedAt:              &sessStartTime,
			CompletedAt:            &sessStartTime,
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessStartTime,
		}

		db.Create(&exec1)
		db.Create(&exec2)
		db.Create(&exec3)
		db.Create(&exec4)
		db.Create(&exec5)
		db.Create(&exec6)
		db.Create(&exec7)
		db.Create(&exec8)

		// Seed Measurements
		meas7 := model.TestMeasurement{
			ID:               "MEAS-001-01",
			TestExecutionID:  exec7.ID,
			SampleIndex:      1,
			MeasuredValue:    &tPcb,
			Unit:             "degC",
			CalculatedResult: "FAILED",
			RecordedBy:       "Marcus Vance",
			RecordedAt:       sessStartTime,
			CreatedAt:        sessStartTime,
		}
		db.Create(&meas7)

		// Seed Engineering Findings
		finding1 := model.EngineeringFinding{
			ID:                     "FND-2024-001",
			PrototypeBuildID:       "PROTO-2024-001",
			PrototypeTestSessionID: strPtr(session1.ID),
			TestExecutionID:        strPtr(exec7.ID),
			Code:                   "FND-THM-001",
			Title:                  "Step-Down Inductor L1 Exceeds 70°C Under Max Compute Load",
			Description:            "During 60-minute continuous benchmark, L1 inductor temperature rose to 74.8°C, exceeding the 70.0°C specification due to inadequate PCB thermal copper area.",
			Category:               "THERMAL",
			Severity:               model.FindingSeverityHigh,
			FindingDisposition:     model.FindingDispDesignChange,
			ChangeCandidateStatus:  model.ChangeCandidateCandidate,
			RecommendedChangeScope: model.ChangeScopePCBLayout,
			RootCause:              "Thermal copper pour on inner layer 2 was necked down near the mounting hole, restricting heat dissipation.",
			ContainmentAction:      "Added temporary thermal gap pad to aluminum chassis during laboratory EVT testing.",
			ResolutionNotes:        "Enlarge L1 footprint and increase thermal via array in REV-B layout.",
			Status:                 "OPEN",
			ReportedBy:             "Marcus Vance",
			AssignedTo:             "Elena Rostova",
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessDoneTime,
		}
		finding2 := model.EngineeringFinding{
			ID:                     "FND-2024-002",
			PrototypeBuildID:       "PROTO-2024-001",
			PrototypeTestSessionID: strPtr(session1.ID),
			Code:                   "FND-ELC-002",
			Title:                  "I2C Bus Rise Time Marginal at 400 kHz",
			Description:            "I2C pull-up resistors (4.7k) exhibit 450ns rise time on SCL line under full sensor bus loading.",
			Category:               "ELECTRICAL",
			Severity:               model.FindingSeverityMedium,
			FindingDisposition:     model.FindingDispBenchRework,
			ChangeCandidateStatus:  model.ChangeCandidateNone,
			RecommendedChangeScope: model.ChangeScopeSchematic,
			RootCause:              "Bus capacitance increased with addition of environmental sensor.",
			ContainmentAction:      "Reworked R12 and R13 pull-up resistors to 2.2k on Build #1.",
			ResolutionNotes:        "Bench rework verified rise time reduced to 180ns.",
			Status:                 "CLOSED",
			ReportedBy:             "Marcus Vance",
			AssignedTo:             "Marcus Vance",
			CreatedAt:              sessStartTime,
			UpdatedAt:              sessDoneTime,
		}
		db.Create(&finding1)
		db.Create(&finding2)

		// Seed Validation Decision
		valDec1 := model.PrototypeValidationDecision{
			ID:                     "VAL-2024-001",
			PrototypeTestSessionID: session1.ID,
			Decision:               model.ValidationDecisionConditionallyValidated,
			Justification:          "Approved for bench and lab RF verification. Full environmental chamber sign-off deferred to REV-B after L1 thermal copper enlargement.",
			DecisionSummary:        "EVT Session 1 achieved 75% pass rate (6 passed, 1 failed, 1 skipped). High severity thermal finding FND-THM-001 dispositioned as Design Change Candidate for REV-B.",
			DecidedBy:              "Sarah Jenkins",
			IsCurrent:              true,
			CreatedAt:              sessDoneTime,
		}
		db.Create(&valDec1)
	}

	log.Println("[INFO] Phase 5 Testing & Validation Control seeded successfully.")
}

func SeedPhase6(db *gorm.DB) {
	// 1. Seed Phase 6 RBAC Permissions
	phase6Permissions := []model.Permission{
		{ID: "PERM-CHG-VIEW", Module: "changes", Action: "view", Code: "changes.view", Description: "View Engineering Changes Dashboard and Registries"},
		{ID: "PERM-ECR-CREATE", Module: "ecr", Action: "create", Code: "ecr.create", Description: "Create and submit Engineering Change Requests"},
		{ID: "PERM-ECR-EDIT", Module: "ecr", Action: "edit", Code: "ecr.edit", Description: "Edit draft Engineering Change Requests and change items"},
		{ID: "PERM-ECR-REVIEW", Module: "ecr", Action: "review", Code: "ecr.review", Description: "Submit technical domain reviews on ECRs"},
		{ID: "PERM-ECR-APPROVE", Module: "ecr", Action: "approve", Code: "ecr.approve", Description: "Approve or reject ECRs in CCB governance chain"},
		{ID: "PERM-ECO-CREATE", Module: "eco", Action: "create", Code: "eco.create", Description: "Create Engineering Change Orders from approved ECRs"},
		{ID: "PERM-ECO-APPROVE", Module: "eco", Action: "approve", Code: "eco.approve", Description: "Authorize and approve ECOs for implementation"},
		{ID: "PERM-ECO-IMPL", Module: "eco", Action: "implement", Code: "eco.implement", Description: "Execute ECO baseline increment and BOM cloning engine"},
		{ID: "PERM-ECO-VERIFY", Module: "eco", Action: "verify", Code: "eco.verify", Description: "Verify ECO completion and link regression test sessions"},
	}

	for _, p := range phase6Permissions {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
		// Assign to Super Admin
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-ADMIN", PermissionID: p.ID})
		// Assign to R&D Manager
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-RNDM", PermissionID: p.ID})
		// Assign view to Viewer
		if p.Action == "view" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-VIEW", PermissionID: p.ID})
		}
		// Assign to Hardware Engineer
		if p.Code == "changes.view" || p.Code == "ecr.create" || p.Code == "ecr.edit" || p.Code == "ecr.review" || p.Code == "eco.implement" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-HWENG", PermissionID: p.ID})
		}
		// Assign to Firmware Engineer
		if p.Code == "changes.view" || p.Code == "ecr.create" || p.Code == "ecr.edit" || p.Code == "ecr.review" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-FWENG", PermissionID: p.ID})
		}
		// Assign to Purchasing
		if p.Code == "changes.view" || p.Code == "ecr.review" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{RoleID: "ROLE-PURCH", PermissionID: p.ID})
		}
	}

	// 2. Seed Master ECR & ECO Datasets
	var ecrCount int64
	db.Model(&model.EngineeringChangeRequest{}).Count(&ecrCount)
	if ecrCount == 0 {
		now := time.Now()
		ecrDate := now.Add(-48 * time.Hour)
		targetDate := now.Add(30 * 24 * time.Hour)

		// Seed ECR-001 (Sourced from Phase 5 Thermal Finding FND-2024-001)
		ecr1 := model.EngineeringChangeRequest{
			ID:                       "ECR-2026-001",
			Code:                     "ECR-2026-001",
			Title:                    "LDO Voltage Regulator Thermal Relief & Ground Plane Rework",
			Origin:                   "TEST_FINDING",
			SourceFindingID:          strPtr("FND-2024-001"),
			ProductID:                "PRD-2024-001",
			ProductVersionID:         "VER-001",
			SourceHardwareRevisionID: "REV-001",
			ChangeType:               "BOM_COMPONENT",
			Priority:                 "CRITICAL",
			Description:              "Prototype Build #1 thermal chamber tests registered 128°C under max continuous compute load. The 3.3V LDO regulator (U3) requires replacement with an automotive-grade, high-thermal-efficiency LDO and solid ground plane stitching vias to prevent thermal throttling.",
			BusinessReason:           "Product operating temperature specification is -40°C to +85°C ambient. Under 65°C ambient chamber test, junction temp exceeded 125°C threshold.",
			TechnicalJustification:   "Replacing TI TLV73333 with TI LP5907 with exposed thermal pad (150°C Tj max) and connecting Pin 2 GND directly to L2 solid plane with 6x 0.3mm thermal vias.",
			Status:                   "APPROVED",
			RequestedBy:              "Alex Chen",
			AssignedLead:             "Marcus Vance",
			EstimatedCost:            4500.00,
			TargetReleaseDate:        &targetDate,
			CreatedAt:                ecrDate,
			UpdatedAt:                now,
		}
		db.Create(&ecr1)

		// Seed Change Items for ECR-001
		items := []model.ECRChangeItem{
			{
				ID:                     "CHG-001-01",
				ECRID:                  ecr1.ID,
				Sequence:               1,
				ChangeType:             "BOM_COMPONENT",
				Title:                  "Upgrade 3.3V LDO to LP5907 with Exposed Thermal Pad",
				AffectedReference:      "U3",
				AffectedComponentID:    strPtr("CMP-003"),
				CurrentState:           "TLV73333PDBVR (SOT-23-5, Tj Max 125°C, $0.18/unit)",
				ProposedState:          "LP5907MFX-3.3 (SOT-23-5 with exposed thermal pad, Tj Max 150°C, $0.24/unit)",
				TechnicalJustification: "Allows higher power dissipation and drops junction temperature by ~44°C under peak load.",
				RiskLevel:              "HIGH",
				Notes:                  "Footprint compatible with PCB layout revisions.",
				CreatedAt:              ecrDate,
				UpdatedAt:              ecrDate,
			},
			{
				ID:                     "CHG-001-02",
				ECRID:                  ecr1.ID,
				Sequence:               2,
				ChangeType:             "PCB_LAYOUT",
				Title:                  "Enlarge L2 Ground Plane Thermal Relief & Add Stitching Vias",
				AffectedReference:      "L2 Inner Ground Plane",
				CurrentState:           "Standard 4-spoke thermal relief spoke width 0.25mm on LDO GND pin.",
				ProposedState:          "Solid copper connection on Pin 2 with 6x 0.3mm thermal vias tied to bottom enclosure heatsink plane.",
				TechnicalJustification: "Enhances heat transfer from IC die directly to PCB copper plane.",
				RiskLevel:              "MEDIUM",
				Notes:                  "Requires minor Gerber re-spin on Layer 2 and Layer 4.",
				CreatedAt:              ecrDate,
				UpdatedAt:              ecrDate,
			},
		}
		db.Create(&items)

		// Seed Impact Analyses for ECR-001
		impacts := []model.ECRImpactAnalysis{
			{
				ID:                  "IMP-001-01",
				ECRID:               ecr1.ID,
				Domain:              "ELECTRICAL",
				ImpactLevel:         "LOW",
				Description:         "Output ripple voltage drops from 45mVpp to 18mVpp. PSRR improves by +14dB at 100kHz.",
				Mitigation:          "Maintain existing 10uF ceramic output capacitor footprint.",
				CostDelta:           0.06,
				ScheduleImpactWeeks: 0.0,
				ScrapDisposition:    "RUN_OUT",
				ReviewerID:          "USR-003",
				ReviewerName:        "Marcus Vance",
				ReviewStatus:        "REVIEWED",
				CreatedAt:           ecrDate,
				UpdatedAt:           now,
			},
			{
				ID:                  "IMP-001-02",
				ECRID:               ecr1.ID,
				Domain:              "THERMAL",
				ImpactLevel:         "CRITICAL",
				Description:         "Thermal simulation confirms max junction temp drops from 128°C to 84°C under 100% duty cycle @ 65°C ambient.",
				Mitigation:          "Validated with thermal chamber profile modeling.",
				CostDelta:           0.00,
				ScheduleImpactWeeks: 0.5,
				ScrapDisposition:    "SCRAP_ALL",
				ReviewerID:          "USR-003",
				ReviewerName:        "Marcus Vance",
				ReviewStatus:        "REVIEWED",
				CreatedAt:           ecrDate,
				UpdatedAt:           now,
			},
			{
				ID:                  "IMP-001-03",
				ECRID:               ecr1.ID,
				Domain:              "MECHANICAL",
				ImpactLevel:         "NOT_AFFECTED",
				Description:         "Enclosure clearances, mounting posts, and connector openings are 100% unaffected.",
				Mitigation:          "Zero enclosure re-tooling required.",
				CostDelta:           0.00,
				ScheduleImpactWeeks: 0.0,
				ScrapDisposition:    "RUN_OUT",
				ReviewerID:          "USR-001",
				ReviewerName:        "Alex Chen",
				ReviewStatus:        "REVIEWED",
				CreatedAt:           ecrDate,
				UpdatedAt:           now,
			},
			{
				ID:                  "IMP-001-04",
				ECRID:               ecr1.ID,
				Domain:              "FIRMWARE",
				ImpactLevel:         "NOT_AFFECTED",
				Description:         "Purely analog hardware power rail change. Zero software register maps or driver logic modified.",
				Mitigation:          "No firmware changes needed.",
				CostDelta:           0.00,
				ScheduleImpactWeeks: 0.0,
				ScrapDisposition:    "RUN_OUT",
				ReviewerID:          "USR-004",
				ReviewerName:        "Elena Rostova",
				ReviewStatus:        "REVIEWED",
				CreatedAt:           ecrDate,
				UpdatedAt:           now,
			},
			{
				ID:                  "IMP-001-05",
				ECRID:               ecr1.ID,
				Domain:              "BOM_COST",
				ImpactLevel:         "LOW",
				Description:         "Net unit BOM cost increases by +$0.06/unit ($0.18 -> $0.24). Fits well within the $120 MSRP margin model.",
				Mitigation:          "Approved by project lead.",
				CostDelta:           0.06,
				ScheduleImpactWeeks: 0.0,
				ScrapDisposition:    "RUN_OUT",
				ReviewerID:          "USR-005",
				ReviewerName:        "David Tanaka",
				ReviewStatus:        "REVIEWED",
				CreatedAt:           ecrDate,
				UpdatedAt:           now,
			},
		}
		db.Create(&impacts)

		// Seed Reviews for ECR-001
		revDone := now.Add(-24 * time.Hour)
		reviews := []model.ECRReview{
			{
				ID:           "REV-001-01",
				ECRID:        ecr1.ID,
				Discipline:   "HARDWARE",
				ReviewerID:   "USR-003",
				ReviewerName: "Marcus Vance",
				Decision:     "APPROVE",
				Comments:     "Hardware engineering agrees. Thermal simulation proves 44°C junction temp reduction.",
				ReviewedAt:   &revDone,
				CreatedAt:    ecrDate,
				UpdatedAt:    revDone,
			},
			{
				ID:           "REV-001-02",
				ECRID:        ecr1.ID,
				Discipline:   "QA",
				ReviewerID:   "USR-002",
				ReviewerName: "Sarah Jenkins",
				Decision:     "APPROVE",
				Comments:     "Quality and reliability signs off. Environmental re-test mandated for REV-B prototype.",
				ReviewedAt:   &revDone,
				CreatedAt:    ecrDate,
				UpdatedAt:    revDone,
			},
		}
		db.Create(&reviews)

		// Seed Approvals for ECR-001 (Unanimous CCB Approval)
		apprDone := now.Add(-12 * time.Hour)
		approvals := []model.ECRApproval{
			{
				ID:            "APPR-001-01",
				ECRID:         ecr1.ID,
				Stage:         "TIER_1_TECH_LEAD",
				ApproverID:    "USR-001",
				ApproverName:  "Alex Chen",
				ApproverRole:  "Technical Lead",
				Decision:      "APPROVED",
				Justification: "Technical solution is sound and necessary for environmental certification.",
				DecidedAt:     &apprDone,
				CreatedAt:     ecrDate,
				UpdatedAt:     apprDone,
			},
			{
				ID:            "APPR-001-02",
				ECRID:         ecr1.ID,
				Stage:         "TIER_2_ENG_MANAGER",
				ApproverID:    "USR-001",
				ApproverName:  "Alex Chen",
				ApproverRole:  "Engineering Manager",
				Decision:      "APPROVED",
				Justification: "Budget approved for PCB re-spin and prototype fabrication.",
				DecidedAt:     &apprDone,
				CreatedAt:     ecrDate,
				UpdatedAt:     apprDone,
			},
			{
				ID:            "APPR-001-03",
				ECRID:         ecr1.ID,
				Stage:         "TIER_3_QUALITY_DIRECTOR",
				ApproverID:    "USR-002",
				ApproverName:  "Sarah Jenkins",
				ApproverRole:  "Quality Director",
				Decision:      "APPROVED",
				Justification: "Quality and Compliance concurs. Mandatory regression test required upon REV-B assembly.",
				DecidedAt:     &apprDone,
				CreatedAt:     ecrDate,
				UpdatedAt:     apprDone,
			},
		}
		db.Create(&approvals)

		// Seed ECO-001 (Sprouted from Approved ECR-001)
		eco1 := model.EngineeringChangeOrder{
			ID:                        "ECO-2026-001",
			Code:                      "ECO-2026-001",
			ECRID:                     ecr1.ID,
			Title:                     "Implement LDO Thermal Rework Baseline Increment (REV-A0 ➔ Target REV-B0)",
			SourceHardwareRevisionID:  "REV-001",
			SourceBOMID:               strPtr("BOM-REV-001"),
			ImplementationOwner:       "Marcus Vance",
			Status:                    "APPROVED",
			ImplementationNotes:       "Authorized for immediate baseline generation. Target revision will be designated as REV-B0 with cloned BOM and modified U3 component.",
			ScrapDisposition:          "RUN_OUT",
			RequiresRegressionTesting: true,
			ApprovedAt:                &apprDone,
			CreatedBy:                 "Alex Chen",
			CreatedAt:                 apprDone,
			UpdatedAt:                 now,
		}
		db.Create(&eco1)

		// Seed ECR-002 (Direct Engineering Review Origin in UNDER_REVIEW status)
		ecr2 := model.EngineeringChangeRequest{
			ID:                       "ECR-2026-002",
			Code:                     "ECR-2026-002",
			Title:                    "QSPI Flash Memory Second Source Qualification (Winbond W25Q128JV)",
			Origin:                   "SUPPLIER_CHANGE",
			ProductID:                "PRD-2024-001",
			ProductVersionID:         "VER-001",
			SourceHardwareRevisionID: "REV-001",
			ChangeType:               "BOM_COMPONENT",
			Priority:                 "MEDIUM",
			Description:              "Primary Macronix flash lead time extended to 26 weeks. Propose adding Winbond W25Q128JV as approved alternate / drop-in second source on BOM.",
			BusinessReason:           "Supply chain resilience and lead time reduction for volume pilot production.",
			TechnicalJustification:   "Pin-compatible SOP-8 package with identical command set and timing specifications at 104MHz SPI clock.",
			Status:                   "UNDER_REVIEW",
			RequestedBy:              "David Tanaka",
			AssignedLead:             "Elena Rostova",
			EstimatedCost:            1200.00,
			TargetReleaseDate:        &targetDate,
			CreatedAt:                now.Add(-6 * time.Hour),
			UpdatedAt:                now,
		}
		db.Create(&ecr2)

		items2 := []model.ECRChangeItem{
			{
				ID:                     "CHG-002-01",
				ECRID:                  ecr2.ID,
				Sequence:               1,
				ChangeType:             "BOM_COMPONENT",
				Title:                  "Qualify Winbond W25Q128JV as Second Source Flash",
				AffectedReference:      "U4",
				AffectedComponentID:    strPtr("CMP-004"),
				CurrentState:           "MX25L12835FM2I-10G (Macronix 128Mbit SPI Flash, $0.85)",
				ProposedState:          "W25Q128JVSIQ (Winbond 128Mbit SPI Flash, $0.82)",
				TechnicalJustification: "100% pin-compatible JEDEC standard SPI/QSPI NOR flash memory.",
				RiskLevel:              "LOW",
				Notes:                  "Requires firmware SFDP parameter validation.",
				CreatedAt:              now.Add(-6 * time.Hour),
				UpdatedAt:              now,
			},
		}
		db.Create(&items2)
	}

	log.Println("[INFO] Phase 6 Engineering Change Management seeded successfully.")
}

// SeedDataImport seeds RBAC permissions for the cross-cutting Data Import and Migration capability.
func SeedDataImport(db *gorm.DB) {
	importPermissions := []model.Permission{
		{ID: "PERM-DIMP-VIEW", Module: "dataimport", Action: "view", Code: "dataimport.view", Description: "View data import batches and migration status"},
		{ID: "PERM-DIMP-UPLOAD", Module: "dataimport", Action: "upload", Code: "dataimport.upload", Description: "Upload source XLS/XLSX workbooks for staging"},
		{ID: "PERM-DIMP-VALIDATE", Module: "dataimport", Action: "validate", Code: "dataimport.validate", Description: "Run parsing, normalization and rule validation"},
		{ID: "PERM-DIMP-AI", Module: "dataimport", Action: "ai", Code: "dataimport.ai", Description: "Trigger AI classification and cleansing suggestions"},
		{ID: "PERM-DIMP-REVIEW", Module: "dataimport", Action: "review", Code: "dataimport.review", Description: "Triage and review staged import records"},
		{ID: "PERM-DIMP-APPROVE", Module: "dataimport", Action: "approve", Code: "dataimport.approve", Description: "Authorize and approve validated batches for import"},
		{ID: "PERM-DIMP-EXECUTE", Module: "dataimport", Action: "execute", Code: "dataimport.execute", Description: "Commit approved staged rows into Master Data"},
	}

	for _, p := range importPermissions {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
		// Assign to Super Admin (ROLE-ADMIN)
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{
			RoleID:       "ROLE-ADMIN",
			PermissionID: p.ID,
		})
		// Assign to R&D Manager (ROLE-RND-MGR)
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{
			RoleID:       "ROLE-RND-MGR",
			PermissionID: p.ID,
		})
		// Assign to Hardware Lead (ROLE-HW-LEAD) - all except approve/execute
		if p.Code != "dataimport.approve" && p.Code != "dataimport.execute" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{
				RoleID:       "ROLE-HW-LEAD",
				PermissionID: p.ID,
			})
		}
		// Assign to Viewer (ROLE-VIEWER) - view only
		if p.Code == "dataimport.view" {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.RolePermission{
				RoleID:       "ROLE-VIEWER",
				PermissionID: p.ID,
			})
		}
	}
	log.Println("[INFO] Cross-Cutting Platform Data Import RBAC Permissions seeded successfully.")
}

func strPtr(s string) *string {
	return &s
}





