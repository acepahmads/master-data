package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/iot_rd_master?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("--- PRODUCT VERSIONS ---")
	rows, _ := db.Query("SELECT id, product_code, version_number, status, deleted_at FROM product_versions")
	for rows.Next() {
		var id, pc, vn, st string
		var del sql.NullString
		_ = rows.Scan(&id, &pc, &vn, &st, &del)
		fmt.Printf("ID: %s | Prod: %s | Ver: %s | Status: %s | Deleted: %v\n", id, pc, vn, st, del.Valid)
	}
	rows.Close()

	fmt.Println("--- HARDWARE REVISIONS ---")
	rows, _ = db.Query("SELECT id, product_version_id, code, name, status, deleted_at FROM hardware_revisions")
	for rows.Next() {
		var id, pvid, code, name, st string
		var del sql.NullString
		_ = rows.Scan(&id, &pvid, &code, &name, &st, &del)
		fmt.Printf("ID: %s | VerID: %s | Code: %s | Name: %s | Status: %s | Deleted: %v\n", id, pvid, code, name, st, del.Valid)
	}
	rows.Close()
}
