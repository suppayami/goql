package resolver

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/suppayami/goql/schema"
)

func makeCreator(db *sql.DB, table *schema.SQLTableStruct) func(map[string]interface{}) int64 {
	return func(values map[string]interface{}) int64 {
		var sqlTxt string
		fieldStatement := make([]string, 0)
		valueStatement := make([]string, 0)
		sqlTxt = fmt.Sprintf("INSERT INTO %s", table.Name)
		for key, value := range values {
			key = schema.GraphqlToSQLFieldName(key)
			if len(fmt.Sprintf("%v", value)) == 0 {
				continue
			}
			if _, ok := value.(string); ok == true {
				value = fmt.Sprintf("'%v'", value)
			}
			fieldStatement = append(fieldStatement, key)
			valueStatement = append(valueStatement, fmt.Sprintf("%v", value))
		}
		sqlTxt = fmt.Sprintf(
			"%s (%s) VALUES (%s)",
			sqlTxt,
			strings.Join(fieldStatement, ", "),
			strings.Join(valueStatement, ", "),
		)
		result, err := db.Exec(sqlTxt)
		if err != nil {
			log.Fatal(err)
		}
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		return lastInsertId
	}
}
