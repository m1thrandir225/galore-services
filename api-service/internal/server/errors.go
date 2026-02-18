package server

import "errors"

var (
	ErrInvalidGlobalConfig = errors.New("invalid global config")
	ErrInvalidStore        = errors.New("invalid store")
	ErrInvalidCache        = errors.New("invalid cache")
	ErrInvalidServerConfig = errors.New("invalid server config")
	ErrInvalidHandler      = errors.New("invalid handler")
	ErrInvalidTokenMaker   = errors.New("invalid token maker")
)
