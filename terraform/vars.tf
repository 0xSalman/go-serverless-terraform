variable "ENV" {
  type = string
}

variable "global" {
  type = map(string)
  default = {
    aws_account = "064097596383"
    region      = "us-east-1"
    company     = "rethesis"

    repository_folder = "./repository"
  }
}
