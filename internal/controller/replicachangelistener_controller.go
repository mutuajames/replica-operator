/*
Copyright 2023 mutuajames.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1 "k8s.io/api/apps/v1"
)

// ReplicaChangeListenerReconciler reconciles a ReplicaChangeListener object
type ReplicaChangeListenerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ReplicaChangeListenerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

func (r *ReplicaChangeListenerReconciler) FilterEvents() predicate.Predicate {
	return predicate.Funcs{CreateFunc: func(event event.CreateEvent) bool {
		return false
	}, UpdateFunc: func(updateEvent event.UpdateEvent) bool {
		return r.CheckForChangeInReplicaNumber(updateEvent)
	}, DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
		return false
	}, GenericFunc: func(genericEvent event.GenericEvent) bool {
		return false
	}}
}

func (r *ReplicaChangeListenerReconciler) CheckForChangeInReplicaNumber(event event.UpdateEvent) bool {
	newDeployment := event.ObjectNew.(*appsv1.Deployment)
	oldDeployment := event.ObjectOld.(*appsv1.Deployment)

	if *newDeployment.Spec.Replicas != *oldDeployment.Spec.Replicas {
		fmt.Printf("Replica number for %s changed from %d to %d\n", newDeployment.Name, *oldDeployment.Spec.Replicas, *newDeployment.Spec.Replicas)
		return true
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReplicaChangeListenerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}, builder.WithPredicates(r.FilterEvents())).
		Complete(r)
}
