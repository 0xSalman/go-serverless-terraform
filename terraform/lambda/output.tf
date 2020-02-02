output "verification_link" {
  value = {
    name = aws_lambda_function.verification_link.function_name
    arn  = aws_lambda_function.verification_link.arn
  }
}

output "clone_user" {
  value = {
    name = aws_lambda_function.clone_user.function_name
    arn  = aws_lambda_function.clone_user.arn
  }
}
