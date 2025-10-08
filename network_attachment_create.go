package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) AttachNetworkToHost(
	name string,
	comment string,
	description string,
	hostID HostID,
	networkID NetworkID,
	nicName string,
	retries ...RetryStrategy,
) (result NetworkAttachment, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("attaching network %s to host %s on nic %s", networkID, hostID, nicName),
		o.logger,
		retries,
		func() error {
			hostService := o.conn.SystemService().HostsService().HostService(string(hostID))
			nicsService := hostService.NicsService()
			nicsResponse, err := nicsService.List().Send()
			if err != nil {
				return err
			}
			var nicID string
			nics, ok := nicsResponse.Nics()
			if !ok {
				return fmt.Errorf("no nics returned from host %s", hostID)
			}
			for _, nic := range nics.Slice() {
				if nic.MustName() == nicName {
					nicID = nic.MustId()
					break
				}
			}
			if nicID == "" {
				return fmt.Errorf("nic %s not found on host %s", nicName, hostID)
			}
			networkAttachmentBuilder := ovirtsdk.NewNetworkAttachmentBuilder()
			networkAttachmentBuilder.Name(name)
			networkAttachmentBuilder.Comment(comment)
			networkAttachmentBuilder.Description(description)
			networkAttachmentBuilder.Host(ovirtsdk.NewHostBuilder().Id(string(hostID)).MustBuild())
			networkAttachmentBuilder.Network(ovirtsdk.NewNetworkBuilder().Id(string(networkID)).MustBuild())
			networkAttachmentBuilder.HostNic(ovirtsdk.NewHostNicBuilder().Id(nicID).MustBuild())
			req := hostService.NetworkAttachmentsService().Add()
			response, err := req.Attachment(networkAttachmentBuilder.MustBuild()).Send()
			if err != nil {
				return err
			}
			sdkAttachment, ok := response.Attachment()
			if !ok {
				return newFieldNotFound("response from network attachment creation", "attachment")
			}
			result, err = convertSDKNetworkAttachment(sdkAttachment, o)
			return err
		})
	return result, err
}

func (m *mockClient) AttachNetworkToHost(
	name string,
	comment string,
	description string,
	hostID HostID,
	networkID NetworkID,
	nicName string,
	retries ...RetryStrategy,
) (result NetworkAttachment, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	id := NetworkAttachmentID(m.GenerateUUID())
	m.networkAttachment[id] = &networkAttachment{}
	return m.networkAttachment[id], nil
}
