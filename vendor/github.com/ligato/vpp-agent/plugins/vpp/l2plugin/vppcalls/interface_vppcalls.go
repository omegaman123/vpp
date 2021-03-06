// Copyright (c) 2017 Cisco and/or its affiliates.
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

package vppcalls

import (
	"fmt"
	"time"

	"github.com/ligato/cn-infra/logging"
	l2ba "github.com/ligato/vpp-agent/plugins/vpp/binapi/l2"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	"github.com/ligato/vpp-agent/plugins/vpp/model/l2"
)

// SetInterfacesToBridgeDomain implements bridge domain handler.
func (handler *BridgeDomainVppHandler) SetInterfacesToBridgeDomain(bdName string, bdIdx uint32, bdIfs []*l2.BridgeDomains_BridgeDomain_Interfaces,
	swIfIndices ifaceidx.SwIfIndex) (ifs []string, wasErr error) {
	defer func(t time.Time) {
		handler.stopwatch.TimeLog(l2ba.SwInterfaceSetL2Bridge{}).LogTimeEntry(time.Since(t))
	}(time.Now())

	if len(bdIfs) == 0 {
		handler.log.Debugf("Bridge domain %v has no new interface to set", bdName)
		return nil, nil
	}

	for _, bdIf := range bdIfs {
		// Verify that interface exists, otherwise skip it.
		ifIdx, _, found := swIfIndices.LookupIdx(bdIf.Name)
		if !found {
			handler.log.Debugf("Required bridge domain %v interface %v not found", bdName, bdIf.Name)
			continue
		}
		if err := handler.addDelInterfaceToBridgeDomain(bdName, bdIdx, bdIf, ifIdx, true); err != nil {
			wasErr = err
			handler.log.Error(wasErr)
		} else {
			handler.log.WithFields(logging.Fields{"Interface": bdIf.Name, "BD": bdName}).Debug("Interface set to bridge domain")
			ifs = append(ifs, bdIf.Name)
		}
	}

	return ifs, wasErr
}

// UnsetInterfacesFromBridgeDomain implements bridge domain handler.
func (handler *BridgeDomainVppHandler) UnsetInterfacesFromBridgeDomain(bdName string, bdIdx uint32, bdIfs []*l2.BridgeDomains_BridgeDomain_Interfaces,
	swIfIndices ifaceidx.SwIfIndex) (ifs []string, wasErr error) {

	defer func(t time.Time) {
		handler.stopwatch.TimeLog(l2ba.SwInterfaceSetL2Bridge{}).LogTimeEntry(time.Since(t))
	}(time.Now())

	if len(bdIfs) == 0 {
		handler.log.Debugf("Bridge domain %v has no obsolete interface to unset", bdName)
		return nil, nil
	}

	for _, bdIf := range bdIfs {
		// Verify that interface exists, otherwise skip it.
		ifIdx, _, found := swIfIndices.LookupIdx(bdIf.Name)
		if !found {
			handler.log.Debugf("Required bridge domain %v interface %v not found", bdName, bdIf.Name)
			continue
		}
		if err := handler.addDelInterfaceToBridgeDomain(bdName, bdIdx, bdIf, ifIdx, false); err != nil {
			wasErr = err
			handler.log.Error(wasErr)
		} else {
			handler.log.WithFields(logging.Fields{"Interface": bdIf.Name, "BD": bdName}).Debug("Interface unset from bridge domain")
			ifs = append(ifs, bdIf.Name)
		}
	}

	return ifs, wasErr
}

func (handler *BridgeDomainVppHandler) addDelInterfaceToBridgeDomain(bdName string, bdIdx uint32, bdIf *l2.BridgeDomains_BridgeDomain_Interfaces,
	ifIdx uint32, add bool) error {
	req := &l2ba.SwInterfaceSetL2Bridge{
		BdID:        bdIdx,
		RxSwIfIndex: ifIdx,
		Shg:         uint8(bdIf.SplitHorizonGroup),
		Enable:      boolToUint(add),
	}
	// Set as BVI.
	if bdIf.BridgedVirtualInterface {
		req.Bvi = 1
		handler.log.Debugf("Interface %v set as BVI", bdIf.Name)
	}
	reply := &l2ba.SwInterfaceSetL2BridgeReply{}

	if err := handler.callsChannel.SendRequest(req).ReceiveReply(reply); err != nil {
		return fmt.Errorf("error while assigning/removing interface %v to bd %v: %v", bdIf.Name, bdName, err)
	} else if reply.Retval != 0 {
		return fmt.Errorf("%s returned %d while assigning/removing interface %v (idx %v) to bd %v",
			reply.GetMessageName(), reply.Retval, bdIf.Name, ifIdx, bdName)
	}

	return nil
}
