package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadS3FileInFile(url string) (*os.File, int64) {
	bucket, item := parseS3Url(url)

	file, err := os.Create(tempS3DownloadFileName)
	shared.CheckErrorAndFatal("Unable to create temporary file for S3 download", err)

	downloadedByteCount := downloadS3FileUsingAwsSdk(bucket, item, file)

	fileInfo, err := file.Stat()
	shared.CheckErrorAndCleanAndFatal("", err, func() {
		os.Remove(tempS3DownloadFileName)
	})
	if downloadedByteCount != fileInfo.Size() {
		os.Remove(tempS3DownloadFileName)
		log.Fatalln("Downloaded byte count does not match the files byte count")
	}

	size := fileInfo.Size()
	shared.VerbosePrintln("Successfully downloaded S3 file")

	return file, size
}

func downloadS3FileUsingAwsSdk(bucket, item string, file *os.File) int64 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessId, awsSecretAccessToken, ""),
	})
	shared.CheckErrorAndCleanAndFatal("Unable to create session with AWS", err, func() {
		os.Remove(tempS3DownloadFileName)
	})

	shared.VerbosePrintln("AWS session created")

	downloader := s3manager.NewDownloader(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}

	shared.VerbosePrintln(fmt.Sprintf("Downloading item '%s' from bucket '%s'", item, bucket))
	downloadedByteCount, err := downloader.Download(file, input)
	shared.CheckErrorAndCleanAndFatal("Unable to download file from s3", err, func() {
		os.Remove(tempS3DownloadFileName)
	})

	return downloadedByteCount
}

func parseS3Url(urlString string) (string, string) {
	log.Println("url:", urlString)
	components, err := url.Parse(urlString)
	if err != nil {
		shared.CheckErrorAndFatal("Unable to parse S3 URL", err)
	}

	log.Println("Components:", components)
	return components.Host, components.Path[1:]
}
