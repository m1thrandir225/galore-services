package users

import "mime/multipart"

type UserID struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateUserInformationRequest struct {
	Email     string                `form:"email" json:"email" binding:"required"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url" json:"avatar_url" binding:"omitempty"`
	Name      string                `form:"name" json:"name" binding:"required"`
	Birthday  string                `form:"birthday" json:"birthday" binding:"required"`
}
type UpdateUserEmailNotificationsRequest struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

type UpdateUserPushNotificationsRequest struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}
