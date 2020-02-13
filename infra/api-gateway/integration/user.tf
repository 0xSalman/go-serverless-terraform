resource "aws_api_gateway_integration" "user_integrations" {
  depends_on              = [var.user_depends_on]
  count                   = length(var.user["ids"])
  rest_api_id             = var.user["api_id"]
  resource_id             = var.user["ids"][count.index]
  http_method             = var.user["methods"][count.index]
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.global["region"]}:lambda:path/2015-03-31/functions/${var.user["lambda"]["arn"]}/invocations"
}

resource "aws_lambda_permission" "user" {
  depends_on    = [aws_api_gateway_integration.user_cors_integrations]
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.user["lambda"]["name"]
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.global["region"]}:${var.global["account"]}:${var.user["api_id"]}/*/*"
}

#resource "aws_api_gateway_method_settings" "user" {
#  rest_api_id = aws_api_gateway_rest_api.user.id
#  stage_name  = aws_api_gateway_deployment.user.stage_name
#  method_path = "$*/*"
#
#  settings {
#    metrics_enabled = false
#    logging_level   = "INFO"
#  }
#}

resource "aws_api_gateway_method" "user_cors_methods" {
  depends_on    = [aws_api_gateway_integration.user_integrations]
  count         = length(var.user["cors_ids"])
  rest_api_id   = var.user["api_id"]
  resource_id   = var.user["cors_ids"][count.index]
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "user_cors_resps" {
  depends_on = [
    aws_api_gateway_integration.user_integrations,
    aws_api_gateway_method.user_cors_methods
  ]
  count       = length(var.user["cors_ids"])
  rest_api_id = var.user["api_id"]
  resource_id = var.user["cors_ids"][count.index]
  http_method = "OPTIONS"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration" "user_cors_integrations" {
  depends_on = [
    aws_api_gateway_integration.user_integrations,
    aws_api_gateway_method.user_cors_methods
  ]
  count       = length(var.user["cors_ids"])
  rest_api_id = var.user["api_id"]
  resource_id = var.user["cors_ids"][count.index]
  http_method = "OPTIONS"
  type        = "MOCK"

  passthrough_behavior = "WHEN_NO_TEMPLATES"
  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }
}

resource "aws_api_gateway_integration_response" "user_cors_integration_resps" {
  depends_on = [
    aws_api_gateway_integration.user_integrations,
    aws_api_gateway_method.user_cors_methods,
    aws_api_gateway_integration.user_cors_integrations
  ]
  count       = length(var.user["cors_ids"])
  rest_api_id = var.user["api_id"]
  resource_id = var.user["cors_ids"][count.index]
  http_method = "OPTIONS"
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'POST,GET,PUT,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
}

resource "aws_api_gateway_deployment" "user" {
  depends_on = [
    aws_api_gateway_integration.user_integrations,
    aws_api_gateway_integration_response.user_cors_integration_resps
  ]
  rest_api_id = var.user["api_id"]
  stage_name  = var.env

  # this is a workaround to force API deployment
  # more info at: https://github.com/terraform-providers/terraform-provider-aws/issues/162
  variables = {
    deployed_at = timestamp()
    #    trigger_hash = sha1(join(",", [
    #      jsonencode(aws_api_gateway_method.get_user_by_id),
    #      jsonencode(aws_api_gateway_integration.get_user_by_id)
    #    ]))
  }

  #  lifecycle {
  #    create_before_destroy = true
  #  }
}
