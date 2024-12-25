output "website_url" {
  description = "The URL of the CCTray XML feed"
  value       = "http://${aws_s3_bucket_website_configuration.ccxml.website_endpoint}/${var.key}"
}

output "lambda_function_arn" {
  description = "The ARN of the Lambda function"
  value       = aws_lambda_function.ccxml.arn
}

output "lambda_function_name" {
  description = "The name of the Lambda function"
  value       = aws_lambda_function.ccxml.function_name
}

output "s3_bucket_name" {
  description = "The name of the S3 bucket"
  value       = aws_s3_bucket.ccxml.id
}
