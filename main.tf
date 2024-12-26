terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0"
    }
  }
}

locals {
  lambda_src_path = "${path.module}/aws-codepipeline-ccxml"
  binary_path     = "${path.module}/dist"
  event_pattern = {
    source      = ["aws.codepipeline"]
    detail-type = ["CodePipeline Stage Execution State Change"]
  }
}

resource "aws_s3_bucket" "ccxml" {
  bucket = var.bucket
  tags   = var.tags
}

resource "aws_s3_bucket_public_access_block" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "ccxml" {
  depends_on = [
    aws_s3_bucket_public_access_block.ccxml,
    aws_s3_bucket_ownership_controls.ccxml,
  ]

  bucket = aws_s3_bucket.ccxml.id
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "ccxml" {
  bucket = aws_s3_bucket.ccxml.id

  index_document {
    suffix = var.key
  }
}

resource "null_resource" "cctest" {
  triggers = {
    source_code_hash = filebase64sha256("${local.lambda_src_path}/main.go")
    uuid             = uuid()
  }

  provisioner "local-exec" {
    command = "cd ${local.lambda_src_path} && mkdir -p dist && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/bootstrap"
  }
}

data "archive_file" "lambda_zip" {
  depends_on  = [null_resource.cctest]
  type        = "zip"
  source_file = "${local.lambda_src_path}/dist/bootstrap"
  output_path = "${path.module}/package.zip"
}

resource "aws_lambda_function" "ccxml" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.function_name
  handler          = "bootstrap"
  description      = "Handler that responds to CodePipeline events by updating a CCTray XML feed"
  memory_size      = var.memory_size
  timeout          = var.timeout
  runtime          = "provided.al2"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  role            = aws_iam_role.ccxml.arn
  tags            = var.tags

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
  name = var.function_name
  tags = var.tags

  assume_role_policy = data.aws_iam_policy_document.ccxml_assume_role_policy.json
}

data "aws_iam_policy_document" "ccxml_role_policy" {
  statement {
    effect = "Allow"
    actions = [
      "codepipeline:ListPipelines",
      "codepipeline:GetPipelineState",
    ]
    resources = ["*"]
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
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:*:*:log-group:/aws/lambda/${var.function_name}:*",
    ]
  }
}

resource "aws_iam_role_policy" "ccxml" {
  role   = aws_iam_role.ccxml.id
  policy = data.aws_iam_policy_document.ccxml_role_policy.json
}

resource "aws_cloudwatch_event_rule" "ccxml" {
  name        = var.function_name
  description = "Rule that matches CodePipeline State Execution State Changes"
  tags        = var.tags

  event_pattern = jsonencode(local.event_pattern)
}

resource "aws_cloudwatch_event_target" "ccxml" {
  target_id = var.function_name
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

output "website" {
  value = "http://${aws_s3_bucket_website_configuration.ccxml.website_endpoint}/${var.key}"
}
