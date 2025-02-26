package app_errors

import "errors"

var (
	ServerDoesNotExist   = errors.New("server does not exist")
	ReceiverDoesNotExist = errors.New("receiver does not exist")
)
