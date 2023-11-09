package model

type InvestigationStatus string

const (
	PiStatusInProgress   InvestigationStatus = "InProgress"
	PiStatusProceed                          = "Proceed"
	PiStatusDoNotProceed                     = "DoNotProceed"
	PiStatusCancelled                        = "Cancelled"
)
