variable "env" {
  type = string
}

variable "global" {
  type = map(string)
}

variable "user" {
  type = any
}

variable "user_depends_on" {
  type    = any
  default = null
}