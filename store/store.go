package store

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type Store struct {
	credentials *StorageCredentials
	session     *session.Session
	storage     *s3.S3
}

func New(credentials *StorageCredentials) *Store {
	return &Store{
		credentials: credentials,
	}
}

func (store *Store) Open() error {
	credentials := credentials.NewStaticCredentials(store.credentials.Key, store.credentials.Secret, "")
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(store.credentials.Region),
		Credentials: credentials,
	})

	if err != nil {
		return err
	}
	store.session = session

	s3Session := s3.New(session)
	store.storage = s3Session
	return nil
}

func (store *Store) GetBuckets() error {
	result, err := store.storage.ListBuckets(nil)
	if err != nil {
		return err
	}
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
	return nil
}

func (store *Store) GetLatestEpoch() (int, error) {
	var epoch int = 0
	ctx := context.Background()
	objects := []string{}

	epochFilesPath := store.credentials.Prefix + "/"

	err := store.storage.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(store.credentials.Bucket),
		Prefix: aws.String(epochFilesPath),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {

		for _, o := range p.Contents {
			key := aws.StringValue(o.Key)

			objects = append(objects, key)
			objectPathSegments := strings.Split(key, "/")
			if len(objectPathSegments) == 2 {
				storedEpoch := objectPathSegments[1]
				readEpoch, err := strconv.Atoi(storedEpoch)
				if err != nil {
					fmt.Println("Detected a non-epoch file")
				}
				if epoch < readEpoch {
					epoch = readEpoch
				}
			}
		}
		return true // continue paging
	})
	if err != nil {
		return epoch, err
	}

	// fmt.Println("Objects in bucket:", objects)
	return epoch, nil
}

func (store *Store) PutEpoch(epoch int, nodeId uuid.UUID) error {
	epochFilePath := store.credentials.Prefix + "/" + strconv.Itoa(epoch)
	req, _ := store.storage.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(store.credentials.Bucket),
		Key:    aws.String(epochFilePath),
		Body:   strings.NewReader(nodeId.String()),
	})
	req.HTTPRequest.Header.Add("If-None-Match", "*")

	err := req.Send()
	if err != nil {
		fmt.Println("Error putting file")
		return err
	}

	if req.Error != nil {
		if aerr, ok := req.Error.(awserr.Error); ok {
			if aerr.Code() == request.CanceledErrorCode {
				// If the SDK can determine the request or retry delay was canceled
				// by a context the CanceledErrorCode error code will be returned.
				fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", req.Error)
			}
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", req.Error)
		}
		return err
	}
	return nil
}
