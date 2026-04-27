package postgres

import "github.com/densmart/smart-stream/internal/adapters/db"

var PGErrorMap = map[string]db.DBErrorCode{
	"23505": db.DBErrorDuplicateUniqueIndex,
	"23503": db.DBErrorWrongInput,
	"23502": db.DBErrorWrongInput,
	"23514": db.DBErrorColumnValueLessThanZero,
}
