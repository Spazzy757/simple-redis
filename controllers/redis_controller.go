/*
Copyright 2023.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	simplev1 "github.com/spazzy757/simple-redis/api/v1"
	iredis "github.com/spazzy757/simple-redis/internal/redis"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=simple.simple.redis,resources=redis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=simple.simple.redis,resources=redis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=simple.simple.redis,resources=redis/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=v1,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=service/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	log.Info("starting reconciliation")

	// fetch redis resource
	var sr simplev1.Redis
	if err := r.Get(ctx, req.NamespacedName, &sr); err != nil {
		log.Error(err, "unable to fetch redis")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var errors error
	if err := r.updateStatus(ctx, sr, simplev1.StatusPending); err != nil {
		errors = multierror.Append(errors, err)
	}
	if err := r.reconcileMasterDeploy(ctx, req, sr); err != nil {
		log.V(1).Error(err, "failed reconciling master deployment")
		errors = multierror.Append(errors, err)
	}

	if err := r.reconcileMasterSvc(ctx, req, sr); err != nil {
		log.V(1).Error(err, "failed reconciling master service")
		errors = multierror.Append(errors, err)
	}

	if err := r.reconcileReplicaDeploy(ctx, req, sr); err != nil {
		log.V(1).Error(err, "failed reconciling replica deployment")
		errors = multierror.Append(errors, err)
	}

	// used to update redis status and handle various errors that could occur in
	// isolation

	if errors != nil {
		if err := r.updateStatus(ctx, sr, simplev1.StatusFailed); err != nil {
			errors = multierror.Append(errors, err)
		}
	} else {
		if err := r.updateStatus(ctx, sr, simplev1.StatusSuccess); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	log.Info("finished reconciliation")
	return ctrl.Result{}, errors
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&simplev1.Redis{}).
		Owns(&v1.Pod{}).
		Complete(r)
}

// updateStatus used to update the status of the Redis instance
func (r *RedisReconciler) updateStatus(ctx context.Context, sr simplev1.Redis, status simplev1.Status) error {
	sr.Status.Status = status
	if err := r.Status().Update(ctx, &sr); err != nil {
		return err
	}
	return nil
}

// reconcileMasterDeploy used to reconcile the master redis instance deployment
func (r *RedisReconciler) reconcileMasterDeploy(ctx context.Context, req ctrl.Request, sr simplev1.Redis) error {
	// master has a single replica for now as multi master would be a future
	// iteration
	// TODO allow multi master setup
	deploy := iredis.GenerateRedisDeploy(sr.Name, req.Namespace, "master", 1)
	if err := controllerutil.SetControllerReference(&sr, deploy, r.Scheme); err != nil {
		return err
	}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Deployment",
		Version: "v1",
	})
	lookup := types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}
	err := r.Get(ctx, lookup, u)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(ctx, deploy)
	} else if err == nil {
		err = r.Update(ctx, deploy)
	}
	return err
}

// reconcileMasterSvc used to reconcile the master redis instance service
func (r *RedisReconciler) reconcileMasterSvc(ctx context.Context, req ctrl.Request, sr simplev1.Redis) error {
	svc := iredis.GenerateRedisSvc(sr.Name, req.Namespace, "master")
	if err := controllerutil.SetControllerReference(&sr, svc, r.Scheme); err != nil {
		return err
	}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Service",
		Version: "v1",
	})
	lookup := types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}
	err := r.Get(ctx, lookup, u)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(ctx, svc)
	} else if err == nil {
		err = r.Update(ctx, svc)
	}
	return err
}

// reconcileReplicaDeploy used to reconcile the master redis instance deployment
func (r *RedisReconciler) reconcileReplicaDeploy(ctx context.Context, req ctrl.Request, sr simplev1.Redis) error {
	replicas := sr.Spec.ClusterSize - 1
	// in the case that cluster size only has 1 or less instance
	// we would set replicas to 0
	if replicas < 0 {
		replicas = 0
	}
	deploy := iredis.GenerateRedisDeploy(sr.Name, req.Namespace, "replica", replicas)
	if err := controllerutil.SetControllerReference(&sr, deploy, r.Scheme); err != nil {
		return err
	}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Deployment",
		Version: "v1",
	})
	lookup := types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}
	err := r.Get(ctx, lookup, u)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(ctx, deploy)
	} else if err == nil {
		err = r.Update(ctx, deploy)
	}
	return err
}
