package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	var (
		containerID string
		interval    time.Duration
	)
	flag.StringVar(&containerID, "c", "", "Container ID")
	flag.DurationVar(&interval, "i", 10*time.Second, "Polling interval")
	flag.Parse()

	ctx := context.Background()

	fmt.Println("unix_time,ram_usage,ram_limit")
	printStats(ctx, cli, containerID, time.Now())

	ticker := time.NewTicker(interval)

	for now := range ticker.C {
		printStats(ctx, cli, containerID, now)
	}
}

func printStats(ctx context.Context, cli *client.Client, containerID string, now time.Time) {
	resp, err := cli.ContainerStatsOneShot(ctx, containerID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var stats types.Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		panic(err)
	}
	// Compute memory stats in MB.
	var (
		limit = float64(stats.MemoryStats.Limit) / 1024 / 1024
		usage = float64(stats.MemoryStats.Usage) / 1024 / 1024
	)
	fmt.Printf("%d,%f,%f\n", now.Unix(), usage, limit)
}
