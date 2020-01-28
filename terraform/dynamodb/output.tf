output "user_table" {
  value = aws_dynamodb_table.user.name
}

output "conversation_table" {
  value = aws_dynamodb_table.conversation.name
}
