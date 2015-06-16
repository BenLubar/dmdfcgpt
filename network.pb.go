// Code generated by protoc-gen-go.
// source: network.proto
// DO NOT EDIT!

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	network.proto

It has these top-level messages:
	Packet
*/
package main

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Packet struct {
	// streams started by the client are odd. streams started by the
	// server are even. responses use the same stream as the request.
	Stream           *uint64 `protobuf:"varint,1,req,name=stream" json:"stream,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}

func (m *Packet) GetStream() uint64 {
	if m != nil && m.Stream != nil {
		return *m.Stream
	}
	return 0
}

func init() {
}
