# Golang Assignment

## Console commands

```bash
# Terminal 1
# clone app repo and dependencies
go get github.com/Galish/golang-assignment

# open app directory
cd ${GOPATH}/src/github.com/Galish/golang-assignment

# install NodeJS dependencies
npm i

# run amqp, redis-server, tor, polipo and consul
npm run env
```

### Production

```bash
# Terminal 2: run Indexer microservice
${GOPATH}/bin/golang-assignment -service=indexer

# Terminal 3: run Crawler microservice
${GOPATH}/bin/golang-assignment -service=crawler

# Terminal 4:  run Frontend microservice
${GOPATH}/bin/golang-assignment -service=frontend

#Terminal 5: run client app at @localhost:8000
npm run client
```

### Development

```bash
# Run microservice
go run main.go -service=indexer/crawler/frontend

# build and install program
go install

# Run client app in development mode at @localhost:8000
npm run client:dev

# Create client app production build
npm run client:build
```
