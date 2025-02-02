package model

type Check struct {
	Schedule *string
	Command  string
	Recover  *string
}

type CheckResult struct {
	CheckOutput   string
	CheckError    *error
	RecoverOutput *string
	RecoverError  *error
	RecheckOutput *string
	RecheckError  *error
}
