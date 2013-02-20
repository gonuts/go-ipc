package ipc

import (
	"encoding/gob"
)

type GobChannel struct {
	Channel
	gob.Encoder
	gob.Decoder
}

func NewGobChannel(name string) *GobChannel {
	return new(GobChannel)
}
