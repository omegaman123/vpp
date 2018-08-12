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
//

package cache

import (
	"github.com/contiv/vpp/plugins/crd/cache/telemetrymodel"
	"github.com/ligato/cn-infra/logging"
)

const subnetmask = "/24"
const vppVNI = 10

// ContivTelemetryCache is used for a in-memory storage of K8s State data
// The cache processes K8s State data updates and RESYNC events through Update()
// and Resync() APIs, respectively.
// The cache allows to get notified about changes via convenient callbacks.
type ContivTelemetryCache struct {
	Deps
	Synced bool
	// todo - here add the maps you have in your db implementation
	VppCache  VppCache
	K8sCache  *K8sCache
	Processor Processor
	Report    map[string][]string
}

// Deps lists dependencies of PolicyCache.
type Deps struct {
	Log logging.Logger
}

// Init initializes policy cache.
func (ctc *ContivTelemetryCache) Init() error {
	// todo - here initialize your maps
	ctc.VppCache = NewVppCache(ctc.Log)
	ctc.K8sCache = NewK8sCache(ctc.Log)
	ctc.Log.Infof("ContivTelemetryCache has been initialized")
	return nil
}

// ListAllVppNodes returns node data for all nodes in the cache.
func (ctc *ContivTelemetryCache) ListAllVppNodes() []*telemetrymodel.Node {
	nodeList := ctc.VppCache.RetrieveAllNodes()
	return nodeList
}

// LookupVppNodes return node data for nodes that match a node name passed
// to the function in the node names slice.
func (ctc *ContivTelemetryCache) LookupVppNodes(nodenames []string) []*telemetrymodel.Node {
	nodeslice := make([]*telemetrymodel.Node, 0)
	for _, name := range nodenames {
		node, err := ctc.VppCache.RetrieveNode(name)
		if err != nil {
			continue
		}
		nodeslice = append(nodeslice, node)
	}
	return nodeslice
}

// DeleteVppNode deletes from the cache those nodes that match a node name passed
// to the function in the node names slice.
func (ctc *ContivTelemetryCache) DeleteVppNode(nodenames []string) {
	for _, str := range nodenames {
		node, err := ctc.VppCache.RetrieveNode(str)
		if err != nil {
			ctc.Log.Error(err)
			return
		}
		ctc.VppCache.DeleteNode(node.Name)
	}
}

// AddVppNode will add a vpp node to the Contiv Telemetry cache with
// the given parameters.
func (ctc *ContivTelemetryCache) AddVppNode(ID uint32, nodeName, IPAdr, ManIPAdr string) error {
	return ctc.VppCache.CreateNode(ID, nodeName, IPAdr, ManIPAdr)
}

// ClearCache with clear all Contiv Telemetry cache data except for the
// data discovered from etcd updates.
func (ctc *ContivTelemetryCache) ClearCache() {
	ctc.VppCache.ClearCache()
	// TODO: clear k8s cache
	ctc.Report = make(map[string][]string)
}

// ReinitializeCache completely re-initializes the Contiv Telemetry cache,
// clearing all data, including discovered vpp and k8s nodes and discovered
// k8s pods.
func (ctc *ContivTelemetryCache) ReinitializeCache() {
	ctc.VppCache.ReinitializeCache()
	// TODO: re-initialize k8s cache
	ctc.Report = make(map[string][]string)
}

func (ctc *ContivTelemetryCache) logErrAndAppendToNodeReport(nodeName string, errString string) {
	ctc.appendToNodeReport(nodeName, errString)
	ctc.Log.Errorf(errString)
}

func (ctc *ContivTelemetryCache) appendToNodeReport(nodeName string, errString string) {
	ctc.Report[nodeName] = append(ctc.Report[nodeName], errString)
}
