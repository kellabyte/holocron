package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/kellabyte/holocron/logging"
	"github.com/kellabyte/holocron/node"
	"github.com/kellabyte/holocron/store"
	"github.com/rs/zerolog"
)

func main() {
	logger := logging.CreateLogger()
	key := flag.String("key", "", "S3 key")
	secret := flag.String("secret", "", "S3 secret")
	region := flag.String("region", "", "S3 region")
	bucket := flag.String("bucket", "", "S3 bucket name")
	prefix := flag.String("prefix", "", "S3 prefix")

	flag.Parse()

	if key == nil || *key == "" {
		panic("S3 key must be provided.")
	}
	if secret == nil || *key == "" {
		panic("S3 secret must be provided.")
	}
	if region == nil || *key == "" {
		panic("S3 region must be provided.")
	}
	if bucket == nil || *key == "" {
		panic("S3 bucket must be provided.")
	}
	if prefix == nil || *key == "" {
		panic("S3 prefix must be provided.")
	}

	storageCredentials := &store.StorageCredentials{
		Key:    *key,
		Secret: *secret,
		Region: *region,
		Bucket: *bucket,
		Prefix: *prefix,
	}

	var nodes = make(map[uuid.UUID]*node.Node, 0)

	// Create node 1.
	node1, err := createNode("01924444-62e3-7152-ae06-3f10d2708fce", storageCredentials, logger)
	if err != nil {
		nodes[node1.NodeId] = node1
	}
	ctx1 := context.Background()
	node1.Start(ctx1)

	// Create node 2.
	node2, err := createNode("01924528-5f9c-7814-a442-d499b66c65f4", storageCredentials, logger)
	if err != nil {
		nodes[node1.NodeId] = node1
	}
	ctx2 := context.Background()
	node2.Start(ctx2)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
	ctx1.Done()
	ctx2.Done()
}

func createNode(id string, credentials *store.StorageCredentials, logger zerolog.Logger) (*node.Node, error) {
	storage := store.New(credentials)
	err := storage.Open()
	if err != nil {
		panic("Unable to open storage")
	}

	nodeId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	nodeLogger := logging.CreateNodeLogger(logger, nodeId)
	var node1 = node.New(nodeId, storage, nodeLogger)
	return node1, nil
}
