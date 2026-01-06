package ovirtclient

import ovirtsdk "github.com/ovirt/go-ovirt"

// NetworkAttachmentID is a unique identifier for a network attachment.
type NetworkAttachmentID string

type NetworkAttachmentClient interface {
	// AttachNetworkToHost attaches a network to a host, returning the resulting network attachment.
	AttachNetworkToHost(hostID HostID, networkID NetworkID, hostNicID HostNICID, retries ...RetryStrategy) (NetworkAttachment, error)
	// DetachNetworkFromHost detaches a network from a host.
	DetachNetworkFromHost(id NetworkAttachmentID, hostID HostID, hostNicID HostNICID, retries ...RetryStrategy) error
	// GetNetworkAttachment retrieves a specific network attachment from a host.
	GetNetworkAttachment(id NetworkAttachmentID, hostID HostID, hostNicID HostNICID, retries ...RetryStrategy) (NetworkAttachment, error)
	// NetworkAttachmentList lists all network attachments from a host.
	NetworkAttachmentList(hostID HostID, retries ...RetryStrategy) ([]NetworkAttachment, error)
}

type NetworkAttachmentData interface {
	ID() NetworkAttachmentID
	HostID() HostID
	NetworkID() NetworkID
	HostNICID() HostNICID
}

type NetworkAttachment interface {
	NetworkAttachmentData
	// Host fetches the host associated with this network attachment. This is a network call and may be slow.
	Host(retries ...RetryStrategy) (Host, error)
	// Network fetches the network associated with this network attachment. This is a network call and may be slow.
	Network(retries ...RetryStrategy) (Network, error)
}

func convertSDKNetworkAttachment(sdkObject *ovirtsdk.NetworkAttachment, client Client) (NetworkAttachment, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "ID")
	}
	host, ok := sdkObject.Host()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "host")
	}
	hostID, ok := host.Id()
	if !ok {
		return nil, newFieldNotFound("Host on NetworkAttachment", "ID")
	}
	network, ok := sdkObject.Network()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "network")
	}
	networkID, ok := network.Id()
	if !ok {
		return nil, newFieldNotFound("Network on NetworkAttachment", "ID")
	}
	hostNic, ok := sdkObject.HostNic()
	if !ok {
		return nil, newFieldNotFound("HostNic on NetworkAttachment", "hostNic")
	}
	hostNicId, ok := hostNic.Id()
	if !ok {
		return nil, newFieldNotFound("HostNic on NetworkAttachment", "ID")
	}
	return &networkAttachment{
		client:    client,
		id:        NetworkAttachmentID(id),
		hostID:    HostID(hostID),
		networkID: NetworkID(networkID),
		hostNicID: HostNICID(hostNicId),
	}, nil
}

type networkAttachment struct {
	client    Client
	id        NetworkAttachmentID
	hostID    HostID
	networkID NetworkID
	hostNicID HostNICID
}

func (n networkAttachment) ID() NetworkAttachmentID {
	return n.id
}

func (n networkAttachment) HostID() HostID {
	return n.hostID
}

func (n networkAttachment) NetworkID() NetworkID {
	return n.networkID
}

func (n networkAttachment) HostNICID() HostNICID {
	return n.hostNicID
}

func (n networkAttachment) Host(retries ...RetryStrategy) (Host, error) {
	return n.client.GetHost(n.hostID, retries...)
}

func (n networkAttachment) Network(retries ...RetryStrategy) (Network, error) {
	return n.client.GetNetwork(n.networkID, retries...)
}
