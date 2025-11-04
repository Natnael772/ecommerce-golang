package pagination

import (
	"net/http"
	"strconv"
)

// Generic pagination that any app could use
type Pagination struct {
    Page    int `json:"page"`
    PerPage int `json:"perPage"`
    Total   int `json:"total"`
}

func (p Pagination) Offset() int {
    return (p.Page - 1) * p.PerPage
}

func New(page, perPage int) Pagination {
    if page <= 0 {
        page = 1
    }
    if perPage <= 0 {
        perPage = 10
    }

    return Pagination{
        Page:    page,
        PerPage: perPage,
    }
}

// GetPaginationParams extracts pagination parameters from the HTTP request
func GetPaginationParams(r *http.Request) (page int, perPage int) {
	// Default values (can be adjusted as needed)
	page = 1
	perPage = 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if val, err := strconv.Atoi(pp); err == nil && val > 0 {
			perPage = val
		}
	}

	return
}