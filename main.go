package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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
		os.Exit(2)
	}

	newCache, err := cache.New(evictionPolicy, capacity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: initialization error\n", err)
		os.Exit(2)
	}

	injectCache := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "cache", newCache)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}

	http.Handle("/test", injectCache(http.HandlerFunc(testHandler)))
	http.ListenAndServe(":8080", nil)
}

func getCache(ctx context.Context) cache.Cache {
	if cache, ok := ctx.Value("cache").(cache.Cache); ok {
		return cache
	}
	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	cache.Set("hello", "world")
	value, err := cache.Get("hello")
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, value)
}
