build:
	@go build -o bin/cashe

run: build
	@./bin/cashe --id=raft0

follower1: build
	@./bin/cashe --listenaddr :3001 --leaderaddr :3000 --raftaddr=:4001 --id=raft1

follower2: build
	@./bin/cashe --listenaddr :3002 --leaderaddr :3000 --raftaddr=:4002 --id=raft2

ct: 
	@go build -o bin/client client/main.go
	@./bin/client

test:
	@go test -v ./...