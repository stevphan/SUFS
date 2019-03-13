# Group Project Notes

## Team Members

* Brian Radebaugh
* Grant Ludwig
* Steven Phan

## Running CLI

* export tokens for S3 Bucket Access
* run `./CLI <action> ...`
  * see DesignDoc.pdf to see how to use the CLI

## Runing the Name Node

* launch ec2 (t3.micro) (project AMI)
* export GOPATH
* copy latest project files to ec2
* go build
* run `./NameNode`

## Running a Data Node

* launch ec2 (t3.micro) (project AMI)
* export GOPATH
* copy latest project files to ec2
* go build
* run `./DataNode <name_node_ip>:8080 <directory_of_block_storage>`

## Libraries Used

* AWS SDK
  * `go get github.com/aws/aws-sdk-go`

## Demo

* small file url: `s3://reviews-radebaug/project/phoenix.png`
* big file url: `s3://reviews-radebaug/project/earth.png`
