package main

import (
	"fmt"
	"os/exec"
)

func executeSSHCommand(node, cmdArg string) ([]byte, error) {
	return exec.Command("/bin/sh", "-c", fmt.Sprintf("ssh -oStrictHostKeyChecking=no root@%s '%s'", node, cmdArg)).CombinedOutput()
}

func executeSCPCommand(source, dest string) ([]byte, error) {
	return exec.Command("scp", "-oStrictHostKeyChecking=no", source, dest).CombinedOutput()
}

func executeLocalCommand(command string, args []string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}

func executeKubernetesCommand(args []string) ([]byte, error) {
	return exec.Command("kubectl", append([]string{"--kubeconfig", "./admin.conf"}, args...)...).CombinedOutput()
}
