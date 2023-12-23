package example

import (
	"expamle/authorization"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	ObjectTag string

	// StorageService represents method signatures for api storage service endpoint.
	// so any object that stratifying this interface can be used as storage service for api endpoint.
	StorageService interface {
		Upload(input *UploadInput, ctx *gin.Context) *BaseResult
	}

	// StorageRepository represents method signatures for storage service domain repository.
	// so any object that stratifying this interface can be used as storage service domain repository
	StorageRepository interface {
		PutObject(bucketName, objectName, path, ct, ext string, tag ObjectTag) (*BaseResult, error)
	}

	// UploadInput is DTO for transferring file upload request params.
	UploadInput struct {
		Claims      *authorization.Claims `json:"-"`
		Tag         string                `json:"tag"`
		Size        int64                 `json:"size"`
		ContentType string                `json:"content_type"`
		FileName    string                `json:"file_name"`
	}

	// DownloadInput is DTO for transferring file download request params.
	DownloadInput struct {
		Token     string `json:"-"`
		AccountId string `json:"account_id"`
		Tag       string `json:"tag"`
		Id        string `json:"id"`
		Size      uint   `json:"-"`
	}

	// StorageEntity is entity of storage file object item info.
	//
	//  swagger:model StorageEntity
	StorageEntity struct {
		gorm.Model

		// Account id of object owner
		Account uint `json:"account" gorm:"column:account;index"`

		// Tag for stored object to manage access permission
		Tag string `json:"tag" gorm:"column:tag;index"`

		// Buket name of stored object that
		Bucket string `json:"bucket" gorm:"column:bucket;index"`

		// ContentType of stored object
		ContentType string `json:"content_type" gorm:"column:content_type;index"`

		// FileName of stored object for back on download request
		FileName string `json:"file_name" gorm:"column:file_name;index"`
	}
)

const (
	PublicTag  ObjectTag = "public"
	PrivateTag ObjectTag = "private"
)
