package main

import (
	"flag"
	"fmt"
	"os"

	cache "github.com/infamous55/go-zestful/cache"
)

func main() {
	var capacity uint64
	var evictionPolicy cache.EvictionPolicy

	flag.Uint64Var(&capacity, "capacity", 0, "set the capacity of the cache")
	flag.Var(&evictionPolicy, "eviction-policy", "set the eviction policy of the cache (LRU or LFU)")

	flag.Parse()

	if capacity == 0 || evictionPolicy == "" {
		flag.Usage()
		return
	}

	_, err := cache.New(evictionPolicy, capacity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: initialization error\n", err)
		os.Exit(1)
	}
}
