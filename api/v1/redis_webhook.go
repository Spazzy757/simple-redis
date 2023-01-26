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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var redislog = logf.Log.WithName("redis-resource")

func (r *Redis) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-simple-simple-redis-v1-redis,mutating=true,failurePolicy=fail,sideEffects=None,groups=simple.simple.redis,resources=redis,verbs=create;update,versions=v1,name=mredis.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Redis{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Redis) Default() {
	redislog.Info("default", "name", r.Name)

	// defaults redis logs level to notice
	if r.Spec.LogLevel == "" {
		r.Spec.LogLevel = RLogLevelNotice
	}

	// dont currently support scaling down to 0
	if r.Spec.ClusterSize == 0 {
		r.Spec.ClusterSize = 1
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-simple-simple-redis-v1-redis,mutating=false,failurePolicy=fail,sideEffects=None,groups=simple.simple.redis,resources=redis,verbs=create;update,versions=v1,name=vredis.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Redis{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateCreate() error {
	redislog.Info("validate create", "name", r.Name)

	return r.validateRedis()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateUpdate(old runtime.Object) error {
	redislog.Info("validate update", "name", r.Name)
	return r.validateRedis()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateDelete() error {
	redislog.Info("validate delete", "name", r.Name)
	// we dont currently need any validation of deletion
	return nil
}

// validateRedis used to run spec through validation
func (r *Redis) validateRedis() error {
	var allErrs field.ErrorList
	if err := r.validateClusterSize(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateLogLevel(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "simple", Kind: "Redis"},
		r.Name,
		allErrs,
	)
}

// validateClusterSize used to validate that cluster size is never set to below
// 0
func (r *Redis) validateClusterSize() *field.Error {
	if r.Spec.ClusterSize <= 0 {
		return field.Invalid(
			field.NewPath("spec").Child("clusterSize"),
			r.Name,
			"invalid cluster size",
		)
	}
	return nil
}

// validateLogLevel used to validate that log level is one of the
// correct levels
func (r *Redis) validateLogLevel() *field.Error {
	switch r.Spec.LogLevel {
	case RLogLevelDebug, RLogLevelNotice, RLogLevelWarning, RLogLevelVerbose:
		return nil
	}
	return field.Invalid(
		field.NewPath("spec").Child("clusterSize"),
		r.Name,
		"logLevel needs to be one of [debug,notice,verbose,warning]",
	)
}
