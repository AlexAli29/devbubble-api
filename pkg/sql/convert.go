package sql

import "github.com/jackc/pgx/v5/pgtype"

func UUIDToString(uuid pgtype.UUID) string {
	userIdValue, _ := uuid.Value()
	userIdString, _ := userIdValue.(string)

	return userIdString
}
