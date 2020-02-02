resource "aws_iam_policy" "log" {
  name        = "log-full-access"
  description = "Give full access to CloudWatch logs"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:*"
      ],
      "Resource": "*"
    }
  ]
}
  EOF
}

resource "aws_iam_policy" "user_table_write" {
  name        = "user-table-${var.env}-write"
  description = "Give write access to user table"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:BatchWrite*",
        "dynamodb:PutItem",
        "dynamodb:Update*"
      ],
      "Resource": "${var.user_table["arn"]}"
    }
  ]
}
  EOF
}

resource "aws_iam_role" "verification_link_role" {
  name               = "${var.verification_link["name"]}-${var.env}-executor"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
  EOF
}

resource "aws_iam_role_policy_attachment" "verification_link_log" {
  policy_arn = aws_iam_policy.log.arn
  role       = aws_iam_role.verification_link_role.name
}

resource "aws_lambda_function" "verification_link" {
  publish       = true
  function_name = "${var.verification_link["name"]}-${var.env}"
  role          = aws_iam_role.verification_link_role.arn
  handler       = var.verification_link["name"]
  runtime       = "go1.x"
  filename      = "${var.folder}/${var.verification_link["name"]}-${var.verification_link["version"]}.zip"

  environment {
    variables = {
      verification_url = "${var.website_url}/verify-email"
    }
  }
}

resource "aws_iam_role" "clone_user_role" {
  name               = "${var.clone_user["name"]}-${var.env}-executor"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
  EOF
}

resource "aws_iam_role_policy_attachment" "clone_user_log" {
  policy_arn = aws_iam_policy.log.arn
  role       = aws_iam_role.clone_user_role.name
}

resource "aws_iam_role_policy_attachment" "clone_user_user_table" {
  policy_arn = aws_iam_policy.user_table_write.arn
  role       = aws_iam_role.clone_user_role.name
}

resource "aws_lambda_function" "clone_user" {
  publish       = true
  function_name = "${var.clone_user["name"]}-${var.env}"
  role          = aws_iam_role.clone_user_role.arn
  handler       = var.clone_user["name"]
  runtime       = "go1.x"
  filename      = "${var.folder}/${var.clone_user["name"]}-${var.clone_user["version"]}.zip"

  environment {
    variables = {
      user_table = var.user_table["name"]
    }
  }
}