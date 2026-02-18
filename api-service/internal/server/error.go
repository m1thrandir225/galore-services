package httpserver

import "errors"

var (
	ErrInvalidGlobalConfig = errors.New("invalid global config")
	ErrInvalidStore        = errors.New("invalid store")
	ErrInvalidCache        = errors.New("invalid cache")
	ErrInvalidServerConfig = errors.New("invalid server config")
	ErrInvalidHandler      = errors.New("invalid handler")
)
