package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Network", func() {
	var (
		server     *testServer
		client     *Client
		networkSpec *NetworkCreateSpec
	)

	BeforeEach(func() {
		server, client = testSetup()
		networkSpec = &NetworkCreateSpec{
			Name: randomString(10, "go-sdk-network-"),
			PortGroups: []string{"portGroup"},
		}
	})

	AfterEach(func() {
		cleanNetworks(client)
		server.Close()
	})

	Describe("CreateDeleteNetwork", func() {
		It("Network create and delete succeeds", func() {
			mockTask := createMockTask("CREATE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Networks.Create(networkSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_NETWORK"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Networks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_NETWORK"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetNetwork", func() {
		It("Get network succeeds", func() {
			mockTask := createMockTask("CREATE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err := client.Networks.Create(networkSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, Network{Name: networkSpec.Name})
			network, err := client.Networks.Get(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(network).ShouldNot(BeNil())
			Expect(network.Name).Should(Equal(networkSpec.Name))

			mockTask = createMockTask("DELETE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Networks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		It("GetAll Network succeeds", func() {
			mockTask := createMockTask("CREATE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err := client.Networks.Create(networkSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &Networks{[]Network{Network{Name: networkSpec.Name}}})
			networks, err := client.Networks.GetAll(&NetworkGetOptions{})

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(networks).ShouldNot(BeNil())

			var found bool
			for _, network := range networks.Items {
				if network.Name == networkSpec.Name && network.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Networks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		It("GetAll Network with optional name succeeds", func() {
			mockTask := createMockTask("CREATE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err := client.Networks.Create(networkSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &Networks{[]Network{Network{Name: networkSpec.Name}}})
			networks, err := client.Networks.GetAll(&NetworkGetOptions{Name: networkSpec.Name})

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(networks).ShouldNot(BeNil())

			var found bool
			for _, network := range networks.Items {
				if network.Name == networkSpec.Name && network.ID == task.Entity.ID {
					found = true
					break
				}
			}
			Expect(len(networks.Items)).Should(Equal(1))
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_NETWORK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Networks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})
})
