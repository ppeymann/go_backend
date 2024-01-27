package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	example "github.com/ppeymann/go_backend"
	"github.com/ppeymann/go_backend/utils"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// storageRepository implements mml_be.StorageRepository interface
type storageRepository struct {
	pg       *gorm.DB
	minio    *minio.Client
	opts     example.StorageOption
	database string
	table    string
	secret   string
}

// NewStorageRepository creates new accountRepository instance and fill its properties with specific args
func NewStorageRepository(pg *gorm.DB, opts example.StorageOption, database, secret string) (example.StorageRepository, error) {
	client, err := minio.New(opts.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.User, opts.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &storageRepository{
		pg:       pg,
		minio:    client,
		opts:     opts,
		database: database,
		table:    "storage",
		secret:   secret,
	}, nil
}

func (r *storageRepository) PutObject(bucketName, objectName, path, ct, ext string, tag example.ObjectTag) (*example.BaseResult, error) {
	aid, err := strconv.Atoi(bucketName)
	if err != nil {
		return nil, example.ErrUserPrincipalsNotFount
	}

	bn := fmt.Sprintf("user-%s-bucket", bucketName)
	fileId := ksuid.New().String()
	file := &example.StorageEntity{
		Model:       gorm.Model{},
		Account:     uint(aid),
		Tag:         objectName,
		Bucket:      bn,
		ContentType: ct,
		FileName:    fmt.Sprintf("%s%s", fileId, ext),
	}

	err = r.pg.Transaction(func(tx *gorm.DB) error {
		if res := r.Model().Create(file).Error; res != nil {
			str := res.(*pgconn.PgError).Message
			if strings.Contains(str, "duplicate key value") {
				return errors.New("account with specific id already exists")
			}
			return res
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	data := fmt.Sprintf(`{"account_id":"%d","tag":"%s","id":"%d"}`, aid, objectName, file.ID)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	err = r.minio.MakeBucket(ctx, bn, minio.MakeBucketOptions{
		Region: r.opts.Region,
	})
	if err != nil {
		exists, existErr := r.minio.BucketExists(ctx, bn)
		if existErr != nil && !exists {
			return nil, err
		}
	}

	info, err := r.minio.FPutObject(ctx, bn, file.FileName, path, minio.PutObjectOptions{
		UserTags:    map[string]string{"target": string(tag)},
		ContentType: ct,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(info)

	cypher, err := utils.EncryptText(data, r.secret)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return &example.BaseResult{
		Result: cypher,
	}, err

}

func (r *storageRepository) Migrate() error {
	return r.pg.AutoMigrate(&example.StorageEntity{})
}

func (r *storageRepository) Name() string {
	return r.table
}

func (r *storageRepository) Model() *gorm.DB {
	return r.pg.Model(&example.StorageEntity{})
}
