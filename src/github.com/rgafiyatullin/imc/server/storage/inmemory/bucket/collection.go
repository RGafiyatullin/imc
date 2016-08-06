package bucket

// The two collections defined here (KV and TTL) provide a simple interface for the Key-Value storage with per-key TTL (msec-wise).
//
// KV associates a given key with the value along with the some deadline.
// Deadline is an ephemeral point in time: it's measured in ticks provided from the outside for several reasons.
// 1. time.Now's resolution is default for the OS: usually microseconds; this is too fine for us - therefore potentially costy;
// 2. dependency on time.Now introduces a major side effect into a collection; enough said.
//
// The collections are not meant to be thread-safe: they are supposed to be accessed sequentially.
// In order to scale out -- use multiple collections as shards over the keyspace.

import (
	"bytes"
	"container/list"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
	"github.com/steveyen/gtreap"
	"math/rand"
)

const ValidThruInfinity = -1 // A special value signalling infinite TTL

type KV interface {
	// Returns nillable KVEntry if there is one associated with the given key
	Get(key string) (KVEntry, bool)

	// Creates and stores a new KVEntry associated with the given key
	Set(key string, value data.Value, validThru int64)

	Del(key string)

	Size() int

	Keys() *gtreap.Treap
}

type KVEntry interface {
	ValidThru() int64
	Value() data.Value
}

type TTL interface {
	SetTTL(k string, deadline int64)
	FetchTimedOut(now int64) (string, bool)
	Size() int
}

// KVEntry implementation

type kventry struct {
	validThru int64
	value     data.Value
}

func NewKVEntry(value data.Value, validThru int64) KVEntry {
	entry := new(kventry)
	entry.value = value
	entry.validThru = validThru
	return entry
}

func (this *kventry) ValidThru() int64 {
	return this.validThru
}

func (this *kventry) Value() data.Value {
	return this.value
}

// KV implementation: native map

func kvTreapStrOrdering(a, b interface{}) int {
	return bytes.Compare([]byte(a.(string)), []byte(b.(string)))
}

type kv struct {
	storage map[string]KVEntry
	keys    *gtreap.Treap
}

func NewKV() KV {
	kv := new(kv)
	kv.storage = make(map[string]KVEntry)
	kv.keys = gtreap.NewTreap(kvTreapStrOrdering)
	return kv
}

func (this *kv) Size() int {
	return len(this.storage)
}

func (this *kv) Keys() *gtreap.Treap {
	return this.keys
}

func (this *kv) Get(k string) (KVEntry, bool) {
	kve, found := this.storage[k]
	return kve, found
}

func (this *kv) Set(key string, value data.Value, validThru int64) {
	entry := NewKVEntry(value, validThru)
	this.storage[key] = entry
	this.keys = this.keys.Upsert(key, rand.Int())
}

func (this *kv) Del(key string) {
	delete(this.storage, key)
	this.keys = this.keys.Delete(key)
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

func (this *ttl) Size() int {
	return this.entries.Len()
}

func (this *ttl) SetTTL(k string, deadline int64) {
	_, found := this.keys[k]
	if found {
		if deadline == ValidThruInfinity {
			this.keyRm(k)
			delete(this.keys, k)
		} else {
			this.keyInsAndRm(k, deadline)
		}
	} else {
		if deadline != ValidThruInfinity {
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
