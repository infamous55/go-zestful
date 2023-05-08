package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/infamous55/go-zestful/api"
	"github.com/infamous55/go-zestful/cache"
)

type timeToLive struct {
	value time.Duration
	isSet bool
}

func (ttl *timeToLive) Set(value string) error {
	result, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse error")
	}
	ttl.value = result
	ttl.isSet = true
	return nil
}

func (ttl *timeToLive) String() string {
	return fmt.Sprint(ttl.value)
}

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
	defaultTtl     timeToLive
	secret         string
	port           portNumber
}

func parseOptions() options {
	opt := options{}

	flag.Uint64Var(&opt.capacity, "capacity", 0, "set the capacity of the cache")
	flag.Var(&opt.evictionPolicy, "eviction-policy", "set the eviction policy of the cache (LRU or LFU)")
	flag.Var(&opt.defaultTtl, "default-ttl", "set the default time-to-live")
	flag.StringVar(&opt.secret, "secret", "", "set the authorization secret")
	flag.Var(&opt.port, "port", "set the port number for the web server")

	flag.Parse()

	envCapacity := os.Getenv("ZESTFUL_CAPACITY")
	if opt.capacity == 0 && envCapacity != "" {
		capacity, err := strconv.ParseUint(envCapacity, 10, 64)
		if err == nil {
			opt.capacity = capacity
		}
	}

	envEvictionPolicy := os.Getenv("ZESTFUL_EVICTION_POLICY")
	if opt.evictionPolicy == "" && (envEvictionPolicy == "LRU" || envEvictionPolicy == "LFU") {
		opt.evictionPolicy = cache.EvictionPolicy(envEvictionPolicy)
	}

	envDefaultTtl := os.Getenv("ZESTFUL_DEFAULT_TTL")
	if !opt.defaultTtl.isSet && envDefaultTtl != "" {
		defaultTtl, err := time.ParseDuration(envDefaultTtl)
		if err == nil {
			opt.defaultTtl.Set(fmt.Sprint(defaultTtl))
		}
	}

	envSecret := os.Getenv("ZESTFUL_SECRET")
	if opt.secret == "" && envSecret != "" {
		opt.secret = envSecret
	}

	envPort := os.Getenv("ZESTFUL_PORT")
	if opt.port == 0 && envPort != "" {
		port, err := strconv.ParseUint(envPort, 10, 64)
		if err == nil && port < 65535 {
			opt.port = portNumber(port)
		}
	}

	if opt.capacity == 0 || opt.evictionPolicy == "" || opt.port == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if opt.secret == "" {
		fmt.Fprint(os.Stderr, "missing value for secret: parse error\n")
		flag.Usage()
		os.Exit(2)
	}
	if !opt.defaultTtl.isSet {
		fmt.Fprint(os.Stderr, "missing value for default-ttl: parse error\n")
		flag.Usage()
		os.Exit(2)
	}

	return opt
}

func main() {
	opt := parseOptions()

	newCache, err := cache.New(opt.capacity, opt.evictionPolicy, opt.defaultTtl.value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: initialization error\n", err)
		os.Exit(2)
	}
	go newCache.DeleteExpired(5 * time.Minute)

	logger := log.New(os.Stdout, "", log.Default().Flags())
	router := mux.NewRouter()
	loggingMiddleware := api.GenerateLoggingMiddleware(logger)
	router.Use(loggingMiddleware)

	itemsRouter := router.PathPrefix("/items").Subrouter()
	authRouter := router.PathPrefix("/auth").Subrouter()
	cacheRouter := router.PathPrefix("/cache").Subrouter()
	keyValue := randomString(32)

	api.RegisterItemsHandlers(itemsRouter)
	authMiddleware := api.GenerateAuthMiddleware([]byte(keyValue))
	cacheMiddleware := api.GenerateCacheMiddleware(newCache)
	itemsRouter.Use(authMiddleware)
	itemsRouter.Use(cacheMiddleware)

	api.RegisterAuthHandlers(authRouter, opt.secret, []byte(keyValue))

	api.RegisterCacheHandlers(cacheRouter)
	cacheRouter.Use(authMiddleware)
	cacheRouter.Use(cacheMiddleware)

	address := fmt.Sprintf(":%v", opt.port)
	fmt.Printf("started on port %v\n", opt.port)
	http.ListenAndServe(address, router)
}
