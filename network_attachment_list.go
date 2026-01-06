package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) NetworkAttachmentList(
	hostID HostID,
	retries ...RetryStrategy,
) (result []NetworkAttachment, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting network attachments from host %s", hostID),
		o.logger,
		retries,
		func() error {
			hostService := o.conn.SystemService().HostsService().HostService(string(hostID))
			sdkAttachment, err := hostService.NetworkAttachmentsService().List().Send()
			if err != nil {
				return err
			}
			attachments, ok := sdkAttachment.Attachments()
			if !ok {
				return newFieldNotFound("response from getting network attachments from host", "attachments")
			}
			result = make([]NetworkAttachment, len(attachments.Slice()))
			for i, attachment := range attachments.Slice() {
				result[i], err = convertSDKNetworkAttachment(attachment, o)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return result, err
}

func (m *mockClient) NetworkAttachmentList(
	hostID HostID,
	retries ...RetryStrategy,
) (result []NetworkAttachment, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil, nil
}
