package iperf

import "github.com/google/uuid"

func NewClient() *Client {
	json := true
	proto := Protocol(PROTO_TCP)
	time := 10
	length := "128KB"
	streams := 1
	c := &Client{
		JSON: &json,
		Proto: &proto,
		TimeSec: &time,
		Length: &length,
		Streams: &streams,
	}
	c.Id = uuid.New().String()
	return c
}

type Client struct {
	SharedOptions
	Host string
	Proto *Protocol
	Bandwidth *string
	TimeSec *int
	Bytes *string
	BlockCount *string
	Length *string
	Streams *int
	Reverse *bool
	Window *string
	MSS *int
	NoDelay *bool
	Version4 *bool
	Version6 *bool
	TOS *int
	ZeroCopy *bool
	OmitSec *int
	Prefix *string
	JSON *bool
}
