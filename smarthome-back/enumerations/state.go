package enumerations

type State int

const (
	PENDING State = iota
	ACCEPTED
	DECLINED
)
