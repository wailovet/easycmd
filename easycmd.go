package easycmd

func EasyCmd(name string, arg ...string) *Pty {
	return NewPty(name, arg...)
}
