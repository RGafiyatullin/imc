
# In-memory cache. Trial task. Juno

## Requirements

```
Simple implementation of Redis-like in-memory cache

Desired features:
* Key-value storage with string, lists, dict support
* Per-key TTL
* Operations:
 ** Get
 ** Set
 ** Update
 ** Remove
 ** Keys
* Custom operations(Get i element on list, get value by key from dict, etc)
* Golang API client
* Telnet-like/HTTP-like API protocol
* Provide some tests, API spec, deployment docs without full coverage, just a few cases and some examples of telnet/http calls to the server.
* Optional features:
 ** persistence to disk/db
 ** scaling(on server-side or on client-side, up to you)
 ** auth
 ** perfomance tests
```

## Solution

The solution has been called `imcd` (in-memory cache daemon).

### Protocol

It's been decided to implement Redis-protocol for the following reasons:
* no need to invent own text-protocol;
* no need to write "Golang API client" as there is a fair amount of those (tested with [Radix](http://redis.io/clients#go));
* no need to write a custom perf-test utility (redis-benchmark does the job).

### Data model

`imcd` is a key-value storage: mapping string-keys into the values of the following types:
* strings;
* lists;
* hash-tables.

A TTL can be associated with a key. TTL precision can be changed compile-time (25ms seems to be okay).

### Storage

The keyspace is divided evenly between a ring of so called buckets. CRC32 is used as the hash-function.

All operations against any particular bucket are serialized (except for the KEYS operation).

Each bucket carries the following low-level data structures to store the keys, values and expiries:

* Go's native `map[string]KVEntry` -- associate keys with the values;
* `list.List/*[TTLEntry]*/` -- ordered list of the keys having non-infinity expiry time;
* `gtreap.Treap/*[string]*/` -- a persistent (immutable) `set` used to satisfy potentially quite a long-lasting KEYS-operation without blocking the bucket.

### Actor model

In order to deal with concurrency painlessly Actor model (goroutine-backed) has been chosen with the following hierarchy:

```
imcd
|_metrics
|_storage_sup
| \_metronome
| \_ring_mgr
|  |_bucket#0
|  |_bucket#1
|  |...
|  \_bucket#n
\_net_srv_listener
  |_acceptor#0
  | |_connection#0
  | |_...
  | \_connection#n
  |_...
  \_acceptor#n
    \...
```

### Deployment

The easiest way to deploy `imcd` is to use carefully brewed [Docker image](https://hub.docker.com/r/rgafiyatullin/imcd/) .

Nevertheless in case of necessity the `imcd` can be launched directly without much fuss. 
The following OS environment variables can be used to configure the service:
* `IMCD_GRAPHITE_ADDR_PLAINTEXT` - default `""`
* `IMCD_GRAPHITE_PREFIX` - default `""`
* `IMCD_STORAGE_RING_SIZE` - default `32`
* `IMCD_NET_BIND` - default `":6379"`
* `IMCD_NET_ACCEPTORS_COUNT` - default `1`
* `IMCD_PASSWORD` - default `""`

###### Grafana dashboard
For a dashboard like that ![](https://github.com/RGafiyatullin/imc/raw/master/doc/grafana-dashboard.png) , use the JSON file located [here](https://github.com/RGafiyatullin/imc/blob/master/doc/grafana-dashboard.json) (replace `imcd-2` with the one of your own).


### Build from source

```
rgmbp:imc rg [] $ ./go-get-deps.sh && make && file bin/imcd
go install src/github.com/rgafiyatullin/imc/server/imcd.go
bin/imcd: Mach-O 64-bit executable x86_64
rgmbp:imc rg [] $
```

### GoLang Client

The following existing client works with `imcd`:
```
go get github.com/mediocregopher/radix.v2/redis
```

### CLI client

Standard Redis client can be used to access `imcd`:
```
rgmbp:imcd rg [] $ redis-cli -h $(docker-machine ip dev) -p 16379
redis 192.168.99.100:16379> SET existent "Hello there!"
OK
redis 192.168.99.100:16379> EXPIRE existent 120
(integer) 1
redis 192.168.99.100:16379> TTL existent
(integer) 114
redis 192.168.99.100:16379> TTL non-existent
(integer) -2
redis 192.168.99.100:16379> GET non-existent
(nil)
redis 192.168.99.100:16379> GET existent
"Hello there!"
redis 192.168.99.100:16379> TTL existent
(integer) 84
redis 192.168.99.100:16379> HSET commands set "SET <key> <value>"
(integer) 0
redis 192.168.99.100:16379> HSET commands get "GET <key>"
(integer) 0
redis 192.168.99.100:16379> HSET commands del "DEL <key>"
(integer) 0
redis 192.168.99.100:16379> HSET commands "expire" "EXPIRE <key> <seconds>"
(integer) 0
redis 192.168.99.100:16379> HSET commands "pexpire" "PEXPIRE <key> <milliseconds>"
(integer) 1
redis 192.168.99.100:16379> HSET commands "ttl" "TTL <key>"
(integer) 0
redis 192.168.99.100:16379> HSET commands "pttl" "PTTL <key>"
(integer) 0
redis 192.168.99.100:16379> HSET commands "persist" "PERSIST <key>"
(integer) 1
redis 192.168.99.100:16379> GET commands
1) 1) del
   2) "DEL <key>"
2) 1) expire
   2) "EXPIRE <key> <seconds>"
3) 1) pexpire
   2) "PEXPIRE <key> <milliseconds>"
4) 1) ttl
   2) "TTL <key>"
5) 1) pttl
   2) "PTTL <key>"
6) 1) persist
   2) "PERSIST <key>"
7) 1) set
   2) "SET <key> <value>"
8) 1) get
   2) "GET <key>"
redis 192.168.99.100:16379> HKEYS commands
1) set
2) get
3) del
4) expire
5) ttl
6) pttl
7) pexpire
8) persist
redis 192.168.99.100:16379> HGETALL commands
1) "PERSIST <key>"
2) "SET <key> <value>"
3) "GET <key>"
4) "DEL <key>"
5) "EXPIRE <key> <seconds>"
6) "TTL <key>"
7) "PTTL <key>"
8) "PEXPIRE <key> <milliseconds>"
redis 192.168.99.100:16379>
```

### Performance

A standard `redis-benchmark` utility has been used.

#### `imcd` (using 4 CPUs)

```
rgmbp:rgafiyatullin rg [dev] $ redis-benchmark -c 50 -n 1000000 -d 1024 -t get,set
====== SET ======
  1000000 requests completed in 12.19 seconds
  50 parallel clients
  1024 bytes payload
  keep alive: 1

95.30% <= 1 milliseconds
99.54% <= 2 milliseconds
99.96% <= 3 milliseconds
99.98% <= 4 milliseconds
99.98% <= 5 milliseconds
99.99% <= 6 milliseconds
99.99% <= 7 milliseconds
99.99% <= 12 milliseconds
100.00% <= 39 milliseconds
100.00% <= 40 milliseconds
100.00% <= 40 milliseconds
82068.12 requests per second

====== GET ======
  1000000 requests completed in 15.18 seconds
  50 parallel clients
  1024 bytes payload
  keep alive: 1

97.91% <= 1 milliseconds
99.67% <= 2 milliseconds
99.98% <= 3 milliseconds
99.99% <= 4 milliseconds
99.99% <= 5 milliseconds
100.00% <= 6 milliseconds
100.00% <= 7 milliseconds
100.00% <= 7 milliseconds
65889.17 requests per second
```

#### `redis` (using a single CPU)

```
rgmbp:rgafiyatullin rg [dev] $ redis-benchmark -c 50 -n 1000000 -d 1024 -t get,set
====== SET ======
  1000000 requests completed in 12.41 seconds
  50 parallel clients
  1024 bytes payload
  keep alive: 1

99.76% <= 1 milliseconds
99.99% <= 2 milliseconds
100.00% <= 2 milliseconds
80573.68 requests per second

====== GET ======
  1000000 requests completed in 12.05 seconds
  50 parallel clients
  1024 bytes payload
  keep alive: 1

99.96% <= 1 milliseconds
100.00% <= 2 milliseconds
82980.67 requests per second
```

## Unmet optional requirements:

* "persistence to disk/db";

 Persistence can be done synchronously (when the bucket itself writes into the persistent storage). That may worsen latency significantly.
 Or it can be done asynchronously (when the bucket just passes a message to the goroutine(s) which in its(their) turn write into the persistent storage).
 The latter option seems more feasible yet it requires either deep-copying of all the key-value tuples passed to the persister or using immutable data structures.
 I would prefer immutable data structures but the time is up, thus we have what we have.

* "scaling(on server-side or on client-side, up to you)."

 Not done. Since it's a *cache*, the data stored in it is not overly valuable. So we do not need to do any kind of replication. 
 In this case there is no need for the communication between the nodes of `imcd` (per chance we decide we need such).
 Therefore I would insist on relying on the client-side in scaling-out.

* "automatic tests."

 Out of time. Sorry. :(

## Conclusion

If you need a kick-ass fast in-memory key-value cache with TTL, optional disk-persistence, basic authentication, support for lists, hashtables, hyperloglogs and strings as the values -- you want to use `Redis`.

Yet I insist that `imcd` isn't that bad at all: `Redis` has evolved for 7 years by now; `imcd` had just several days. :)




