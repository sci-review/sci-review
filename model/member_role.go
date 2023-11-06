package model

type MemberRole string

const (
	MemberOwner    MemberRole = "MemberOwner"
	MemberAdmin               = "MemberAdmin"
	MemberReviewer            = "MemberReviewer"
)
