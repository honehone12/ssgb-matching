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

func MakeWaitingMessage() StatusMessage {
	return StatusMessage{
		Status: StatusWaitng,
	}
}

func MakeMatchedMessage(gsip gsip.GSIP) StatusMessage {
	return StatusMessage{
		Status: StatusMatched,
		Gsip:   gsip,
	}
}
