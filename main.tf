variable "bucket" {
  description = "The bucket that will contain the feed"
}

variable "key" {
  description = "The key within the bucket that will contain the feed"
  default     = "cc.xml"
}

resource "aws_s3_bucket" "ccxml" {
  bucket        = var.bucket
  force_destroy = true
}

resource "aws_s3_bucket_website_configuration" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id

  index_document {
    suffix = "cc.xml"
  }
}

resource "aws_s3_bucket_versioning" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id
  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_ownership_controls" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "ccxml" {
  depends_on = [aws_s3_bucket_ownership_controls.ccxml]
  bucket     = aws_s3_bucket.ccxml.id
  acl        = "public-read"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

output "website" {
  value = "http://${aws_s3_bucket_website_configuration.ccxml.website_endpoint}/${var.key}"
}


resource "null_resource" "cctest" {
  provisioner "local-exec" {
    working_dir = path.module
    command     = <<EOF
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o bootstrap && \
      chmod +x bootstrap
    EOF
  }

  triggers = {
    force_run = uuid()
  }
}

data "archive_file" "lambda_zip" {
  depends_on  = [null_resource.cctest]
  type        = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "${path.module}/bootstrap.zip"
}

resource "aws_lambda_function" "ccxml" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = "ccxml"
  handler          = "bootstrap"
  description      = "Handler that responds to CodePipeline events by updating a CCTray XML feed"
  memory_size      = 256
  timeout          = 30
  runtime          = "provided.al2"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  role             = aws_iam_role.ccxml.arn

  environment {
    variables = {
      BUCKET = var.bucket
      KEY    = var.key
    }
  }
}

data "aws_iam_policy_document" "ccxml_assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_role" "ccxml" {
  name = "ccxml"

  assume_role_policy = data.aws_iam_policy_document.ccxml_assume_role_policy.json
}

data "aws_iam_policy_document" "ccxml_role_policy" {
  statement {
    effect = "Allow"

    actions = [
      "codepipeline:ListPipelines",
      "codepipeline:GetPipelineState",
    ]

    resources = [
      "*",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:PutObject",
      "s3:PutObjectAcl",
    ]

    resources = [
      "arn:aws:s3:::${var.bucket}/${var.key}",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "logs:*",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy" "ccxml" {
  role = aws_iam_role.ccxml.id

  policy = data.aws_iam_policy_document.ccxml_role_policy.json
}

locals {
  event_pattern = {
    source      = ["aws.codepipeline"]
    detail-type = ["CodePipeline Stage Execution State Change"]
  }

  event_pattern_json = jsonencode(local.event_pattern)
}

resource "aws_cloudwatch_event_rule" "ccxml" {
  name           = "ccxml"
  description    = "Rule that matches CodePipeline State Execution State Changes"
  event_pattern  = local.event_pattern_json
  event_bus_name = "default"
}

resource "aws_cloudwatch_event_target" "ccxml" {
  target_id = "ccxml"
  rule      = aws_cloudwatch_event_rule.ccxml.name
  arn       = aws_lambda_function.ccxml.arn
}

resource "aws_lambda_permission" "ccxml" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ccxml.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.ccxml.arn
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.31.0"
    }
  }
}
