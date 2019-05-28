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
	Buffer    *bytes.Buffer
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

func NewPty(name string, arg ...string) *Pty {
	cmd := exec.Command(name, arg...)
	ptytmp := Pty{
		cmd:    cmd,
		Buffer: bytes.NewBuffer(nil),
	}
	return &ptytmp
}
