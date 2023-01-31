package enum

type Operation string

const (
	ADD     = Operation("ADD")
	USE     = Operation("USE")
	RECEIVE = Operation("RECEIVE")
	SEND    = Operation("SEND")
)
