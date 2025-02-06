package model

import (
	"encoding/json"
	"time"
)

type Check struct {
	Schedule string  `json:"schedule"`
	Command  string  `json:"command"`
	Recover  *string `json:"recover"`
}

func (c *Check) UnmarshalJSON(data []byte) error {
	type Alias Check
	if err := json.Unmarshal(data, &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}); err != nil {
		return err
	}

	if c.Schedule == "" {
		c.Schedule = "* * * * *"
	}

	return nil
}

type CheckResult struct {
	Started   time.Time   `yaml:"started"`
	Completed time.Time   `yaml:"completed"`
	Status    CheckStatus `yaml:"status"`

	CheckOutput string  `yaml:"check_output"`
	CheckError  *string `yaml:"check_error,omitempty"`

	RecoverOutput *string `yaml:"recover_output,omitempty"`
	RecoverError  *string `yaml:"recover_error,omitempty"`

	RecheckOutput *string `yaml:"recheck_output,omitempty"`
	RecheckError  *string `yaml:"recheck_error,omitempty"`
}

type CheckStatus int

const (
	CheckSuccess CheckStatus = iota
	CheckRecovered
	CheckFailed
)

func (s CheckStatus) String() string {
	return [...]string{"CheckSuccess", "CheckRecovered", "CheckFailed"}[s]
}
