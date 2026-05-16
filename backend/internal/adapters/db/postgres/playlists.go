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
	playlistsTable = "playlists"
)

func (dbc *PgDB) CreatePlaylist(data dto.CreatePlaylistDTO) (models.Playlist, *db.DBError) {
	var result models.Playlist

	query, err := createPlaylistQuery(data)
	if err != nil {
		logger.Debugf("[postgres] CreatePlaylist error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] CreatePlaylist error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Name,
		&result.Type, &result.Poster, &result.ParentID, &result.HasChildren); err != nil {
		logger.Debugf("[postgres] CreatePlaylist error scan result: %s", err.Error())
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

	// Update parent's has_children flag if parent_id is set
	if data.ParentID != nil {
		if dbErr := dbc.UpdateParentHasChildren(*data.ParentID); dbErr != nil {
			logger.Debugf("[postgres] CreatePlaylist error update parent has_children: %v", dbErr)
			// Don't fail the operation, just log the error
		}
	}

	return result, nil
}

func (dbc *PgDB) RetrievePlaylist(id string) (models.Playlist, *db.DBError) {
	var result models.Playlist

	query, err := retrievePlaylistQuery(id)
	if err != nil {
		logger.Debugf("[postgres] RetrievePlaylist error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] RetrievePlaylist error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Name,
		&result.Type, &result.Poster, &result.ParentID, &result.HasChildren); err != nil {
		logger.Debugf("[postgres] RetrievePlaylist error scan result: %s", err.Error())
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

func (dbc *PgDB) UpdatePlaylist(id string, data dto.UpdatePlaylistDTO) (models.Playlist, *db.DBError) {
	var result models.Playlist

	// Get current playlist to check if parent_id changed
	var oldParentID *string
	if data.ParentID != nil {
		currentPlaylist, dbErr := dbc.RetrievePlaylist(id)
		if dbErr != nil {
			return result, dbErr
		}
		oldParentID = currentPlaylist.ParentID
	}

	query, err := updatePlaylistQuery(id, data)
	if err != nil {
		logger.Debugf("[postgres] UpdatePlaylist error create query: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] UpdatePlaylist error acquire conn: %s", err.Error())
		return result, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	row := conn.QueryRow(dbc.ctx, query)
	if err = row.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Name,
		&result.Type, &result.Poster, &result.ParentID, &result.HasChildren); err != nil {
		logger.Debugf("[postgres] UpdatePlaylist error scan result: %s", err.Error())
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

	// Update parent's has_children flag if parent_id changed
	if data.ParentID != nil {
		// Update old parent if it exists and is different from new parent
		if oldParentID != nil && (result.ParentID == nil || *oldParentID != *result.ParentID) {
			if dbErr := dbc.UpdateParentHasChildren(*oldParentID); dbErr != nil {
				logger.Debugf("[postgres] UpdatePlaylist error update old parent has_children: %v", dbErr)
			}
		}
		// Update new parent if it exists
		if result.ParentID != nil {
			if dbErr := dbc.UpdateParentHasChildren(*result.ParentID); dbErr != nil {
				logger.Debugf("[postgres] UpdatePlaylist error update new parent has_children: %v", dbErr)
			}
		}
	}

	return result, nil
}

func (dbc *PgDB) DeletePlaylist(id string) *db.DBError {
	query, err := deletePlaylistQuery(id)
	if err != nil {
		logger.Debugf("[postgres] DeletePlaylist error create query: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] DeletePlaylist error acquire conn: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	_, err = conn.Exec(dbc.ctx, query)
	if err != nil {
		logger.Debugf("[postgres] DeletePlaylist error exec query: %s", err.Error())
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

func (dbc *PgDB) SearchPlaylists(filter dto.SearchPlaylistsDTO) ([]models.Playlist, uint64, *db.DBError) {
	var rowsCount uint64
	var results []models.Playlist

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] SearchPlaylists error acquire conn: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	tx, err := conn.BeginTx(dbc.ctx, pgx.TxOptions{})
	if err != nil {
		logger.Debugf("[postgres] SearchPlaylists error begin tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	countQuery, err := countPlaylistsQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchPlaylists error create count query: %s", err.Error())
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
		logger.Debugf("[postgres] SearchPlaylists error scan count result: %s", err.Error())
		if rbErr := tx.Rollback(dbc.ctx); rbErr != nil {
			logger.Errorf("error rollback tx: %s", rbErr.Error())
		}
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	listQuery, err := limitedListPlaylistsQuery(filter)
	if err != nil {
		logger.Debugf("[postgres] SearchPlaylists error create list query: %s", err.Error())
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
		logger.Debugf("[postgres] SearchPlaylists error exec list query: %s", err.Error())
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
		var result models.Playlist
		if err = rows.Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt, &result.Name,
			&result.Type, &result.Poster, &result.ParentID, &result.HasChildren); err != nil {
			logger.Debugf("[postgres] SearchPlaylists error scan list result: %s", err.Error())
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
		logger.Debugf("[postgres] SearchPlaylists error commit tx: %s", err.Error())
		return results, 0, &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	return results, rowsCount, nil
}

func createPlaylistQuery(data dto.CreatePlaylistDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	record := goqu.Record{
		"created_at":   time.Now().UTC(),
		"updated_at":   time.Now().UTC(),
		"name":         data.Name,
		"type":         data.Type,
		"has_children": false,
	}

	if data.Poster != nil {
		record["poster"] = *data.Poster
	}
	if data.ParentID != nil {
		record["parent_id"] = *data.ParentID
	}

	ds := dialect.Insert(playlistsTable).Rows(record).Returning("id", "created_at", "updated_at", "name",
		"type", "poster", "parent_id", "has_children")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] createPlaylistQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func retrievePlaylistQuery(id string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "name", "type", "poster", "parent_id", "has_children").
		From(playlistsTable).Where(goqu.Ex{"id": id})

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] retrievePlaylistQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func updatePlaylistQuery(id string, data dto.UpdatePlaylistDTO) (string, error) {
	dialect := goqu.Dialect("postgres")
	records := goqu.Record{
		"updated_at": time.Now().UTC(),
	}

	if data.Name != nil {
		records["name"] = *data.Name
	}
	if data.Type != nil {
		records["type"] = *data.Type
	}
	if data.Poster != nil {
		records["poster"] = *data.Poster
	}
	if data.ParentID != nil {
		records["parent_id"] = *data.ParentID
	}

	ds := dialect.Update(playlistsTable).Set(records).Where(goqu.Ex{"id": id}).
		Returning("id", "created_at", "updated_at", "name", "type", "poster", "parent_id", "has_children")
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] updatePlaylistQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func deletePlaylistQuery(id string) (string, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Delete(playlistsTable).Where(goqu.Ex{"id": id})
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] deletePlaylistQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func countPlaylistsQuery(data dto.SearchPlaylistsDTO) (string, error) {
	ds := listPlaylistsQuery(data)

	ds = ds.Select(goqu.COUNT("*").As("rows_count"))
	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] countPlaylistsQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func limitedListPlaylistsQuery(data dto.SearchPlaylistsDTO) (string, error) {
	ds := listPlaylistsQuery(data)

	ds = ds.Select("id", "created_at", "updated_at", "name", "type", "poster", "parent_id", "has_children")
	ds = sort(ds, data.Order)
	ds = paginate(ds, data.Limit, data.Offset)

	sql, _, err := ds.ToSQL()
	logger.Debugf("[postgres] limitedListPlaylistsQuery: %s", sql)
	if err != nil {
		return "", err
	}
	return sql, nil
}

func listPlaylistsQuery(data dto.SearchPlaylistsDTO) *goqu.SelectDataset {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select("id", "created_at", "updated_at", "name", "type", "poster", "parent_id", "has_children").From(playlistsTable)

	ds = filterByDatetime(ds, data.CreatedAtFrom, data.CreatedAtTo, "")

	if data.ID != nil {
		ds = ds.Where(goqu.Ex{"id": *data.ID})
	}
	if data.Name != nil {
		ds = ds.Where(goqu.L("name ILIKE ?", "%"+*data.Name+"%"))
	}
	if data.Type != nil {
		ds = ds.Where(goqu.Ex{"type": *data.Type})
	}

	// Filter by parent_id: if IncludeAllHierarchy is true, don't filter by parent_id at all
	if data.IncludeAllHierarchy != nil && *data.IncludeAllHierarchy {
		// Don't apply parent_id filter - load all playlists
	} else if data.ParentID != nil {
		// Filter by specific parent_id
		ds = ds.Where(goqu.Ex{"parent_id": *data.ParentID})
	} else {
		// Default: show only root-level playlists (parent_id IS NULL)
		ds = ds.Where(goqu.Ex{"parent_id": nil})
	}

	return ds
}

// UpdateParentHasChildren updates has_children flag for parent playlist
func (dbc *PgDB) UpdateParentHasChildren(parentID string) *db.DBError {
	// Count children playlists
	dialect := goqu.Dialect("postgres")
	ds := dialect.Select(goqu.COUNT("*").As("count")).
		From(playlistsTable).
		Where(goqu.Ex{"parent_id": parentID})

	countQuery, _, err := ds.ToSQL()
	if err != nil {
		logger.Debugf("[postgres] updateParentHasChildren error create count query: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	conn, err := dbc.pool.Acquire(dbc.ctx)
	if err != nil {
		logger.Debugf("[postgres] updateParentHasChildren error acquire conn: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}
	defer conn.Release()

	var count int
	row := conn.QueryRow(dbc.ctx, countQuery)
	if err = row.Scan(&count); err != nil {
		logger.Debugf("[postgres] updateParentHasChildren error scan count: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	// Update has_children flag
	hasChildren := count > 0
	updateDS := dialect.Update(playlistsTable).
		Set(goqu.Record{"has_children": hasChildren, "updated_at": time.Now().UTC()}).
		Where(goqu.Ex{"id": parentID})

	updateQuery, _, err := updateDS.ToSQL()
	if err != nil {
		logger.Debugf("[postgres] updateParentHasChildren error create update query: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorWrongInput,
			Message: db.DBErrorWrongInput.String(),
		}
	}

	_, err = conn.Exec(dbc.ctx, updateQuery)
	if err != nil {
		logger.Debugf("[postgres] updateParentHasChildren error exec update: %s", err.Error())
		return &db.DBError{
			Code:    db.DBErrorGeneral,
			Message: db.DBErrorGeneral.String(),
		}
	}

	logger.Debugf("[postgres] updateParentHasChildren: parent_id=%s, has_children=%v", parentID, hasChildren)
	return nil
}