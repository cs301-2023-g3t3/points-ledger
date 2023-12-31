data "aws_ssm_parameter" "db_user" {
  name = "db_user"
}

data "aws_ssm_parameter" "db_password" {
  name = "db_password"
}

data "aws_ssm_parameter" "db_url" {
  name = "db_url"
}

data "aws_ssm_parameter" "jwt_secret" {
  name = "jwt_secret"
}

data "aws_iam_role" "existing_role" {
    name = "lambda-execution-role"
}

output "existing_role_arn" {
    value = data.aws_iam_role.existing_role.arn
}
