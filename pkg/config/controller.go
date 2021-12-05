package config

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	"github.com/hongshixing/db_operator/pkg/apis/dbconfig/v1"
	"github.com/hongshixing/db_operator/pkg/builders"

	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	APIVersion      = "v1"
	Kind            = "DbConfig"
	Group           = "api.myit.fun"
	GroupAPIVersion = Group + "/" + APIVersion
)

type DbConfigController struct {
	client.Client
}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}

// OnDelete 当资源被删除时 重建资源
func (dc *DbConfigController) OnDelete(event event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.Object.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupAPIVersion { // deploy被删除 重新回到队列中 重建
			limitingInterface.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: event.Object.GetNamespace(),
					Name:      ref.Name,
				}})
		}
	}
}

// OnUpdate 当资源被更新时 重建资源
func (dc *DbConfigController) OnUpdate(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupAPIVersion { // deploy被删除 重新回到队列中 重建
			fmt.Println("OnUpdate 当资源被更新时 重建资源s")
			limitingInterface.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: event.ObjectNew.GetNamespace(),
					Name:      ref.Name,
				}})
		}
	}
}

func (dc *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = logf.FromContext(ctx)
	db := &v1.DbConfig{}
	if err := dc.Get(ctx, req.NamespacedName, db); err != nil {
		return reconcile.Result{}, err
	}

	builder, err := builders.NewDeployBuilder(db, dc.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = builder.Build(ctx); err != nil { // 创建出deployment
		return reconcile.Result{}, err
	}

	klog.Info(db)

	return reconcile.Result{}, nil
}

func (dc *DbConfigController) InjectClient(c client.Client) error {
	dc.Client = c
	return nil
}
