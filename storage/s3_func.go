package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func UploadFileToS3(s *session.Session, fileHeader *multipart.FileHeader) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	file.Read(buffer)

	// create a unique file name for the file
	// 此处文件名称即为上传到aws bucket的名称，也是文件url路径的一部分。也可以在此处拼接url全路径，把域名部分拼接在前边即可。
	tempFileName := "pictures/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("test-bucket"), // bucket名称，把自己创建的bucket名称替换到此处即可
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, err
}
