package pagination

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