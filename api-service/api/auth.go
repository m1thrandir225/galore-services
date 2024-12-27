package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1thrandir225/galore-services/mail"
	"github.com/m1thrandir225/galore-services/security"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassowrd(requestData.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	imageData, err := util.BytesFromFile(requestData.AvatarUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userId := uuid.New()

	avatarUrl, err := server.storage.UploadFile(imageData, userId.String(), requestData.AvatarUrl.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	otpSecret, err := security.GenerateOTPSecret()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		ID:         userId,
		Email:      requestData.Email,
		Birthday:   dbDate,
		Name:       requestData.Name,
		AvatarUrl:  avatarUrl,
		Password:   hashedPassword,
		HotpSecret: otpSecret,
	}

	newEntry, err := server.store.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.CreateHotpCounter(ctx, db.CreateHotpCounterParams{
		UserID:  newEntry.ID,
		Counter: 0,
	})
	if err != nil {
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
		"subject":       "Welcome to Galore",
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
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing required fields")))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, reqData.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting user")))
		return
	}
	currentUserCounter, err := server.store.GetCurrentCounter(ctx, user.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = server.store.CreateHotpCounter(ctx, db.CreateHotpCounterParams{
				UserID:  user.ID,
				Counter: 0,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error creating hotp counter")))
				return
			}
			currentUserCounter = 0
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting hotp counter")))
		return
	}

	otpCode, err := security.GenerateHOTP(user.HotpSecret, uint64(currentUserCounter))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error while generating otp code")))
		return
	}

	emailTemplate := mail.GeneratePasswordOTPMail(otpCode)
	server.scheduler.EnqueueJob("send_mail", map[string]interface{}{
		"email":         reqData.Email,
		"mail_template": emailTemplate,
		"subject":       "Galore - Password Reset Request",
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
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting user")))
		return
	}

	currentUserCounter, err := server.store.GetCurrentCounter(ctx, user.ID)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			//No counter is found initialize a new one
			err = server.store.CreateHotpCounter(ctx, db.CreateHotpCounterParams{
				UserID:  user.ID,
				Counter: 0,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to generate hotp counter")))
				return
			}
			currentUserCounter = 0
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting hotp counter")))
		return
	}

	isOtpCorrect := false

	for lookAhead := uint64(0); lookAhead < 5; lookAhead++ {
		tC := uint64(currentUserCounter) + lookAhead
		isValid, vaErr := security.ValidateHOTP(user.HotpSecret, reqData.OTP, tC)
		if vaErr != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("otp verification failed")))
			return
		}
		if isValid {
			_, iErr := server.store.IncreaseCounter(ctx, user.ID)
			if iErr != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to increase hotp counter")))
				return
			}
			isOtpCorrect = true
			break
		}
	}

	if !isOtpCorrect {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("otp verification failed")))
		return
	}

	requestValidUntil := time.Now().Add(time.Minute * 5)
	pgValidUntil := pgtype.Timestamptz{
		Time:             requestValidUntil,
		InfinityModifier: 0,
		Valid:            true,
	}

	arg := db.CreateResetPasswordRequestParams{
		ValidUntil: pgValidUntil,
		UserID:     user.ID,
	}
	passwordChangeRequest, err := server.store.CreateResetPasswordRequest(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error creating reset password request")))
		return
	}
	currentTime := time.Now()

	if currentTime.After(requestValidUntil.In(currentTime.Location())) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("reset password request has expired")))
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("getting reset password request")))
		return
	}

	user, err := server.store.GetUser(ctx, resetPasswordRequest.UserID)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("required user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting user")))
		return
	}

	currentTime := time.Now().UTC()

	validUntilTime := resetPasswordRequest.ValidUntil.Time.In(currentTime.Location())

	if currentTime.After(validUntilTime) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("reset password request has expired")))
		return
	}

	if reqData.NewPassword != reqData.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid password")))
		return
	}

	newPassword, err := util.HashPassowrd(reqData.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error while hashing password")))
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: newPassword,
	}
	err = server.store.UpdateUserPassword(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error updating password")))
		return
	}

	pgTypeTimeStamp := pgtype.Timestamptz{
		Time:             time.Now(),
		Valid:            true,
		InfinityModifier: 0,
	}

	updatePasswordRequestArg := db.UpdateResetPasswordRequestParams{
		ID:            resetPasswordRequest.ID,
		PasswordReset: true,
		ValidUntil:    pgTypeTimeStamp,
	}
	_, err = server.store.UpdateResetPasswordRequest(ctx, updatePasswordRequestArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error updating password reset request")))
		return
	}

	template := mail.GeneratePasswordResetSuccessfullyMail()

	sendMailArgs := map[string]interface{}{
		"email":         user.Email,
		"mail_template": template,
		"subject":       "Your Password Was Changed",
	}

	server.scheduler.EnqueueJob("send_mail", sendMailArgs)

	ctx.Status(http.StatusOK)
}
