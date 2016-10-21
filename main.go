package main

import "flag"

var (
	cleanUp   = flag.Bool("cleanup", false, "removes all dropletes of the cluster")
	provision = flag.Bool("provision", false, "provisions a cluster like specified in cluster.yaml")
	tokenFile = flag.String("token", ".token", "path to the file containing the API token")
)

func main() {
	flag.Parse()
	cluster := readClusterConfig()
	client := createNewDOClient(*tokenFile)

	if *cleanUp {
		deprovisionCluster(cluster, client)
	}

	if *provision {
		provisionCluster(cluster, client)
	}
}
