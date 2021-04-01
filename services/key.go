package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/infrmods/xbus/utils"
)

var rValidName = regexp.MustCompile(`(?i)^[a-z][a-z0-9_.-]{5,}$`)
var rValidService = regexp.MustCompile(`(?i)^[a-z][a-z0-9_.-]{5,}:[a-z0-9][a-z0-9_.-]*$`)
var rValidZone = regexp.MustCompile(`(?i)^[a-z0-9][a-z0-9_-]{3,}$`)
var rValidExt = regexp.MustCompile(`(?i)^[a-z0-9][a-z0-9_-]{3,16}$`)

func checkName(name string) error {
	if !rValidName.MatchString(name) {
		return utils.NewError(utils.EcodeInvalidName, "")
	}
	return nil
}

func checkService(service string) error {
	if !rValidService.MatchString(service) {
		return utils.NewError(utils.EcodeInvalidService, "")
	}
	return nil
}

func checkServiceZone(service, zone string) error {
	if !rValidService.MatchString(service) {
		return utils.NewError(utils.EcodeInvalidService, "")
	}
	if !rValidZone.MatchString(zone) {
		return utils.NewError(utils.EcodeInvalidZone, "")
	}
	return nil
}

var rValidAddress = regexp.MustCompile(`(?i)^[a-z0-9:_.-]+$`)

func (ctrl *ServiceCtrl) checkAddress(addr string) error {
	if addr == "" {
		return utils.NewError(utils.EcodeInvalidEndpoint, "missing address")
	}
	if !rValidAddress.MatchString(addr) {
		return utils.NewError(utils.EcodeInvalidAddress, "")
	}
	if ctrl.config.isAddressBanned(addr) {
		return utils.NewError(utils.EcodeInvalidAddress, "banned")
	}
	return nil
}

func (ctrl *ServiceCtrl) serviceEntryPrefix(name string) string {
	return fmt.Sprintf("%s/%s/", ctrl.config.KeyPrefix, name)
}

func (ctrl *ServiceCtrl) serviceZoneKey(service string, zone string) string {
	return fmt.Sprintf("%s/%s", service, zone)
}

const serviceDescNodeKey = "desc"

func (ctrl *ServiceCtrl) serviceDescKey(service, zone string) string {
	return fmt.Sprintf("%s/%s/%s/desc", ctrl.config.KeyPrefix, service, zone)
}

func (ctrl *ServiceCtrl) isServiceDescKey(key string) bool {
	return strings.HasSuffix(key, "/"+serviceDescNodeKey)
}

const serviceKeyNodePrefix = "node_"

func (ctrl *ServiceCtrl) serviceNodeKey(service, zone, addr string) string {
	return fmt.Sprintf("%s/%s/%s/node_%s", ctrl.config.KeyPrefix, service, zone, addr)
}

func (ctrl *ServiceCtrl) serviceDescNotifyKey(service, zone string) string {
	return fmt.Sprintf("%s-descs/%s/%s", ctrl.config.KeyPrefix, zone, service)
}

func (ctrl *ServiceCtrl) serviceM5NotifyKey(service, zone string) string {
	return fmt.Sprintf("%s-md5s/%s/%s", ctrl.config.KeyPrefix, zone, service)
}

func (ctrl *ServiceCtrl) serviceM5NotifyPrefix(zone string) string {
	if zone != "" {
		return fmt.Sprintf("%s-md5s/%s/", ctrl.config.KeyPrefix, zone)
	}
	return fmt.Sprintf("%s-md5s/", ctrl.config.KeyPrefix)
}

func (ctrl *ServiceCtrl) serviceDescNotifyKeyPrefix(zone string) string {
	if zone != "" {
		return fmt.Sprintf("%s-descs/%s/", ctrl.config.KeyPrefix, zone)
	}
	return fmt.Sprintf("%s-descs/", ctrl.config.KeyPrefix)
}

type serviceDescKey struct {
	service string
	zone    string
}

func (ctrl *ServiceCtrl) splitServiceDescNotifyKey(key string) *serviceDescKey {
	prefix := ctrl.config.KeyPrefix + "-descs/"
	if strings.HasPrefix(key, prefix) {
		parts := strings.Split(key[len(prefix):], "/")
		if len(parts) == 2 {
			return &serviceDescKey{zone: parts[0], service: parts[1]}
		}
	}
	return nil
}

func (ctrl *ServiceCtrl) splitServiceM5NotifyKey(key string) *serviceDescKey {
	prefix := ctrl.config.KeyPrefix + "-md5s/"
	if strings.HasPrefix(key, prefix) {
		parts := strings.Split(key[len(prefix):], "/")
		if len(parts) == 2 {
			return &serviceDescKey{zone: parts[0], service: parts[1]}
		}
	}
	return nil
}
