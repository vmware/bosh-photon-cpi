package photon

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for projects API.
type ProjectsAPI struct {
	client *Client
}

// Options for GetDisks API.
type DiskGetOptions struct {
	Name string `urlParam:"name"`
}

// Options for GetDisks API.
type VmGetOptions struct {
	Name string `urlParam:"name"`
}

var projectUrl string = "/projects/"

// Deletes the project with specified ID. Any VMs, disks, etc., owned by the project must be deleted first.
func (api *ProjectsAPI) Delete(projectID string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+projectUrl+projectID, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Creates a disk on the specified project.
func (api *ProjectsAPI) CreateDisk(projectID string, spec *DiskCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+projectUrl+projectID+"/disks", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets disks for project with the specified ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *ProjectsAPI) GetDisks(projectID string, options *DiskGetOptions) (result *DiskList, err error) {
	uri := api.client.Endpoint + projectUrl + projectID + "/disks"
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &DiskList{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Creates a VM on the specified project.
func (api *ProjectsAPI) CreateVM(projectID string, spec *VmCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+projectUrl+projectID+"/vms", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets all tasks with the specified project ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *ProjectsAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + projectUrl + id + "/tasks"
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &TaskList{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets vms for project with the specified ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *ProjectsAPI) GetVMs(projectID string, options *VmGetOptions) (result *VMs, err error) {
	uri := api.client.Endpoint + projectUrl + projectID + "/vms"
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &VMs{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Creates a cluster on the specified project.
func (api *ProjectsAPI) CreateCluster(projectID string, spec *ClusterCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+projectUrl+projectID+"/clusters", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets clusters for project with the specified ID
func (api *ProjectsAPI) GetClusters(projectID string) (result *Clusters, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+projectUrl+projectID+"/clusters", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Clusters{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets the project with a specified ID.
func (api *ProjectsAPI) Get(id string) (project *ProjectCompact, err error) {
	res, err := rest.Get(api.client.httpClient, api.getEntityUrl(id), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	project = &ProjectCompact{}
	err = json.NewDecoder(res.Body).Decode(project)
	return
}

// Set security groups for this project, overwriting any existing ones.
func (api *ProjectsAPI) SetSecurityGroups(projectID string, securityGroups *SecurityGroups) (task *Task, err error) {
	return setSecurityGroups(api.client, api.getEntityUrl(projectID), securityGroups)
}

func (api *ProjectsAPI) getEntityUrl(id string) (url string) {
	return api.client.Endpoint + projectUrl + id
}