package log

// RollType file rolling type
type RollType int

const (
	RollNONE = iota
	RollByFileSize
	RollByDay
)
