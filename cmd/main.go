package main

import (
	"github.com/overnilx/udptun"
	"path/filepath"
	"sync"
)

func main() {
	paths, err := filepath.Glob("*.json")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, path := range paths {

		wg.Add(1)

		server, err := udptun.NewServer(path)
		if err != nil {
			panic(err)
		}

		go func(s *udptun.Server) {

			err := s.Run()
			if err != nil {
				panic(err)
			}

		}(server)
	}

	wg.Wait()
}