# 🚀 **Cashe - Blazing Fast Distributed Key-Value Store implemented with Raft and Golang**

Cashe is a distributed key-value store database that leverages the Raft consensus algorithm and is implemented using the Go programming language. It provides a highly available and fault-tolerant solution for storing and retrieving data across a cluster of nodes.

# **Features**
### ✨ **High Availability**: Cashe ensures that data remains available even in the event of node failures, ensuring continuous service uptime.

### 🔒 **Consistency**: The Raft consensus algorithm employed by Cashe ensures strong consistency across all nodes in the cluster.

### 📈 **Scalability**: Cashe is designed to scale horizontally by adding more nodes to the cluster, allowing it to handle increasing workloads with ease.

### 🔑 **Key-Value Store**: Cashe stores data in a simple key-value format, providing efficient access to stored information.

### ⚡️ **Fast and Efficient**: The Go programming language used in Cashe ensures high performance and efficient resource utilization.


### 🌐 **Distributed Architecture**: Cashe distributes data across multiple nodes, allowing for better load distribution and fault tolerance.


## **Getting Started**

1. Install Dependencies
```shell
make install
```
2. Start a master node 
```shell
./bin/cashe --id=raft0 --listenaddr=:3000  --raftaddr=:4000
```
3. Start second and third node
```shell
./bin/cashe --listenaddr :3001 --leaderaddr :3000 --raftaddr=:4001 --id=raft1

./bin/cashe --listenaddr :3002 --leaderaddr :3000 --raftaddr=:4002 --id=raft2
```
