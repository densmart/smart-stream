package web

import (
	"regexp"
	"strconv"
	"strings"

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
	RawURL string
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
		RawURL: data.RawURL,
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

func (p *Paginator) ToLinkHeader() string {
	var link []string

	if p.MakeNextLink() != "" {
		link = append(link, "<"+p.MakeNextLink()+">; rel=\"next\"")
	}
	if p.MakePrevLink() != "" {
		link = append(link, "<"+p.MakePrevLink()+">; rel=\"prev\"")
	}
	return strings.Join(link, ",\n")
}

func (p *Paginator) GetOffset() uint {
	return p.Offset
}

func (p *Paginator) GetLimit() uint {
	return p.Limit
}

func (p *Paginator) GetOrder() string {
	return p.Order
}

func (p *Paginator) GetOrderDit() string {
	if strings.HasPrefix(p.Order, "-") {
		return "desc"
	}
	return "asc"
}

func (p *Paginator) GetTotalPages() uint {
	pages, mod := divmod(int(p.Totals), int(p.Limit))
	if mod >= 0 {
		pages++
	}
	return uint(pages)
}

func (p *Paginator) MakeNextLink() string {
	totalPages := int(p.GetTotalPages())
	r := regexp.MustCompile("offset=([0-9]+)")
	if r.MatchString(p.RawURL) {
		currentOffset, err := strconv.Atoi(r.FindStringSubmatch(p.RawURL)[1])
		if err != nil {
			return ""
		}
		nextOffset := currentOffset + int(p.Limit)
		if uint64(nextOffset) > p.Totals {
			return ""
		}
		return r.ReplaceAllString(p.RawURL, "offset="+strconv.Itoa(nextOffset))
	} else if totalPages > 1 {
		return p.RawURL + "&offset=" + strconv.Itoa(int(p.Limit))
	}
	return ""
}

func (p *Paginator) MakePrevLink() string {
	r := regexp.MustCompile("&offset=([0-9]+)")
	if r.MatchString(p.RawURL) {
		currentOffset, err := strconv.Atoi(r.FindStringSubmatch(p.RawURL)[1])
		if err != nil {
			return ""
		}
		nextOffset := currentOffset - int(p.Limit)
		if int64(nextOffset) < 0 {
			return ""
		}
		return r.ReplaceAllString(p.RawURL, "&offset="+strconv.Itoa(nextOffset))
	}

	return ""
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
