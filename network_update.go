package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) UpdateNetwork(
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
		fmt.Sprintf("updating network %s", name),
		o.logger,
		retries,
		func() error {
			networkBuilder := ovirtsdk.NewNetworkBuilder()
			networkBuilder.Name(name)
			networkBuilder.DataCenter(ovirtsdk.NewDataCenterBuilder().Id(string(dataCenterId)).MustBuild())
			networkBuilder.Description(description)
			networkBuilder.Comment(comment)
			networkBuilder.Vlan(ovirtsdk.NewVlanBuilder().Id(int64(vlanID)).MustBuild())
			req := o.conn.SystemService().NetworksService().NetworkService(name).Update()
			response, err := req.Network(networkBuilder.MustBuild()).Send()
			if err != nil {
				return err
			}
			network, ok := response.Network()
			if !ok {
				return newFieldNotFound("response from network update", "network")
			}
			result, err = convertSDKNetwork(network, o)
			return err
		})
	return result, err
}
