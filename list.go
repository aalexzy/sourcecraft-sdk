package sourcecraft

import (
	"fmt"
	"net/url"
)

// ListOptions contains pagination and filtering options for list operations.
// It supports token-based pagination, page size control, filtering, and sorting.
type ListOptions struct {
	// PageToken is the token for retrieving the next page of results.
	// Leave empty for the first page.
	PageToken string

	// PageSize specifies the maximum number of items to return per page.
	// If zero or negative, the server's default page size will be used.
	PageSize int

	// Filter is a filter expression to narrow down the results.
	// The syntax depends on the specific API endpoint being called.
	Filter string

	// SortBy specifies the field(s) to sort results by.
	// The format depends on the specific API endpoint being called.
	SortBy string
}

func (o ListOptions) getURLQuery() url.Values {
	query := make(url.Values)
	if len(o.PageToken) > 0 {
		query.Add("page_token", fmt.Sprintf("%s", o.PageToken))
	}
	if o.PageSize > 0 {
		query.Add("page_size", fmt.Sprintf("%d", o.PageSize))
	}
	if len(o.Filter) > 0 {
		query.Add("filter", o.Filter)
	}
	if len(o.SortBy) > 0 {
		query.Add("sort_by", o.SortBy)
	}
	return query
}

func (o *ListOptions) setDefaults() {
	if len(o.PageToken) == 0 {
		o.PageSize = 0
		return
	}
}
