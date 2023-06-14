build:
	@go build -o bin/cashe

install:
	@go mod download

run: build
	@./bin/cashe --id=raft0

follower1: build
	@./bin/cashe --listenaddr :3001 --leaderaddr :3000 --raftaddr=:4001 --id=raft1

follower2: build
	@./bin/cashe --listenaddr :3002 --leaderaddr :3000 --raftaddr=:4002 --id=raft2

cl:
	@go run client/runtest/main.go

test:
	@go test -v ./...