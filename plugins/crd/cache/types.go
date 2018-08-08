// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"github.com/contiv/vpp/plugins/ksr/model/node"
	pod2 "github.com/contiv/vpp/plugins/ksr/model/pod"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/vpp-agent/plugins/vpp/model/interfaces"
)

// here goes different cache types
//Update this whenever a new DTO type is added.
const numDTOs = 5

//Node is a struct to hold all relevant information of a kubernetes node.
//It is populated with various information such as the interfaces and L2Fibs
//as well as the name and IP Addresses.
type Node struct {
	ID                uint32
	IPAdr             string
	ManIPAdr          string
	Name              string
	NodeLiveness      *NodeLiveness
	NodeInterfaces    map[int]NodeInterface
	NodeBridgeDomains map[int]NodeBridgeDomain
	NodeL2Fibs        map[string]NodeL2FibEntry
	NodeTelemetry     map[string]NodeTelemetry
	NodeIPArp         []NodeIPArpEntry
	report            []string
	podMap            map[string]*pod2.Pod
}

//Cache holds various maps which all take different keys but point to the same underlying value.
type Cache struct {
	nMap       map[string]*Node
	loopIPMap  map[string]*Node
	gigEIPMap  map[string]*Node
	loopMACMap map[string]*Node
	k8sNodeMap map[string]*node.Node
	hostIPMap  map[string]*Node
	podMap     map[string]*pod2.Pod
	report     []string

	logger logging.Logger
}

//NodeLiveness holds the unmarshalled node liveness JSON data
type NodeLiveness struct {
	BuildVersion string `json:"build_version"`
	BuildDate    string `json:"build_date"`
	State        uint32 `json:"state"`
	StartTime    uint32 `json:"start_time"`
	LastChange   uint32 `json:"last_change"`
	LastUpdate   uint32 `json:"last_update"`
	CommitHash   string `json:"commit_hash"`
}

// NodeDTO holds generic node information to be sent over a channel and associated with a name in the cache.
type NodeDTO struct {
	NodeName string
	NodeInfo interface{}
	err      error
}

type nodeInterfaces map[int]NodeInterface
type nodeBridgeDomains map[int]NodeBridgeDomain
type nodeL2FibTable map[string]NodeL2FibEntry
type nodeTelemetries map[string]NodeTelemetry
type nodeIPArpTable []NodeIPArpEntry

//NodeTelemetry holds the unmarshalled node telemetry JSON data
type NodeTelemetry struct {
	Command string   `json:"command"`
	Output  []output `json:"output"`
}

type output struct {
	command string
	output  []outputEntry
}

type outputEntry struct {
	nodeName string
	count    int
	reason   string
}

//NodeL2FibEntry holds unmarshalled L2Fib JSON data
type NodeL2FibEntry struct {
	BridgeDomainIdx          uint32 `json:"bridge_domain_idx"`
	OutgoingInterfaceSwIfIdx uint32 `json:"outgoing_interface_sw_if_idx"`
	PhysAddress              string `json:"phys_address"`
	StaticConfig             bool   `json:"static_config"`
	BridgedVirtualInterface  bool   `json:"bridged_virtual_interface"`
}

//NodeInterface holds unmarshalled Interface JSON data
type NodeInterface struct {
	VppInternalName string                   `json:"vpp_internal_name"`
	Name            string                   `json:"name"`
	IfType          interfaces.InterfaceType `json:"type,omitempty"`
	Enabled         bool                     `json:"enabled,omitempty"`
	PhysAddress     string                   `json:"phys_address,omitempty"`
	Mtu             uint32                   `json:"mtu,omitempty"`
	Vxlan           vxlan                    `json:"vxlan,omitempty"`
	IPAddresses     []string                 `json:"ip_addresses,omitempty"`
	Tap             tap                      `json:"tap,omitempty"`
}

type vxlan struct {
	SrcAddress string `json:"src_address"`
	DstAddress string `json:"dst_address"`
	Vni        uint32 `json:"vni"`
}

//NodeIPArpEntry holds unmarshalled IP ARP data
type NodeIPArpEntry struct {
	Interface  uint32 `json:"interface"`
	IPAddress  string `json:"IPAddress"`
	MacAddress string `json:"MacAddress"`
	Static     bool   `json:"Static"`
}

type tap struct {
	Version    uint32 `json:"version"`
	HostIfName string `json:"host_if_name"`
}

//NodeBridgeDomain holds the unmarshalled bridge domain data.
type NodeBridgeDomain struct {
	Interfaces []bdinterfaces `json:"interfaces"`
	Name       string         `json:"name"`
	Forward    bool           `json:"forward"`
}

type bdinterfaces struct {
	SwIfIndex uint32 `json:"sw_if_index"`
}