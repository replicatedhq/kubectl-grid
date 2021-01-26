package kubectl

import (
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
)

func run(cmd *exec.Cmd) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "failed to create stdout reader")
	}
	defer stdoutReader.Close()

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "failed to create stderr reader")
	}
	defer stderrReader.Close()

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start kubectl")
	}

	stdout, _ := ioutil.ReadAll(stdoutReader)
	stderr, _ := ioutil.ReadAll(stderrReader)

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// maybe log outputs instead of this "cleverness"
			if len(stderr) > 0 {
				return errors.New(string(stderr))
			} else if len(stdout) > 0 {
				return errors.New(string(stdout))
			}
			return errors.Wrap(err, "failed to run kubectl")
		}
	}

	return nil
}
