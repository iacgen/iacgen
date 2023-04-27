package tfaws

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const ecsTemplate = `resource "aws_security_group" "{{.AppName}}_sg" {
	name        = "{{.AppName}}-${var.project}-${var.env}-${var.launch_type}-ecs-sg"
	vpc_id      = module.vpc.id
  
	ingress {
	  protocol        = "tcp"
	  from_port       = {{.ContainerPort}}
	  to_port         = {{.ContainerPort}}
	  security_groups = [module.lb_sg.id]
	}
  
	egress {
	  protocol    = "-1"
	  from_port   = 0
	  to_port     = 0
	  cidr_blocks = ["0.0.0.0/0"]
	}
  }

module "{{.AppName}}_task_definition" {
	source         = "./modules/ecs/task_definition"
	name           = "{{.AppName}}-${var.project}-${var.env}-ecs-task-def"
	launch_type    = [var.launch_type]
	network_mode   = var.network_mode
	cpu            = var.cpu
	memory         = var.memory
	execution_role = module.task_execution_role.arn
	definitions = templatefile("definitions/container_definition.json", {
	  repository_url  = "{{.Image}}"
	  definition_name = "{{.AppName}}-${var.project}-${var.env}-${var.launch_type}"
	  container_port  = {{.ContainerPort}}
	  host_port       = {{.ContainerPort}}
	})
  }
  
  module "{{.AppName}}_ecs_service" {
	source          = "./modules/ecs/service"
	name            = "{{.AppName}}-${var.project}-${var.env}-${var.launch_type}-ecs-service"
	cluster         = module.ecs_cluster.id
	task_definition = module.{{.AppName}}_task_definition.arn
	desired_count   = var.desired_tasks
	launch_type     = var.launch_type
	lb_target_group = module.{{.AppName}}_target_group.arn
	container_name  = "{{.AppName}}-${var.project}-${var.env}-${var.launch_type}"
	container_port  = {{.ContainerPort}}
	network_config  = [
		{
		  subnets         = module.vpc.private_subnet_ids
		  public_ip       = "false"
		  security_groups = [aws_security_group.{{.AppName}}_sg.id]
		}
	  ]
	http_listener   = module.{{.AppName}}_http_listener.arn
	namespace 		= var.project
	dns_name 		= "{{.AppName}}"
  }
  
  module "{{.AppName}}_target_group" {
	source      = "./modules/alb/target_group"
	name        = "{{.AppName}}-${var.project}-${var.env}-tg"
	port        = {{.ContainerPort}}
	protocol    = var.tg_protocol
	target_type = var.tg_type
	vpc_id      = module.vpc.id
  }
  
  module "{{.AppName}}_http_listener" {
	source      = "./modules/alb/lb_listener"
	lb_arn      = module.load_balancer.arn
	port        = var.lb_listener_port
	protocol    = var.lb_listener_protocol
	action_type = var.http_listener_action
	tg_arn      = module.{{.AppName}}_target_group.arn
  }`

type TfAwsECSTemplate struct {
	AppName       string `json:"app_name"`
	Image         string `json:"image"`
	ContainerPort int    `json:"container_port"`
}

type TfAwsECSTemplates []TfAwsECSTemplate

func (g *TfAwsECSTemplate) RenderTemplate(basedir string) error {
	// create new template
	t, err := template.New("ecstemplate").Parse(ecsTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse ECS template: %w", err)
	}

	// render template
	file, err := os.Create(filepath.Join(basedir, fmt.Sprintf("%s.tf", g.AppName)))
	if err != nil {
		return fmt.Errorf("failed to create file for ecs template: %w", err)
	}
	if err := t.Execute(file, g); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}
