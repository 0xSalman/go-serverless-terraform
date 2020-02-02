variable "env" {
  type = string
}

variable "global" {
  type = map(string)
}

variable "verification_link_lambda" {
  type = map(string)
}

variable "clone_user_lambda" {
  type = map(string)
}
