package resolver

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/suppayami/goql/schema"
)

func makeReader(db *sql.DB, table *schema.SQLTableStruct) func(map[string]interface{}) []map[string]string {
	return func(wheres map[string]interface{}) []map[string]string {
		var sqlTxt string
		whereStatement := make([]string, 0)
		rows := make([]map[string]string, 0)
		sqlTxt = fmt.Sprintf("SELECT * FROM %s", table.Name)
		for key, value := range wheres {
			if len(fmt.Sprintf("%v", value)) == 0 {
				continue
			}
			if _, ok := value.(string); ok == true {
				value = fmt.Sprintf("'%v'", value)
			}
			if !strings.EqualFold(key, "first") && !strings.EqualFold(key, "offset") {
				whereStatement = append(whereStatement, fmt.Sprintf("%s=%v", key, value))
			}
		}
		if len(whereStatement) > 0 {
			sqlTxt = fmt.Sprintf("%s WHERE", sqlTxt)
			sqlTxt = fmt.Sprintf("%s %s", sqlTxt, strings.Join(whereStatement, " AND "))
		}
		if first, ok := wheres["first"]; ok {
			sqlTxt = fmt.Sprintf("%s LIMIT %v", sqlTxt, first)
			if offset, ok := wheres["offset"]; ok {
				sqlTxt = fmt.Sprintf("%s OFFSET %v", sqlTxt, offset)
			}
		}
		fmt.Println(sqlTxt)
		sqlRows, err := db.Query(sqlTxt)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlRows.Close()
		cols, err := sqlRows.Columns()
		if err != nil {
			log.Fatal(err)
		}
		for sqlRows.Next() {
			columns := make([]sql.NullString, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}
			if err := sqlRows.Scan(columnPointers...); err != nil {
				log.Fatal(err)
			}
			m := make(map[string]string)
			for i, colName := range cols {
				val := columnPointers[i].(*sql.NullString)
				m[schema.SQLToGraphqlFieldName(colName)] = val.String
			}
			rows = append(rows, m)
		}
		return rows
	}
}
