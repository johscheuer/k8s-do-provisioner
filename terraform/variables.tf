variable "do_token" {}

variable "ssh_key_id" {}

variable "do_image" {
    default = "ubuntu-16-04-x64"
}

variable "user" {
    default = "anonymous"
}

variable "node_count" {
    default = 1
}

variable "region" {
  default   = "fra1"
}

variable "size" {
  default   = "8gb"
}

variable "token" {
  default   = "a7e9da.7776e834bd816af8"
}
