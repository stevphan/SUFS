copyComponentsToEc2 () {
    if [ -z $1 ]; then
        echo "Must pass in EC2 IP"
        return -1
    fi

    echo "- Copying CLI..."
    scp -r -i "C:/Users/Grant Ludwig/Documents/AWSkey.pem" "C:/Users/Grant Ludwig/Desktop/Cloud-Computing-Project/src/CLI/*.go" ec2-user@$1:/home/ec2-user/project_repo/src/CLI/

    echo "- Copying Data Node..."
    scp -r -i "C:/Users/Grant Ludwig/Documents/AWSkey.pem" "C:/Users/Grant Ludwig/Desktop/Cloud-Computing-Project/src/DataNode/*.go" ec2-user@$1:/home/ec2-user/project_repo/src/DataNode/

    echo "- Copying Name Node..."
    scp -r -i "C:/Users/Grant Ludwig/Documents/AWSkey.pem" "C:/Users/Grant Ludwig/Desktop/Cloud-Computing-Project/src/NameNode/*.go" ec2-user@$1:/home/ec2-user/project_repo/src/NameNode/

    echo "- Copying shared..."
    scp -r -i "C:/Users/Grant Ludwig/Documents/AWSkey.pem" "C:/Users/Grant Ludwig/Desktop/Cloud-Computing-Project/src/shared/*.go" ec2-user@$1:/home/ec2-user/project_repo/src/shared/

    echo "- Finished"
}

copyComponentsToEc2 "$@"
