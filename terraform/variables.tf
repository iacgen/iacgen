variable "region" {
  description = "AWS region in which the project needs to be setup"
}

variable "project" {
  description = "Name of application"
}

variable "env" {
  description = "Environment type (dev, qa, stage, prod)"
}

# load balancer variables
variable "lb_type" {
  default     = "application"
  description = "Load Balancer type. Can be application, network or classic"
}

variable "lb_listener_port" {
  default     = 80
  description = "Load balancer's listening port"
}

variable "lb_listener_protocol" {
  default     = "HTTP"
  description = "Load balancer's protocol"
}

variable "http_listener_action" {
  default     = "forward"
  description = "Listener's action. Can be forward, redirect, fixed-response, authenticate-cognito or authenticate-oidc"
}

variable "tg_protocol" {
  default = "HTTP"
  description = "Target group's protocol"
}

variable "tg_type" {
  default = "ip"
  description = "Type of the target group. Can be ip or instance"
}

# security group variables
variable "sg_type_ingress" {
  default = "ingress"
}

variable "sg_type_egress" {
  default = "egress"
}

variable "all_cidr_block" {
  default = "0.0.0.0/0"
}

variable "tcp_protocol" {
  default = "tcp"
}

# ecs variables
variable "launch_type" {
  default = "FARGATE"
  description = "ECS launch type. Either EC2 or Fargate"
}

variable "ecr_policy" {
  default = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"
}

variable "ecs_task_policy" {
  default = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

variable "task_role_service" {
  default = "ecs-tasks.amazonaws.com"
}

variable "desired_tasks" {
  default = 2
}

variable "network_mode" {
  default = "awsvpc"
}

variable "cpu" {

}

variable "memory" {

}