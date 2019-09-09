package easycmd

import (
	"bytes"
	"os"
	"os/exec"
)

type Pty struct {
	cmd       *exec.Cmd
	ptyFile   *os.File
	ptyWriter func(p []byte) (int, error)
	ptyEnd    func()
	ptyError  func(err interface{})
	Buffer    *bytes.Buffer
	IsWinUtf8 bool
}

type ptyWriter struct {
	callback func(p []byte)
}

func (that *ptyWriter) Write(p []byte) (int, error) {
	that.callback(p)
	return len(p), nil
}

func (that *Pty) Write(p []byte) (int, error) {
	n, err := that.ptyWriter(p)
	return n, err
}

func (that *Pty) SetEventEnd(c func()) {
	that.ptyEnd = c
}

func (that *Pty) Close(c func()) {
	that.cmd.Process.Kill()
}

func (that *Pty) SetEventError(c func(err interface{})) {
	that.ptyError = c
}

func NewPty(name string, arg ...string) *Pty {
	cmd := exec.Command(name, arg...)
	ptytmp := Pty{
		cmd:    cmd,
		Buffer: bytes.NewBuffer(nil),
	}
	return &ptytmp
}
