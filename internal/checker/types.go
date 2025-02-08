package checker

import (
	"time"
)

type Check struct {
	Schedule string  `yaml:"schedule"`
	Command  string  `yaml:"command"`
	Recover  *string `yaml:"recover"`
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

type Result struct {
	Command string  `yaml:"command"`
	Recover *string `yaml:"recover"`

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
