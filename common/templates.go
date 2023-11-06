package common

import "sci-review/model"

type PageData struct {
	Title   string
	Active  string
	Message string
	Errors  []ErrorResponse
	User    *model.Principal
}
