package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/Minio"
	"net/url"
	"path"
)

func (r *Repository) deleteImageFromMinio(objectName string) error {
	// Удаление объекта из MinIO
	filename, err := extractFilenameFromURL(objectName)
	err = r.Minio.RemoveObject(Minio.BucketName, filename)
	if err != nil {
		return err
	}

	return nil
}

func extractFilenameFromURL(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	// Используйте path.Base для извлечения последней части пути, которая будет являться именем файла
	filename := path.Base(u.Path)

	return filename, nil
}
