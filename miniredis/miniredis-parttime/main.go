package main

import (
	"log"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
)

const maxReq = 10

func main() {
	s := miniredis.NewMiniRedis()

	// Optionally set some keys your code expects:
	s.Set("foo", "bar")
	s.HSet("some", "other", "key")

	// NOTE: enable to be accessed from other containers
	if err := s.StartAddr("0.0.0.0:6379"); err != nil {
		panic(err)
	}
	defer s.Close()

	log.Printf("miniredis serves on %s\n", s.Addr())

	done := make(chan struct{})
	defer close(done)

	s.Server().SetPreHook(mergeHooks(
		mockInfoHook,
		limitRequestHook(done),
	))

	for {
		select {
		case <-time.After(999999 * time.Hour):
			return
		case <-done:
			return
		}
	}
}

func mergeHooks(hooks ...server.Hook) server.Hook {
	return func(c *server.Peer, cmd string, args ...string) bool {
		for _, h := range hooks {
			// NOTE: return value of h
			// 		 true:  hook already worked and no need to run command
			//       false: next hook should be applied
			if h(c, cmd, args...) {
				return true
			}
		}

		return false
	}
}

func limitRequestHook(done chan struct{}) server.Hook {
	nReq := 0

	return func(c *server.Peer, cmd string, args ...string) bool {
		nReq++

		if nReq > maxReq {
			c.WriteError("MiniRedis went home.")
			log.Println("Bye!")
			done <- struct{}{}
			return true
		}

		log.Printf("Request: %d/%d\n", nReq, maxReq)
		return false
	}
}

func mockInfoHook(c *server.Peer, cmd string, args ...string) bool {
	// NOTE: mock `INFO replication` required for Dapr sidecar initialization
	// 		 (miniredis does not handle `INFO` command)
	if cmd != "INFO" || len(args) < 1 {
		return false
	}

	if args[0] != "replication" {
		return false
	}

	mockLines := []string{
		"# Replication",
		"role:master",
		"connected_slaves:0",
		"master_failover_state:no-failover",
		"master_replid:b9bca6c53f5f6e52047e05566897add8b3f3c662",
		"master_replid2:0000000000000000000000000000000000000000",
		"master_repl_offset:0",
		"second_repl_offset:-1",
		"repl_backlog_active:0",
		"repl_backlog_size:1048576",
		"repl_backlog_first_byte_offset:0",
		"repl_backlog_histlen:0",
		"",
	}

	c.WriteLen(len(mockLines))
	for _, l := range mockLines {
		c.WriteBulk(l)
	}
	c.Flush()

	return false
}
