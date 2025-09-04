package s3Client

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/common/constants"
)

type S3Client struct {
	Sentry sentry.SentryLayer
}

func NewS3Client() *S3Client {
	return &S3Client{
		Sentry: sentry.InitializeSentry(),
	}
}

func InitializeS3Client() S3ClientInterface {
	return NewS3Client()
}

type S3ClientInterface interface {
	UploadMultipartFile(file multipart.File, fileHeader *multipart.FileHeader, fileName string) (string, string)
	UploadLocalFile(file *os.File, fileName string, bucket string) (string, string)
	UploadFileFromString(base64File string, fileName string) (string, string)
}

func ConnectAwsS3() (*session.Session, error) {
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(constants.Region),
			Credentials: credentials.NewStaticCredentials(
				constants.S3AccessKeyID,
				constants.S3SecretAccessKey,
				"",
			),
		})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *S3Client) UploadMultipartFile(file multipart.File, fileHeader *multipart.FileHeader, fileName string) (string, string) {
	if file == nil {
		return "", ""
	}
	session, err := ConnectAwsS3()
	if err != nil {
		return "", ""
	}
	s3Svc := s3.New(session)
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	size := fileHeader.Size
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return "", ""
	}

	uploadParams := &s3manager.UploadInput{
		Bucket: &constants.Bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(buffer),
	}

	_, err = uploader.Upload(uploadParams)

	if err != nil {
		return "", ""
	}
	filepath := "https://" + constants.Bucket + "." + "s3-" + constants.Region + ".amazonaws.com/" + fileName
	return filepath, fileName
}

func (s *S3Client) UploadLocalFile(file *os.File, fileName string, bucket string) (string, string) {
	session, err := ConnectAwsS3()
	if err != nil {
		return "", ""
	}
	s3Svc := s3.New(session)
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	if file == nil {
		return "", ""
	}

	uploadParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &fileName,
		Body:   file,
	}
	//upload to the s3 bucket
	_, err = uploader.Upload(uploadParams)

	if err != nil {
		return "", ""
	}
	filepath := "https://" + bucket + "." + "s3-" + constants.Region + ".amazonaws.com/" + fileName
	return filepath, fileName
}

func (s *S3Client) UploadFileFromString(base64File string, fileName string) (string, string) {
	session, err := ConnectAwsS3()
	if err != nil {
		return "", ""
	}
	s3Svc := s3.New(session)
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	data, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return "", ""
	}

	uploadParams := &s3manager.UploadInput{
		Bucket: &constants.Bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(data),
	}
	//upload to the s3 bucket
	_, errUp := uploader.Upload(uploadParams)

	if errUp != nil {
		return "", ""
	}
	filepath := "https://" + constants.Bucket + "." + "s3-" + constants.Region + ".amazonaws.com/" + fileName
	return filepath, fileName

}
