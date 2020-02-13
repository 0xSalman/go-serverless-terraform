output "userfiles" {
  value = {
    arn  = aws_s3_bucket.userfiles.arn
    id   = aws_s3_bucket.userfiles.id
    name = aws_s3_bucket.userfiles.bucket
  }
}
