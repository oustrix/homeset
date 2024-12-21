package store

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oustrix/homeset/internal/models"
)

const (
	UsersTable = "users"
)

var (
	UsersTableColumns = getColumns(models.User{}, UsersTable)
)

func getColumns(model any, tableName string) string {
	r := reflect.Indirect(reflect.ValueOf(model)).Type()

	cols := make([]string, 0, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		colName := r.Field(i).Tag.Get("db")
		if colName == "" || colName == "-" {
			continue
		}

		cols = append(cols, fmt.Sprintf("%s.%s", tableName, colName))
	}

	return strings.Join(cols, ",")
}
