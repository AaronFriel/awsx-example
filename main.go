package main

import (
	aws_ecs "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := aws_ecs.NewCluster(ctx, "cluster", nil)
		if err != nil {
			return err
		}
		lb, err := lb.NewApplicationLoadBalancer(ctx, "lb", nil)
		if err != nil {
			return err
		}
		_, err = ecs.NewFargateService(ctx, "nginx", &ecs.FargateServiceArgs{
			Cluster: cluster.Arn,
			TaskDefinitionArgs: &ecs.FargateServiceTaskDefinitionArgs{
				Container: &ecs.TaskDefinitionContainerDefinitionArgs{
					Image:  pulumi.StringPtr("nginx:latest"),
					Cpu:    pulumi.IntPtr(512),
					Memory: pulumi.IntPtr(128),
					PortMappings: ecs.TaskDefinitionPortMappingArray{
						ecs.TaskDefinitionPortMappingArgs{
							ContainerPort: pulumi.IntPtr(80),
							TargetGroup:   lb.DefaultTargetGroup,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("dnsName", lb.LoadBalancer.DnsName())
		return nil
	})
}
