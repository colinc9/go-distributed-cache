import aws_cdk as core
import aws_cdk.assertions as assertions

from app_infra.app_infra_stack import AppInfraStack

# example tests. To run these tests, uncomment this file along with the example
# resource in app_infra/app_infra_stack.py
def test_sqs_queue_created():
    app = core.App()
    stack = AppInfraStack(app, "app-infra")
    template = assertions.Template.from_stack(stack)

#     template.has_resource_properties("AWS::SQS::Queue", {
#         "VisibilityTimeout": 300
#     })
