# Kubernetes Digital Ocean Cluster Provisioner

This is a PoC and uses Kubeadm and the Digital Ocean API.

## Usage

Create a file `.token` with your DO access token.

```bash
$ cp cluster.yaml.example cluster.yaml

$ ./k8s-do-provisioner --help                                                                                                                                                                                                                          ±[●][master] 23:51:12
Usage of ./k8s-do-provisioner:
  -cleanup
    	removes all dropletes of the cluster
  -provision
    	provisions a cluster like specified in cluster.yaml
  -token string
    	path to the file containing the API token (default ".token")
```

### Create a cluster

```bash
./k8s-do-provisior --provision
```

## Remove a cluster

```bash
./k8s-do-provisior --provision
```

## Configure kubectl

```bash
$ export KUBECONFIG="${KUBECONFIG}:$(pwd)/admin.conf"

$ kubectl config use-context admin@kubernetes
```

# Todo

- [ ] Allow multiple clusters