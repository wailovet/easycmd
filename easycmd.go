package easycmd

func EasyCmd(name string, arg ...string) *Pty {
	return NewPty(name, arg...)
}

func EasyCmdNotPty(name string, arg ...string) *Cmd {
	return NewCmd(name, arg...)
}
