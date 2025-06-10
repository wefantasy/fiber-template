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

func IsSuccess(err error) bool {
	if err == nil {
		return true
	}
	e := Error(err.Error())
	return errors.Is(e, Nil)
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
	AuthFailed               Error = "AuthFailed"
	UsernameOrPasswordFailed Error = "UsernameOrPasswordFailed"
	TokenGenerateFailed      Error = "TokenGenerateFailed"

	// 用户侧错误
	ParamError Error = "ParamError"

	// 三方问题
	ExternalError Error = "ExternalError"
)
