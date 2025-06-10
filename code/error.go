package code

import (
	"errors"
)

type Error string

func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	return string(e)
}

func (e Error) ToError() error {
	return errors.New(e.String())
}

func (e Error) IsNil() bool {
	return errors.Is(e, Nil)
}

func (e Error) IsNotNil() bool {
	return !e.IsNil()
}

func ParseError(err string) Error {
	return Error(err)
}

const (
	// 无错误
	Nil Error = ""

	// 服务侧错误
	ServerError         Error = "ServerError"
	PasswordCryptFailed Error = "PasswordCryptFailed"
	JsonMarshalFailed   Error = "JsonMarshalFailed"
	JsonUnmarshalFailed Error = "JsonUnmarshalFailed"

	// 数据库错误
	DatabaseError         Error = "DatabaseError"
	RedisGetDataFailed    Error = "RedisGetDataFailed"
	RedisSetDataFailed    Error = "RedisSetDataFailed"
	RedisDeleteDataFailed Error = "RedisDeleteDataFailed"
	RedisKeyNotExist      Error = "RedisKeyNotExist"

	// 认证模块错误
	UsernameOrPasswordFailed Error = "UsernameOrPasswordFailed"
	TokenGenerateFailed      Error = "TokenGenerateFailed"

	// 用户侧错误
	ParamError Error = "ParamError"

	// 三方问题
	ExternalError Error = "ExternalError"
)
