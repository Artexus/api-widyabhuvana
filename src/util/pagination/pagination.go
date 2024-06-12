package pagination

type Pagination struct {
	Page   int `json:"page" form:"page"`
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

func (p *Pagination) Paginate() {
	if p.Page == 0 || p.Limit == 0 {
		p.Page = 1
		p.Limit = 10
		p.Offset = 0
	} else {
		p.Offset = (p.Page - 1) * p.Limit
	}
}
