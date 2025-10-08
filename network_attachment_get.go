package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetNetworkAttachment(
	id NetworkAttachmentID,
	hostID HostID,
	nicName string,
	retries ...RetryStrategy,
) (result NetworkAttachment, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting network attachment %s from host %s on nic %s", id, hostID, nicName),
		o.logger,
		retries,
		func() error {
			hostService := o.conn.SystemService().HostsService().HostService(string(hostID))
			sdkAttachment, err := hostService.NetworkAttachmentsService().AttachmentService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			attachment, ok := sdkAttachment.Attachment()
			if !ok {
				return newFieldNotFound("response from getting network attachment from host", "attachment")
			}
			result, err = convertSDKNetworkAttachment(attachment, o)
			return err
		})
	return result, err
}

func (m *mockClient) GetNetworkAttachment(
	id NetworkAttachmentID,
	hostID HostID,
	nicName string,
	retries ...RetryStrategy,
) (result NetworkAttachment, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil, nil
}
