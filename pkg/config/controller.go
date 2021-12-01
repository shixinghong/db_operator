package config

import (
	"context"
	"github.com/hongshixing/db-operator/pkg/apis/dbconfig/v1"
	"k8s.io/klog/v2"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DbConfigController struct {
	client.Client
}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}

func (a *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = logf.FromContext(ctx)
	db := &v1.DbConfig{}
	err := a.Get(ctx, req.NamespacedName, db)
	if err != nil {
		return reconcile.Result{}, err
	}

	klog.Info(db)

	return reconcile.Result{}, nil
}

func (a *DbConfigController) InjectClient(c client.Client) error {
	a.Client = c
	return nil
}
