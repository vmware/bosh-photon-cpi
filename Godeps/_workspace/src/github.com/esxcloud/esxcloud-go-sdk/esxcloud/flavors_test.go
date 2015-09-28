package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flavor", func() {
	var (
		server     *testServer
		client     *Client
		flavorSpec *FlavorCreateSpec
	)

	BeforeEach(func() {
		server, client = testSetup()
		flavorSpec = &FlavorCreateSpec{
			Name: randomString(10, "go-sdk-flavor-"),
			Kind: "vm",
			Cost: []QuotaLineItem{QuotaLineItem{"GB", 16, "vm.memory"}},
		}
	})

	AfterEach(func() {
		cleanFlavors(client)
		server.Close()
	})

	Describe("CreateGetAndDeleteFlavor", func() {
		It("Flavor create and delete succeeds", func() {
			mockTask := createMockTask("CREATE_FLAVOR", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			flavorSpec := &FlavorCreateSpec{
				Name: randomString(10, "go-sdk-flavor-"),
				Kind: "vm",
				Cost: []QuotaLineItem{QuotaLineItem{"GB", 16, "vm.memory"}},
			}
			task, err := client.Flavors.Create(flavorSpec)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_FLAVOR"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_FLAVOR", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Flavors.Delete(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_FLAVOR"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetFlavor", func() {
		var (
			flavorName string
			flavorID   string
		)

		BeforeEach(func() {
			flavorName, flavorID = createFlavor(server, client)
		})

		It("Get flavor succeeds", func() {
			server.SetResponseJson(200, Flavor{Name: flavorName})
			flavor, err := client.Flavors.Get(flavorID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(flavor).ShouldNot(BeNil())
			Expect(flavor.ID).Should(Equal(flavorID))
			Expect(flavor.Name).Should(Equal(flavorName))

			mockTask := createMockTask("DELETE_FLAVOR", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			_, err = client.Flavors.Delete(flavorID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		It("Get all flavor succeeds", func() {
			server.SetResponseJson(200, &FlavorList{[]Flavor{Flavor{Name: flavorName}}})
			flavorList, err := client.Flavors.GetAll(&FlavorGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(flavorList).ShouldNot(BeNil())

			var found bool
			for _, flavor := range flavorList.Items {
				if flavor.Name == flavorName && flavor.ID == flavorID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask := createMockTask("DELETE_FLAVOR", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			_, err = client.Flavors.Delete(flavorID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetTasks", func() {
		It("GetTasks returns a completed task", func() {
			mockTask := createMockTask("CREATE_FLAVOR", "COMPLETED")
			mockTask.Entity.ID = "mock-task-id"
			server.SetResponseJson(200, mockTask)

			task, err := client.Flavors.Create(flavorSpec)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
			taskList, err := client.Flavors.GetTasks(task.Entity.ID, &TaskGetOptions{})

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(taskList).ShouldNot(BeNil())
			Expect(taskList.Items).Should(ContainElement(*task))

			mockTask = createMockTask("DELETE_FLAVOR", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			_, err = client.Flavors.Delete(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})
})
