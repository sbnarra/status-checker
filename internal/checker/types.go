package checker

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
	type Defaults Check
	var defaults = Defaults{
		Schedule: "* * * * *",
	}
	if err := json.Unmarshal(data, &defaults); err != nil {
		return err
	}
	*c = Check(defaults)
	return nil
}

func (c *Check) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Defaults Check
	var defaults = Defaults{
		Schedule: "* * * * *",
	}
	if err := unmarshal(&defaults); err != nil {
		return err
	}
	*c = Check(defaults)
	return nil
}

type CheckResult struct {
	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed"`
	Status    string    `json:"status"`

	CheckOutput string  `json:"check_output"`
	CheckError  *string `json:"check_error,omitempty"`

	RecoverOutput *string `json:"recover_output,omitempty"`
	RecoverError  *string `json:"recover_error,omitempty"`

	RecheckOutput *string `json:"recheck_output,omitempty"`
	RecheckError  *string `json:"recheck_error,omitempty"`
}
