package scalingo

import "strconv"

type PaginationOpts struct {
	Page    int
	PerPage int
}

func (opts PaginationOpts) ToMap() map[string]string {
	return map[string]string{
		"page":     strconv.Itoa(opts.Page),
		"per_page": strconv.Itoa(opts.PerPage),
	}
}

type PaginationMeta struct {
	PrevPage    int `json:"prev_page"`
	CurrentPage int `json:"current_page"`
	NextPage    int `json:"next_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}
