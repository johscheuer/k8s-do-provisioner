# Kubernetes Digital Ocean Cluster Provisioner

This is a PoC and uses Kubeadm and the Digital Ocean API.

## Usage

Create a file `.token` with your DO access token.

```bash
$ cp cluster.yaml.example cluster.yaml

$ ./k8s-do-provisioner
```

## Configure kubectl

```bash
$ export KUBECONFIG="${KUBECONFIG}:$(pwd)/admin.conf"
$ kubectl config use-context admin@kubernetes
```
