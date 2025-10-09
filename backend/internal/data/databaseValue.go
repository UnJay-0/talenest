package data

import (
	"fmt"
	"strconv"
	"time"
)

type Value interface {
	GetValueString() string
}

type IntValue struct {
	val int
}

func NewIntValue(val int) *IntValue {
	return &IntValue{
		val: val,
	}
}

func (value *IntValue) GetValueString() string {
	return strconv.Itoa(value.val)
}

func (value *IntValue) String() string {
	return value.GetValueString()
}

type FloatValue struct {
	val float32
}

func NewFloatValue(val float32) *FloatValue {
	return &FloatValue{
		val: val,
	}
}

func (value *FloatValue) GetValueString() string {
	return fmt.Sprintf("%f", value.val)
}

func (value *FloatValue) String() string {
	return value.GetValueString()
}

type StringValue struct {
	val string
}

func NewStringValue(val string) *StringValue {
	return &StringValue{
		val: val,
	}
}

func (value *StringValue) GetValueString() string {
	return fmt.Sprintf("'%s'", value.val)
}

func (value *StringValue) String() string {
	return value.GetValueString()
}

type TokenValue struct {
	token string
}

func NewTokenValue(token string) *TokenValue {
	return &TokenValue{
		token: token,
	}
}

func (value *TokenValue) GetValueString() string {
	return value.token
}

func (value *TokenValue) String() string {
	return value.GetValueString()
}

type TimeValue struct {
	val time.Time
}

func NewTimeValue(val time.Time) *TimeValue {
	return &TimeValue{
		val: val,
	}
}

func (value *TimeValue) GetValueString() string {
	return fmt.Sprintf("'%s'", value.val.Format("2006-01-02 15:04:05"))
}

func (value *TimeValue) String() string {
	return value.GetValueString()
}
