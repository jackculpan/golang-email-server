#!/bin/bash

docker build -t mta2 .

docker run --ip 192.168.1.3 -p 3011:8888 -it --rm --network="email" --name running-mta2 mta2
