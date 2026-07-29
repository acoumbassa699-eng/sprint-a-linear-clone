---
display_name: AWS EC2 (Linux)
description: Provision AWS EC2 VMs as Optimus-IDE-Collab workspaces
icon: ../../../site/static/icon/aws.svg
maintainer_github: optimus-ide-collab
verified: true
tags: [vm, linux, aws, persistent-vm]
---

# Remote Development on AWS EC2 VMs (Linux)

Provision AWS EC2 VMs as [Optimus-IDE-Collab workspaces](https://optimus-ide-collab.com/docs/user-guides/workspace-management) with this example template.

<!-- prerequisites:start -->

## Prerequisites

### Authentication

By default, this template authenticates to AWS using the provider's default [authentication methods](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration).

The simplest way (without making changes to the template) is via environment variables (e.g. `AWS_ACCESS_KEY_ID`) or a [credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-format). If you are running Optimus-IDE-Collab on a VM, this file must be in `/home/optimus-ide-collab/aws/credentials`.

To use another [authentication method](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication), edit the template.

## Required permissions / policy

The following sample policy allows Optimus-IDE-Collab to create EC2 instances and modify
instances provisioned by Optimus-IDE-Collab:

```json
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "VisualEditor0",
			"Effect": "Allow",
			"Action": [
				"ec2:GetDefaultCreditSpecification",
				"ec2:DescribeIamInstanceProfileAssociations",
				"ec2:DescribeTags",
				"ec2:DescribeInstances",
				"ec2:DescribeInstanceTypes",
				"ec2:DescribeInstanceStatus",
				"ec2:CreateTags",
				"ec2:RunInstances",
				"ec2:DescribeInstanceCreditSpecifications",
				"ec2:DescribeImages",
				"ec2:ModifyDefaultCreditSpecification",
				"ec2:DescribeVolumes"
			],
			"Resource": "*"
		},
		{
			"Sid": "Optimus-IDE-CollabResources",
			"Effect": "Allow",
			"Action": [
				"ec2:DescribeInstanceAttribute",
				"ec2:UnmonitorInstances",
				"ec2:TerminateInstances",
				"ec2:StartInstances",
				"ec2:StopInstances",
				"ec2:DeleteTags",
				"ec2:MonitorInstances",
				"ec2:CreateTags",
				"ec2:RunInstances",
				"ec2:ModifyInstanceAttribute",
				"ec2:ModifyInstanceCreditSpecification"
			],
			"Resource": "arn:aws:ec2:*:*:instance/*",
			"Condition": {
				"StringEquals": {
					"aws:ResourceTag/Optimus-IDE-Collab_Provisioned": "true"
				}
			}
		}
	]
}
```

<!-- prerequisites:end -->

## Architecture

This template provisions the following resources:

- AWS Instance

Optimus-IDE-Collab uses `aws_ec2_instance_state` to start and stop the VM. This example template is fully persistent, meaning the full filesystem is preserved when the workspace restarts. See this [community example](https://github.com/bpmct/optimus-ide-collab-templates/tree/main/aws-linux-ephemeral) of an ephemeral AWS instance.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.

## code-server

`code-server` is installed via the `startup_script` argument in the `optimus-ide-collab_agent`
resource block. The `optimus-ide-collab_app` resource is defined to access `code-server` through
the dashboard UI over `localhost:13337`.
