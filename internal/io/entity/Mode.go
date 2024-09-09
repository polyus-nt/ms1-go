package entity

type Mode int64

const (
	ModeBlank Mode = iota
	ModeRun
	ModeProg
	ModeConf
	ModeErr
)