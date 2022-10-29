import aws_cdk as cdk
import aws_cdk.aws_ecr as ecr
import aws_cdk.aws_ec2 as ec2
import aws_cdk.aws_ecs as ecs
import aws_cdk.aws_iam as iam
import aws_cdk.aws_ecs_patterns as ecs_patterns
from constructs import Construct

class AppInfraStack(cdk.Stack):

    def __init__(self, scope: cdk.App, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        err_repository = ecr.Repository(self,
            "app-repository",
            repository_name="app-repository")

        vpc = ec2.Vpc(self, "app-vpc", max_azs = 3)

        cluster = ecs.Cluster(self, "app-cluster", cluster_name = "app-cluster", vpc = vpc)

        execution_role = iam.Role(self, "app-execution-role",
            assumed_by=iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
            role_name="app-execution-role")

        execution_role.add_to_policy(iam.PolicyStatement(
            effect = iam.Effect.ALLOW,
            resources=["*"],
            actions=[
                "ecr:GetAuthorizationToken",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
            ]
        ))

        task_image_options = ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
            family="app-task-definition",
            execution_role=execution_role,
            image=ecs.ContainerImage.from_registry("amazon/amazon-ecs-sample"),
            container_name="app",
            container_port=8080,
            log_driver=ecs.LogDrivers.aws_logs(stream_prefix="app-container")
        )

        ecs_patterns.ApplicationLoadBalancedFargateService(self, "app-service",
            cluster=cluster,
            service_name="app-service",
            desired_count=2,
            task_image_options=task_image_options,
            public_load_balancer=True)

        task_role = iam.Role(self, "app-task-role", assumed_by=iam.ServicePrincipal("ecs-tasks.amazonaws.com"), role_name="app-task-role")

        task_role.add_to_policy(iam.PolicyStatement(effect=iam.Effect.ALLOW, resources=["*"], actions=["ecs:*", "ec2:DescribeNetworkInterfaces"]))