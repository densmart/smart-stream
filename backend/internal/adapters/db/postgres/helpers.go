package postgres

import (
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
)

const (
	defaultLimit = 500
)

func filterByDatetime(ds *goqu.SelectDataset, createdAtFrom, createdAtTo time.Time, prefix string) *goqu.SelectDataset {
	if !createdAtFrom.IsZero() {
		ds = ds.Where(goqu.I(prefix + "created_at").Gte(createdAtFrom))
	}
	if !createdAtTo.IsZero() {
		ds = ds.Where(goqu.I(prefix + "created_at").Lte(createdAtTo))
	}
	return ds
}

func sort(ds *goqu.SelectDataset, sortBy *string) *goqu.SelectDataset {
	if sortBy != nil {
		sortField := *sortBy
		if strings.HasPrefix(sortField, "-") {
			sortField = strings.TrimPrefix(sortField, "-")
			ds = ds.Order(goqu.I(sortField).Desc())
		} else {
			ds = ds.Order(goqu.I(sortField).Asc())
		}
	} else {
		ds = ds.Order(goqu.I("created_at").Desc())
	}
	return ds
}

func paginate(ds *goqu.SelectDataset, limit, offset *uint) *goqu.SelectDataset {
	if limit != nil {
		ds = ds.Limit(*limit)
	} else {
		ds = ds.Limit(defaultLimit)
	}
	if offset != nil {
		ds = ds.Offset(*offset)
	}
	return ds
}
