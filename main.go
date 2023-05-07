package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/infamous55/go-zestful/api"
	"github.com/infamous55/go-zestful/cache"
)

type portNumber uint16

func (p *portNumber) Set(value string) error {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil || result > 65535 {
		return fmt.Errorf("parse error")
	}
	*p = portNumber(result)
	return nil
}

func (p *portNumber) String() string {
	return fmt.Sprint(*p)
}

func main() {
	var capacity uint64
	var evictionPolicy cache.EvictionPolicy
	var defaultTtl cache.TimeToLive
	var secret string
	var port portNumber

	flag.Uint64Var(&capacity, "capacity", 0, "set the capacity of the cache")
	flag.Var(&evictionPolicy, "eviction-policy", "set the eviction policy of the cache (LRU or LFU)")
	flag.Var(&defaultTtl, "default-ttl", "set the default time-to-live")
	flag.StringVar(&secret, "secret", "", "set the authorization secret")
	flag.Var(&port, "port", "set the port number for the web server")

	flag.Parse()

	if capacity == 0 || evictionPolicy == "" || port == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if secret == "" {
		if os.Getenv("ZESTFUL_SECRET") != "" {
			secret = os.Getenv("ZESTFUL_SECRET")
		} else {
			fmt.Fprint(os.Stderr, "missing value for secret: initialization error\n")
			flag.Usage()
			os.Exit(2)
		}
	}

	newCache, err := cache.New(capacity, evictionPolicy, defaultTtl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: initialization error\n", err)
		os.Exit(2)
	}
	go newCache.DeleteExpired(5 * time.Minute)

	router := mux.NewRouter()
	itemsRouter := router.PathPrefix("/items").Subrouter()
	authRouter := router.PathPrefix("/auth").Subrouter()
	cacheRouter := router.PathPrefix("/cache").Subrouter()
	keyValue := randomString(32)

	api.RegisterItemsHandlers(itemsRouter)
	authMiddleware := api.GenerateAuthMiddleware([]byte(keyValue))
	cacheMiddleware := api.GenerateCacheMiddleware(newCache)
	itemsRouter.Use(cacheMiddleware)
	itemsRouter.Use(authMiddleware)

	api.RegisterAuthHandlers(authRouter, secret, []byte(keyValue))

	api.RegisterCacheHandlers(cacheRouter)

	address := fmt.Sprintf(":%v", port)
	http.ListenAndServe(address, router)
}
