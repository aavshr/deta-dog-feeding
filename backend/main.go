package main

import (
	"log"
)

func main() {
	s, err := NewServer(&Config{
		Address: ":8080",
		CodeStoreName: "mermaidCodeStore",
		CodeStoreKey: nil,
	})
	if err != nil {
		log.Fatalf("failed to get new server: %s\n", err)
	}
	log.Fatal(s.Start())
}