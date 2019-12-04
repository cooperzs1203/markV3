package mface

type MMessage interface {
	MsgID() string
	ConnID() string

	Marshal() []byte
}
