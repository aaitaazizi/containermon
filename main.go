package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	containerName = flag.String("container", "", "Name or ID of the container to monitor")
	outputFormat  = flag.String("output-format", "json", "Output format: json | csv")
	interval      = flag.Int("interval", 5, "collection interval (in seconds)")

	previousTotalUsage uint64
	previousTime       time.Time
)

func main() {
	flag.Parse()
	if *containerName == "" {
		fmt.Println("--container flag is required")
		flag.Usage()
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	stats := getStats(cli)
	start := time.Now()
	startUsage := stats.CPUStats.CPUUsage.TotalUsage
	previousTotalUsage = stats.CPUStats.CPUUsage.TotalUsage
	previousTime = start
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)

	containerStats := func(f *os.File) {
		for range ticker.C {
			stats = getStats(cli)
			now := time.Now()
			elapsed := now.Sub(start)
			intervalElapsed := now.Sub(previousTime)
			printStats(stats, now, elapsed, intervalElapsed, startUsage, f)
			previousTotalUsage = stats.CPUStats.CPUUsage.TotalUsage
			previousTime = now
		}
	}

	if *outputFormat == "csv" {
		f, _ := os.Create("./stats.csv")
		defer f.Close()
		fmt.Println("ts,timeElapsed,cpuTimeElapsed,percentCPUSinceStart,percentCPUThisInterval,memoryUsageKiB")
		fmt.Fprintf(f, "ts,timeElapsed,cpuTimeElapsed,percentCPUSinceStart,percentCPUThisInterval,memoryUsageKiB")
		fmt.Fprintf(f, "\n")
		containerStats(f)
	} else {
		f, _ := os.Create("./stats.json")
		defer f.Close()
		containerStats(f)
	}
}

func getStats(cli *client.Client) *types.StatsJSON {
	resp, err := cli.ContainerStats(context.TODO(), *containerName, false)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	stats := &types.StatsJSON{}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, stats)
	if err != nil {
		panic(err)
	}

	return stats
}

func printStats(stats *types.StatsJSON, now time.Time, elapsed time.Duration, intervalElapsed time.Duration, startUsage uint64, f *os.File) {
	ts := now.UTC().Format(time.RFC3339)
	timeElapsed := elapsed.Seconds()
	// cpu time in seconds
	cpuTimeElapsed := float64(stats.CPUStats.CPUUsage.TotalUsage-startUsage) / 1000000000
	percentCPUSinceStart := float64(stats.CPUStats.CPUUsage.TotalUsage-startUsage) / float64(elapsed.Nanoseconds()) * 100
	percentCPUThisInterval := float64(stats.CPUStats.CPUUsage.TotalUsage-previousTotalUsage) / float64(intervalElapsed.Nanoseconds()) * 100

	if *outputFormat == "csv" {
		// csv
		// ts,timeElapsed,cpuTimeElapsed,percentCPUSinceStart,percentCPUThisInterval,memoryUsageKiB
		fmt.Printf("%s,%.2f,%.2f,%.2f,%.2f,%.1f\n",
			ts,
			timeElapsed,
			cpuTimeElapsed,
			percentCPUSinceStart,
			percentCPUThisInterval,
			float64(stats.MemoryStats.Usage)/1024)
		fmt.Fprintf(f, "%s,%.2f,%.2f,%.2f,%.2f,%.1f\n",
			ts,
			timeElapsed,
			cpuTimeElapsed,
			percentCPUSinceStart,
			percentCPUThisInterval,
			float64(stats.MemoryStats.Usage)/1024)

	} else {
		// json
		fmt.Printf(`{"ts":"%s","timeElapsed":%.2f,"cpuTimeElapsed":%.2f,"percentCPUSinceStart":%.2f,"percentCPUThisInterval":%.2f,"memoryUsageKiB":%.1f}`,
			ts,
			timeElapsed,
			cpuTimeElapsed,
			percentCPUSinceStart,
			percentCPUThisInterval,
			float64(stats.MemoryStats.Usage)/1024)
		fmt.Fprintf(f, `{"ts":"%s","timeElapsed":%.2f,"cpuTimeElapsed":%.2f,"percentCPUSinceStart":%.2f,"percentCPUThisInterval":%.2f,"memoryUsageKiB":%.1f}`,
			ts,
			timeElapsed,
			cpuTimeElapsed,
			percentCPUSinceStart,
			percentCPUThisInterval,
			float64(stats.MemoryStats.Usage)/1024)
		fmt.Fprintf(f, "\n")
		fmt.Println()
	}
}
