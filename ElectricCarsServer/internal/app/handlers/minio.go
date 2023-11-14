package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/Minio"
	"ElectricCarsServer/ElectricCarsServer/internal/app/utils"
	"fmt"
	"github.com/minio/minio-go"
	"mime/multipart"
)

func (h *Handler) createImageInMinio(file *multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := header.Filename
	if errName := utils.GenerateUniqueName(&objectName); errName != nil {
		return "", errName
	}

	if _, err := h.Minio.PutObject("electric-cars-server", objectName, *file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", Minio.MinioHost, Minio.BucketName, objectName), nil
}
