package main

import (
	"fmt"
	"os"
	"status/internal/checker"
	"status/internal/model"
	"status/internal/util"
)

func main() {
	var checkPath string
	if len(os.Args) > 1 {
		checkPath = os.Args[1]
	} else {
		checkPath = "checks.json"
	}

	if checker, err := checker.New(checkPath, onChecked); err != nil {
		panic(err)
	} else {
		checker.Start()
		util.WaitUserTermination()
		checker.Stop()
	}
}

func onChecked(name string, check model.Check, result model.CheckResult) {
	if result.CheckError != nil {
		fmt.Printf("check '%s' failed: %s", name, result.CheckOutput)
	}

	if result.RecoverError != nil {
		fmt.Printf("recover '%s' failed: %s", name, result.CheckOutput)
	}

	if result.RecheckError != nil {
		fmt.Printf("re-check '%s' failed: %s", name, result.CheckOutput)
	}

	// notify Slack
	// update db via API
}
