resource "aws_iam_role" "cloudwatch" {
  name = "api-gateway-cloudwatch"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "cloudwatch" {
  name = "default"
  role = aws_iam_role.cloudwatch.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents",
                "logs:GetLogEvents",
                "logs:FilterLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_api_gateway_account" "main" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch.arn
}

resource "aws_api_gateway_rest_api" "user" {
  name        = "user-${var.env}"
  description = "User CRUD api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  minimum_compression_size = 0
}

resource "aws_api_gateway_request_validator" "user" {
  name                        = "user"
  rest_api_id                 = aws_api_gateway_rest_api.user.id
  validate_request_body       = true
  validate_request_parameters = true
}

resource "aws_api_gateway_resource" "user_path" {
  rest_api_id = aws_api_gateway_rest_api.user.id
  parent_id   = aws_api_gateway_rest_api.user.root_resource_id
  path_part   = "users"
}

resource "aws_api_gateway_resource" "user_path_with_id" {
  rest_api_id = aws_api_gateway_rest_api.user.id
  parent_id   = aws_api_gateway_resource.user_path.id
  path_part   = "{id}"
}

resource "aws_api_gateway_method" "get_user_by_id" {
  rest_api_id = aws_api_gateway_rest_api.user.id
  resource_id = aws_api_gateway_resource.user_path_with_id.id
  http_method = "GET"

  authorization = "AWS_IAM"

  request_parameters = {
    "method.request.path.id" = true
  }
}

resource "aws_api_gateway_model" "update_user_request" {
  rest_api_id  = aws_api_gateway_rest_api.user.id
  name         = "UpdateUserRequest"
  description  = "validate update user request"
  content_type = "application/json"

  schema = <<EOF
{
  "title" : "UpdateUserRequest",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    },
    "phoneNumber": {
      "type": "string"
    },
    "country": {
      "type": "string"
    },
    "city": {
      "type": "string"
    },
    "postalCode": {
      "type": "string"
    },
    "activelyLooking": {
      "type": "boolean"
    },
    "resumes": {
      "type": "object",
      "properties": {
        "key": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "primary": {
          "type": "boolean"
        },
        "uploaded": {
          "type": "number"
        },
        "index": {
          "type": "number"
        }
      }
    },
    "coverLetters": {
      "type": "object",
      "properties": {
        "key": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "primary": {
          "type": "boolean"
        },
        "uploaded": {
          "type": "number"
        },
        "index": {
          "type": "number"
        }
      }
    },
    "transcripts": {
      "type": "object",
      "properties": {
        "key": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "primary": {
          "type": "boolean"
        },
        "uploaded": {
          "type": "number"
        },
        "index": {
          "type": "number"
        }
      }
    },
    "title": {
      "type": "string",
      "enum": ["Professor", "Associate Professor", "Assistant Professor", "Adjunct Professor", "Emeritus Professor"]
    },
    "school": {
      "type": "string",
      "enum": ["McMaster University", "Queen's University", "Ryerson University", "University of Guelph", "University of Toronto", "University of Waterloo", "Western University", "Wilfrid Laurier University", "York University"]
    },
    "department": {
      "type": "string",
      "enum": ["Biology", "Chemistry", "Computer Science", "Engineering", "Math", "Physics"]
    },
    "acceptingApplicants": {
      "type": "boolean"
    },
    "studentCriteria": {
      "type": "string",
      "enum": ["Publications", "GPA", "Industry Experience", "Well Rounded", "Local", "International", "Extra Curriculars"]
    }
  }
}
EOF
}

resource "aws_api_gateway_method" "update_user_by_id" {
  rest_api_id = aws_api_gateway_rest_api.user.id
  resource_id = aws_api_gateway_resource.user_path_with_id.id
  http_method = "PUT"

  authorization = "AWS_IAM"

  request_parameters = {
    "method.request.path.id" = true
  }

  request_validator_id = aws_api_gateway_request_validator.user.id
  request_models = {
    "application/json" = aws_api_gateway_model.update_user_request.name
  }
}

module "integration" {
  source = "./integration"

  env    = var.env
  global = var.global

  user_depends_on = [aws_api_gateway_method.get_user_by_id]

  user = {
    api_id = aws_api_gateway_rest_api.user.id
    lambda = var.user_lambda
    ids = [
      aws_api_gateway_resource.user_path_with_id.id,
      aws_api_gateway_resource.user_path_with_id.id
    ]
    methods = [
      aws_api_gateway_method.get_user_by_id.http_method,
      aws_api_gateway_method.update_user_by_id.http_method,
    ]
    cors_ids = [aws_api_gateway_resource.user_path_with_id.id]
  }
}