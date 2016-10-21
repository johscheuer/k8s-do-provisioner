package main

import (
	"flag"
	"os"
)

var (
	cleanUp = flag.Bool("cleanup", false, "removes all dropletes of the cluster")
)

func main() {
	flag.Parse()
	cluster := readClusterConfig()
	client := createNewDOClient()

	if *cleanUp {
		deprovision(cluster, client)
		os.Exit(0)
	}

	provision(cluster, client)
}
