package common

type ProblemDetail struct {
	Type     string  `json:"type"`
	Title    string  `json:"title"`
	Status   int     `json:"status"`
	Detail   string  `json:"detail"`
	Instance string  `json:"instance"`
	Fields   []Field `json:"fields"`
}

func NewProblemDetail(title string, status int) *ProblemDetail {
	return &ProblemDetail{
		Type:     "",
		Title:    title,
		Status:   status,
		Detail:   title,
		Instance: "",
		Fields:   []Field{},
	}
}

func InvalidJson() *ProblemDetail {
	return &ProblemDetail{
		Type:     "",
		Title:    "Invalid Json format",
		Status:   400,
		Detail:   "Invalid Json format",
		Instance: "",
		Fields:   []Field{},
	}
}

func ProblemWithErrors(errors []Field) *ProblemDetail {
	return &ProblemDetail{
		Type:     "",
		Title:    "Bad Request",
		Status:   400,
		Detail:   "Bad Request",
		Instance: "",
		Fields:   errors,
	}

}
