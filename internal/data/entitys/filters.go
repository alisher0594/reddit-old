package entitys

import (
	"github.com/alisher0594/reddit-old/internal/validator"
	"math"
	"strings"
)

const promotedCount = 2

// Filters ...
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// SortColumn ...
func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// SortDirection ...
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Limit ...
func (f Filters) Limit() int {
	return f.PageSize
}

// Offset ...
func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// Validate ...
func (f *Filters) Validate(v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// PromotedPerPage ...
func (f Filters) PromotedPerPage() int {
	return promotedCount
}

// PromotedOffset ...
func (f Filters) PromotedOffset() int {
	return (f.Page - 1) * promotedCount
}

// Metadata ...
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// CalculateMetadata ...
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
