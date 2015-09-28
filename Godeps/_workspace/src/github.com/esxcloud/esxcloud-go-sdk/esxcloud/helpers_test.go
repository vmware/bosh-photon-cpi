package esxcloud

import (
	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func toJson(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		// Since this method is only for testing, don't return
		// any errors, just panic.
		panic("Error serializing struct into JSON")
	}
	// json.Marshal returns []byte, convert to string
	return string(res[:])
}

func hasStep(task *Task, operation, state string) bool {
	for _, step := range task.Steps {
		if step.State == state && step.Operation == operation {
			return true
		}
	}
	return false
}

func createTenant(server *testServer, client *Client) string {
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
	task, err := client.Tenants.Create(tenantSpec)
	GinkgoT().Log(err)
	Expect(err).Should(BeNil())
	return task.Entity.ID
}

// Checks the list of tenants and deletes the ones created by go-sdk
func cleanTenants(client *Client) {
	tenants, err := client.Tenants.GetAll()
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, tenant := range tenants.Items {
		if strings.HasPrefix(tenant.Name, "go-sdk-tenant-") {
			cleanProjects(client, tenant.ID)
			_, err := client.Tenants.Delete(tenant.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func createResTicket(server *testServer, client *Client, tenantID string) string {
	resTicketName := randomString(10)
	spec := &ResourceTicketCreateSpec{
		Name:   resTicketName,
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
	}
	mockTask := createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	_, err := client.Tenants.CreateResourceTicket(tenantID, spec)
	GinkgoT().Log(err)
	Expect(err).Should(BeNil())
	return resTicketName
}

func createProject(server *testServer, client *Client, tenantID string, resName string) string {
	mockTask := createMockTask("CREATE_PROJECT", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	projSpec := &ProjectCreateSpec{
		ResourceTicket: ResourceTicketReservation{
			resName,
			[]QuotaLineItem{QuotaLineItem{"GB", 2, "vm.memory"}},
		},
		Name: randomString(10, "go-sdk-project-"),
	}
	task, err := client.Tenants.CreateProject(tenantID, projSpec)
	GinkgoT().Log(err)
	Expect(err).Should(BeNil())
	return task.Entity.ID
}

// Checks the projects for the tenant and deletes ones created by go-sdk
func cleanProjects(client *Client, tenantID string) {
	projList, err := client.Tenants.GetProjects(tenantID, &ProjectGetOptions{})
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, proj := range projList.Items {
		if strings.HasPrefix(proj.Name, "go-sdk-project-") {
			_, err := client.Projects.Delete(proj.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

// Returns flavorName, flavorID
func createFlavor(server *testServer, client *Client) (string, string) {
	mockTask := createMockTask("CREATE_FLAVOR", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	flavorName := randomString(10, "go-sdk-flavor-")
	flavorSpec := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "persistent-disk.cost"}},
		"persistent-disk",
		flavorName,
	}
	task, err := client.Flavors.Create(flavorSpec)
	GinkgoT().Log(err)
	Expect(err).Should(BeNil())
	return flavorName, task.Entity.ID
}

func cleanFlavors(client *Client) {
	flavorList, err := client.Flavors.GetAll(&FlavorGetOptions{})
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, flavor := range flavorList.Items {
		if strings.HasPrefix(flavor.Name, "go-sdk-flavor-") {
			_, err := client.Flavors.Delete(flavor.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func cleanDisks(client *Client, projID string) {
	diskList, err := client.Projects.GetDisks(projID, &DiskGetOptions{})
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, disk := range diskList.Items {
		if strings.HasPrefix(disk.Name, "go-sdk-disk-") {
			task, err := client.Disks.Delete(disk.ID)
			task, err = client.Tasks.Wait(task.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func createImage(server *testServer, client *Client) string {
	mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
	server.SetResponseJson(200, mockTask)

	// create image from file
	imagePath := "../testdata/tty_tiny.ova"
	task, err := client.Images.CreateFromFile(imagePath, &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
	task, err = client.Tasks.Wait(task.ID)

	GinkgoT().Log(err)
	Expect(err).Should(BeNil())

	return task.Entity.ID
}

func cleanImages(client *Client) {
	imageList, err := client.Images.GetAll()
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, image := range imageList.Items {
		if image.Name == "tty_tiny.ova" {
			task, err := client.Images.Delete(image.ID)
			task, err = client.Tasks.Wait(task.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func cleanVMs(client *Client, projID string) {
	vmList, err := client.Projects.GetVMs(projID, &VmGetOptions{})
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, vm := range vmList.Items {
		if strings.HasPrefix(vm.Name, "go-sdk-vm-") {
			task, err := client.VMs.Delete(vm.ID)
			task, err = client.Tasks.Wait(task.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func cleanHosts(client *Client) {
	hostList, err := client.Hosts.GetAll()
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, host := range hostList.Items {
		if host.Metadata != nil {
			if val, ok := host.Metadata["Test"]; ok && val == "go-sdk-host"{
				task, err := client.Hosts.Delete(host.ID)
				task, err = client.Tasks.Wait(task.ID)
				if err != nil {
					GinkgoT().Log(err)
				}
			}
		}
	}
}

func cleanNetworks(client *Client) {
	networks, err := client.Networks.GetAll(&NetworkGetOptions{})
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, network := range networks.Items {
		if strings.HasPrefix(network.Name, "go-sdk-network-") {
			task, err := client.VMs.Delete(network.ID)
			task, err = client.Tasks.Wait(task.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}

func cleanClusters(client *Client, projID string) {
	clusters, err := client.Projects.GetClusters(projID)
	if err != nil {
		GinkgoT().Log(err)
	}
	for _, cluster := range clusters.Items {
		if strings.HasPrefix(cluster.Name, "go-sdk-cluster-") {
			task, err := client.Clusters.Delete(cluster.ID)
			task, err = client.Tasks.Wait(task.ID)
			if err != nil {
				GinkgoT().Log(err)
			}
		}
	}
}
