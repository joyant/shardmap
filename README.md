# shardmap

Write from memory after read [shardmap](https://github.com/tidwall/shardmap).

# Performance

```cassandraql
go version go1.12.1 darwin/amd64

     number of cpus: 8
     number of keys: 1000000
            keysize: 10
        random seed: 1571653203921315000

-- sync.Map --
set: 1,000,000 ops over 8 threads in 1383ms, 723,045/sec, 1383 ns/op
get: 1,000,000 ops over 8 threads in 597ms, 1,675,285/sec, 596 ns/op
rng:       100 ops over 8 threads in 1745ms, 57/sec, 17446015 ns/op
del: 1,000,000 ops over 8 threads in 388ms, 2,574,752/sec, 388 ns/op

-- stdlib map --
set: 1,000,000 ops over 8 threads in 782ms, 1,278,189/sec, 782 ns/op
get: 1,000,000 ops over 8 threads in 55ms, 18,117,572/sec, 55 ns/op
rng:       100 ops over 8 threads in 457ms, 218/sec, 4574736 ns/op
del: 1,000,000 ops over 8 threads in 410ms, 2,440,056/sec, 409 ns/op

-- github.com/joyant/shardmap --
set: 1,000,000 ops over 8 threads in 153ms, 6,529,058/sec, 153 ns/op
get: 1,000,000 ops over 8 threads in 47ms, 21,184,919/sec, 47 ns/op
rng:       100 ops over 8 threads in 358ms, 279/sec, 3577649 ns/op
del: 1,000,000 ops over 8 threads in 66ms, 15,105,070/sec, 66 ns/op
```