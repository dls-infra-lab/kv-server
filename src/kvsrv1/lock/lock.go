package lock

import (
	"kv-server/kvsrv1/rpc"
	"kv-server/kvtest1"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck kvtest.IKVClerk
	lockname string
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// This interface supports multiple locks by means of the
// lockname argument; locks with different names should be
// independent.
func MakeLock(ck kvtest.IKVClerk, lockname string) *Lock {
	lk := &Lock{ck: ck, lockname: lockname}
	
	// adding lock to kv store for client to access it
	// "0" = lock free, "1" = lock acquired
	ck.Put(lockname, "0", 0)
	return lk
}

func (lk *Lock) Acquire() {
	for {
		flag, version, getErr := lk.ck.Get(lk.lockname)
		if flag == "0" && getErr == rpc.OK {
			putErr := lk.ck.Put(lk.lockname, "1", version)
			if putErr == rpc.OK {
				return
			}
		}
	}
}

func (lk *Lock) Release() {
	if flag, version, getErr := lk.ck.Get(lk.lockname); flag == "1" && getErr == rpc.OK {
		lk.ck.Put(lk.lockname, "0", version)
	}
}
