#!/bin/bash

docker build -t mta .

docker run --ip 192.168.1.8 -p 3001:8888 -it --rm --network="email" --name running-mta mta
