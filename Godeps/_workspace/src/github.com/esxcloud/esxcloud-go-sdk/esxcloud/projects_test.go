package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	var (
		server     *testServer
		client     *Client
		tenantID   string
		resName    string
		projID     string
		flavorName string
		flavorID   string
	)

	BeforeEach(func() {
		server, client = testSetup()
		tenantID = createTenant(server, client)
		resName = createResTicket(server, client, tenantID)
		projID = createProject(server, client, tenantID, resName)
		flavorName, flavorID = createFlavor(server, client)

	})

	AfterEach(func() {
		cleanDisks(client, projID)
		cleanFlavors(client)
		cleanTenants(client)
		server.Close()
	})

	Describe("GetProjectTasks", func() {
		It("GetTasks returns a completed task", func() {
			mockTask := createMockTask("CREATE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			diskSpec := &DiskCreateSpec{
				Flavor:     flavorName,
				Kind:       "persistent-disk",
				CapacityGB: 2,
				Name:       randomString(10, "go-sdk-disk-"),
			}

			task, err := client.Projects.CreateDisk(projID, diskSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
			taskList, err := client.Projects.GetTasks(projID, &TaskGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(taskList).ShouldNot(BeNil())
			Expect(taskList.Items).Should(ContainElement(*task))

			// Clean disk
			mockTask = createMockTask("DELETE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Disks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetProjectDisks", func() {
		It("GetAll returns disk", func()  {
			mockTask := createMockTask("CREATE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			diskSpec := &DiskCreateSpec{
				Flavor:     flavorName,
				Kind:       "persistent-disk",
				CapacityGB: 2,
				Name:       randomString(10, "go-sdk-disk-"),
			}

			task, err := client.Projects.CreateDisk(projID, diskSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			diskMock := PersistentDisk{
				Name:       diskSpec.Name,
				Flavor:     diskSpec.Flavor,
				CapacityGB: diskSpec.CapacityGB,
				Kind:       diskSpec.Kind,
			}
			server.SetResponseJson(200, &DiskList{[]PersistentDisk{diskMock}})
			diskList, err := client.Projects.GetDisks(projID, &DiskGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(diskList).ShouldNot(BeNil())

			var found bool
			for _, disk := range diskList.Items {
				if disk.Name == diskSpec.Name && disk.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Disks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetProjectVms", func() {
		var (
			imageID      string
			flavorSpec   *FlavorCreateSpec
			vmFlavorSpec *FlavorCreateSpec
		)

		BeforeEach(func() {
			imageID = createImage(server, client)
			flavorSpec = &FlavorCreateSpec{
				[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
				"ephemeral-disk",
				randomString(10, "go-sdk-flavor-"),
			}

			_, err := client.Flavors.Create(flavorSpec)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			vmFlavorSpec = &FlavorCreateSpec{
				Name: randomString(10, "go-sdk-flavor-"),
				Kind: "vm",
				Cost: []QuotaLineItem{
					QuotaLineItem{"GB", 2, "vm.memory"},
					QuotaLineItem{"COUNT", 4, "vm.cpu"},
				},
			}
			_, err = client.Flavors.Create(vmFlavorSpec)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		AfterEach(func() {
			cleanVMs(client, projID)
		})

		It("GetAll returns vm", func() {
			mockTask := createMockTask("CREATE_VM", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			vmSpec := &VmCreateSpec{
				Flavor:        vmFlavorSpec.Name,
				SourceImageID: imageID,
				AttachedDisks: []AttachedDisk{
					AttachedDisk{
						CapacityGB: 1,
						Flavor:     flavorSpec.Name,
						Kind:       "ephemeral-disk",
						Name:       randomString(10),
						State:      "STARTED",
						BootDisk:   true,
					},
				},
				Name: randomString(10, "go-sdk-vm-"),
			}

			task, err := client.Projects.CreateVM(projID, vmSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			mockVm := VM{Name: vmSpec.Name}
			server.SetResponseJson(200, &VMs{[]VM{mockVm}})
			vmList, err := client.Projects.GetVMs(projID, &VmGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(vmList).ShouldNot(BeNil())

			var found bool
			for _, vm := range vmList.Items {
				if vm.Name == vmSpec.Name && vm.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_VM", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.VMs.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetProjectClusters", func() {
		It("GetAll returns cluster", func() {
			mockTask := createMockTask("CREATE_CLUSTER", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			clusterSpec := &ClusterCreateSpec{
				Name: randomString(10, "go-sdk-cluster-"),
				Type: "KUBERNETES",
				SlaveCount: 50,
				ExtendedProperties: map[string]string{},
			}

			task, err := client.Projects.CreateCluster(projID, clusterSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			mockCluster := Cluster{Name: clusterSpec.Name}
			server.SetResponseJson(200, &Clusters{[]Cluster{mockCluster}})
			clusterList, err := client.Projects.GetClusters(projID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(clusterList).ShouldNot(BeNil())

			var found bool
			for _, cluster := range clusterList.Items {
				if cluster.Name == clusterSpec.Name && cluster.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_CLUSTER", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Clusters.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})
})
