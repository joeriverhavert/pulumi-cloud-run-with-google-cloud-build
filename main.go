// TO DO
// Add authication, our group only
// Revision traffic split

package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//Load godotenv
		err := godotenv.Load()
		if err != nil {
			ctx.Log.Error("No .env file found!", nil)
		}

		// Dotenv variables
		service_name := os.Getenv("PULUMI_CLOUDRUN_SERVICE_NAME")
		service_location := os.Getenv("PULUMI_CLOUDRUN_SERVICE_LOCATION")
		service_container_image := os.Getenv("PULUMI_CLOUDRUN_SERVICE_CONTAINER_IMAGE")
		service_container_port, _ := strconv.Atoi(os.Getenv("PULUMI_CLOUDRUN_SERVICE_CONTAINER_PORT"))
		service_container_port_protocol := os.Getenv("PULUMI_CLOUDRUN_SERVICE_CONTAINER_PORT_PROTOCOL")

		// Current service
		var traffic cloudrun.ServiceTrafficArray
		currentService, err := cloudrun.LookupService(ctx, &cloudrun.LookupServiceArgs{
			Name:     service_name,
			Location: service_location,
		})

		if err != nil && strings.Contains(err.Error(), "not found") {
			ctx.Log.Info("No services are found, all traffic will go to the new revision.", nil)

			//traffic split 100%
			traffic = cloudrun.ServiceTrafficArray{
				&cloudrun.ServiceTrafficArgs{
					Percent:        pulumi.Int(100),
					LatestRevision: pulumi.Bool(true),
				},
			}
		} else if err != nil {
			return err
		} else {
			ctx.Log.Info("Detected existing revision. Routing 30% of traffic to the new revision, 70% to the previous one.", nil)
			// traffic split 70&/30%
			currentRev := currentService.Statuses[0].LatestCreatedRevisionName
			traffic = cloudrun.ServiceTrafficArray{
				&cloudrun.ServiceTrafficArgs{
					Percent:        pulumi.Int(30),
					LatestRevision: pulumi.Bool(true),
				},
				&cloudrun.ServiceTrafficArgs{
					Percent:      pulumi.Int(70),
					RevisionName: pulumi.String(currentRev),
				},
			}

		}

		// Create new Cloud Run service
		run, err := cloudrun.NewService(ctx, service_name, &cloudrun.ServiceArgs{
			Name:     pulumi.String(service_name),
			Location: pulumi.String(service_location),
			Template: cloudrun.ServiceTemplateArgs{
				Spec: cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: pulumi.String(service_container_image),
							Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
								&cloudrun.ServiceTemplateSpecContainerPortArgs{
									ContainerPort: pulumi.Int(service_container_port),
									Protocol:      pulumi.String(service_container_port_protocol),
								},
							},
							StartupProbe: &cloudrun.ServiceTemplateSpecContainerStartupProbeArgs{
								PeriodSeconds:    pulumi.Int(5),
								FailureThreshold: pulumi.Int(3),
								TcpSocket: &cloudrun.ServiceTemplateSpecContainerStartupProbeTcpSocketArgs{
									Port: pulumi.Int(service_container_port),
								},
							},
						},
					},
				},
			},
			Traffics: traffic,
		})
		if err != nil {
			return err
		}

		noauth, err := organizations.LookupIAMPolicy(ctx, &organizations.LookupIAMPolicyArgs{
			Bindings: []organizations.GetIAMPolicyBinding{
				{
					Role: "roles/run.invoker",
					Members: []string{
						"allUsers",
					},
				},
			},
		}, nil)
		if err != nil {
			return err
		}

		_, err = cloudrun.NewIamPolicy(ctx, "noauth", &cloudrun.IamPolicyArgs{
			Location:   run.Location,
			Project:    run.Project,
			Service:    run.Name,
			PolicyData: pulumi.String(noauth.PolicyData),
		})

		// Export the DNS name of the bucket
		ctx.Export("Cloud run service", run.Name)
		return nil
	})
}
