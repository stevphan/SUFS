# Group Project Notes

## Running CLI

* export tokens for S3 Bucket Access
* run cli from anywhere

## Runing the Name Node

* launch ec2 (t3.micro) (project AMI)
* export GOPATH
* copy latest project files to ec2
* go build
* run `./namenode`

## Running a Data Node

* launch ec2 (t3.micro) (project AMI)
* export GOPATH
* copy latest project files to ec2
* go build
* run `./datanode <name_node_addr>:8080 <directory_of_block_storage>`

## Notes

[Google Drive Documents](https://drive.google.com/drive/folders/1fM6cTPVd33H-8DGVdshsqFIiVyhb7iQa)
