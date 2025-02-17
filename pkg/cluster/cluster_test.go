package cluster

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCluster_NewSnapshot(t *testing.T) {
	tests := []struct {
		name string
		want *Snapshot
	}{
		{
			name: "NewSnapshot",
			want: &Snapshot{
				Cluster: *mockCluster,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mockCluster
			got := c.NewSnapshot()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.NewSnapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClusterName(t *testing.T) {
	type args struct {
		nodes []v1.Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "getClusterName",
			args: args{
				nodes: []v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "node1",
							Labels: map[string]string{
								"Cluster": "test",
							},
						},
					},
				},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClusterName(tt.args.nodes); got != tt.want {
				t.Errorf("getClusterName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAwsCreds(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "NewAwsCreds",
			args: args{
				region: "us-east-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAwsCreds(tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAwsCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
