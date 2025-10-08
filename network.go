package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Network" -n "network" -T NetworkID

// NetworkID is the UUID if the network.
type NetworkID string

// NetworkClient describes the functions related to oVirt networks.
//
// See https://www.ovirt.org/documentation/administration_guide/#chap-Logical_Networks for details.
type NetworkClient interface {
	// GetNetwork returns a single network based on its ID.
	GetNetwork(id NetworkID, retries ...RetryStrategy) (Network, error)
	// ListNetworks returns all networks on the oVirt engine.
	ListNetworks(retries ...RetryStrategy) ([]Network, error)
	// CreateNetwork creates a new network in the specified data center with the given parameters.
	CreateNetwork(dataCenterId DatacenterID, name string, description string, comment string, vlanID int, retries ...RetryStrategy) (Network, error)
	// RemoveNetwork removes the specified network.
	RemoveNetwork(id NetworkID, retries ...RetryStrategy) error
	// UpdateNetwork updates the specified network with the given parameters.
	UpdateNetwork(id NetworkID, dataCenterId DatacenterID, name string, description string, comment string, vlanID int, retries ...RetryStrategy) (Network, error)
}

// NetworkData is the core of Network, providing only the data access functions, but not the client
// functions.
type NetworkData interface {
	// ID returns the auto-generated identifier for this network.
	ID() NetworkID
	// Name returns the user-give nname for this network.
	Name() string
	// DatacenterID is the identifier of the datacenter object.
	DatacenterID() DatacenterID
	// Description returns the description of this network.
	Description() string
	// Comment returns the comment of this network.
	Comment() string
	// VlanID returns the VLAN ID of this network.
	VlanID() int
}

// Network is the interface defining the fields for networks.
type Network interface {
	NetworkData

	// Datacenter fetches the datacenter associated with this network. This is a network call and may be slow.
	Datacenter(retries ...RetryStrategy) (Datacenter, error)
}

func CreateNetworkParams()

func convertSDKNetwork(sdkObject *ovirtsdk4.Network, client *oVirtClient) (Network, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("network", "id")
	}
	name, ok := sdkObject.Name()
	if !ok {
		return nil, newFieldNotFound("network", name)
	}
	dc, ok := sdkObject.DataCenter()
	if !ok {
		return nil, newFieldNotFound("network", "datacenter")
	}
	dcID, ok := dc.Id()
	if !ok {
		return nil, newFieldNotFound("datacenter on network", "ID")
	}
	description, ok := sdkObject.Description()
	if !ok {
		description = ""
	}
	comment, ok := sdkObject.Comment()
	if !ok {
		comment = ""
	}
	vlan, ok := sdkObject.Vlan()
	var vlanID int
	if ok && vlan != nil {
		vlanID64, ok := vlan.Id()
		if ok {
			vlanID = int(vlanID64)
		}
	} else {
		vlanID = 0
	}
	return &network{
		client:      client,
		id:          NetworkID(id),
		name:        name,
		dcID:        DatacenterID(dcID),
		description: description,
		comment:     comment,
		vlanID:      vlanID,
	}, nil
}

type network struct {
	client Client

	id   NetworkID
	name string
	dcID DatacenterID
	// description is an optional description of the network.
	description string
	// comment is an optional comment for the network.
	comment string
	// vlanID is the VLAN ID of the network. 0 means no VLAN.
	vlanID int
}

func (n network) ID() NetworkID {
	return n.id
}

func (n network) Name() string {
	return n.name
}

func (n network) DatacenterID() DatacenterID {
	return n.dcID
}

func (n network) Datacenter(retries ...RetryStrategy) (Datacenter, error) {
	return n.client.GetDatacenter(n.dcID, retries...)
}

func (n network) Description() string {
	return n.description
}

func (n network) Comment() string {
	return n.comment
}

func (n network) VlanID() int {
	return n.vlanID
}
