package main

import (
	"fmt"
	gee_cache_jz "gee-cache-jz"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	gee_cache_jz.NewGroup("scores", 2<<10, gee_cache_jz.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}

			return nil, fmt.Errorf("%s not exist in DB", key)
		}))

	addr := "localhost:9999"
	peers := gee_cache_jz.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
