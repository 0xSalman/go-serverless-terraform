variable "ENV" {
  type = string
}

variable "WEBSITE_URL" {
  type = string
}

variable "VERIFICATION_LINK_VERSION" {
  type = string
}

variable "CLONE_USER_VERSION" {
  type = string
}

variable "global" {
  type = map(string)
  default = {
    account  = "064097596383"
    region   = "us-east-1"
    app_name = "rethesis"

    ses_email_arn = "arn:aws:ses:us-east-1:064097596383:identity/noreply@rethesis.com"
  }
}
