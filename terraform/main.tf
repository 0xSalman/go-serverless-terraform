terraform {
  backend "s3" {
    bucket  = "rethesis-terraform"
    key     = "terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

module "Dynamo" {
  source = "./dynamodb"

  ENV = var.ENV
}
