package repo

import (
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	logger *logrus.Logger
	db     *gorm.DB
	Minio  *minio.Client
}

func NewRepository(dsn string, log *logrus.Logger, m *minio.Client) (*Repository, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:     gormDB,
		logger: log,
		Minio:  m,
	}, nil
}
