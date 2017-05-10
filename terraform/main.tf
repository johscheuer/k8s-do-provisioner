# Configure the DigitalOcean Provider
provider "digitalocean" {
    token               = "${var.do_token}"
}

data "template_file" "master_userdata" {
    template = "${file("terraform/templates/master.tpl")}"
    vars {
        region          = "${var.region}"
        token           = "${var.token}"
    }
}

data "template_file" "node_userdata" {
    template = "${file("terraform/templates/node.tpl")}"
    vars {
        region          = "${var.region}"
        token           = "${var.token}"
        master          = "${digitalocean_droplet.master.ipv4_address_private}"
    }
}

resource "digitalocean_droplet" "master" {
    image               = "${var.do_image}"
    name                = "${var.user}-master"
    private_networking  = true
    region              = "${var.region}"
    size                = "${var.size}"
    ssh_keys            = ["${var.ssh_key_id}"]
    tags                = ["${digitalocean_tag.kubernetes.id}", "${digitalocean_tag.user.id}"]
    user_data           = "${data.template_file.master_userdata.rendered}"
}

resource "digitalocean_droplet" "node" {
    image               = "${var.do_image}"
    count               = "${var.node_count}"
    name                = "${var.user}-node-${count.index}"
    private_networking  = true
    region              = "${var.region}"
    size                = "${var.size}"
    ssh_keys            = ["${var.ssh_key_id}"]
    tags                = ["${digitalocean_tag.kubernetes.id}", "${digitalocean_tag.user.id}"]
    user_data           = "${data.template_file.node_userdata.rendered}"
    depends_on          = ["digitalocean_droplet.master"]
}

resource "digitalocean_tag" "kubernetes" {
    name                = "kubernetes"
}

resource "digitalocean_tag" "user" {
    name                = "${var.user}"
}

output "master" {
    value               = "${digitalocean_droplet.master.ipv4_address}"
}

output "master_private" {
    value               = "${digitalocean_droplet.master.ipv4_address_private}"
}

output "nodes" {
    value = [
      "${digitalocean_droplet.node.*.ipv4_address}",
    ]
}

output "private_ips" {
    value = [
      "${digitalocean_droplet.node.*.ipv4_address_private}",
    ]
}
