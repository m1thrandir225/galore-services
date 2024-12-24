package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/pgvector/pgvector-go"
)

type CreateOrUpdateCocktailRequest struct {
	Image        *multipart.FileHeader `form:"file" binding:"required"`
	Name         string                `form:"name" binding:"required"`
	Ingredients  string                `form:"ingredients" binding:"required"`
	Glass        string                `form:"glass" binding:"required"`
	Instructions string                `form:"instructions" binding:"required"`
	IsAlcoholic  bool                  `form:"isAlcoholic" binding:"required"`
}

func (server *Server) createCocktail(ctx *gin.Context) {
	var requestData CreateOrUpdateCocktailRequest
	var isAlcoholic pgtype.Bool
	var ingredientDto dto.IngredientDto

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := isAlcoholic.Scan(requestData.IsAlcoholic)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Unmarshal ingredients and instructions to dto objects
	if err = json.Unmarshal([]byte(requestData.Ingredients), &ingredientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
		return
	}

	// Unload the bytes from the uploaded file
	fileData, err := util.BytesFromFile(requestData.Image)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	/*
	* Get the id of the currently logged-in user to use as a name for the folder that the uploaded file will be placed in.
	 */
	data, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	payload := data.(*token.Payload)

	// Upload the file to the public/user_id/file
	fileName, err := server.storage.UploadFile(fileData, payload.ID.String(), requestData.Image.Filename)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	nameEmbedding, err := server.embedding.GenerateEmbedding(requestData.Name)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateCocktailParams{
		Name:         requestData.Name,
		Image:        fileName,
		Glass:        requestData.Glass,
		Ingredients:  ingredientDto,
		Instructions: requestData.Instructions,
		IsAlcoholic:  isAlcoholic,
		Embedding:    pgvector.NewVector(nameEmbedding),
	}

	cocktail, err := server.store.CreateCocktail(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Categorize the current cocktail
	err = server.categorizer.CategorizeCocktail(cocktail)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktail)
}

func (server *Server) createCocktailMigrate(ctx *gin.Context) {
	var requestData CreateOrUpdateCocktailRequest
	var isAlcoholic pgtype.Bool
	var ingredientDto dto.IngredientDto

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := isAlcoholic.Scan(requestData.IsAlcoholic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Unmarshal ingredients and instructions to dto objects
	if err = json.Unmarshal([]byte(requestData.Ingredients), &ingredientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
		return
	}

	// Unload the bytes from the uploaded file
	fileData, err := util.BytesFromFile(requestData.Image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Upload the file to the public/cocktails/file
	fileName, err := server.storage.UploadFile(fileData, "cocktails", requestData.Image.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	nameEmbedding, err := server.embedding.GenerateEmbedding(requestData.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateCocktailParams{
		Name:         requestData.Name,
		Image:        fileName,
		Glass:        requestData.Glass,
		Ingredients:  ingredientDto,
		Instructions: requestData.Instructions,
		IsAlcoholic:  isAlcoholic,
		Embedding:    pgvector.NewVector(nameEmbedding),
	}

	cocktail, err := server.store.CreateCocktail(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Categorize the current cocktail
	err = server.categorizer.CategorizeCocktail(cocktail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktail)
}

/*
* To delete a cocktail we first must:
* 1.  Get the id
* 2. Get the cocktail
* 3. Delete the image from any storage
* 4. Delete the cocktail itself
 */
func (server *Server) deleteCocktail(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// 1. Get cocktail id
	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// 2. Get the cocktail data
	cocktail, err := server.store.GetCocktail(ctx, cocktailId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// 3. Delete the associated image
	err = server.storage.DeleteFile(cocktail.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// 4. Delete the cocktail itself
	err = server.store.DeleteCocktail(ctx, cocktailId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

func (server *Server) getCocktails(ctx *gin.Context) {
	searchQuery := ctx.DefaultQuery("search", "")

	cocktails, err := server.store.SearchCocktails(ctx, searchQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktails)
}

func (server *Server) getCocktail(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cachedCocktailJson, cacheErr := server.cache.GetItem(ctx, cocktailId.String())
	if cacheErr == nil {
		// Cache Hit, this was recently accessed
		var cachedCocktailData db.Cocktail
		err = json.Unmarshal([]byte(cachedCocktailJson), &cachedCocktailData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		cachedCocktailData.Image = util.UrlFromFilePath(server.config.HTTPServerAddress, cachedCocktailData.Image)
		ctx.JSON(http.StatusOK, cachedCocktailData)
	}

	cocktail, err := server.store.GetCocktail(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Cache Missed we save the cocktail to the cache
	if cacheErr != nil {
		cocktailJson, jsonMarshalError := json.Marshal(cocktail)
		if jsonMarshalError != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(jsonMarshalError))
			return
		}

		saveCacheError := server.cache.SaveItem(ctx, cocktail.ID.String(), string(cocktailJson))
		if saveCacheError != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(saveCacheError))
			return
		}
	}

	cocktail.Image = util.UrlFromFilePath(server.config.HTTPServerAddress, cocktail.Image)

	ctx.JSON(http.StatusOK, cocktail)
}

func (server *Server) updateCocktail(ctx *gin.Context) {
	var uriData UriId
	var requestData CreateOrUpdateCocktailRequest
	var isAlcoholic pgtype.Bool
	var ingredientDto dto.IngredientDto

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Unmarshal ingredients and instructions to dto objects
	if err := json.Unmarshal([]byte(requestData.Ingredients), &ingredientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
		return
	}
	err := isAlcoholic.Scan(requestData.IsAlcoholic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cocktail, err := server.store.GetCocktail(ctx, cocktailId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newImageData, err := util.BytesFromFile(requestData.Image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newFilePath, err := server.storage.ReplaceFile(cocktail.Image, newImageData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.UpdateCocktailParams{
		ID:           cocktailId,
		Name:         requestData.Name,
		Instructions: requestData.Instructions,
		Ingredients:  ingredientDto,
		Image:        newFilePath,
		Glass:        requestData.Glass,
		IsAlcoholic:  isAlcoholic,
	}

	updated, err := server.store.UpdateCocktail(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updated.Image = util.UrlFromFilePath(server.config.HTTPServerAddress, updated.Image)

	ctx.JSON(http.StatusOK, updated)
}

func (server *Server) getCategoriesForCocktail(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	categories, err := server.store.GetCategoriesForCocktail(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) getSimilarCocktails(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cocktails, err := server.store.GetCocktailAndSimilar(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, cocktails)
}
