# Kubernetes Digital Ocean Cluster Provisioner

This is a PoC that uses `terraform`, `DigitalOcean` and `kubeadm`

## Prerequisites

- [jq](https://stedolan.github.io/jq/)
- [doctl](https://github.com/digitalocean/doctl)

## Usage

### Setup

If you installed `doctl` and initialized it, your `access-token` is available in `$HOME/.config/doctl/config.yaml`.

Be sure that your SSH key is imported in DigitalOcean:

```
doctl compute ssh-key import my-key --public-key-file $HOME/.ssh/id_rsa.pub
```

Initalize all variables used by the terraform setup:

```bash
echo "do_token=\"$(grep "access-token" $HOME/.config/doctl/config.yaml | sed 's/access-token: //g')\"" > terraform.tfvars
echo "ssh_key_id=\"$(doctl compute ssh-key get $(ssh-keygen -E md5 -lf "$HOME/.ssh/id_rsa.pub" | awk '{ print $2 }' | sed -e "s/^MD5://") -o json | jq '.[-1].id')\"" >> terraform.tfvars
echo "user=\"$(logname)\"" >> terraform.tfvars
```

Additional variables:

```bash
echo "node_count=10" >> terraform.tfvars
echo "size=\"18gb\"" >> terraform.tfvars
echo "token=\"a7e9da.7776e834bd816af8\"" >> terraform.tfvars
```

### Cluster creation

First check that everything works:

```bash
terraform plan -var-file="terraform.tfvars" terraform
```

If the plan command exists successful we can run apply:

```bash
terraform apply -var-file="terraform.tfvars" terraform
```

If you want to destroy the cluster simply run:

```bash
terraform destoy -var-file="terraform.tfvars" terraform
```

## Next steps

To interact with the cluster we need to fetch the config:

```bash
scp root@$(terraform output -json | jq -r '.master.value'):/etc/kubernetes/admin.conf .
export KUBECONFIG=$(pwd)/admin.conf
kubectl get nodes
```

Now we need to install a [pod network](https://kubernetes.io/docs/concepts/cluster-administration/addons/) (example for Calico):

```bash
kubectl apply -f http://docs.projectcalico.org/v2.1/getting-started/kubernetes/installation/hosted/kubeadm/1.6/calico.yaml
```

And also install the Kubernetes dashboard:

```bash
kubectl apply -f  https://git.io/kube-dashboard
```

If you want to access the dashboard:

```bash
kubectl proxy > /dev/null &
# On OSX
open http://localhost:8001/ui
```

# Todo

- [ ] Allow multiple clusters
- [ ] Kubefed
