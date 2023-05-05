package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

func defaultTokenFunc() (string, error) {
	return cliConfig.GetToken()
}

func RunCommandOrDie(app string, args ...string) {
	cmd := exec.Command(app, args...)
	log.Debugf("running command: %v", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		cobra.CheckErr(fmt.Sprint(err) + ": " + string(output))
	}
	log.Infof(string(output))
}
