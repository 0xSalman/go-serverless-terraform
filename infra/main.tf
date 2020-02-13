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

module "S3" {
  source = "./s3"

  env    = var.ENV
  global = var.global
}

module "Lambda" {
  source = "./lambda"

  env         = var.ENV
  folder      = "./repository"
  website_url = var.WEBSITE_URL
  user_table  = module.Dynamo.user

  verification_link = {
    name    = "verification-link"
    version = var.VERIFICATION_LINK_VERSION
  }

  clone_user = {
    name    = "clone-user"
    version = var.CLONE_USER_VERSION
  }

  user = {
    name    = "user"
    version = var.USER_VERSION
  }
}

module "Cognito" {
  source = "./cognito"

  env    = var.ENV
  global = var.global

  verification_link_lambda = module.Lambda.verification_link
  clone_user_lambda        = module.Lambda.clone_user
  userfiles_bucket_name    = module.S3.userfiles["name"]
}

module "ApiGateway" {
  source = "./api-gateway"

  env         = var.ENV
  global      = var.global
  user_lambda = module.Lambda.user
}
