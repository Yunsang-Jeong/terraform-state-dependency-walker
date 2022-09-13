package awsApi

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errors "github.com/pkg/errors"
)

func stringSuffixFilter(target string, suffixes []string) bool {
	flag := false

	for _, suffix := range suffixes {
		if strings.HasSuffix(target, suffix) {
			flag = true
		}
	}

	return flag
}

func GetFilteredBucketObjectList(awsConfig aws.Config, bucketName string, objectSuffixFilter []string) (*[]string, error) {
	client := s3.NewFromConfig(awsConfig)

	paginator := s3.NewListObjectsV2Paginator(client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		},
	)

	objectList := []string{}

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "[awsAPi:GetFilteredBucketObjectList]fail to paginate of ListObjectV2")
		}

		for _, obj := range page.Contents {
			if !stringSuffixFilter(*obj.Key, objectSuffixFilter) {
				continue
			}

			objectList = append(objectList, *obj.Key)
		}
	}

	return &objectList, nil
}

func DownloadBucketObjectToBuffer(awsConfig aws.Config, bucketName string, objectName string) ([]byte, error) {
	client := s3.NewFromConfig(awsConfig)

	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(client)
	_, err := downloader.Download(context.TODO(),
		buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[awsAPi:DownloadBucketObjectToBuffer] fail to download object(%s) from bucket(%s)", objectName, bucketName))
	}

	return buffer.Bytes(), nil
}
