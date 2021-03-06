package objectstorage

import (
	"bytes"
	"io"
	"net/http"

	"github.com/minio/minio-go"
	"github.com/nkhang/pluto/pkg/logger"
)

type minioClient struct {
	client *minio.Client
}

func NewMinioClient(endpoint, accessKey, secretKey string, ssl bool) (*minioClient, error) {
	client, err := minio.New(endpoint, accessKey, secretKey, ssl)
	if err != nil {
		return nil, err
	}
	return &minioClient{
		client: client,
	}, nil
}

func (c *minioClient) Put(collection, filename string, reader io.Reader, size int64, contentType string) (int64, error) {
	return c.client.PutObject(collection, filename, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
}

func (c *minioClient) PutImage(collection, filename string, reader io.Reader, size int64) (int64, error) {
	if ok, _ := c.client.BucketExists(collection); !ok {
		logger.Infof("making bucket %s", collection)
		err := c.client.MakeBucket(collection, "ap-southeast-1")
		if err != nil {
			panic(err)
		}
	}
	b := make([]byte, size)
	n, err := reader.Read(b)
	if err != nil || n == 0 {
		return int64(n), err
	} 
	contentType := http.DetectContentType(b)
	return c.Put(collection, filename, bytes.NewReader(b), size, contentType)
}
