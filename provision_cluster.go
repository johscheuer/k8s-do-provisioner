package main

import (
	"fmt"
	"strings"
	"sync"
)

func provisionRepo(nodes []*Node) {
	//TODO install latest docker from docker

	var wg sync.WaitGroup
	wg.Add(len(nodes))

	for _, node := range nodes {
		go func(node *Node) {
			fmt.Printf("Provision Node %s\n", node.Name)
			// TODO error checking

			executeSCPCommand("./files/kubernetes.repo", fmt.Sprintf("root@%s:/etc/yum.repos.d/kubernetes.repo", node.publicIP))
			executeSSHCommand(node.publicIP, "setenforce 0")
			executeSSHCommand(node.publicIP, "yum install -y docker kubelet kubeadm kubectl kubernetes-cni")
			executeSSHCommand(node.publicIP, "systemctl enable docker && systemctl start docker")
			executeSSHCommand(node.publicIP, "systemctl enable kubelet && systemctl start kubelet")
			wg.Done()
		}(node)
	}

	wg.Wait()
}

func provisionMaster(master *Node) string {
	fmt.Println("Initialize Kubernetes Master")
	// TODO error checking
	out, _ := executeSSHCommand(master.publicIP, fmt.Sprintf("kubeadm init --api-advertise-addresses=%s", master.publicIP))

	token := ""
	for _, s := range strings.Split(string(out), "\n") {
		if strings.Contains(s, "kubeadm join") {
			token = s
		}
	}

	fmt.Println(token)
	return token
}

func provisionNode(tokenString string, nodes []*Node) {
	// TODO error checking
	var wg sync.WaitGroup
	wg.Add(len(nodes))

	for _, node := range nodes {
		go func(node *Node) {
			fmt.Printf("Provision Kubernetes Node %s\n", node.Name)
			executeSSHCommand(node.publicIP, fmt.Sprintf("%s", tokenString))
			wg.Done()
		}(node)
	}

	wg.Wait()
}

func copyAdminConf(master *Node) {
	executeSCPCommand(fmt.Sprintf("root@%s:/etc/kubernetes/admin.conf", master.publicIP), "./admin.conf")
	out, _ := executeKubernetesCommand([]string{"get", "nodes"})
	fmt.Println(string(out))
	out, _ = executeKubernetesCommand([]string{"create", "-f", "https://raw.githubusercontent.com/tigera/canal/master/k8s-install/kubeadm/canal.yaml"})
	fmt.Println(string(out))
	out, _ = executeKubernetesCommand([]string{"create", "-f", "https://rawgit.com/kubernetes/dashboard/master/src/deploy/kubernetes-dashboard.yaml"})
	fmt.Println(string(out))

	fmt.Println("Use \"kubectl --kubeconfig admin.conf\" to interact with your new cluster")
}
