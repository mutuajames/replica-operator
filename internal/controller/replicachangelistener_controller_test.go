package controller

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func Test_CheckForChangeInReplicaNumber(t *testing.T) {
	r := &ReplicaChangeListenerReconciler{}

	type args struct {
		updateEvent event.UpdateEvent
	}

	var oldReplicaCount *int32 = new(int32)
	var newReplicaCount *int32 = new(int32)

	*newReplicaCount = 2
	*oldReplicaCount = 3

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Should return false if replica number has not changed", args: args{
			updateEvent: event.UpdateEvent{
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: newReplicaCount,
					},
				},
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: newReplicaCount,
					},
				},
			},
		}, want: false},
		{name: "Should return true if replica number has changed", args: args{
			updateEvent: event.UpdateEvent{
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: newReplicaCount,
					},
				},
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: newReplicaCount,
					},
				},
			},
		}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.CheckForChangeInReplicaNumber(tt.args.updateEvent); got != tt.want {
				t.Errorf("filterEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
