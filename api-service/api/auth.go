package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1thrandir225/galore-services/mail"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
	otp "github.com/pquerna/otp/totp"
)

type registerUserRequest struct {
	Name      string                `form:"name" json:"name" binding:"required"`
	Email     string                `form:"email" json:"email" binding:"required,email"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url" json:"avatar_url" binding:"required"`
	Password  string                `form:"password" json:"password" binding:"required"`
	Birthday  string                `form:"birthday" json:"birthday" binding:"required"`
}

type registerUserResponse struct {
	User                  db.CreateUserRow `json:"user"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	AccessToken           string           `json:"access_token"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	SessionID             uuid.UUID        `json:"session_id"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	SessionID    string `json:"session_id" binding:"required,uuid"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type loginUserResponse struct {
	User                  db.CreateUserRow `json:"user"`
	RefreshToken          string           `json:"refresh_token"`
	AccessToken           string           `json:"access_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	SessionID             uuid.UUID        `json:"session_id"`
}

type logoutRequest struct {
	SessionID uuid.UUID `json:"session_id" binding:"required"`
}

type logoutResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ForgotPasswordResponse struct {
	OTP        string    `json:"otp"`
	ValidUntil time.Time `json:"valid_until"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	ResetPasswordRequest db.ResetPasswordRequest `json:"reset_password_request"`
	Email                string                  `json:"email"`
}

type ResetPasswordRequest struct {
	PasswordRequestId string `json:"password_request_id" binding:"required"`
	NewPassword       string `json:"new_password" binding:"required"`
	ConfirmPassword   string `json:"confirm_password" binding:"required"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var requestData registerUserRequest

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	dbDate, err := util.TimeToDbDate(requestData.Birthday)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassowrd(requestData.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	imageData, err := util.BytesFromFile(requestData.AvatarUrl)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userId := uuid.New()

	avatarUrl, err := server.storage.UploadFile(imageData, userId.String(), requestData.AvatarUrl.Filename)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	args := db.CreateUserParams{
		ID:        userId,
		Email:     requestData.Email,
		Birthday:  dbDate,
		Name:      requestData.Name,
		AvatarUrl: avatarUrl,
		Password:  hashedPassword,
	}

	newEntry, err := server.store.CreateUser(ctx, args)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(newEntry.ID, server.config.AccessTokenDuration)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(newEntry.ID, server.config.RefreshTokenDuration)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Email:        args.Email,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	b, err := json.Marshal(session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	server.cache.SaveItem(ctx, args.Email, string(b))

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	template := mail.GenerateWelcomeMail(newEntry.Email)
	sendMailArgs := map[string]interface{}{
		"email":         newEntry.Email,
		"mail_template": template,
	}

	server.scheduler.EnqueueJob("send_mail", sendMailArgs)

	if err != nil {
		log.Print(errorResponse(err))
	}

	response := registerUserResponse{
		User:                  newEntry,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var requestData loginUserRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ComparePassword(user.Password, requestData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Email:        requestData.Email,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		SessionID: session.ID,
		User: db.CreateUserRow{
			ID:                        user.ID,
			Name:                      user.Name,
			Email:                     user.Email,
			AvatarUrl:                 user.AvatarUrl,
			Birthday:                  user.Birthday,
			EnabledPushNotifications:  user.EnabledPushNotifications,
			EnabledEmailNotifications: user.EnabledPushNotifications,
			CreatedAt:                 user.CreatedAt,
		},
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) logout(ctx *gin.Context) {
	// This will be an auth only route so you would need to send a acces token, and as a specific meter need to send the session id
	// so when you logout the session can be deleted or maybe set it to blocked, let it be for the time being just be deleted
	// TODO: set the session to blocked instead of deleting it so the user can see his own previous sessions.

	var requestData logoutRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	_, err := server.store.InvalidateSession(ctx, requestData.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// No need to return anything
	ctx.JSON(http.StatusOK, logoutResponse{
		Message: "You have been logged out",
	})
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var requestData refreshTokenRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionId, err := uuid.Parse(requestData.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(requestData.RefreshToken)

	if err != nil {
		_, _ = server.store.InvalidateSession(ctx, sessionId)
		//No need to handle the error as if there is an error it is still unauthorized.

		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	newToken, _, err := server.tokenMaker.CreateToken(payload.UserId, server.config.AccessTokenDuration)

	ctx.JSON(http.StatusOK, refreshTokenResponse{
		AccessToken: newToken,
	})
}

func (server *Server) forgotPassword(ctx *gin.Context) {
	var reqData ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUserByEmail(ctx, reqData.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	codeValidUntil := time.Now().Add(time.Minute * 5)
	otpCode, err := otp.GenerateCode(server.config.TOTPSecret, codeValidUntil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	emailTemplate := mail.GeneratePasswordOTPMail(otpCode)
	server.scheduler.EnqueueJob("send_mail", map[string]interface{}{
		"email":         reqData.Email,
		"mail_template": emailTemplate,
	})

	ctx.Status(http.StatusOK)
}

func (server *Server) verifyOTP(ctx *gin.Context) {
	var reqData VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, reqData.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	isValid := otp.Validate(reqData.OTP, server.config.TOTPSecret)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	requestValidUntil := time.Now().Add(time.Minute * 5)

	var pgRequestValidUntil pgtype.Timestamp
	err = pgRequestValidUntil.Scan(requestValidUntil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateResetPasswordRequestParams{
		ValidUntil: pgRequestValidUntil,
		UserID:     user.ID,
	}
	passwordChangeRequest, err := server.store.CreateResetPasswordRequest(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, VerifyOTPResponse{
		Email:                user.Email,
		ResetPasswordRequest: passwordChangeRequest,
	})
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var reqData ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	resetPasswordRequestId, err := uuid.Parse(reqData.PasswordRequestId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resetPasswordRequest, err := server.store.GetResetPasswordRequest(ctx, resetPasswordRequestId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reset password request not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if resetPasswordRequest.ValidUntil.Time.After(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if reqData.NewPassword != reqData.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newPassword, err := util.HashPassowrd(reqData.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.UpdateUserPasswordParams{
		ID:       resetPasswordRequest.UserID,
		Password: newPassword,
	}
	err = server.store.UpdateUserPassword(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
