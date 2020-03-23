#!/bin/bash

docker build -t bbs .

docker run --ip 192.168.1.9 -p 3002:8888 -it --rm --network="email" --name running-bbs bbs
