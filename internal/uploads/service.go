package uploads

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type Service struct {
	bucket        string
	publicBaseURL string
}

func NewService(bucket, publicBaseURL string) (*Service, error) {
	return &Service{bucket: bucket, publicBaseURL: strings.TrimRight(publicBaseURL, "/")}, nil
}

func (s *Service) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("bucket GCS não configurado")
	}
	defer file.Close()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("credenciais do GCP inválidas ou ausentes: %w", err)
	}
	defer client.Close()

	ext := strings.ToLower(path.Ext(header.Filename))
	objectName := fmt.Sprintf("%s/%s%s", strings.Trim(folder, "/"), uuid.NewString(), ext)
	writer := client.Bucket(s.bucket).Object(objectName).NewWriter(ctx)
	writer.ContentType = header.Header.Get("Content-Type")
	writer.CacheControl = "public, max-age=31536000"
	writer.Created = time.Now()

	if _, err := io.Copy(writer, file); err != nil {
		_ = writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	if s.publicBaseURL != "" {
		return s.publicBaseURL + "/" + objectName, nil
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, objectName), nil
}
