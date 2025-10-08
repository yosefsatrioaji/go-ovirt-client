package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateNetwork(
	dataCenterId DatacenterID,
	name string,
	description string,
	comment string,
	vlanID int,
	retries ...RetryStrategy,
) (result Network, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	if err := validateNetworkCreationParameters(name, dataCenterId); err != nil {
		return nil, err
	}
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("creating network %s", name),
		o.logger,
		retries,
		func() error {
			networkBuilder := ovirtsdk.NewNetworkBuilder()
			networkBuilder.Name(name)
			networkBuilder.DataCenter(ovirtsdk.NewDataCenterBuilder().Id(string(dataCenterId)).MustBuild())
			networkBuilder.Description(description)
			networkBuilder.Comment(comment)
			networkBuilder.Vlan(ovirtsdk.NewVlanBuilder().Id(int64(vlanID)).MustBuild())
			req := o.conn.SystemService().NetworksService().Add()
			response, err := req.Network(networkBuilder.MustBuild()).Send()
			if err != nil {
				return err
			}
			network, ok := response.Network()
			if !ok {
				return newFieldNotFound("response from network creation", "network")
			}
			result, err = convertSDKNetwork(network, o)
			return err
		})
	return result, err
}

func validateNetworkCreationParameters(name string, dataCenterId DatacenterID) error {
	if name == "" {
		return newError(EBadArgument, "name cannot be empty for Network creation")
	}
	if dataCenterId == "" {
		return newError(EBadArgument, "Datacenter ID cannot be empty for Datacenter creation")
	}
	return nil
}
