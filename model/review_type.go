package model

type ReviewType string

const (
	SystematicReview ReviewType = "SystematicReview"
	ScopingReview               = "ScopingReview"
	RapidReview                 = "RapidReview"
)
