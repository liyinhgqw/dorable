# Doracle

## Overview
Doracle is a distributed fault-tolerant **timestamp oracle** server. For single node oracle server, please refer to [Oracle](https://github.com/liyinhgqw/oracle). The fault-tolerant protocol is based on [Raft](https://github.com/goraft/raft), which is an alternative to Paxos.

## Usage

* Install
```
    go get github.com/liyinhgqw/doracle
```

* Run doracle servers
```
    go run cmd/drun.go -p 4001 ./node.1     // node.1
    go run cmd/drun.go -p 4002 -join localhost:4001 ./node.2     // node.1
    go run cmd/drun.go -p 4003 -join localhost:4001 ./node.3     // node.3
```

* For help information
```
    go run cmd/drun.go
```

* Use client

We implemented golang client stub, which is thread-safe.
```go
    client, err := doracle.NewClient([]string{
        "localhost:4001", "localhost:4002", "localhost:4003"})
    if err != nil {
        log.Fatalln(err)
    }
    if ts, err := client.TS(); err != nil {
        log.Println('ts error')
    } else {
        ...
    }
```

## Performance
Currently, qps is 0.6 million per second. Promising optimization might include:

* Use tcp ``rpc`` instead of http.

## Lessons Learned

* Close http response ``Body`` after reading it
* Do not reuse ``Reader`` or ``Writer`` (e.g., ``bytes.Buffer``) because after read or write, the offset is changed.

