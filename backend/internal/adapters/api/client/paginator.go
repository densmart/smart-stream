package client

import (
	"github.com/densmart/smart-stream/internal/adapters/dto"
)

const (
	defaultPageLimit = 100
	maxPageLimit     = 500
)

type Paginator struct {
	Limit  uint
	Offset uint
	Order  string
	Totals uint64
}

func NewPaginator(data dto.BaseSearchRequestDTO, totals uint64) *Paginator {
	limit := uint(defaultPageLimit)
	if data.Limit != nil {
		limit = *data.Limit
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}
	offset := uint(0)
	if data.Offset != nil {
		offset = *data.Offset
	}

	order := ""
	if data.Order == nil {
		order = "-created_at"
	} else {
		order = *data.Order
	}

	return &Paginator{
		Limit:  limit,
		Offset: offset,
		Order:  order,
		Totals: totals,
	}
}

func (p *Paginator) ToRepresentation() PaginationInfo {
	return PaginationInfo{
		Total: p.Totals,
		Pages: p.GetTotalPages(),
	}
}

func (p *Paginator) GetTotalPages() uint {
	pages, mod := divmod(int(p.Totals), int(p.Limit))
	if mod >= 0 {
		pages++
	}
	return uint(pages)
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}