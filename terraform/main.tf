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

  env = var.ENV
}

module "Lambda" {
  source = "./lambda"

  env         = var.ENV
  folder      = "./repository"
  website_url = var.WEBSITE_URL
  verification_link = {
    name    = "verification-link"
    version = var.VERIFICATION_LINK_VERSION
  }
  clone_user = {
    name    = "clone-user"
    version = var.CLONE_USER_VERSION
  }
  user_table = module.Dynamo.user
}

module "Cognito" {
  source = "./cognito"

  env    = var.ENV
  global = var.global

  verification_link_lambda = module.Lambda.verification_link
  clone_user_lambda        = module.Lambda.clone_user
}
