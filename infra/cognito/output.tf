output "userpool" {
  value = {
    arn = aws_cognito_user_pool.main.arn,
    id  = aws_cognito_user_pool.main.id
  }
}

