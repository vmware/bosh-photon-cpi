package esxcloud

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Represents stateless context needed to call esxcloud APIs.
type Client struct {
	options         ClientOptions
	httpClient      *http.Client
	Endpoint        string
	Status          *StatusAPI
	Tenants         *TenantsAPI
	Tasks           *TasksAPI
	Projects        *ProjectsAPI
	Flavors         *FlavorsAPI
	Images          *ImagesAPI
	Disks           *DisksAPI
	VMs             *VmAPI
	Hosts           *HostsAPI
	Deployments     *DeploymentsAPI
	ResourceTickets *ResourceTicketsAPI
	Networks        *NetworksAPI
	Clusters        *ClustersAPI
}

// Options for Client
type ClientOptions struct {
	// When using the Tasks.Wait APIs, defines the duration of how long
	// the SDK should continue to poll the server. Default is 30 minutes.
	// TasksAPI.WaitTimeout() can be used to specify timeout on
	// individual calls.
	TaskPollTimeout time.Duration

	// Whether or not to ignore any TLS errors when talking to esxcloud,
	// false by default.
	IgnoreCertificate bool

	// For tasks APIs, defines the delay between each polling attempt.
	// Default is 100 milliseconds.
	taskPollDelay time.Duration

	// For tasks APIs, defines the number of retries to make in the event
	// of an error. Default is 3.
	taskRetryCount int

	// AccessToken for user authentication. Default is empty.
	Token string
}

// Creates a new ESXCloud client with specified options. If options
// is nil, default options will be used.
func NewClient(endpoint string, options *ClientOptions) (c *Client) {
	defaultOptions := &ClientOptions{
		TaskPollTimeout:   30 * time.Minute,
		taskPollDelay:     100 * time.Millisecond,
		taskRetryCount:    3,
		Token:             "",
		IgnoreCertificate: false,
	}
	if options != nil {
		if options.TaskPollTimeout != 0 {
			defaultOptions.TaskPollTimeout = options.TaskPollTimeout
		}
		if options.taskPollDelay != 0 {
			defaultOptions.taskPollDelay = options.taskPollDelay
		}
		if options.taskRetryCount != 0 {
			defaultOptions.taskRetryCount = options.taskRetryCount
		}
		if options.Token != "" {
			defaultOptions.Token = options.Token
		}
		defaultOptions.IgnoreCertificate = options.IgnoreCertificate
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: defaultOptions.IgnoreCertificate},
	}
	c = &Client{Endpoint: endpoint, httpClient: &http.Client{Transport: tr}}
	// Ensure a copy of options is made, rather than using a pointer
	// which may change out from underneath if misused by the caller.
	c.options = *defaultOptions
	c.Status = &StatusAPI{c}
	c.Tenants = &TenantsAPI{c}
	c.Tasks = &TasksAPI{c}
	c.Projects = &ProjectsAPI{c}
	c.Flavors = &FlavorsAPI{c}
	c.Images = &ImagesAPI{c}
	c.Disks = &DisksAPI{c}
	c.VMs = &VmAPI{c}
	c.Hosts = &HostsAPI{c}
	c.Deployments = &DeploymentsAPI{c}
	c.ResourceTickets = &ResourceTicketsAPI{c}
	c.Networks = &NetworksAPI{c}
	c.Clusters = &ClustersAPI{c}
	return
}

// Creates a new ESXCloud client with specified options and http.Client.
// Useful for functional testing where http calls must be mocked out.
// If options is nil, default options will be used.
func NewTestClient(endpoint string, options *ClientOptions, httpClient *http.Client) (c *Client) {
	c = NewClient(endpoint, options)
	c.httpClient = httpClient
	return
}
