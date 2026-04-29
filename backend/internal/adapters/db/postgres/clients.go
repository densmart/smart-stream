package postgres

import (
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/pkg/logger"
)

const (
	clientsTable = "clients"
)

func (dbc *PgDB) CreateClient(data dto.CreateClientDTO) (models.Client, *db.DBError) {
	var result models.Client

	query, err := createClientQuery(data)
	if err != nil {
		logger.Debugf("[postgres] CreateClient error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] CreateClient error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Login, &result.Password, &result.Email,
		&result.IsActive, &result.CurrentVersion, &result.LastLoginAt); err != nil {
		logger.Debugf("[postgres] CreateClient error scan result: %s", err.Error())
		code := db.DBErrorGeneral
		if pgErr, ok := err.(*pgconn.PgError); ok {
			code, ok = PGErrorMap[pgErr.Code]
			if !ok {
				code = db.DBErrorGeneral
			}
		}
		return result, &db.DBError{
			Code:    code,
			Message: code.String(),
		}
	}

	return result, nil
}

func (dbc *PgDB) RetrieveClient(id string) (models.Client, *db.DBError) {
	var result models.Client

	query, err := retrieveClientQuery(&id, nil)
	if err != nil {
		logger.Debugf("[postgres] RetrieveClient error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] RetrieveClient error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Login, &result.Password, &result.Email,
		&result.IsActive, &result.CurrentVersion, &result.LastLoginAt); err != nil {
		logger.Debugf("[postgres] RetrieveClient error scan result: %s", err.Error())
		code := db.DBErrorGeneral
		if pgErr, ok := err.(*pgconn.PgError); ok {
			code, ok = PGErrorMap[pgErr.Code]
			if !ok {
				code = db.DBErrorGeneral
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			code = db.DBErrorNotFound
		}
		return result, &db.DBError{
			Code:    code,
			Message: code.String(),
		}
	}

	return result, nil
}

func (dbc *PgDB) RetrieveClientByLogin(login string) (models.Client, *db.DBError) {
	var result models.Client

	query, err := retrieveClientQuery(nil, &login)
	if err != nil {
		logger.Debugf("[postgres] RetrieveClientByLogin error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] RetrieveClientByLogin error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Login, &result.Password, &result.Email,
		&result.IsActive, &result.CurrentVersion, &result.LastLoginAt); err != nil {
		logger.Debugf("[postgres] RetrieveClientByLogin error scan result: %s", err.Error())
		code := db.DBErrorGeneral
		if pgErr, ok := err.(*pgconn.PgError); ok {
			code, ok = PGErrorMap[pgErr.Code]
			if !ok {
				code = db.DBErrorGeneral
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			code = db.DBErrorNotFound
		}
		return result, &db.DBError{
			Code:    code,
			Message: code.String(),
		}
	}

	return result, nil
}

func (dbc *PgDB) UpdateClient(id string, data dto.UpdateClientDTO) (models.Client, *db.DBError) {
	var result models.Client

	query, err := updateClientQuery(id, data)
	if err != nil {
		logger.Debugf("[postgres] UpdateClient error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] UpdateClient error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Login, &result.Password, &result.Email,
		&result.IsActive, &result.CurrentVersion, &result.LastLoginAt); err != nil {
		logger.Debugf("[postgres] UpdateClient error scan result: %s", err.Error())
		code := db.DBErrorGeneral
		if pgErr, ok := err.(*pgconn.PgError); ok {
			code, ok = PGErrorMap[pgErr.Code]
			if !ok {
				code = db.DBErrorGeneral
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			code = db.DBErrorNotFound
		}
		return result, &db.DBError{
			Code:    code,
			Message: code.String(),
		}
	}

	return result, nil
}

func (dbc *PgDB) DeleteClient(id string) *db.DBError {
	query, err := deleteClientQuery(id)
	if err != nil {
		logger.Debugf("[postgres] DeleteClient error create query: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] DeleteClient error acquire conn: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	_, err = conn.Exec(dbc.ctx, query)
	if err != nil {
		logger.Debugf("[postgres] DeleteClient error exec query: %s", err.Error())
		code := db.DBErrorGeneral
		if pgErr, ok := err.(*pgconn.PgError); ok {
			code, ok = PGErrorMap[pgErr.Code]
			if !ok {
				code = db.DBErrorGeneral
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			code = db.DBErrorNotFound
		}
		return &db.DBError{
			Code:    code,
			Message: code.String(),
		}
	}

	return nil
}

func (dbc *PgDB) SearchClients(filter dto.SearchClientsDTO) ([]models.Client, uint64, *db.DBError) {
	var rowsCount uint64
	var results []models.Client

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] SearchClients error acquire conn: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	tx, err := conn.BeginTx(dbc.ctx, pgx.TxOptions{})
	if err != nil {
		logger.Debugf("[postgres] SearchClients error begin tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	countQuery, err := countClientsQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchClients error create count query: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	row := tx.QueryRow(dbc.ctx, countQuery)
	if err = row.Scan(&rowsCount); err != nil {
		logger.Debugf("[postgres] SearchClients error scan count result: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	listQuery, err := limitedListClientsQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchClients error create list query: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	rows, err := tx.Query(dbc.ctx, listQuery)
	if err != nil {
		logger.Debugf("[postgres] SearchClients error exec list query: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		var result models.Client
		if err = rows.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Login, &result.Password, &result.Email,
			&result.IsActive, &result.CurrentVersion, &result.LastLoginAt); err != nil {
			logger.Debugf("[postgres] SearchClients error scan list result: %s", err.Error())
			if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
				logger.Errorf("error rollback tx: %s", rbErr.Error())
			}
			return results, 0, &db.DBError{
				Code:    db.DBErrorGeneral,
				Message: db.DBErrorGeneral.String(),
			}
		}
		results = append(results, result)
	}

	if err = tx.Commit(dbc.ctx); err != nil {
		logger.Debugf("[postgres] SearchClients error commit tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	return results, rowsCount, nil
}

func createClientQuery(data dto.CreateClientDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	record := goqu.Record{
		"created_at": time.Now().UTC(),
		"updated_at": time.Now().UTC(),
		"login":      data.Login,
		"password":   data.Password,
		"is_active":  data.IsActive,
	}

	if data.Email != nil {
		record["email"] = *data.Email
	}

	ds := dialect.Insert(clientsTable).Rows(record).Returning("id", "created_at", "updated_at", "login",
		"password", "email", "is_active", "current_version", "last_login_at")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] createClientQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func retrieveClientQuery(id *string, login *string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "login", "password", "email", "is_active",
		"current_version", "last_login_at").From(clientsTable)

	if id != nil {
		ds = ds.Where(goqu.Ex{"id": *id})
	} else if login != nil {
		ds = ds.Where(goqu.Ex{"login": *login})
	}

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] retrieveClientQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func updateClientQuery(id string, data dto.UpdateClientDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	records := goqu.Record{
		"updated_at": time.Now().UTC(),
	}

	if data.Login != nil {
		records["login"] = *data.Login
	}
	if data.Password != nil {
		records["password"] = *data.Password
	}
	if data.Email != nil {
		records["email"] = *data.Email
	}
	if data.IsActive != nil {
		records["is_active"] = *data.IsActive
	}
	if data.LastLoginAt != nil {
		if data.LastLoginAt.IsZero() {
			records["last_login_at"] = nil
		} else {
			records["last_login_at"] = *data.LastLoginAt
		}
	}
	if data.CurrentVersion != nil {
		records["current_version"] = *data.CurrentVersion
	}

	ds := dialect.Update(clientsTable).Set(records).Where(goqu.Ex{"id": id}).Returning("id", "created_at",
		"updated_at", "login", "password", "email", "is_active", "current_version", "last_login_at")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] updateClientQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func deleteClientQuery(id string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Delete(clientsTable).Where(goqu.Ex{"id": id})
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] deleteClientQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func countClientsQuery(data dto.SearchClientsDTO) (string, error) {
	ds := listClientsQuery(data)

	ds = ds.Select(goqu.COUNT("*").As("rows_count"))
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] countClientsQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func limitedListClientsQuery(data dto.SearchClientsDTO) (string, error) {
	ds := listClientsQuery(data)

	ds = ds.Select("id", "created_at", "updated_at", "login", "password", "email", "is_active",
		"current_version", "last_login_at")
	ds = sort(ds, data.Order)
	ds = paginate(ds, data.Limit, data.Offset)

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] limitedListClientsQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func listClientsQuery(data dto.SearchClientsDTO) *goqu.SelectDataset {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "login", "password", "email", "is_active",
		"current_version", "last_login_at").From(clientsTable)

	ds = filterByDatetime(ds, data.CreatedAtFrom, data.CreatedAtTo, "")

	if data.ID != nil {
		ds = ds.Where(goqu.Ex{"id": *data.ID})
	}
	if data.Login != nil {
		ds = ds.Where(goqu.Ex{"login": *data.Login})
	}
	if data.Email != nil {
		ds = ds.Where(goqu.Ex{"email": *data.Email})
	}
	if data.IsActive != nil {
		ds = ds.Where(goqu.Ex{"is_active": *data.IsActive})
	}
	if data.LastLoginAtFrom != nil {
		ds = ds.Where(goqu.Ex{"last_login_at": goqu.Op{"gte": *data.LastLoginAtFrom}})
	}
	if data.LastLoginAtTo != nil {
		ds = ds.Where(goqu.Ex{"last_login_at": goqu.Op{"lte": *data.LastLoginAtTo}})
	}
	if data.CurrentVersion != nil {
		ds = ds.Where(goqu.Ex{"current_version": *data.CurrentVersion})
	}

	return ds
}
