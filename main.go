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

type options struct {
	capacity       uint64
	evictionPolicy cache.EvictionPolicy
	defaultTtl     cache.TimeToLive
	secret         string
	port           portNumber
}

func parseOptions(opt *options) {
	flag.Uint64Var(&opt.capacity, "capacity", 0, "set the capacity of the cache")
	flag.Var(&opt.evictionPolicy, "eviction-policy", "set the eviction policy of the cache (LRU or LFU)")
	flag.Var(&opt.defaultTtl, "default-ttl", "set the default time-to-live")
	flag.StringVar(&opt.secret, "secret", "", "set the authorization secret")
	flag.Var(&opt.port, "port", "set the port number for the web server")

	flag.Parse()
}

func main() {
	opt := options{}
	parseOptions(&opt)

	if opt.capacity == 0 || opt.evictionPolicy == "" || opt.port == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if opt.secret == "" {
		fmt.Fprint(os.Stderr, "missing value for secret: initialization error\n")
		flag.Usage()
		os.Exit(2)
	}

	newCache, err := cache.New(opt.capacity, opt.evictionPolicy, opt.defaultTtl)
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

	api.RegisterAuthHandlers(authRouter, opt.secret, []byte(keyValue))

	api.RegisterCacheHandlers(cacheRouter)

	address := fmt.Sprintf(":%v", opt.port)
	http.ListenAndServe(address, router)
}
