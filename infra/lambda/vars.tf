variable "env" {
  type = string
}

variable "folder" {
  type = string
}

variable "website_url" {
  type = string
}

variable "verification_link" {
  type = map(string)
}

variable "clone_user" {
  type = map(string)
}

variable "user" {
  type = map(string)
}

variable "user_table" {
  type = map(string)
}