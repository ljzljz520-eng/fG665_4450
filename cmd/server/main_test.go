package main

import "testing"

func TestMainConfig(t *testing.T) {
	if (Config{DB: "x", Addr: ":1"}).DB != "x" {
		t.Fatal("config")
	}
}
