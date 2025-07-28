package ovirtclient_test

import (
	"fmt"

	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

// The following example demonstrates how to create a virtual machine. It is set up
// using the test helper, but can be easily modified to use the client directly.
func ExampleVMClient_create() {
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// This is the cluster the VM will be created on.
	clusterID := helper.GetClusterID()
	// Use the blank template as a starting point.
	templateID := helper.GetBlankTemplateID()
	// Set the VM name
	name := "test-vm"
	// Create the optional parameters.
	params := ovirtclient.CreateVMParams()

	// Create the VM...
	vm, err := client.CreateVM(clusterID, templateID, name, params)
	if err != nil {
		panic(fmt.Sprintf("failed to create VM (%v)", err))
	}

	// ... and then remove it. Alternatively, you could call client.RemoveVM(vm.ID()).
	if err := vm.Remove(); err != nil {
		panic(fmt.Sprintf("failed to remove VM (%v)", err))
	}
}
