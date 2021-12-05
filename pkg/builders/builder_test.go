package builders

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"reflect"
	"testing"
)

func TestDeployBuilder_Build(t *testing.T) {
	builder, err := NewDeployBuilder("test", "default")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(builder.Replicas(10).Build())
}

func TestDeployBuilder_Replicas(t *testing.T) {
	type fields struct {
		Deployment *appsv1.Deployment
	}
	type args struct {
		r int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *DeployBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &DeployBuilder{
				Deployment: tt.fields.Deployment,
			}
			if got := b.Replicas(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replicas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeployBuilder(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    *DeployBuilder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeployBuilder(tt.args.name, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeployBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployBuilder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
