package api

type UriId struct {
	ID string `uri:"id" binding:"required,uuid"`
}
