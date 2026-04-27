package db

const (
	DBErrorGeneral DBErrorCode = iota
	DBErrorWrongInput
	DBErrorNotFound
	DBErrorDuplicateUniqueIndex
	DBErrorColumnValueLessThanZero
)

type DBErrorCode uint

func (dbec DBErrorCode) String() string {
	strings := []string{
		"general DB Error",
		"wrong Input Data",
		"record Not Found",
		"duplicate index entry",
		"column value less than zero",
	}
	if dbec >= DBErrorCode(len(strings)) {
		return "Unknown"
	}
	return strings[dbec]
}

type DBError struct {
	Code    DBErrorCode
	Message string
}

func (dbe *DBError) Error() string {
	return dbe.Message
}
