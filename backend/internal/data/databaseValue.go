package data

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"talenest/backend/internal/utils"
	"time"
)

type Value interface {
	GetValueString() string
	Value() (driver.Value, error)
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

func (value *IntValue) Value() (driver.Value, error) {
	return value.val, nil
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

func (value *FloatValue) Value() (driver.Value, error) {
	return value.val, nil
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

func (value *StringValue) Value() (driver.Value, error) {
	return value.val, nil
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

func (value *TokenValue) Value() (driver.Value, error) {
	return value.token, nil
}

func GetTokens(number int, symbol string) []Value {
	tokens := []Value{}
	for i := 0; i < number; i++ {
		tokens = append(tokens, NewTokenValue(symbol))
	}
	return tokens
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
	return fmt.Sprintf("'%s'", value.val.Format(utils.DATETIME_FORMAT))
}

func (value *TimeValue) String() string {
	return value.GetValueString()
}

func (value *TimeValue) Value() (driver.Value, error) {
	return value.val, nil
}
