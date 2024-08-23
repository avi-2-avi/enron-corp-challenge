resource "aws_vpc" "vpc" {
  cidr_block           = "172.16.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

# Private subnet
resource "aws_subnet" "sub1" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${var.region}a"
}

# Private subnet
resource "aws_subnet" "sub2" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.16.2.0/24"
  availability_zone = "${var.region}b"
}

# Private subnet
resource "aws_subnet" "sub3" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.16.3.0/24"
  availability_zone = "${var.region}c"
}

# Public subnet
resource "aws_subnet" "sub4" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = "172.16.4.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "${var.region}d"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "public_rt_assoc" {
  subnet_id      = aws_subnet.sub4.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_security_group" "frontend_sg" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_security_group_rule" "frontend_http_inbound" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.frontend_sg.id
}

resource "aws_security_group_rule" "frontend_to_backend" {
  type                     = "egress"
  from_port                = 8080
  to_port                  = 8080
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.frontend_sg.id
  security_group_id        = aws_security_group.backend_sg.id
}

resource "aws_security_group" "backend_sg" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_security_group_rule" "backend_inbound" {
  type                     = "ingress"
  from_port                = 8080
  to_port                  = 8080
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.frontend_sg.id
  security_group_id        = aws_security_group.backend_sg.id
}

resource "aws_security_group_rule" "backend_to_zincsearch" {
  type                     = "egress"
  from_port                = 4080
  to_port                  = 4080
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.backend_sg.id
  security_group_id        = aws_security_group.zincsearch_sg.id
}

resource "aws_security_group" "zincsearch_sg" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_security_group_rule" "zincsearch_inbound" {
  type                     = "ingress"
  from_port                = 4080
  to_port                  = 4080
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.backend_sg.id
  security_group_id        = aws_security_group.zincsearch_sg.id
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        },
      },
    ],
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  ]
}

resource "aws_ecs_cluster" "cluster" {
  name = "ecs-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# Frontend
resource "aws_ecs_task_definition" "frontend_task" {
  family                   = "enron-frontend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name      = "frontend"
      image     = "9433462577/enron-front:v5"
      essential = true
      portMappings = [{
        containerPort = 80
        hostPort      = 80
      }]
    }
  ])
}

resource "aws_ecs_service" "frontend_service" {
  name            = "frontend-service"
  cluster         = aws_ecs_cluster.cluster.arn
  task_definition = aws_ecs_task_definition.frontend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.sub4.id]
    security_groups  = [aws_security_group.frontend_sg.id]
    assign_public_ip = true
  }
}

# Backend
resource "aws_ecs_task_definition" "backend_task" {
  family                   = "backend-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name      = "backend"
      image     = "9433462577/enron-back:v5"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
      }]
  }])
}

resource "aws_ecs_service" "backend_service" {
  name            = "backend-service"
  cluster         = aws_ecs_cluster.cluster.arn
  task_definition = aws_ecs_task_definition.backend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.sub1.id, aws_subnet.sub2.id, aws_subnet.sub3.id]
    security_groups = [aws_security_group.backend_sg.id]
  }
}

# ZincSearch
resource "aws_ecs_task_definition" "zincsearch_task" {
  family                   = "zincsearch-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "zincsearch"
      image     = "public.ecr.aws/zinclabs/zincsearch:latest"
      essential = true
      portMappings = [
        {
          containerPort = 4080
          hostPort      = 4080
      }]
      environment = [
        {
          name  = "ZINC_DATA_PATH"
          value = "/data"
        },
        {
          name  = "ZINC_FIRST_ADMIN_USER"
          value = "admin"
        },
        {
          name  = "ZINC_FIRST_ADMIN_PASSWORD"
          value = "Pass123!!!"
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "zincsearch_service" {
  name            = "zincsearch-service"
  cluster         = aws_ecs_cluster.cluster.arn
  task_definition = aws_ecs_task_definition.zincsearch_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.sub1.id, aws_subnet.sub2.id, aws_subnet.sub3.id]
    security_groups = [aws_security_group.zincsearch_sg.id]
  }
}
