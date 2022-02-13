package daos

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func getContextExecutor(tx *sql.Tx) (contextExecutor boil.ContextExecutor) {
	if tx == nil {
		contextExecutor = boil.GetContextDB()
	} else {
		contextExecutor = boil.ContextExecutor(tx)
	}
	return contextExecutor
}

func EncodeCursor(id int) string {
	payload := []byte(fmt.Sprint(id))
	cursor := base64.StdEncoding.EncodeToString(payload)
	return cursor
}

func DecodeCursor(after string) (int, error) {
	if after == "" {
		return 0, nil
	}
	payload, err := base64.StdEncoding.DecodeString(after)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(string(payload))
	if err != nil {
		return 0, err
	}
	return id, nil
}
