package sink

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type GCSStore struct {
	CredentialsFile string
	Bucket          string
	Log             *zap.SugaredLogger
}

func (s *GCSStore) Store(in io.Reader, fileName string) error {
	ctx := context.Background()

	gcsClient, err := storage.NewClient(ctx, option.WithCredentialsFile(s.CredentialsFile))
	if err != nil {
		s.Log.Infof("Error creating new client: %v", err)
		return err
	}

	bucket := gcsClient.Bucket(s.Bucket)

	b := make([]byte, 2048)

	coreFile := bucket.Object(fileName)

	coreFileWriter := coreFile.NewWriter(ctx)
	defer coreFileWriter.Close()

	for {
		read, err := in.Read(b)
		if err != nil && err != io.EOF {
			s.Log.Infof("Couldn't read from stdin: %v \n", err)
			return err
		}

		_, err = coreFileWriter.Write(b)
		if err != nil {
			s.Log.Infof("Couldn't write to core file %s/%s: %v", s.Bucket, fileName, err)
			return err
		}

		if read < 2048 {
			break
		}
	}

	return nil
}
