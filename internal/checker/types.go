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
	Status Status `json:"status"`

	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed"`

	Check   *CmdResult `json:"check"`
	Recover *CmdResult `json:"recover"`
	ReCheck *CmdResult `json:"recheck"`
}

type CmdResult struct {
	Command string `json:"command"`
	Status  Status `json:"status"`

	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed"`

	Output string  `json:"output"`
	Error  *string `json:"error"`
}

type Status string

const (
	StatusRunning   Status = "Running"
	StatusSuccess   Status = "Success"
	StatusRecovered Status = "Recovered"
	StatusFailed    Status = "Failed"
)
