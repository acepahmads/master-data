package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("================================================================================")
	fmt.Println(" [AUDIT 2] DATABASE FORENSIC SCHEMA INSPECTION")
	fmt.Println(" Target: MySQL 8.0 Engine (InnoDB) / Database: iot_rd_master")
	fmt.Println("================================================================================")

	dsn := "root:@tcp(127.0.0.1:3306)/iot_rd_master?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("[FAIL] Database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	tables := []string{"imp_batches", "imp_files", "imp_staged_rows"}
	for _, t := range tables {
		fmt.Printf("\n--- Table: %s ---\n", t)
		var tblName, createStmt string
		err := db.QueryRow("SHOW CREATE TABLE " + t).Scan(&tblName, &createStmt)
		if err != nil {
			fmt.Printf("[FAIL] SHOW CREATE TABLE %s error: %v\n", t, err)
			continue
		}
		fmt.Println(createStmt)

		// Check Indexes & Keys
		rows, err := db.Query(fmt.Sprintf("SHOW INDEX FROM %s", t))
		if err == nil {
			fmt.Println("\nIndexes:")
			for rows.Next() {
				var table, nonUnique, keyName, seqInIndex, colName string
				var collation, cardinality, subPart, packed, null, indexType, comment, indexComment, visible, expression sql.NullString
				_ = rows.Scan(&table, &nonUnique, &keyName, &seqInIndex, &colName, &collation, &cardinality, &subPart, &packed, &null, &indexType, &comment, &indexComment, &visible, &expression)
				fmt.Printf("  - Key: %-25s | Column: %-20s | NonUnique: %s | Type: %s\n", keyName, colName, nonUnique, indexType.String)
			}
			rows.Close()
		}
	}

	// Verify Foreign Keys & Constraints
	fmt.Println("\n--- Information Schema Foreign Key Constraints ---")
	fkQuery := `
		SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = 'iot_rd_master' AND TABLE_NAME IN ('imp_batches', 'imp_files', 'imp_staged_rows') AND REFERENCED_TABLE_NAME IS NOT NULL;
	`
	fkRows, err := db.Query(fkQuery)
	if err != nil {
		fmt.Printf("[FAIL] Querying foreign keys error: %v\n", err)
	} else {
		for fkRows.Next() {
			var tbl, cName, col, refTbl, refCol string
			_ = fkRows.Scan(&tbl, &cName, &col, &refTbl, &refCol)
			fmt.Printf("  [FK] %s (%s) ➔ %s (%s) [Constraint: %s]\n", tbl, col, refTbl, refCol, cName)
		}
		fkRows.Close()
	}
	fmt.Println("================================================================================")
}
