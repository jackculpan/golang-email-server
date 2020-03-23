#!/bin/bash

docker build -t msa .

docker run --ip 192.168.1.7 -p 3000:8888 -it --rm --network="email" --name running-msa msa
