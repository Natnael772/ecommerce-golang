package httputil

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetIntQuery extracts an integer query parameter with default value
func GetIntQuery(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil && intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}

// GetStringQuery extracts a string query parameter with default value
func GetStringQuery(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBoolQuery extracts a boolean query parameter with default value
func GetBoolQuery(r *http.Request, key string, defaultValue bool) bool {
	if value := r.URL.Query().Get(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// URLParamInt gets an integer URL parameter with default value (Chi specific)
func URLParamInt(r *http.Request, key string, defaultValue int) int {
	if value := chi.URLParam(r, key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// URLParamString gets a string URL parameter with default value (Chi specific)
func URLParamString(r *http.Request, key, defaultValue string) string {
	if value := chi.URLParam(r, key); value != "" {
		return value
	}
	return defaultValue
}

// PaginationParams represents common pagination query parameters
type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// ParsePaginationParams extracts pagination parameters from request
func ParsePaginationParams(r *http.Request) PaginationParams {
	return PaginationParams{
		Page:    GetIntQuery(r, "page", 1),
		PerPage: GetIntQuery(r, "per_page", 20),
	}
}

// SearchParams represents common search parameters
type SearchParams struct {
	PaginationParams
	Query string `json:"query,omitempty"`
}

// ParseSearchParams extracts search parameters from request
func ParseSearchParams(r *http.Request) SearchParams {
	return SearchParams{
		PaginationParams: ParsePaginationParams(r),
		Query:            GetStringQuery(r, "query", ""),
	}
}

// SortParams represents common sorting parameters
type SortParams struct {
	SortBy string `json:"sort_by,omitempty"`
	Order  string `json:"order,omitempty"` // asc, desc
}

// ParseSortParams extracts sorting parameters from request
func ParseSortParams(r *http.Request) SortParams {
	order := GetStringQuery(r, "order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	
	return SortParams{
		SortBy: GetStringQuery(r, "sort_by", ""),
		Order:  order,
	}
}