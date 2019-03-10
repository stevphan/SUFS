package main

import (
	"fmt"
	"log"
	"net/url"
	"shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadS3File(url string) []byte {
	bucket, item := parseS3Url(url)

	shared.VerbosePrintln(fmt.Sprintf("Downloading item '%s' from bucket '%s'", item, bucket))

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessId, awsSecretAccessToken, ""),
	})
	if err != nil {
		log.Fatalln("Unable to create session with AWS:", err)
	}

	shared.VerbosePrintln("AWS session created")

	downloader := s3manager.NewDownloader(sess)
	outputBuffer := &aws.WriteAtBuffer{}
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}

	downloadedByteCount, err := downloader.Download(outputBuffer, input)
	if err != nil {
		log.Fatalln("Unable to download S3 bucket item:", err)
	}

	shared.VerbosePrintln(fmt.Sprintf("Successfully downloaded S3 file (%d bytes)", downloadedByteCount))

	return outputBuffer.Bytes()
}

func parseS3Url(urlString string) (string, string) {
	components, err := url.Parse(urlString)
	if err != nil {
		shared.CheckErrorAndFatal("Unable to parse S3 URL", err)
	}

	return components.Host, components.Path[1:]
}
