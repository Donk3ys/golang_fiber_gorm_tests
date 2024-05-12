#!/bin/bash

build=prod
ip="x.x.x.x"
imageName=api-$build
containerName=api-$build
pemFile=.pem

# Build locallly and upload image to the remote server
docker build  -f Dockerfile.stage -t $imageName . 
docker save $imageName > $imageName.tar 
scp -i ~/.ssh/$pemFile $imageName.tar ubuntu@$ip:/home/ubuntu

# Connect to the remote server
ssh -i ~/.ssh/$pemFile ubuntu@$ip << EOF
docker load < $imageName.tar
docker stop $containerName
docker rm $containerName
docker run -d --restart unless-stopped --network local -p 8000:8000 -v /home/ubuntu/public:/app/public -v /home/ubuntu/private:/app/private --name $containerName $imageName:latest
docker image prune -af
EOF

echo "Completed $build build $containerName deploy: ssh -i ~/.ssh/$pemFile ubuntu@$ip"
