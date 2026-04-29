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
	mediaTable = "media"
)

func (dbc *PgDB) CreateMedia(data dto.CreateMediaDTO) (models.Media, *db.DBError) {
	var result models.Media

	query, err := createMediaQuery(data)
	if err != nil {
		logger.Debugf("[postgres] CreateMedia error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] CreateMedia error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.PlaylistID, &result.Name,
		&result.Poster, &result.Format, &result.Path, &result.Duration, &result.Size, &result.Order); err != nil {
		logger.Debugf("[postgres] CreateMedia error scan result: %s", err.Error())
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

func (dbc *PgDB) RetrieveMedia(id string) (models.Media, *db.DBError) {
	var result models.Media

	query, err := retrieveMediaQuery(id)
	if err != nil {
		logger.Debugf("[postgres] RetrieveMedia error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] RetrieveMedia error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.PlaylistID, &result.Name,
		&result.Poster, &result.Format, &result.Path, &result.Duration, &result.Size, &result.Order); err != nil {
		logger.Debugf("[postgres] RetrieveMedia error scan result: %s", err.Error())
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

func (dbc *PgDB) UpdateMedia(id string, data dto.UpdateMediaDTO) (models.Media, *db.DBError) {
	var result models.Media

	query, err := updateMediaQuery(id, data)
	if err != nil {
		logger.Debugf("[postgres] UpdateMedia error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] UpdateMedia error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.PlaylistID, &result.Name,
		&result.Poster, &result.Format, &result.Path, &result.Duration, &result.Size, &result.Order); err != nil {
		logger.Debugf("[postgres] UpdateMedia error scan result: %s", err.Error())
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

func (dbc *PgDB) DeleteMedia(id string) *db.DBError {
	query, err := deleteMediaQuery(id)
	if err != nil {
		logger.Debugf("[postgres] DeleteMedia error create query: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] DeleteMedia error acquire conn: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	_, err = conn.Exec(dbc.ctx, query)
	if err != nil {
		logger.Debugf("[postgres] DeleteMedia error exec query: %s", err.Error())
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

func (dbc *PgDB) SearchMedia(filter dto.SearchMediaDTO) ([]models.Media, uint64, *db.DBError) {
	var rowsCount uint64
	var results []models.Media

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] SearchMedia error acquire conn: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	tx, err := conn.BeginTx(dbc.ctx, pgx.TxOptions{})
	if err != nil {
		logger.Debugf("[postgres] SearchMedia error begin tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	countQuery, err := countMediaQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchMedia error create count query: %s", err.Error())
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
		logger.Debugf("[postgres] SearchMedia error scan count result: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	listQuery, err := limitedListMediaQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchMedia error create list query: %s", err.Error())
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
		logger.Debugf("[postgres] SearchMedia error exec list query: %s", err.Error())
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
		var result models.Media
		if err = rows.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.PlaylistID, &result.Name,
			&result.Poster, &result.Format, &result.Path, &result.Duration, &result.Size, &result.Order); err != nil {
			logger.Debugf("[postgres] SearchMedia error scan list result: %s", err.Error())
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
		logger.Debugf("[postgres] SearchMedia error commit tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	return results, rowsCount, nil
}

func createMediaQuery(data dto.CreateMediaDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	record := goqu.Record{
		"created_at": time.Now().UTC(),
		"updated_at": time.Now().UTC(),
		"name":       data.Name,
		"format":     data.Format,
		"path":       data.Path,
		"duration":   data.Duration,
		"size":       data.Size,
		"order":      0,
	}

	if data.PlaylistID != nil {
		record["playlist_id"] = *data.PlaylistID
	}
	if data.Poster != nil {
		record["poster"] = *data.Poster
	}
	if data.Order != nil {
		record["order"] = *data.Order
	}

	ds := dialect.Insert(mediaTable).Rows(record).Returning("id", "created_at", "updated_at", "playlist_id",
		"name", "poster", "format", "path", "duration", "size", "order")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] createMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func retrieveMediaQuery(id string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "playlist_id", "name", "poster", "format",
		"path", "duration", "size", "order").From(mediaTable).Where(goqu.Ex{"id": id})

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] retrieveMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func updateMediaQuery(id string, data dto.UpdateMediaDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	records := goqu.Record{
		"updated_at": time.Now().UTC(),
	}

	if data.PlaylistID != nil {
		records["playlist_id"] = *data.PlaylistID
	}
	if data.Name != nil {
		records["name"] = *data.Name
	}
	if data.Poster != nil {
		records["poster"] = *data.Poster
	}
	if data.Format != nil {
		records["format"] = *data.Format
	}
	if data.Path != nil {
		records["path"] = *data.Path
	}
	if data.Duration != nil {
		records["duration"] = *data.Duration
	}
	if data.Size != nil {
		records["size"] = *data.Size
	}
	if data.Order != nil {
		records["order"] = *data.Order
	}

	ds := dialect.Update(mediaTable).Set(records).Where(goqu.Ex{"id": id}).Returning("id", "created_at",
		"updated_at", "playlist_id", "name", "poster", "format", "path", "duration", "size", "order")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] updateMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func deleteMediaQuery(id string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Delete(mediaTable).Where(goqu.Ex{"id": id})
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] deleteMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func countMediaQuery(data dto.SearchMediaDTO) (string, error) {
	ds := listMediaQuery(data)

	ds = ds.Select(goqu.COUNT("*").As("rows_count"))
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] countMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func limitedListMediaQuery(data dto.SearchMediaDTO) (string, error) {
	ds := listMediaQuery(data)

	ds = ds.Select("id", "created_at", "updated_at", "playlist_id", "name", "poster", "format",
		"path", "duration", "size", "order")
	ds = sort(ds, data.Order)
	ds = paginate(ds, data.Limit, data.Offset)

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] limitedListMediaQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func listMediaQuery(data dto.SearchMediaDTO) *goqu.SelectDataset {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "playlist_id", "name", "poster", "format",
		"path", "duration", "size", "order").From(mediaTable)

	ds = filterByDatetime(ds, data.CreatedAtFrom, data.CreatedAtTo, "")

	if data.ID != nil {
		ds = ds.Where(goqu.Ex{"id": *data.ID})
	}
	if data.PlaylistID != nil {
		ds = ds.Where(goqu.Ex{"playlist_id": *data.PlaylistID})
	}
	if data.Name != nil {
		ds = ds.Where(goqu.Ex{"name": *data.Name})
	}
	if data.Format != nil {
		ds = ds.Where(goqu.Ex{"format": *data.Format})
	}
	// Фильтр для медиа без плейлистов
	if data.OnlyUnassigned != nil && *data.OnlyUnassigned {
		ds = ds.Where(goqu.Ex{"playlist_id": nil})
	}

	return ds
}