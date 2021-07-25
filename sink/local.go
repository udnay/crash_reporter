package sink

import (
	"io"
	"os"

	"go.uber.org/zap"
)

type LocalStore struct {
	Log *zap.SugaredLogger
}

func (l *LocalStore) Store(in io.Reader, fileName string) error {

	f, err := os.Create(fileName)
	if err != nil {
		l.Log.Errorf("Unable to create %s: %v", fileName, err)
		return err
	}
	defer f.Close()

	b := make([]byte, 2048)

	for {
		read, err := in.Read(b)
		if err != nil && err != io.EOF {
			l.Log.Errorf("Unable to read from in buffer: %v", err)
			return err
		}

		_, err = f.Write(b[:read])
		if err != nil {
			l.Log.Errorf("Unable to write to %s: %v", fileName, err)
			return err
		}

		if read < 2048 {
			break
		}
	}

	return nil
}
