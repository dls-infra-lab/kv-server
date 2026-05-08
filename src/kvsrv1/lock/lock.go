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
	id string
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// This interface supports multiple locks by means of the
// lockname argument; locks with different names should be
// independent.
func MakeLock(ck kvtest.IKVClerk, lockname string) *Lock {
	id := kvtest.RandValue(8)
	// each lock client must have a unique id
	lk := &Lock{ck: ck, lockname: lockname, id: id}
	return lk
}

func (lk *Lock) Acquire() {
	for {
		id, version, getErr := lk.ck.Get(lk.lockname)
		if (getErr == rpc.ErrNoKey)  {
			putErr := lk.ck.Put(lk.lockname, lk.id, 0)
			if putErr == rpc.OK {
				return
			}
		} else if (id == "" && getErr == rpc.OK) {
			putErr := lk.ck.Put(lk.lockname, lk.id, version)
			if putErr == rpc.OK {
				return
			}
		}
	}
}

func (lk *Lock) Release() {
	if id, version, getErr := lk.ck.Get(lk.lockname); id == lk.id && getErr == rpc.OK {
		lk.ck.Put(lk.lockname, "", version)
	}
}
