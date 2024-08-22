resource "aws_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
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

  # Allow HTTP access
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow egress to backend
  egress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [aws_security_group.backend_sg.id]
  }
}


resource "aws_security_group" "backend_sg" {
  vpc_id = aws_vpc.vpc.id

  # Allow ingress from frontend
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [aws_security_group.frontend_sg.id]
  }

  # Allow egress to ZincSearch
    egress {
        from_port   = 4080
        to_port     = 4080
        protocol    = "tcp"
        cidr_blocks = [aws_security_group.zincsearch_sg.id]
    }
}

resource "aws_security_group" "zincsearch_sg" {
  vpc_id = aws_vpc.vpc.id

  # Allow ZincSearch port access
  ingress {
    from_port   = 4080
    to_port     = 4080
    protocol    = "tcp"
    cidr_blocks = [aws_security_group.backend_sg.id]
  }
}

resource "aws_ecs_cluster" "cluster" {
  name = "${var.company_name}-ecs-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# Frontend
resource "aws_ecs_task_definition" "frontend_task" {
  family = "frontend-task"

  container_definitions = jsonencode([
    {
      name      = "frontend"
      image     = "TODOOOOO"
      essential = true
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
      }]
  }])
}

resource "aws_ecs_service" "frontend_service" {
  name            = "frontend-service"
  cluster         = aws_ecs_cluster.cluster.arn
  task_definition = aws_ecs_task_definition.frontend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets = [aws_subnet.sub4.id]
    security_groups = [aws_security_group.frontend_sg.id]
  }
}

# Backend
resource "aws_ecs_task_definition" "backend_task" {
  family = "backend-task"

  container_definitions = jsonencode([
    {
      name      = "backend"
      image     = "TODOOOOO"
      essential = true
      portMappings = [
        {
          containerPort = 5173
          hostPort      = 5173
      }]
  }])
}

resource "aws_ecs_service" "backend_service" {
  name           = "backend-service"
    cluster        = aws_ecs_cluster.cluster.arn
    task_definition = aws_ecs_task_definition.backend_task.arn
    desired_count  = 1
    launch_type    = "FARGATE"

    network_configuration {
      subnets = [aws_subnet.sub1.id, aws_subnet.sub2.id, aws_subnet.sub3.id]
      security_groups = [aws_security_group.backend_sg.id]
    }
}

# ZincSearch
resource "aws_ecs_task_definition" "zincsearch_task" {
  family = "zincsearch-task"

  container_definitions = jsonencode([
    {
      name      = "zincsearch"
      image     = "TODOOOOO"
      essential = true
      portMappings = [
        {
          containerPort = 4080
          hostPort      = 4080
      }]
  }])
}

resource "aws_ecs_service" "zincsearch_service" {
  name           = "zincsearch-service"
    cluster        = aws_ecs_cluster.cluster.arn
    task_definition = aws_ecs_task_definition.zincsearch_task.arn
    desired_count  = 1
    launch_type    = "FARGATE"

    network_configuration {
      subnets = [aws_subnet.sub1.id, aws_subnet.sub2.id, aws_subnet.sub3.id]
      security_groups = [aws_security_group.zincsearch_sg.id]
    }
}