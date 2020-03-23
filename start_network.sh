#!/bin/bash
docker network create --subnet 192.168.1.0/24 email

#these bash scripts build and run docker containers for the 'here.com' server (containing MSA, MTA, BBS microservices)
#here/msa/run.shh
#here/mta/run.shh
#here/bbs/run.shh

#These files below will activate a second MSA and MTA microservice, which can be used to send emails to the 'here.com' server
#there/msa/run.shh
#there/mta/run.shh

#once set up, you will be able to send messages between the two servers.
#please note that both microservices need to be running in order to receive emails between the two (otherwise the API calls are missed)
#in order to run the auto-updater for the MTA's send functions, the dockerfile/MTA programmes need to be running first.
#to enable to the auto-updaters, the bash scripts are called run-updater.sh (in here and there / mta)


#main commands for here.com server from within docker
#docker exec running-bbs curl http://192.168.1.7:8888/emails
#docker exec running-bbs curl http://192.168.1.8:8888/send

#main commands for there.com server from within docker
#docker exec running-bbs curl http://192.168.1.2:8888/emails
#docker exec running-bbs curl http://192.168.1.3:8888/send

#main commands for BBS server from within docker
#docker exec running-bbs curl http://192.168.1.9:8888/servers

#main commands for here.com server from localhost
#curl http://localhost:3000/emails
#curl http://localhost:3001/send

#main commands for there.com server from localhost
#curl http://localhost:3010/emails
#curl http://localhost:3011/send

#main commands for BBS server from localhost
#curl http://localhost:3003/servers
