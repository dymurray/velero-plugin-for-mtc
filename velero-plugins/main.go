package main

import (
	"github.com/migtools/velero-plugin-for-mtc/velero-plugins/common"
	"github.com/sirupsen/logrus"
	veleroplugin "github.com/vmware-tanzu/velero/pkg/plugin/framework"
)

func main() {
	veleroplugin.NewServer().
		RegisterRestoreItemAction("mtc.openshift.io/01-common-restore-plugin", newCommonRestorePlugin).
		Serve()
}

func newCommonRestorePlugin(logger logrus.FieldLogger) (interface{}, error) {
	return &common.RestorePlugin{Log: logger}, nil
}
