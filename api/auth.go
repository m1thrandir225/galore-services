package api

func (server *Server) verifyAuthCookie(cookie string) bool {
	payload, err := server.tokenMaker.VerifyToken(cookie)
	if err != nil {
		return false
	}
	err = payload.Valid()

	if err != nil {
		return false
	}

	return true
}
