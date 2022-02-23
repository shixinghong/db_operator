package main

import (
	"github.com/hongshixing/db_operator/pkg/controller"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/hongshixing/db_operator/pkg/apis/dbconfig/v1"
	"github.com/hongshixing/db_operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
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

	dbConfigController := controller.NewDbConfigController()
	if err = builder.ControllerManagedBy(mgr).
		For(&v1.DbConfig{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}},
			handler.Funcs{
				DeleteFunc: dbConfigController.OnDelete,
				UpdateFunc: dbConfigController.OnUpdate,
			}).
		Complete(dbConfigController); err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}

}
