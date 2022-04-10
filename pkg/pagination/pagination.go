package pagination

import "math"

type Query struct {
	Size    int    `json:"-" form:"size"`
	Page    int    `json:"page,omitempty" form:"page"`
}

func (q *Query) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

func (q *Query) GetLimit() int {
	return q.Size
}

func (q *Query) GetPage() int {
	return q.Page
}

func (q *Query) GetSize() int {
	return q.Size
}

func GetTotalPages(totalCount int64, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
