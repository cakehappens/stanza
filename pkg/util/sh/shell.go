package sh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Waiter func(cmd *exec.Cmd) error

type WritePipeHandler func(writer io.WriteCloser) error

type ReadPipeHandler func(reader io.ReadCloser) error

// For use when y
type MultiReadPipeHandler func(readers ...io.ReadCloser) error

type RunOptions struct {
	Args   []string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
	Stdin io.Reader
	Waiter Waiter
}

type RunOption func(opts *RunOptions)

// Merges os.Environment variables to Env map passed to command
func OptionWithOsEnv(opts *RunOptions) error {
	for _, item := range os.Environ() {
		split := strings.SplitN(item, "=", 2)
		opts.Env[split[0]] = split[1]
	}
	return nil
}

// Run provides an easy to use API on top of RunLow
// By default, nothing is done to stdout, stderr, stdin
// so if you want to consume them, you must bind to them using Waiter func
// If Waiter remains unset, it defaults to cmd.Wait()
func Run(ctx context.Context, cmd string, options ...RunOption) error {
	runOpts := &RunOptions{}

	for _, optFn := range options {
		optFn(runOpts)
	}

	if runOpts.Waiter == nil {
		runOpts.Waiter = func(cmd *exec.Cmd) error {
			return cmd.Wait()
		}
	}

	if runOpts.Env == nil {
		runOpts.Env = make(map[string]string)
	}

	if runOpts.Args == nil {
		runOpts.Args = make([]string, 0)
	}

	cmdEnv := make([]string, len(runOpts.Env))

	for k, v := range runOpts.Env {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}

	code, err, ran := RunLow(
		ctx,
		cmdEnv,
		runOpts.Waiter,
		runOpts.Stdout,
		runOpts.Stderr,
		runOpts.Stdin,
		cmd,
		runOpts.Args...)

	if !ran {
		return fmt.Errorf("command failed to run: %w", err)
	}

	if err != nil {
		return fmt.Errorf("command ran and exited with code %d: %w", code, err)
	}

	return nil
}

// Run is the lowest level function to execute an external command
// since this function uses cmd.Start, it does not natively block
// you must provide that functionality
// No environment variables are explicitly passed to the command besides those provided
// Environment variables provided must be in the form of KEY=VALUE
func RunLow(
	ctx context.Context,
	env []string,
	waiter Waiter,
	stdout io.Writer,
	stderr io.Writer,
	stdin io.Reader,
	cmd string,
	args ...string,
	) (code int, err error, ran bool) {

	c := exec.CommandContext(ctx, cmd, args...)
	c.Env = env

	c.Stderr = stderr
	c.Stdout = stdout
	c.Stdin = stdin

	err = c.Start()
	if err != nil {
		return -1, err, false
	}

	err = waiter(c)

	if err == nil {
		return 0, nil, true
	}

	if err, ok := err.(*exec.ExitError); ok {
		return err.ExitCode(), err, true
	}

	return -1, err, false
}
