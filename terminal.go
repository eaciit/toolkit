package toolkit

import (
	"os/exec"
)

func RunCommand(cmd string, parm ...string) (string, error) {
	cmdOut, err := exec.Command(cmd, parm...).Output()
	if err == nil {
		return string(cmdOut), nil
	} else {
		return "", err
	}
}
