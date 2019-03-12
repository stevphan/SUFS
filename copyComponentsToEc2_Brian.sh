copyComponentsToEc2 () {
    if [ -z $1 ]; then
        echo "Must pass in EC2 IP"
        return -1
    fi

    echo "- Copying CLI..."
    scp -r -i ~/.ssh/aws-seattleu-cloud-computing-1.pem ~/Documents/Seattle\ University/2019\ \(1\ Winter\)\ CPSC-5910-02\ Cloud\ Computing/project_repo/src/CLI/*.go ec2-user@$1:/home/ec2-user/project_repo/src/CLI/

    echo "- Copying Data Node..."
    scp -r -i ~/.ssh/aws-seattleu-cloud-computing-1.pem ~/Documents/Seattle\ University/2019\ \(1\ Winter\)\ CPSC-5910-02\ Cloud\ Computing/project_repo/src/DataNode/*.go ec2-user@$1:/home/ec2-user/project_repo/src/DataNode/

    echo "- Copying Name Node..."
    scp -r -i ~/.ssh/aws-seattleu-cloud-computing-1.pem ~/Documents/Seattle\ University/2019\ \(1\ Winter\)\ CPSC-5910-02\ Cloud\ Computing/project_repo/src/NameNode/*.go ec2-user@$1:/home/ec2-user/project_repo/src/NameNode/

    echo "- Copying shared..."
    scp -r -i ~/.ssh/aws-seattleu-cloud-computing-1.pem ~/Documents/Seattle\ University/2019\ \(1\ Winter\)\ CPSC-5910-02\ Cloud\ Computing/project_repo/src/shared/*.go ec2-user@$1:/home/ec2-user/project_repo/src/shared/

    echo "- Finished"
}

copyComponentsToEc2 "$@"
