package kvsrv

import (
	"log"
	"sync"

	"kv-server/kvsrv1/rpc"
	"kv-server/labrpc"
	"kv-server/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type MapValue struct {
	Value   string
	Version rpc.Tversion
}

type KVServer struct {
	mu sync.Mutex
	kvMap map[string]MapValue
}

func MakeKVServer() *KVServer {
	kv := &KVServer{kvMap: map[string]MapValue{}}
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if res, ok := kv.kvMap[args.Key]; ok {
		reply.Value = res.Value
		reply.Version = res.Version
		reply.Err = rpc.OK
	} else {
		reply.Err = rpc.ErrNoKey
	}

}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	res, ok := kv.kvMap[args.Key]
	if (!ok && args.Version == 0) {
		kv.kvMap[args.Key] = MapValue{
			Value: args.Value,
			Version: args.Version + 1,
		}
		reply.Err = rpc.OK
	} else if (ok && res.Version == args.Version){
		kv.kvMap[args.Key] = MapValue{
			Value: args.Value,
			Version: args.Version + 1,
		}
		reply.Err = rpc.OK
	} else if (ok && res.Version != args.Version) {
		reply.Err = rpc.ErrVersion
	} else {
		reply.Err = rpc.ErrNoKey
	}
}



// You can ignore all arguments; they are for replicated KVservers
func StartKVServer(tc *tester.TesterClnt, ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []any {
	kv := MakeKVServer()
	return []any{kv}
}
