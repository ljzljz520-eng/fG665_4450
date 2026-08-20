package main

import (
	"flag"
	"fmt"
	"instrumentarchive/api"
	"instrumentarchive/service"
	"instrumentarchive/store"
	"log"
	"net/http"
)

type Config struct {
	DB   string
	Addr string
}

func ParseConfig() Config {
	db := flag.String("db", "instrument.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	return Config{DB: *db, Addr: *addr}
}
func main() {
	cfg := ParseConfig()
	s, e := store.Open(cfg.DB)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	svc := service.New(s, service.StaticClock{"demo"})
	fmt.Println("instrument archive listening", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, api.New(svc).Routes()))
}
