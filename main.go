package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

// Credentials
var projectID = "hidden-howl-252922"                                         // Your ProjectID
var jsonPath = filepath.Join(os.Getenv("HOME"), "mypersonalgcpjsonkey.json") // path to your JSON file

// GoogleCloudClient is a generic wrapper for talking with individual services inside Google Cloud
// such as Cloud Resource Manager, IAM, Services, Billing and DNS
type GoogleCloudClient struct {
	// Structs from Google library
	Resource *cloudresourcemanager.Service
	IAM      *iam.Service
	Service  *servicemanagement.APIService
	Billing  *cloudbilling.APIService
	DNS      *dns.Service

	// Required user input
	ProjectID string
	JSONPath  string
}

// NewGoogleCloudClient returns a pointer to the `GoogleCloudClient` instance
func NewGoogleCloudClient(projectID string, json string) (*GoogleCloudClient, error) {
	ctx := context.Background()

	// Client for Cloud Resource Manager
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with Cloud Resource Manager Service: %v", err)
	}

	// Client for IAM
	iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the IAM Service: %v", err)
	}

	// Client for Service Infrastructure Manager
	servicemanagementService, err := servicemanagement.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Service Management Service: %v", err)
	}

	// Client for Cloud Billing
	cloudbillingService, err := cloudbilling.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Cloud Billing Account: %v", err)
	}

	// Client for Google Cloud DNS API
	dnsService, err := dns.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Cloud DNS: %v", err)
	}

	return &GoogleCloudClient{
		Resource:  cloudresourcemanagerService,
		IAM:       iamService,
		Service:   servicemanagementService,
		Billing:   cloudbillingService,
		DNS:       dnsService,
		ProjectID: projectID,
		JSONPath:  json,
	}, nil
}

// ListProjects lists the Projects of a GCP service account and returns an error
func (c *GoogleCloudClient) ListProjects() (*cloudresourcemanager.ListProjectsResponse, error) {
	projectsList, err := c.Resource.Projects.List().Do()
	if err != nil {
		return nil, err
	}
	return projectsList, nil
}

// GetProject returns a project from GCP
func (c *GoogleCloudClient) GetProject(projectID string) (*cloudresourcemanager.Project, error) {
	project, err := c.Resource.Projects.Get(projectID).Do()
	if err != nil {
		return nil, err
	}
	return project, nil
}

// DeleteProject deletes a project from GCP
func (c *GoogleCloudClient) DeleteProject(projectID string) (*cloudresourcemanager.Empty, error) {
	project, err := c.Resource.Projects.Delete(projectID).Do()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func main() {
	gcpClient, err := NewGoogleCloudClient(projectID, jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gcpClient.Resource.BasePath)
	fmt.Println(gcpClient.IAM.BasePath)
	fmt.Println(gcpClient.Service.BasePath)
	fmt.Println(gcpClient.Billing.BasePath)
	fmt.Println(gcpClient.DNS.BasePath)

	// Get a List of Projects
	resp, err := gcpClient.ListProjects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range resp.Projects {
		fmt.Printf("Project Name: %s\tProjectID: %s\r\n", project.Name, project.ProjectId)
	}

	// Get a Project
	resp1, err := gcpClient.GetProject("hidden-howl-252922")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp1.Name)

	fmt.Printf("The lifecycle state of %s project is %s:", resp1.Name, resp1.LifecycleState)

	// Delete a Project
	// resp2, err := gcpClient.DeleteProject("hidden-howl-252922")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(resp2.ServerResponse)

}
