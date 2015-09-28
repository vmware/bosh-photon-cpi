package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	var (
		server *testServer
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		cleanTenants(client)
		server.Close()
	})

	Describe("CreateAndDeleteTenant", func() {
		It("Tenant create and delete succeeds", func() {
			mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
			task, err := client.Tenants.Create(tenantSpec)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_TENANT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_TENANT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Tenants.Delete(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_TENANT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})

		It("Tenant create fails", func() {
			tenantSpec := &TenantCreateSpec{}
			task, err := client.Tenants.Create(tenantSpec)

			Expect(err).ShouldNot(BeNil())
			Expect(task).Should(BeNil())
		})
	})

	Describe("GetTenant", func() {
		It("Get returns a tenant ID", func() {
			mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			tenantName := randomString(10, "go-sdk-tenant-")
			tenantSpec := &TenantCreateSpec{Name: tenantName}
			task, err := client.Tenants.Create(tenantSpec)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_TENANT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			server.SetResponseJson(200, &Tenants{[]Tenant{Tenant{Name: tenantName}}})
			tenants, err := client.Tenants.GetAll()

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(tenants).ShouldNot(BeNil())

			var found bool
			for _, tenant := range tenants.Items {
				if tenant.Name == tenantName && tenant.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_TENANT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			_, err = client.Tenants.Delete(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetTenantTasks", func() {
		var (
			option string
		)

		Context("no extra options for GetTask", func() {
			BeforeEach(func() {
				option = ""
			})

			It("GetTasks returns a completed task", func() {
				mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
				mockTask.Entity.ID = "mock-task-id"
				server.SetResponseJson(200, mockTask)
				tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
				task, err := client.Tenants.Create(tenantSpec)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
				Expect(task).ShouldNot(BeNil())
				Expect(task.Operation).Should(Equal("CREATE_TENANT"))
				Expect(task.State).Should(Equal("COMPLETED"))

				server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
				taskList, err := client.Tenants.GetTasks(task.Entity.ID, &TaskGetOptions{State: option})
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
				Expect(taskList).ShouldNot(BeNil())
				Expect(taskList.Items).Should(ContainElement(*task))

				mockTask = createMockTask("DELETE_TENANT", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				_, err = client.Tenants.Delete(task.Entity.ID)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
			})
		})

		Context("Searching COMPLETED state for GetTask", func() {
			BeforeEach(func() {
				option = "COMPLETED"
			})

			It("GetTasks returns a completed task", func() {
				mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
				mockTask.Entity.ID = "mock-task-id"
				server.SetResponseJson(200, mockTask)
				tenantSpec := &TenantCreateSpec{Name: randomString(10, "go-sdk-tenant-")}
				task, err := client.Tenants.Create(tenantSpec)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
				Expect(task).ShouldNot(BeNil())
				Expect(task.Operation).Should(Equal("CREATE_TENANT"))
				Expect(task.State).Should(Equal("COMPLETED"))

				server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
				taskList, err := client.Tenants.GetTasks(task.Entity.ID, &TaskGetOptions{State: option})
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
				Expect(taskList).ShouldNot(BeNil())
				Expect(taskList.Items).Should(ContainElement(*task))

				mockTask = createMockTask("DELETE_TENANT", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				_, err = client.Tenants.Delete(task.Entity.ID)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
			})
		})

	})
})

var _ = Describe("ResourceTicket", func() {
	var (
		server   *testServer
		client   *Client
		tenantID string
	)

	BeforeEach(func() {
		server, client = testSetup()
		tenantID = createTenant(server, client)
	})

	AfterEach(func() {
		cleanTenants(client)
		server.Close()
	})

	Describe("CreateAndGetResourceTicket", func() {
		It("Resource ticket create and get succeeds", func() {
			mockTask := createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			spec := &ResourceTicketCreateSpec{
				Name:   randomString(10),
				Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
			}
			task, err := client.Tenants.CreateResourceTicket(tenantID, spec)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_RESOURCE_TICKET"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockResList := ResourceList{[]ResourceTicket{ResourceTicket{TenantId: tenantID, Name: spec.Name, Limits: spec.Limits}}}
			server.SetResponseJson(200, mockResList)
			resList, err := client.Tenants.GetResourceTickets(tenantID, &ResourceTicketGetOptions{spec.Name})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(resList).ShouldNot(BeNil())

			var found bool
			for _, res := range resList.Items {
				if res.Name == spec.Name && res.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())
		})
	})
})

var _ = Describe("Project", func() {
	var (
		server   *testServer
		client   *Client
		tenantID string
		resName  string
	)

	BeforeEach(func() {
		server, client = testSetup()
		tenantID = createTenant(server, client)
		resName = createResTicket(server, client, tenantID)
	})

	AfterEach(func() {
		cleanTenants(client)
		server.Close()
	})

	Describe("CreateGetAndDeleteProject", func() {
		It("Project create and delete succeeds", func() {
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
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_PROJECT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockProjects := ProjectList{[]ProjectCompact{ProjectCompact{Name: projSpec.Name}}}
			server.SetResponseJson(200, mockProjects)
			projList, err := client.Tenants.GetProjects(tenantID, &ProjectGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(projList).ShouldNot(BeNil())

			var found bool
			for _, proj := range projList.Items {
				if proj.Name == projSpec.Name && proj.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_PROJECT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Projects.Delete(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_PROJECT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})
})
