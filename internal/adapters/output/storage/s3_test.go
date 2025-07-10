package storage_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/storage"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type s3TestSuite struct {
	suite.Suite

	ctx          context.Context
	mockS3Client *mocks.MockS3Client
	s3Adapter    *storage.S3Client
}

func (suite *s3TestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.mockS3Client = mocks.NewMockS3Client(ctrl)
	suite.s3Adapter = storage.NewS3Adapter(suite.mockS3Client, "bucket-videos-entrada", "bucket-videos-saida")
}

func (suite *s3TestSuite) AfterTest(_, _ string) {
	_ = os.Remove("local/video.mp4")
	_ = os.RemoveAll("local")
}

func Test_S3TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(s3TestSuite))
}

func (suite *s3TestSuite) Test_DownloadFile() {
	st := suite.T()

	st.Run("should download file successfully", func(t *testing.T) {
		objectKey := "video.mp4"

		suite.mockS3Client.EXPECT().
			GetObject(gomock.Any(), gomock.Any()).
			Return(&s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("conteudo do arquivo"))),
			}, nil)

		downloadedFile, err := suite.s3Adapter.DownloadFile(suite.ctx, objectKey)
		tempDir := os.TempDir()
		suite.True(strings.HasPrefix(downloadedFile.Path, tempDir))
		suite.NoError(err)
		suite.NotEmpty(downloadedFile.Path)
		suite.Contains(downloadedFile.Path, "video-")
		suite.Contains(downloadedFile.Path, ".tmp")
	})

	st.Run("should return error when S3 GetObject fails", func(t *testing.T) {
		objectKey := "video.mp4"

		suite.mockS3Client.EXPECT().
			GetObject(gomock.Any(), gomock.Any()).
			Return(nil, io.ErrUnexpectedEOF)

		downloadedFile, err := suite.s3Adapter.DownloadFile(suite.ctx, objectKey)
		suite.Error(err)
		suite.Nil(downloadedFile)
	})

	st.Run("Should return error when failed to copy S3 content", func(t *testing.T) {
		objectKey := "video.mp4"

		errorReader := io.NopCloser(&errorOnRead{})

		suite.mockS3Client.EXPECT().
			GetObject(gomock.Any(), gomock.Any()).
			Return(&s3.GetObjectOutput{
				Body: errorReader,
			}, nil)

		downloadedFile, err := suite.s3Adapter.DownloadFile(suite.ctx, objectKey)
		suite.Error(err)
		suite.Nil(downloadedFile)
	})
}

func (suite *s3TestSuite) Test_UploadFile() {
	st := suite.T()

	st.Run("should upload file successfully", func(t *testing.T) {
		localFilePath := "local/video.mp4"
		objectKey := "upload/video.mp4"

		err := os.MkdirAll("local", 0755)
		suite.NoError(err)

		file, err := os.Create(localFilePath)
		suite.NoError(err)
		_, _ = file.Write([]byte("conteudo de teste"))
		file.Close()

		suite.mockS3Client.EXPECT().
			PutObject(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		err = suite.s3Adapter.UploadFile(suite.ctx, localFilePath, objectKey)
		suite.NoError(err)
	})

	st.Run("should return error when S3 PutObject fails", func(t *testing.T) {
		localFilePath := "local/video.mp4"
		objectKey := "upload/video.mp4"

		err := os.MkdirAll("local", 0755)
		suite.NoError(err)

		file, err := os.Create(localFilePath)
		suite.NoError(err)
		_, _ = file.Write([]byte("conteudo de teste"))
		file.Close()

		suite.mockS3Client.EXPECT().
			PutObject(gomock.Any(), gomock.Any()).
			Return(nil, io.ErrUnexpectedEOF)

		err = suite.s3Adapter.UploadFile(suite.ctx, localFilePath, objectKey)
		suite.Error(err)
	})
}

type errorOnRead struct{}

func (e *errorOnRead) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
