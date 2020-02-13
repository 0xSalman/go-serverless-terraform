# TODO add group: https://www.terraform.io/docs/providers/aws/r/cognito_user_group.html

resource "aws_cognito_user_pool" "main" {
  name = "main-${var.env}"

  auto_verified_attributes = ["email"]
  username_attributes      = ["email"]

  schema {
    attribute_data_type = "String"
    name                = "email"
    required            = true
    mutable             = true

    string_attribute_constraints {
      min_length = 5
      max_length = 300
    }
  }

  schema {
    attribute_data_type = "String"
    # cognito does not allow `required=true` for custom attributes
    # so use nickname for group to enforce required
    name     = "nickname"
    required = true
    mutable  = true
    string_attribute_constraints {
      min_length = 7
      max_length = 20
    }
  }

  email_configuration {
    source_arn            = var.global["ses_email_arn"]
    email_sending_account = "DEVELOPER"
  }

  password_policy {
    minimum_length    = 8
    require_uppercase = true
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
  }

  lambda_config {
    custom_message    = var.verification_link_lambda["arn"]
    post_confirmation = var.clone_user_lambda["arn"]
  }

  tags = {
    Stage = var.env
  }
}

resource "aws_lambda_permission" "invoke_verification_link" {
  statement_id  = "AllowExecutionFromCognito"
  action        = "lambda:InvokeFunction"
  function_name = var.verification_link_lambda["name"]
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = "arn:aws:cognito-idp:${var.global["region"]}:${var.global["account"]}:userpool/${aws_cognito_user_pool.main.id}"
}

resource "aws_lambda_permission" "invoke_user_clone" {
  statement_id  = "AllowExecutionFromCognito"
  action        = "lambda:InvokeFunction"
  function_name = var.clone_user_lambda["name"]
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = "arn:aws:cognito-idp:${var.global["region"]}:${var.global["account"]}:userpool/${aws_cognito_user_pool.main.id}"
}

resource "aws_cognito_user_pool_client" "web_client" {
  name         = "web-client"
  user_pool_id = aws_cognito_user_pool.main.id
}

resource "aws_cognito_user_pool_domain" "main" {
  domain       = "${var.env}-${var.global["app_name"]}"
  user_pool_id = aws_cognito_user_pool.main.id
}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name = "main ${var.env}"

  allow_unauthenticated_identities = false

  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.web_client.id
    provider_name           = "cognito-idp.${var.global["region"]}.amazonaws.com/${aws_cognito_user_pool.main.id}"
    server_side_token_check = false
  }
}

resource "aws_iam_role" "main_authenticated" {
  name = "authenticated-${var.env}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.main.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "authenticated"
        }
      }
    }
  ]
}
EOF

  tags = {
    Stage   = var.env
    Service = "cognito"
  }
}

resource "aws_iam_policy" "userfiles_full_access" {
  name        = "userfiles-${var.env}-full-access"
  description = "Give read and write access to cognito users to userfiles S3 bucket"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::${var.userfiles_bucket_name}/protected/$${cognito-identity.amazonaws.com:sub}/*",
        "arn:aws:s3:::${var.userfiles_bucket_name}/private/$${cognito-identity.amazonaws.com:sub}/*"
      ],
      "Effect": "Allow"
    },
    {
      "Condition": {
        "StringLike": {
          "s3:prefix": [
            "protected/",
            "protected/*",
            "private/$${cognito-identity.amazonaws.com:sub}/",
            "private/$${cognito-identity.amazonaws.com:sub}/*"
          ]
        }
      },
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::${var.userfiles_bucket_name}"
      ],
      "Effect": "Allow"
    }
  ]
}
  EOF
}

resource "aws_iam_role_policy_attachment" "userfiles_full_access" {
  policy_arn = aws_iam_policy.userfiles_full_access.arn
  role       = aws_iam_role.main_authenticated.name
}

resource "aws_iam_policy" "api_execute" {
  name        = "api-${var.env}-execute"
  description = "Give APi gateway execute access to cognito users"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "execute-api:Invoke"
      ],
      "Resource": [
        "arn:aws:execute-api:${var.global["region"]}:${var.global["account"]}:*/*"
      ],
      "Effect": "Allow"
    }
  ]
}
  EOF
}

resource "aws_iam_role_policy_attachment" "api_execute" {
  policy_arn = aws_iam_policy.api_execute.arn
  role       = aws_iam_role.main_authenticated.name
}

resource "aws_iam_role" "main_unauthenticated" {
  name = "unauthenticated-${var.env}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.main.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "unauthenticated"
        }
      }
    }
  ]
}
EOF

  tags = {
    Stage   = var.env
    Service = "cognito"
  }
}

resource "aws_cognito_identity_pool_roles_attachment" "main" {
  identity_pool_id = aws_cognito_identity_pool.main.id

  roles = {
    authenticated   = aws_iam_role.main_authenticated.arn
    unauthenticated = aws_iam_role.main_unauthenticated.arn
  }
}
