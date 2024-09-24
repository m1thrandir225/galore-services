package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
)

/*
*
This is a universal function that should be called when we want to upload a file
*/
func (server *Server) uploadFile(ctx *gin.Context) (bool, error) {
	file, err := ctx.FormFile("file")

	if err != nil {
		return false, err
	}
	err = ctx.SaveUploadedFile(file, "/public")

	if err != nil {
		return false, err
	}

	return true, nil
}

func (server *Server) deleteFile(ctx *gin.Context) (bool, error) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		return false, err
	}

	_, err := os.Open(uriData.ID)

	if errors.Is(err, os.ErrNotExist) {
		return false, os.ErrNotExist
	}
	err = os.Remove("/public/" + uriData.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

/*
* I don't think the following two functions are needed since the files are gonna be placed inside
* the public folder and the URL can be hand crafted to get the file on the client side.
 */
func (server *Server) getUserFiles(ctx *gin.Context) {}

func (server *Server) getUserFile(ctx *gin.Context) {
	//var uriData userIdFileIdUri
	//
	//if err := ctx.ShouldBindUri(&uriData); err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//	return
	//}
	//
}
