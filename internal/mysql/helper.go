package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mklepium/chats/internal/util"
)

func GetAllTables(db *sql.DB, dbName string) {
	query := `SHOW TABLES FROM ` + dbName + `;`
	rows, err := db.Query(query)
	util.CheckErr(err)

	tables := make([]string, 0)
	for rows.Next() {
		var table sql.NullString
		err = rows.Scan(&table)
		util.CheckErr(err)
		tables = append(tables, table.String)
	}
	rows.Close()
	fmt.Println("Tables in the database:")
	for _, table := range tables {
		fmt.Println(table)
	}
}

func GetAllDatabases(db *sql.DB) {
	query := `SHOW DATABASES;`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			log.Fatal(err)
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Databases in the MySQL server:")
	for _, database := range databases {
		fmt.Println(database)
	}
}

func GetStructureForTable(db *sql.DB, dbName string, tableName string) {
	query := `DESCRIBE ` + dbName + `.` + tableName + `;`
	rows, err := db.Query(query)
	util.CheckErr(err)

	var (
		fieldName, fieldType, nullAllowed, key, defaultValue, extraInfo sql.NullString
	)

	columns := make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&fieldName, &fieldType, &nullAllowed, &key, &defaultValue, &extraInfo)
		util.CheckErr(err)
		columns = append(columns, fieldName.String)
	}
	rows.Close()

	fmt.Printf("Columns in the table %v: \n", tableName)
	for _, column := range columns {
		fmt.Println(column)
	}
}
