package utils

import (
	"strconv"
)

// como obtengo un string debo convertirlo a int32 por que asi lo pide la funcion en el sql.go
func ConvertTo32(s string) (int32, error) {
	id64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id64), nil
}
