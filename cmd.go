package easycmd

import (
	"bytes"
	"os"
	"os/exec"
)

type Cmd struct {
	cmd       *exec.Cmd
	ptyFile   *os.File
	cmdWriter func(p []byte) (int, error)
	cmdEnd    func()
	cmdError  func(err interface{})
	Buffer    *bytes.Buffer
	IsWinUtf8 bool
}

type cmdWriter struct {
	callback func(p []byte)
}

func (that *cmdWriter) Write(p []byte) (int, error) {
	that.callback(p)
	return len(p), nil
}

func (that *Cmd) Write(p []byte) (int, error) {
	n, err := that.cmdWriter(p)
	return n, err
}

func (that *Cmd) SetEventEnd(c func()) {
	that.cmdEnd = c
}

func (that *Cmd) Close() {
	that.cmd.Process.Kill()
}

func (that *Cmd) SetEventError(c func(err interface{})) {
	that.cmdError = c
}

func (that *Cmd) Start(callback func(data []byte)) error {
	var err error

	stdInIn, _ := that.cmd.StdinPipe()

	var out = &WinWriter{
		callback:  callback,
		IsWinUtf8: that.IsWinUtf8,
	}
	that.cmd.Stderr = out
	that.cmd.Stdout = out
	that.cmdWriter = func(p []byte) (i int, e error) {
		return stdInIn.Write(p)
	}
	err = that.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			errs := recover()
			if errs != nil {
				if that.cmdError != nil {
					that.cmdError(errs)
				}
			}
			if that.cmdEnd != nil {
				that.cmdEnd()
			}
		}()

		that.cmd.Wait()
	}()

	return nil
}

func NewCmd(name string, arg ...string) *Cmd {
	cmd := exec.Command(name, arg...)
	cmdtmp := Cmd{
		cmd:    cmd,
		Buffer: bytes.NewBuffer(nil),
	}
	return &cmdtmp
}
