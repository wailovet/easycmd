package easycmd

import (
	"github.com/kr/pty"
	"io"
)

func (that *Pty) Start(callback func(data []byte)) error {
	var err error
	that.ptyFile, err = pty.Start(that.cmd)
	that.ptyWriter = func(b []byte) (int, error) {
		return that.ptyFile.Write(b)
	}

	if err != nil {
		return err
	}

	go func() {

		defer func() {
			errs := recover()
			if errs != nil {
				if that.ptyError != nil {
					that.ptyError(errs)
				}
			}
			if that.ptyEnd != nil {
				that.ptyEnd()
			}
		}()

		_, _ = io.Copy(&ptyWriter{
			callback: func(p []byte) {
				that.Buffer.Write(p)
				callback(p)
			},
		}, that.ptyFile)

		that.ptyFile.Close()
	}()

	return nil
}
