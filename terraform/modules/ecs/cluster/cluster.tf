resource "aws_service_discovery_http_namespace" "http"{
  name = var.name
}

resource "aws_ecs_cluster" "cluster" {
  name = var.name
  service_connect_defaults {
    namespace = aws_service_discovery_http_namespace.http.arn
  }
}
