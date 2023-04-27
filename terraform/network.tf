#################### create network infrastructure ############

# create vpc, subnets
module "vpc" {
  source     = "./modules/network/vpc"
  name       = "${var.project}-${var.env}-vpc"
  cidr_block = "10.0.0.0/16"
}

# create internet gateway in vpc
module "internet_gateway" {
  source = "./modules/network/internet_gateway"
  vpc_id = module.vpc.id
  name   = "${var.project}-${var.env}-ig"
}

# create nat gateway in public subnet
module "nat_gateway" {
  source    = "./modules/network/nat_gateway"
  subnet_id = element(module.vpc.public_subnet_ids, 1)
  name      = "${var.project}-${var.env}-natgw"
}

# create public route table
module "public_route_table" {
  source  = "./modules/network/route_table"
  vpc_id  = module.vpc.id
  name    = "${var.project}-${var.env}-public-table"
  subnets = module.vpc.public_subnet_ids
}

# all outgoing traffic in public subnet
# to be directed to internet gateway
module "public_internet_gw_route" {
  source           = "./modules/network/route"
  route_table_id   = module.public_route_table.id
  destination_cidr = "0.0.0.0/0"
  gateway_id       = module.internet_gateway.id
}

# create private route table
module "private_route_table" {
  source  = "./modules/network/route_table"
  vpc_id  = module.vpc.id
  name    = "${var.project}-${var.env}-private-table"
  subnets = module.vpc.private_subnet_ids
}

# all outgoing traffic in private subnet
# to be directed to nat gateway
module "private_nat_gw_route" {
  source           = "./modules/network/route"
  route_table_id   = module.private_route_table.id
  destination_cidr = "0.0.0.0/0"
  gateway_id       = module.nat_gateway.id
}

#################### load balancer ################
# create security group for load balancer
module "lb_sg" {
  source = "./modules/network/sg"
  name   = "${var.project}-${var.env}-lb-sg"
  vpc_id = module.vpc.id
}

# create ingress rule
module "lb_ingress_http_rule" {
  source         = "./modules/network/sg_rule_cidr"
  type           = var.sg_type_ingress
  from_port      = 80
  to_port        = 80
  cidr_block     = [var.all_cidr_block]
  protocol       = var.tcp_protocol
  security_group = module.lb_sg.id
}

# create egress rule
module "lb_egress_sg_rule" {
  source         = "./modules/network/sg_rule_cidr"
  type           = var.sg_type_egress
  from_port      = 0
  to_port        = 0
  cidr_block     = [var.all_cidr_block]
  protocol       = var.tcp_protocol
  security_group = module.lb_sg.id
}

# create a load balancer in public subnets
# attach security group 
module "load_balancer" {
  source     = "./modules/alb/load_balancer"
  name       = "${var.project}-${var.env}-lb"
  lb_type    = var.lb_type
  lb_sg      = [module.lb_sg.id]
  lb_subnets = module.vpc.public_subnet_ids
}

####################### ECS #################
# create ecs task execution role
module "task_execution_role" {
  source  = "./modules/ecs/iam"
  name    = "${var.project}-${var.env}-execution-role"
  service = var.task_role_service
}

module "execution_role_ecr" {
  source     = "./modules/ecs/role_policy_attachment"
  role       = module.task_execution_role.name
  policy_arn = var.ecr_policy
}

module "execution_role_ecs" {
  source     = "./modules/ecs/role_policy_attachment"
  role       = module.task_execution_role.name
  policy_arn = var.ecs_task_policy
}

module "ecs_cluster" {
  source = "./modules/ecs/cluster"
  name   = "${var.project}-${var.env}-ecs-cluster"
}