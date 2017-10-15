// Package s3fs implements fs.FS for s3 object storage
package s3fs

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/fs"
)

var sess *session.Session
var uploader *s3manager.Uploader
var downloader *s3manager.Downloader
var client = &http.Client{}

func init() {
	var err error

	prov := &credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     config.AWSAccessKeyID(),
			SecretAccessKey: config.AWSSecretKey(),
		},
	}

	creds := credentials.NewCredentials(prov)

	sess = session.Must(session.NewSession(&aws.Config{
		Logger:      config.Logger,
		Region:      aws.String(config.AWSRegion()),
		Endpoint:    config.AWSBaseURL(),
		HTTPClient:  client,
		Credentials: creds,
	}))

	uploader = s3manager.NewUploader(sess)
	downloader = s3manager.NewDownloader(sess)
}

// FS is the entrypoint for the s3fs backed fs.FS
type FS struct {
	Bucket string
}

// New will return a new fs.FS for use with s3 object storage
func New() FS {
	return FS{config.AWSBucket()}
}

// Get the file from the s3fs bucket at path
// TODO: Find a way to not use tmp files here, this is slow and scales poorly.
func (s3fs FS) Get(path string) (io.ReadCloser, error) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		tmp.Close()
		return nil, err
	}

	_, err = downloader.Download(tmp, &s3.GetObjectInput{
		Bucket: aws.String(s3fs.Bucket),
		Key:    aws.String(path),
	})

	return tmp, err
}

// Save will save the file to the s3fs bucket
func (s3fs FS) Save(file *os.File) (string, error) {
	fn, err := fs.GenUniqueFileName(file)
	if err != nil {
		return "", err
	}

	canonicalURL := "/" + fn

	// Retrieve the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	length := strconv.Itoa(len(content))
	reader := bytes.NewReader(content)

	req, err := http.NewRequest("PUT", s3fs.url()+"/"+fn, reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Length", length)
	req.Header.Set("Authorization")

	return "", nil
}
