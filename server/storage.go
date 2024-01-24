package server

import (
	example "expamle"
	"github.com/gin-gonic/gin"
)

type storageHandler struct {
	service example.StorageService
	config  *example.Configuration
}

func (s *Server) InitStorageHandlers(svc example.StorageService, config *example.Configuration) {
	s.storage = svc
	handler := storageHandler{
		service: svc,
		config:  config,
	}

	group := s.Router.Group("api/v1/storage")
	{
		group.Use(s.Authenticate(example.AllRoles))
		{
			group.POST("/upload", handler.Upload)
		}
	}
}

// Upload handles uploading files request.
//
// @BasePath			/api/v1/storage
// @Summary				uploading file to storage service
// @Description			upload specified file to storage service with specified properties
// @Tags				Storage
// @Accept				mpfd
// @Produce				json
//
// @Param				file 		formData  	file 	true	"uploading file"
// @Param				tag			path		string	true	     "string enums"   Enums(public)
//
// @Param				Authenticate header string false "authenticate paseto token [Required If AuthMode: paseto]"
//
// @Success				200			{object}	example.BaseResult	 	"always returns status 200 but body contains errors"
// @Router				/storage/upload/{tag}	[post]
// @Security			Authenticate Header
// @Security			Session
func (h *storageHandler) Upload(ctx *gin.Context) {
	input := &example.UploadInput{
		Tag:  ctx.Param("tag"),
		Size: 0,
	}

	res := h.service.Upload(input, ctx)
	ctx.JSON(res.Status, res)
}
