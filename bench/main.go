package main

import (
	"fmt"
	"github.com/joyant/shardmap"
	"github.com/tidwall/lotsa"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func randKey(rnd *rand.Rand, n int) string {
	s := make([]byte, n)
	rnd.Read(s)
	for i := 0; i < n; i++ {
		s[i] = 'a' + (s[i] % 26)
	}
	return string(s)
}

func main()  {
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	N := 1000000
	K := 10

	fmt.Printf("\n")
	fmt.Printf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("\n")
	fmt.Printf("     number of cpus: %d\n", runtime.NumCPU())
	fmt.Printf("     number of keys: %d\n", N)
	fmt.Printf("            keysize: %d\n", K)
	fmt.Printf("        random seed: %d\n", seed)

	fmt.Printf("\n")

	keysm := make(map[string]interface{}, N)
	for len(keysm) < N {
		keysm[randKey(rnd, K)] = true
	}
	keys := make([]string, 0, N)
	for key := range keysm {
		keys = append(keys, key)
	}

	lotsa.Output = os.Stdout

	var mu sync.RWMutex

	println("-- sync.Map --")
	var sm sync.Map
	print("set: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		sm.Store(keys[i], i)
	})

	print("get: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		v, _ := sm.Load(keys[i])
		if v.(int) != i {
			panic("bad news")
		}
	})

	print("rng:       ")
	lotsa.Ops(100, runtime.NumCPU(), func(i, thread int) {
		sm.Range(func(key, value interface{}) bool {
			return true
		})
	})

	print("del: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		mu.Lock()
		sm.Delete(keys[i])
		mu.Unlock()
	})
	println()

	println("-- stdlib map --")
	m := make(map[string]interface{})

	print("set: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		mu.Lock()
		m[keys[i]] = i
		mu.Unlock()
	})

	print("get: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		mu.RLock()
		v := m[keys[i]]
		mu.RUnlock()
		if v.(int) != i {
			panic("bad news")
		}
	})

	print("rng:       ")
	lotsa.Ops(100, runtime.NumCPU(), func(i, thread int) {
		mu.RLock()
		for _, v := range m {
			if v == nil {
				panic("bad news")
			}
		}
		mu.RUnlock()
	})

	print("del: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		mu.Lock()
		delete(m, keys[i])
		mu.Unlock()
	})

	println()

	println("-- github.com/joyant/shardmap --")
	var com shardmap.Map

	print("set: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		com.Set(keys[i], i)
	})

	print("get: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		v, _ := com.Get(keys[i])
		if v.(int) != i {
			panic("bad news")
		}
	})

	print("rng:       ")
	lotsa.Ops(100, runtime.NumCPU(), func(i, thread int) {
		com.Range(func(key string, value interface{}) bool {
			return true
		})
	})

	print("del: ")
	lotsa.Ops(N, runtime.NumCPU(), func(i, thread int) {
		com.Delete(keys[i])
	})

	println()
}
