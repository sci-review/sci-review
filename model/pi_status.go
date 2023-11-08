package model

type PiStatus string

const (
	PiStatusInProgress   PiStatus = "InProgress"
	PiStatusProceed               = "Proceed"
	PiStatusDoNotProceed          = "DoNotProceed"
	PiStatusCancelled             = "Cancelled"
)
