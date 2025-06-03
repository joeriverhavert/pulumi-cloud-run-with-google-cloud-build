# Pulumi GCP Go: Example Cloud Run Application with Cloud Build

This example project provisions a Google Cloud Run application using Pulumi, Go and Cloud build. It demonstrates how to define infrastructure as code for serverless deployments in GCP.

## Features
	•	Integration of the Pulumi GCP provider within a Go-based Pulumi program
	•	Deployment of a Cloud Run service with configurable runtime parameters
	•	Support for multiple environments using Pulumi stacks
	•	CI/CD via Cloud Build

## Providers
	•	Google Cloud Platform via the Pulumi GCP SDK for Go (github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp)

## Prerequisites
	•	Go 1.20 or later installed
	•	A Google Cloud account with billing enabled
	•	GCP credentials configured for Pulumi (for example, via gcloud auth application-default login)


## Project Layout
```bash
.
├── cloudbuild.yaml     # Cloud Build configuration file defining CI/CD pipeline steps
├── Dockerfile            # Docker image build instructions for the Cloud Run service
├── go.mod                # Go module definition for dependency management
├── go.sum                # Checksums for Go module dependencies to ensure integrity
├── hyperspace-demo-app   # Directory containing your Go application source code
├── main.go               # Pulumi program that provisions Cloud Run infrastructure and resources
├── Pulumi.sbx.yaml       # Pulumi stack configuration file for the ‘sbx’ (sandbox) environment
├── Pulumi.yaml           # Pulumi project configuration defining project metadata
└── README.md             # This documentation file describing project usage and structure
```
