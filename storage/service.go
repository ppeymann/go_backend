package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	example "github.com/ppeymann/go_backend"
)

type storageService struct {
	opts example.StorageOption
	repo example.StorageRepository
}

func NewService(opts example.StorageOption, repo example.StorageRepository) example.StorageService {
	return &storageService{
		opts: opts,
		repo: repo,
	}
}

func (s *storageService) Upload(input *example.UploadInput, ctx *gin.Context) *example.BaseResult {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{ErrBadUploadRequest.Error()},
		}
	}

	tempName := fmt.Sprintf("upload_%d_%s_%s", input.Claims.Subject)

	resFile, err := os.Create("./data/" + tempName)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	fex := filepath.Ext(resFile.Name())
	file, err := formFile.Open()
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	_, err = io.Copy(resFile, file)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	defer func() {
		_ = resFile.Close()
		_ = os.Remove(resFile.Name())
	}()

	fn := resFile.Name()
	result, err := s.repo.PutObject(strconv.Itoa(int(input.Claims.Subject)), input.Tag, fn, input.ContentType, fex, example.ObjectTag(input.Tag))
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	return result
}
