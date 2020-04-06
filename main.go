package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type msSQLColumns struct {
	columnName string
	dataType   string
}

func main() {

	// TODO: Use a connection string from the command line
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", "NOT_YET")

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Cannot connect to database.")
	}
	defer db.Close()

	sqlString := `
		SELECT
			SchemaName = c.table_schema,
			TableName = c.table_name,
			ColumnName = c.column_name,
			DataType = data_type,
			IsNullable = c.is_nullable,
			CharacterMaxLength = c.character_maximum_length,
			NumericPrecision = c.numeric_precision
		FROM
			information_schema.columns c
			INNER JOIN information_schema.tables t ON c.table_name = t.table_name
			AND c.table_schema = t.table_schema
			AND t.table_type = 'BASE TABLE'
		ORDER BY
			SchemaName,
			TableName,
			ordinal_position
	`
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		schemaName         string
		tableName          string
		columnName         string
		dataType           string
		isNullable         sql.NullString
		characterMaxLength sql.NullInt64
		numericPrecision   sql.NullInt64
	)

	schemaMap := map[string][]msSQLColumns{}

	for rows.Next() {
		err := rows.Scan(&schemaName, &tableName, &columnName, &dataType, &isNullable, &characterMaxLength, &numericPrecision)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tableName, columnName, dataType, isNullable, characterMaxLength, numericPrecision)

		schemaMap[tableName] = append(schemaMap[tableName], msSQLColumns{columnName, dataType})
	}

	fmt.Println(schemaMap)
}
