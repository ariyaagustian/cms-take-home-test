package minio

import (
	"bytes"
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Minio  *minio.Client
	Bucket string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) *Client {
	m, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("gagal init MinIO: %v", err)
	}

	// Cek dan buat bucket jika belum ada
	ctx := context.Background()
	exists, err := m.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("gagal cek bucket: %v", err)
	}
	if !exists {
		if err := m.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("gagal buat bucket: %v", err)
		}
		log.Printf("Bucket %s berhasil dibuat", bucket)
	}

	return &Client{
		Minio:  m,
		Bucket: bucket,
	}
}

func (c *Client) Upload(ctx context.Context, objectName string, contentType string, data []byte) error {
	reader := bytes.NewReader(data)

	_, err := c.Minio.PutObject(ctx, c.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) GenerateSignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values) // bisa diisi content-disposition atau lainnya
	url, err := c.Minio.PresignedGetObject(ctx, c.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	err := c.Minio.RemoveObject(ctx, c.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Gagal hapus file di MinIO: %v", err)
	}
	return err
}
