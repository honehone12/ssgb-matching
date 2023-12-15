package messages

import "ssgb-matching/gsip"

const (
	StatusError byte = iota
	StatusWaitng
	StatusMatched
)

type StatusMessage struct {
	Status byte
	Gsip   gsip.GSIP
}
