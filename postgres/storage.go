package postgres

import (
	example "expamle"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
		table:    "storage_entities",
		secret:   secret,
	}, nil
}

func (r *storageRepository) PutObject(bucketName, objectName, path, ct, ext string, tag example.ObjectTag) (*example.BaseResult, error) {
	return nil, nil
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
