package main

import (
	"errors"
	"fmt"
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	"github.com/esxcloud/bosh-esxcloud-cpi/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Dispatch", func() {
	var (
		ctx        *types.Context
		configPath string
	)

	AfterEach(func() {
		if configPath != "" {
			os.Remove(configPath)
		}
	})

	It("returns a valid bosh JSON response given valid arguments", func() {
		actions := map[string]types.ActionFn{
			"create_vm": createVM,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

		Expect(res.Result).Should(Equal("fake-vm-id"))
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when given an invalid argument", func() {
		actions := map[string]types.ActionFn{
			"create_vm": createVM,
		}
		args := []interface{}{5}
		res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(types.CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when function errors", func() {
		actions := map[string]types.ActionFn{
			"create_vm": createVmError,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(types.CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when function panics", func() {
		actions := map[string]types.ActionFn{
			"create_vm": createVmPanic,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(types.CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when method not implemented", func() {
		actions := map[string]types.ActionFn{}
		args := []interface{}{"fake-agent-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(types.NotImplementedError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("loads JSON config correctly", func() {
		configFile, err := ioutil.TempFile("", "bosh-esxcloud-cpi-config")
		if err != nil {
			panic(err)
		}
		configPath = configFile.Name()
		jsonConfig := `{"ESXCloud":{"APIFE":{"Hostname":"none", "Port": 123}}}`
		configFile.WriteString(jsonConfig)

		context, err := loadConfig(configPath)
		expectedURL := fmt.Sprintf("http://%s:%d", "none", 123)
		Expect(context.ECClient.Endpoint).Should(Equal(expectedURL))
		Expect(err).Should(BeNil())
	})
})

func createVM(ctx *types.Context, args []interface{}) (result interface{}, err error) {
	_, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where agent_id should be")
	}
	return "fake-vm-id", nil
}

func createVmError(ctx *types.Context, args []interface{}) (result interface{}, err error) {
	return nil, errors.New("error occurred")
}

func createVmPanic(ctx *types.Context, args []interface{}) (result interface{}, err error) {
	panic("oh no!")
}
