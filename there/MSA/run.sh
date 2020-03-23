#!/bin/bash

docker build -t msa2 .

docker run --ip 192.168.1.2 -p 3010:8888 -it --rm --network="email" --name running-msa2 msa2
