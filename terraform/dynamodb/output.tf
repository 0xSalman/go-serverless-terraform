output "user" {
  value = {
    name = aws_dynamodb_table.user.name
    id   = aws_dynamodb_table.user.id
    arn  = aws_dynamodb_table.user.arn
  }
}

output "conversation" {
  value = {
    name = aws_dynamodb_table.conversation.name
    id   = aws_dynamodb_table.conversation.id
    arn  = aws_dynamodb_table.conversation.arn
  }
}
