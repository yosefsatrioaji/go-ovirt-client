package ovirtclient

import ovirtsdk "github.com/ovirt/go-ovirt"

// NetworkAttachmentID is a unique identifier for a network attachment.
type NetworkAttachmentID string

type NetworkAttachmentClient interface {
	// AttachNetworkToHost attaches a network to a host, returning the resulting network attachment.
	AttachNetworkToHost(comment string, description string, hostID HostID, networkID NetworkID, nicName string, retries ...RetryStrategy) (NetworkAttachment, error)
	// DetachNetworkFromHost detaches a network from a host.
	DetachNetworkFromHost(id NetworkAttachmentID, hostID HostID, nicName string, retries ...RetryStrategy) error
	// GetNetworkAttachment retrieves a specific network attachment from a host.
	GetNetworkAttachment(id NetworkAttachmentID, hostID HostID, nicName string, retries ...RetryStrategy) (NetworkAttachment, error)
}

type NetworkAttachmentData interface {
	Comment() string
	Description() string
	ID() NetworkAttachmentID
	HostID() HostID
	NetworkID() NetworkID
	NicName() string
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
	comment, ok := sdkObject.Comment()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "comment")
	}
	description, ok := sdkObject.Description()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "description")
	}
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
		return nil, newFieldNotFound("NetworkAttachment", "host nic")
	}
	nicName, ok := hostNic.Name()
	if !ok {
		return nil, newFieldNotFound("NetworkAttachment", "host nic name")
	}
	return &networkAttachment{
		client:      client,
		comment:     comment,
		description: description,
		id:          NetworkAttachmentID(id),
		hostID:      HostID(hostID),
		networkID:   NetworkID(networkID),
		nicName:     nicName,
	}, nil
}

type networkAttachment struct {
	client      Client
	comment     string
	description string
	id          NetworkAttachmentID
	hostID      HostID
	networkID   NetworkID
	nicName     string
}

func (n networkAttachment) ID() NetworkAttachmentID {
	return n.id
}

func (n networkAttachment) Comment() string {
	return n.comment
}

func (n networkAttachment) Description() string {
	return n.description
}

func (n networkAttachment) HostID() HostID {
	return n.hostID
}

func (n networkAttachment) NetworkID() NetworkID {
	return n.networkID
}

func (n networkAttachment) NicName() string {
	return n.nicName
}

func (n networkAttachment) Host(retries ...RetryStrategy) (Host, error) {
	return n.client.GetHost(n.hostID, retries...)
}

func (n networkAttachment) Network(retries ...RetryStrategy) (Network, error) {
	return n.client.GetNetwork(n.networkID, retries...)
}
