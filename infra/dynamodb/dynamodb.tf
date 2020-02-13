resource "aws_dynamodb_table" "user" {
  name         = "user-${var.env}"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"
  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "conversation" {
  name         = "conversation-${var.env}"
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "pkey"
  range_key = "skey"
  attribute {
    name = "pkey"
    type = "S"
  }
  attribute {
    name = "skey"
    type = "N"
  }
}