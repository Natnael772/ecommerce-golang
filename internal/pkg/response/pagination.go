package response

import (
	"net/http"
	"strconv"
)

// PaginationMeta creates Meta for paginated responses
func PaginationMeta(r *http.Request, page, perPage, total int) Meta {
	pages := calculatePages(total, perPage)
	
	return Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
		Next:    buildNextURL(r, page, perPage, pages),
		Prev:    buildPrevURL(r, page, perPage),
	}
}

// PaginatedSuccess creates a success response with pagination meta
func PaginatedSuccess[T any](w http.ResponseWriter, r *http.Request, data T, page, perPage, total int) {
	meta := PaginationMeta(r, page, perPage, total)
	SuccessWithMeta(w, http.StatusOK, data, meta)
}

func calculatePages(total, perPage int) int {
	if total == 0 || perPage == 0 {
		return 0
	}
	pages := total / perPage
	if total%perPage > 0 {
		pages++
	}
	return pages
}

func buildNextURL(r *http.Request, page, perPage, pages int) string {
	if page >= pages {
		return ""
	}
	return buildPageURL(r, page+1, perPage)
}

func buildPrevURL(r *http.Request, page, perPage int) string {
	if page <= 1 {
		return ""
	}
	return buildPageURL(r, page-1, perPage)
}

func buildPageURL(r *http.Request, page, perPage int) string {
	query := r.URL.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("per_page", strconv.Itoa(perPage))
	return r.URL.Path + "?" + query.Encode()
}