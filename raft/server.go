package rf

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

type RaftServer struct {
	*raft.Raft
}

func New(serverID, port string) *RaftServer {
	var (
		cfg           = raft.DefaultConfig()
		fsm           = &raft.MockFSM{}
		logStore      = raft.NewInmemStore()
		stableStore   = raft.NewInmemStore()
		snapShotStore = raft.NewInmemSnapshotStore()
		timeout       = time.Second * 10
	)
	// os.Getenv("RAFT_MASTER_NODE_ID")

	cfg.LocalID = raft.ServerID(serverID)
	addr := fmt.Sprintf("127.0.0.1%s", port)
	tr, err := raft.NewTCPTransport(addr, nil, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	server := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(cfg.LocalID),
		Address:  tr.LocalAddr(),
	}

	serverConfig := raft.Configuration{
		Servers: []raft.Server{server},
	}
	r, err := raft.NewRaft(cfg, fsm, stableStore, logStore, snapShotStore, tr)
	if err != nil {
		log.Fatal(err)
	}
	if err := r.BootstrapCluster(serverConfig).Error(); err != nil {
		log.Fatal(err)
	}
	return &RaftServer{Raft: r}
}
