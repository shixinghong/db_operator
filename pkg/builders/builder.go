package builders

import (
	"bytes"
	"context"
	"strconv"
	"text/template"

	v1 "github.com/hongshixing/db_operator/pkg/apis/dbconfig/v1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const deployTemp = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: old-{{ .Name }}
  namespace: {{ .Namespace}}
spec:
  selector:
    matchLabels:
      app: old-{{ .Namespace}}-{{ .Name }}
  replicas: 1
  template:
    metadata:
      labels:
        app: old-{{ .Namespace}}-{{ .Name }}
        version: v1
    spec:
      containers:
        - name: old-{{ .Namespace}}-{{ .Name }}-container
          image: docker.io/shenyisyn/dbcore:v1
          imagePullPolicy: IfNotPresent
          ports:
             - containerPort: 8081
             - containerPort: 8090
`

type DeployBuilder struct {
	Deployment *appsv1.Deployment
	Db         *v1.DbConfig
	client.Client
}

func deployName(name string) string {
	return "old-" + name
}

// NewDeployBuilder build new deploy
func NewDeployBuilder(db *v1.DbConfig, client client.Client) (*DeployBuilder, error) {
	deploy := &appsv1.Deployment{}
	// 判断资源是否已经被创建
	err := client.Get(context.Background(), types.NamespacedName{
		Namespace: db.Namespace,
		Name:      deployName(db.Name), // 注意前缀
	}, deploy)

	if err != nil { // 说明资源没有被创建
		// 创建新的资源
		deploy.Name, deploy.Namespace = db.Name, db.Namespace
		tpl, _ := template.New("deploy").Parse(deployTemp)

		var resTemp bytes.Buffer
		err = tpl.Execute(&resTemp, deploy)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(resTemp.Bytes(), deploy)
		if err != nil {
			return nil, err
		}
	}

	return &DeployBuilder{
		Deployment: deploy,
		Db:         db,
		Client:     client,
	}, nil
}

// Replicas set replicas number
func (b *DeployBuilder) Replicas(r int32) *DeployBuilder {
	b.Deployment.Spec.Replicas = &r
	return b
}

// syncReplicas 同步在crd中设置的replicas
func (b *DeployBuilder) syncReplicas() *DeployBuilder {
	return b.Replicas(b.Db.Spec.Replicas)
}

// setOwnerRef 设置OwnerReferences 在删除crd同时级联删除crd创建的资源
func (b *DeployBuilder) setOwnerRef() *DeployBuilder {
	b.Deployment.OwnerReferences = append(b.Deployment.OwnerReferences, metav1.OwnerReference{
		APIVersion: b.Db.APIVersion,
		Kind:       b.Db.Kind,
		Name:       b.Db.Name,
		UID:        b.Db.UID,
	})
	return b
}

// Build 创建出新的deployment
func (b *DeployBuilder) Build(ctx context.Context) error {
	b.syncReplicas().setOwnerRef()                          // 同步dbconfig里面的replicas属性，同时设置OwnerReferences
	if b.Deployment.ObjectMeta.CreationTimestamp.IsZero() { // 如果是模板创建没有creationTimestamp
		return b.Create(ctx, b.Deployment)
	} else {
		if err := b.Update(ctx, b.Deployment); err != nil {
			return err
		}
		rep := b.Deployment.Status.Replicas
		b.Db.Status.Replicas = rep
		b.Db.Status.Ready = strconv.Itoa(int(rep)) + "/" + strconv.Itoa(int(b.Db.Spec.Replicas))
		//fmt.Println(b.Db.Status.Ready)b
		return b.Status().Update(ctx, b.Db) // 更新status字段中的状态
	}
	//return nil
}
