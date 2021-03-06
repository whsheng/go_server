// Code generated by protoc-gen-go.
// source: test_message.proto
// DO NOT EDIT!

/*
Package Cmd_Test is a generated protocol buffer package.

It is generated from these files:
	test_message.proto

It has these top-level messages:
	TestMessage
	RequestLogin
*/
package Cmd_Test

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type MSG_INDEX int32

const (
	MSG_INDEX_THE_FIRST MSG_INDEX = 1
)

var MSG_INDEX_name = map[int32]string{
	1: "THE_FIRST",
}
var MSG_INDEX_value = map[string]int32{
	"THE_FIRST": 1,
}

func (x MSG_INDEX) Enum() *MSG_INDEX {
	p := new(MSG_INDEX)
	*p = x
	return p
}
func (x MSG_INDEX) String() string {
	return proto.EnumName(MSG_INDEX_name, int32(x))
}
func (x *MSG_INDEX) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MSG_INDEX_value, data, "MSG_INDEX")
	if err != nil {
		return err
	}
	*x = MSG_INDEX(value)
	return nil
}

type TestMessage_MSG_SECOND int32

const (
	TestMessage_THE_SECOND TestMessage_MSG_SECOND = 1
)

var TestMessage_MSG_SECOND_name = map[int32]string{
	1: "THE_SECOND",
}
var TestMessage_MSG_SECOND_value = map[string]int32{
	"THE_SECOND": 1,
}

func (x TestMessage_MSG_SECOND) Enum() *TestMessage_MSG_SECOND {
	p := new(TestMessage_MSG_SECOND)
	*p = x
	return p
}
func (x TestMessage_MSG_SECOND) String() string {
	return proto.EnumName(TestMessage_MSG_SECOND_name, int32(x))
}
func (x *TestMessage_MSG_SECOND) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(TestMessage_MSG_SECOND_value, data, "TestMessage_MSG_SECOND")
	if err != nil {
		return err
	}
	*x = TestMessage_MSG_SECOND(value)
	return nil
}

type RequestLogin_MSG_SECOND int32

const (
	RequestLogin_THE_SECOND RequestLogin_MSG_SECOND = 2
)

var RequestLogin_MSG_SECOND_name = map[int32]string{
	2: "THE_SECOND",
}
var RequestLogin_MSG_SECOND_value = map[string]int32{
	"THE_SECOND": 2,
}

func (x RequestLogin_MSG_SECOND) Enum() *RequestLogin_MSG_SECOND {
	p := new(RequestLogin_MSG_SECOND)
	*p = x
	return p
}
func (x RequestLogin_MSG_SECOND) String() string {
	return proto.EnumName(RequestLogin_MSG_SECOND_name, int32(x))
}
func (x *RequestLogin_MSG_SECOND) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(RequestLogin_MSG_SECOND_value, data, "RequestLogin_MSG_SECOND")
	if err != nil {
		return err
	}
	*x = RequestLogin_MSG_SECOND(value)
	return nil
}

type TestMessage struct {
	FIRST            *MSG_INDEX              `protobuf:"varint,1,opt,enum=Cmd.Test.MSG_INDEX,def=1" json:"FIRST,omitempty"`
	SECOND           *TestMessage_MSG_SECOND `protobuf:"varint,2,opt,enum=Cmd.Test.TestMessage_MSG_SECOND,def=1" json:"SECOND,omitempty"`
	Name             []byte                  `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Age              *int32                  `protobuf:"varint,4,opt,name=age" json:"age,omitempty"`
	Desc             []byte                  `protobuf:"bytes,5,opt,name=desc" json:"desc,omitempty"`
	Count            *int32                  `protobuf:"varint,6,opt,name=count" json:"count,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *TestMessage) Reset()         { *m = TestMessage{} }
func (m *TestMessage) String() string { return proto.CompactTextString(m) }
func (*TestMessage) ProtoMessage()    {}

const Default_TestMessage_FIRST MSG_INDEX = MSG_INDEX_THE_FIRST
const Default_TestMessage_SECOND TestMessage_MSG_SECOND = TestMessage_THE_SECOND

func (m *TestMessage) GetFIRST() MSG_INDEX {
	if m != nil && m.FIRST != nil {
		return *m.FIRST
	}
	return Default_TestMessage_FIRST
}

func (m *TestMessage) GetSECOND() TestMessage_MSG_SECOND {
	if m != nil && m.SECOND != nil {
		return *m.SECOND
	}
	return Default_TestMessage_SECOND
}

func (m *TestMessage) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *TestMessage) GetAge() int32 {
	if m != nil && m.Age != nil {
		return *m.Age
	}
	return 0
}

func (m *TestMessage) GetDesc() []byte {
	if m != nil {
		return m.Desc
	}
	return nil
}

func (m *TestMessage) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

type RequestLogin struct {
	FIRST            *MSG_INDEX               `protobuf:"varint,1,opt,enum=Cmd.Test.MSG_INDEX,def=1" json:"FIRST,omitempty"`
	SECOND           *RequestLogin_MSG_SECOND `protobuf:"varint,2,opt,enum=Cmd.Test.RequestLogin_MSG_SECOND,def=2" json:"SECOND,omitempty"`
	Username         []byte                   `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Password         []byte                   `protobuf:"bytes,4,opt,name=password" json:"password,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *RequestLogin) Reset()         { *m = RequestLogin{} }
func (m *RequestLogin) String() string { return proto.CompactTextString(m) }
func (*RequestLogin) ProtoMessage()    {}

const Default_RequestLogin_FIRST MSG_INDEX = MSG_INDEX_THE_FIRST
const Default_RequestLogin_SECOND RequestLogin_MSG_SECOND = RequestLogin_THE_SECOND

func (m *RequestLogin) GetFIRST() MSG_INDEX {
	if m != nil && m.FIRST != nil {
		return *m.FIRST
	}
	return Default_RequestLogin_FIRST
}

func (m *RequestLogin) GetSECOND() RequestLogin_MSG_SECOND {
	if m != nil && m.SECOND != nil {
		return *m.SECOND
	}
	return Default_RequestLogin_SECOND
}

func (m *RequestLogin) GetUsername() []byte {
	if m != nil {
		return m.Username
	}
	return nil
}

func (m *RequestLogin) GetPassword() []byte {
	if m != nil {
		return m.Password
	}
	return nil
}

func init() {
	proto.RegisterEnum("Cmd.Test.MSG_INDEX", MSG_INDEX_name, MSG_INDEX_value)
	proto.RegisterEnum("Cmd.Test.TestMessage_MSG_SECOND", TestMessage_MSG_SECOND_name, TestMessage_MSG_SECOND_value)
	proto.RegisterEnum("Cmd.Test.RequestLogin_MSG_SECOND", RequestLogin_MSG_SECOND_name, RequestLogin_MSG_SECOND_value)
}
