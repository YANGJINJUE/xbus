package services

import (
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/glog"
	"github.com/infrmods/xbus/utils"
	"strings"
)

func (ctrl *ServiceCtrl) makeService(kvs []*mvccpb.KeyValue) (*Service, error) {
	var service Service
	service.Endpoints = make([]ServiceEndpoint, 0, len(kvs))
	for _, kv := range kvs {
		parts := strings.Split(string(kv.Key), "/")
		name := parts[len(parts)-1]
		if strings.HasPrefix(name, serviceKeyNodePrefix) {
			var endpoint ServiceEndpoint
			if err := json.Unmarshal(kv.Value, &endpoint); err != nil {
				glog.Errorf("unmarshal endpoint fail(%#v): %v", string(kv.Value), err)
				return nil, utils.NewError(utils.EcodeDamagedEndpointValue, "")
			}
			service.Endpoints = append(service.Endpoints, endpoint)
		} else if strings.HasPrefix(name, serviceDescNodeKey) {
			if err := json.Unmarshal(kv.Value, &service.ServiceDesc); err != nil {
				glog.Errorf("invalid desc(%s), unmarshal fail: %v", string(kv.Key), string(kv.Value))
				return nil, utils.NewSystemError("service-data damanged")
			}
		}
	}
	return &service, nil
}
