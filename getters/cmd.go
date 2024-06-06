package getters;

import (
	"os/exec"
	"errors"
)

func Cmd(f *exec.Cmd) error {
	if err := f.Run(); err != nil {
		return errors.New("incorrect command");
	} else {
		return nil;
	}
}