package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v1"

	"github.com/digitalocean/godo"
)

type Cluster struct {
	Master  string
	Nodes   int
	Region  string
	Size    string
	Image   string
	SSHKeys []godo.DropletCreateSSHKey `yaml:"sshkeys"`
}

type Node struct {
	Name      string
	publicIP  string
	privateIP string
}

func (c *Cluster) getNodeNames() []string {
	nodes := []string{}

	for i := 0; i < c.Nodes; i++ {
		nodes = append(nodes, fmt.Sprintf("node-%d", i))
	}

	return nodes
}

func readClusterConfig() *Cluster {
	cluster := &Cluster{}

	clusterContent, err := ioutil.ReadFile("./cluster.yaml")

	if err = yaml.Unmarshal([]byte(clusterContent), cluster); err != nil {
		log.Fatalf("error: %v", err)
	}

	return cluster
}

func provision(cluster *Cluster, client *godo.Client) {
	//TODO not more than 10 droplet at a time :(
	createRequests := &godo.DropletMultiCreateRequest{
		Names:  append([]string{cluster.Master}, cluster.getNodeNames()...),
		Region: cluster.Region,
		Size:   cluster.Size,
		Image: godo.DropletCreateImage{
			Slug: cluster.Image,
		},
		SSHKeys:           cluster.SSHKeys,
		PrivateNetworking: true,
	}

	// TODO PrivateNetworking -> ? do wee need this yes for -> --address=127.0.0.1
	// probably not bad -> Ingress only at master

	newDroplets, _, err := client.Droplets.CreateMultiple(createRequests)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		os.Exit(-1)
	}

	fmt.Println("Wait for creation")
	time.Sleep(60 * time.Second)

	master := &Node{}
	nodes := []*Node{}

	for _, droplet := range newDroplets {
		drop, _, err := client.Droplets.Get(droplet.ID)
		if err != nil {
			fmt.Printf("Something bad happened: %s\n\n", err)
		}

		for drop.Status != "active" {
			//TODO backoff 5 tries
			time.Sleep(5 * time.Second)
			fmt.Println("Wait for creation...")
			drop, _, err = client.Droplets.Get(droplet.ID)
			if err != nil {
				fmt.Printf("Something bad happened: %s\n\n", err)
			}
		}

		pubIP, _ := drop.PublicIPv4()
		privIP, _ := drop.PrivateIPv4()

		if drop.Name == "master" {
			master = &Node{
				Name:      drop.Name,
				publicIP:  pubIP,
				privateIP: privIP,
			}

			continue
		}

		nodes = append(nodes, &Node{
			Name:      drop.Name,
			publicIP:  pubIP,
			privateIP: privIP,
		})
	}

	provisionRepo(append(nodes, master))
	token := provisionMaster(master)
	provisionNode(token, nodes)
	copyAdminConf(master)
}

func deprovision(cluster *Cluster, client *godo.Client) {
	droplets, _, err := client.Droplets.List(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	nodes := cluster.getNodeNames()

	var wg sync.WaitGroup
	wg.Add(len(droplets))

	for _, droplet := range droplets {
		go func(droplet godo.Droplet) {
			delete := false
			if droplet.Name == cluster.Master {
				delete = true
			}

			for _, node := range nodes {
				if delete {
					break
				}

				if droplet.Name == node {
					delete = true
				}
			}

			if delete {
				fmt.Printf("Delete: %s\n", droplet.Name)
				if _, err := client.Droplets.Delete(droplet.ID); err != nil {
					fmt.Println(err)
				}
			}
			wg.Done()
		}(droplet)
	}

	wg.Wait()

	executeLocalCommand("rm", []string{"-f", "./admin.conf"})
}
