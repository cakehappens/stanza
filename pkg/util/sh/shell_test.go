package sh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

// https://npf.io/2015/06/testing-exec-command/

func TestHelperProcess(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, dockerRunResult)
	os.Exit(0)
}

func TestOptionWithOsEnv(t *testing.T) {
	type args struct {
		opts *RunOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OptionWithOsEnv(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("OptionWithOsEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		ctx     context.Context
		cmd     string
		options []RunOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.ctx, tt.args.cmd, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunLow(t *testing.T) {
	type args struct {
		ctx    context.Context
		env    []string
		waiter Waiter
		stdin  io.Reader
		cmd    string
		args   []string
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
		wantCode   int
		wantErr    error
		wantRan    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			gotCode, gotErr, gotRan := RunLow(tt.args.ctx, tt.args.env, tt.args.waiter, stdout, stderr, tt.args.stdin, tt.args.cmd, tt.args.args...)
			if gotStdout := stdout.String(); gotStdout != tt.wantStdout {
				t.Errorf("RunLow() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
			}
			if gotStderr := stderr.String(); gotStderr != tt.wantStderr {
				t.Errorf("RunLow() gotStderr = %v, want %v", gotStderr, tt.wantStderr)
			}
			if gotCode != tt.wantCode {
				t.Errorf("RunLow() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("RunLow() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if gotRan != tt.wantRan {
				t.Errorf("RunLow() gotRan = %v, want %v", gotRan, tt.wantRan)
			}
		})
	}
}