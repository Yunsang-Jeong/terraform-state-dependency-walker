package aws

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

func (a *AWS) DownloadS3ObjectToBuffer(bucketName string, objectSuffixFilter []string) (map[string][]byte, error) {
	bufferMap := map[string][]byte{}

	client := s3.NewFromConfig(a.AWSClientConfig)

	paginator := s3.NewListObjectsV2Paginator(client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		},
	)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "fail to paginate of ListObjectV2")
		}

		for _, obj := range page.Contents {

			if !stringSuffixFilter(*obj.Key, objectSuffixFilter) {
				continue
			}

			buf := manager.NewWriteAtBuffer([]byte{})

			downloader := manager.NewDownloader(client)
			_, err := downloader.Download(context.TODO(),
				buf,
				&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    obj.Key,
				},
			)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("fail to download object(%s) from aws s3 bucket(%s)", *obj.Key, bucketName))
			}

			bufferMap[*obj.Key] = buf.Bytes()
		}
	}

	return bufferMap, nil
}
