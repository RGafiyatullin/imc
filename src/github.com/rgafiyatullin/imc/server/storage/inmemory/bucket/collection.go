// The two collections defined here (KV and TTL) provide a simple interface for the Key-Value storage with per-key TTL (msec-wise).
//
// KV associates a given key with the value along with the some deadline.
// Deadline is an ephemeral point in time: it's measured in ticks provided from the outside for several reasons.
// 1. time.Now's resolution is default for the OS: usually microseconds; this is too fine for us - therefore potentially costy;
// 2. dependency on time.Now introduces a major side effect into a collection; enough said.
//
// The collections are not meant to be thread-safe: they are supposed to be accessed sequentially.
// In order to scale out -- use multiple collections as shards over the keyspace.
package bucket

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
)

type KV interface {
	// Returns nillable KVEntry if there is one associated with the given key
	Get(key string) (KVEntry, bool)

	// Creates and stores a new KVEntry associated with the given key
	Set(key string, value data.Value, validThru int64)

	Del(key string)
}

type KVEntry interface {
	validThru() int64
	value() data.Value
}

type TTL interface {
	SetTTL(k string, deadline int64)
	FetchTimedOut(now int64) (string, bool)
}

// KVEntry implementation

type kventry struct {
	validThru_ int64
	value_     data.Value
}

func NewKVEntry(value data.Value, validThru int64) KVEntry {
	entry := new(kventry)
	entry.value_ = value
	entry.validThru_ = validThru
	return entry
}

func (this *kventry) validThru() int64 {
	return this.validThru_
}

func (this *kventry) value() data.Value {
	return this.value_
}

// KV implementation

type kv struct {
	storage map[string]KVEntry
}

func NewKV() KV {
	kv := new(kv)
	kv.storage = make(map[string]KVEntry)
	return kv
}

func (this *kv) Get(k string) (KVEntry, bool) {
	kve, found := this.storage[k]
	return kve, found
}

func (this *kv) Set(key string, value data.Value, validThru int64) {
	entry := NewKVEntry(value, validThru)
	this.storage[key] = entry
}

func (this *kv) Del(key string) {
	delete(this.storage, key)
}

// TTL implementation
// quite naïve here... yep...
type ttlentry struct {
	key       string
	validThru int64
}

type ttl struct {
	keys    map[string]bool
	entries *list.List
}

func NewTTL() TTL {
	ttl := new(ttl)
	ttl.entries = list.New()
	ttl.keys = make(map[string]bool)
	return ttl
}

func (this *ttl) SetTTL(k string, deadline int64) {
	_, found := this.keys[k]
	if found {
		if deadline == -1 {
			this.keyRm(k)
			delete(this.keys, k)
		} else {
			this.keyInsAndRm(k, deadline)
		}
	} else {
		if deadline != -1 {
			this.keyIns(k, deadline)
			this.keys[k] = true
		}
	}
}

func (this *ttl) keyRm(k string) {
	for elt := this.entries.Front(); elt != nil; elt = elt.Next() {
		entry := elt.Value.(*ttlentry)
		if entry.key == k {
			this.entries.Remove(elt)
			return
		}
	}
}
func (this *ttl) keyInsAndRm(k string, deadline int64) {
	inserted := false
	removed := false

	newEntry := new(ttlentry)
	newEntry.key = k
	newEntry.validThru = deadline

	for elt := this.entries.Front(); elt != nil; elt = elt.Next() {
		entry := elt.Value.(*ttlentry)
		if entry.key == k {
			this.entries.Remove(elt)
			removed = true

			if removed && inserted {
				return
			}
		}
		if entry.validThru >= deadline {
			this.entries.InsertBefore(newEntry, elt)
			inserted = true

			if removed && inserted {
				return
			}
		}
	}

	this.entries.PushBack(newEntry)
}
func (this *ttl) keyIns(k string, deadline int64) {
	newEntry := new(ttlentry)
	newEntry.key = k
	newEntry.validThru = deadline

	for elt := this.entries.Front(); elt != nil; elt = elt.Next() {
		entry := elt.Value.(*ttlentry)

		if entry.validThru >= deadline {
			this.entries.InsertBefore(newEntry, elt)
			return
		}
	}

	this.entries.PushBack(newEntry)
}

func (this *ttl) FetchTimedOut(now int64) (string, bool) {
	if this.entries.Len() == 0 {
		return "", false
	}

	head := this.entries.Front()
	entry := head.Value.(*ttlentry)
	if entry.validThru < now {
		this.entries.Remove(head)
		return entry.key, true
	} else {
		return "", false
	}
}
