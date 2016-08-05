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
	"github.com/rgafiyatullin/imc/protocol/resp/types"
)

type KV interface {
	// Returns nillable KVEntry if there is one associated with the given key
	Get(key string) (KVEntry, bool)

	// Creates and stores a new KVEntry associated with the given key
	Set(key string, value types.BasicType, validThru uint64)

	Del(key string)
}

type KVEntry interface {
	validThru() uint64
	value() types.BasicType
}

type TTL interface {
	SetTTL(k string, deadline uint64)
	FetchTimedOut(now uint64) (string, bool)
}

// KVEntry implementation

type kventry struct {
	validThru_ uint64
	value_     types.BasicType
}

func NewKVEntry(value types.BasicType, validThru uint64) KVEntry {
	entry := new(kventry)
	entry.value_ = value
	entry.validThru_ = validThru
	return entry
}

func (this *kventry) validThru() uint64 {
	return this.validThru_
}

func (this *kventry) value() types.BasicType {
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

func (this *kv) Set(key string, value types.BasicType, validThru uint64) {
	entry := NewKVEntry(value, validThru)
	this.storage[key] = entry
}

func (this *kv) Del(key string) {
	delete(this.storage, key)
}

// TTL implementation

type ttl struct{}

func NewTTL() TTL {
	ttl := new(ttl)
	return ttl
}

func (this *ttl) SetTTL(k string, deadline uint64) {
	// TODO: well, set the TTL
}

func (this *ttl) FetchTimedOut(now uint64) (string, bool) {
	// TODO: fetch a single timed out entry

	return "", false
}