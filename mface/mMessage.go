package mface

type MMessage interface {
	MsgID() string
	RouteID() string
	ConnID() string

	Marshal() []byte
}
