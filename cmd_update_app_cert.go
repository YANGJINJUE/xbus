package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/google/subcommands"
	"github.com/infrmods/xbus/apps"
	"github.com/infrmods/xbus/utils"
)

// KeyCertCmd key cert cmd
type UpdateAppCert struct {
}

// Name cmd name
func (cmd *UpdateAppCert) Name() string {
	return "update-app-cert"
}

// Synopsis cmd synopsis
func (cmd *UpdateAppCert) Synopsis() string {
	return "get app's update-app-cert"
}

// Usage cmd usage
func (cmd *UpdateAppCert) Usage() string {
	return "update-app-cert [OPTIONS] app"
}

// SetFlags cmd set flags
func (cmd *UpdateAppCert) SetFlags(f *flag.FlagSet) {
}

// Execute cmd execute
func (cmd *UpdateAppCert) Execute(_ context.Context, f *flag.FlagSet, v ...interface{}) subcommands.ExitStatus {
	x := NewXBus()
	db := x.NewDB()
	appCtrl := x.NewAppCtrl(db, x.Config.Etcd.NewEtcdClient())
	privKey, err := utils.NewPrivateKey("", 2048)
	if err != nil {
		glog.Errorf("generate private key fail: %v", err)
		return subcommands.ExitFailure
	}
	appList, err := apps.GetAppList(db)
	if err != nil {
		glog.Errorf("update-app-cert fail: %v", err)
		return subcommands.ExitFailure
	}
	for _, app := range appList {
		if _, err := appCtrl.NewApp(&app, privKey, nil, nil, 3650, true); err != nil {
			glog.Errorf("create app fail: %v", err)
			return subcommands.ExitFailure
		}
		glog.Info("app: " + app.Name + " update success")
	}
	glog.Info("update-app-cert all success")
	return subcommands.ExitSuccess
}
