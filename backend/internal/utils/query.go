package utils

type Query struct {
	SearchType string `json:"search_type"`
	Query      struct {
		Term  string `json:"term"`
		Field string `json:"field"`
	} `json:"query,omitempty"`
	From       int      `json:"from"`
	MaxResults int      `json:"max_results"`
	Source     []string `json:"_source"`
	SortFields []string `json:"sort_fields"`
}

func BuildQuery(filter string, page, size int, sort, order string) Query {
	from := (page - 1) * size

	if order == "desc" {
		sort = "-" + sort
	}

	if filter == "" {
		return Query{
			SearchType: "matchall",
			From:       from,
			MaxResults: size,
			Source:     []string{},
			SortFields: []string{sort},
		}
	} else {
		return Query{
			SearchType: "match",
			Query: struct {
				Term  string `json:"term"`
				Field string `json:"field"`
			}{
				Term:  filter,
				Field: "content",
			},
			From:       from,
			MaxResults: size,
			Source:     []string{},
			SortFields: []string{sort},
		}
	}
}
