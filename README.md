## Pre-requisite 
1. brew install grpcurl and brew install golang
2. brew install protobuf 
3. brew install postgresql
4. brew install libpq
5. go get github.com/joho/godotenv/cmd/godotenv
6. brew install golang-migrate
7. Docker 


## Assumptions
Golang is setup on machine to execute program.
Using grpcurl to send shel script in binary format, the content is store in db as string.
I tried using byte[] and got stuck with grpc EOF error. Updated code to send binary encoded string 
using https://www.base64decode.org/ and then convert the binary to string and store
script as string. 


## Execution
1. make postgres createdb migrateup server
2. ./run_take_home.sh
3. make migratedown dropdb



## Notes
1. I have not enabled, reflection on GRPC and that's the reason working with protoset file.
2. I have created a GRPC server to support highly performant and scalable API's and makes use of binary data rather than just text which makes the communication more compact and more efficient.
3. Do to time constraint :
   a. I did not add background jobs.
   b. Did not create docker container.
4. Can improve the code:
    a. Implementing schedule service apart from api server.
    b. Implementing cron schedule for task to execute in background.
    c. Implementing monitoring for functionality. 
    