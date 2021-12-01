package main

import (
	"os"

	"github.com/hongshixing/db-operator/pkg/apis/dbconfig/v1"
	"github.com/hongshixing/db-operator/pkg/config"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	logf.SetLogger(zap.New())

	var log = logf.Log.WithName("dbconfig")
	kubeConfig := config.InitK8S()
	mgr, err := manager.New(kubeConfig, manager.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	if err = v1.SchemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "could not add manager")
		os.Exit(1)
	}

	if err = builder.ControllerManagedBy(mgr).For(&v1.DbConfig{}).Complete(config.NewDbConfigController()); err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}
