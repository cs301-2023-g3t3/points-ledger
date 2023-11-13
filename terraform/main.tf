provider "aws" {
  region = "ap-southeast-1"
}

terraform {
  backend "s3" {
    bucket         = "itsa-terraform-states"
    key            = "points-ledger-terraform.tfstate"
    region         = "ap-southeast-1"
    dynamodb_table = "points-ledger-terraform-state-lock"
  }
}

resource "aws_iam_policy" "lambda_rds_permissions" {
  name = "lambda-rds-permissions"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Statement1",
            "Effect": "Allow",
            "Action": [
                "rds-data:BatchExecuteStatement",
                "rds-data:BeginTransaction",
                "rds-data:CommitTransaction",
                "rds-data:ExecuteSql",
                "rds-data:ExecuteStatement",
                "rds-data:RollbackTransaction",
                "rds-db:connect"
            ],
            "Resource": [
                "arn:aws:rds:ap-southeast-1:345215350058:cluster:itsa-db"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_rds_permissions_attachment" {
  policy_arn = aws_iam_policy.lambda_rds_permissions.arn
  role = aws_iam_role.lambda_role.name
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "lambda.amazonaws.com" },
        Action    = "sts:AssumeRole"
      }
    ]
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AWSLambda_FullAccess"
  ]
}

resource "aws_lambda_function" "this" {
  function_name    = "points-ledger-api"
  runtime          = "go1.x"
  handler          = "main"
  role             = aws_iam_role.lambda_role.arn
  filename         = "./main.zip"
  source_code_hash = filebase64sha256("./main.zip")
  timeout          = 10

  vpc_config {
    subnet_ids         = ["subnet-067ead2dfed1c40f8", "subnet-07722d2cb666693fd"]
    security_group_ids = ["sg-05a445100fc469383"]
  }

  environment {
    variables = {
      ENV         = "lambda"
      DB_HOST     = data.aws_ssm_parameter.db_url.value
      DB_PORT     = 3306
      DB_USER     = data.aws_ssm_parameter.db_user.value
      DB_PASSWORD = data.aws_ssm_parameter.db_password.value
      DB_NAME     = "itsa"
      JWT_SECRET  = data.aws_ssm_parameter.jwt_secret.value
    }
  }
}

resource "aws_api_gateway_resource" "root" {
  rest_api_id = "vm1swtn9ii"
  parent_id   = "9gy5jtm4yf"
  path_part   = "points"
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = "vm1swtn9ii"
  parent_id   = aws_api_gateway_resource.root.id
  path_part   = "{proxy+}"
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:ap-southeast-1:345215350058:vm1swtn9ii/*/*/points/*"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id             = "vm1swtn9ii"
  resource_id             = aws_api_gateway_resource.this.id
  http_method             = aws_api_gateway_method.this.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.this.invoke_arn
}

resource "aws_api_gateway_method" "this" {
  rest_api_id      = "vm1swtn9ii"
  resource_id      = aws_api_gateway_resource.this.id
  http_method      = "ANY"
  authorization    = "CUSTOM"
  authorizer_id    = "kjkxid"
  api_key_required = false
}
