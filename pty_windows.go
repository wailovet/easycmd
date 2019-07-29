package easycmd

import "github.com/axgle/mahonia"

type WinWriter struct {
	callback  func(data []byte)
	IsWinUtf8 bool
}

func (that *WinWriter) Write(p []byte) (n int, err error) {
	if !that.IsWinUtf8 {
		that.callback([]byte(that.ConvertToString(string(p), "gbk", "utf-8")))
	} else {
		that.callback(p)
	}
	return len(p), nil
}

func (that *WinWriter) ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func (that *Pty) Start(callback func(data []byte)) error {
	var err error

	stdInIn, _ := that.cmd.StdinPipe()

	var out = &WinWriter{
		callback:  callback,
		IsWinUtf8: that.IsWinUtf8,
	}
	that.cmd.Stderr = out
	that.cmd.Stdout = out
	that.ptyWriter = func(p []byte) (i int, e error) {
		return stdInIn.Write(p)
	}
	err = that.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		that.cmd.Wait()
		if that.ptyEnd != nil {
			that.ptyEnd()
		}
	}()

	return nil
}
