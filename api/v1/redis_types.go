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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resource status enum
type Status string

const (
	StatusPending Status = "Pending"
	StatusFailed  Status = "Failed"
	StatusSuccess Status = "Success"
)

// redis log levels enum
type RedisLogLevel string

const (
	RLogLevelDebug   RedisLogLevel = "debug"
	RLogLevelVerbose RedisLogLevel = "verbose"
	RLogLevelNotice  RedisLogLevel = "notice"
	RLogLevelWarning RedisLogLevel = "warning"
)

// RedisSpec defines the desired state of Redis
type RedisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ClusterSize determines the amount of redis instances running
	ClusterSize int `json:"clusterSize,omitempty"`

	// LogLevel specifies the redis verbosity level.
	// This can be one of:
	// debug (a lot of information, useful for development/testing)
	// verbose (many rarely useful info, but not a mess like the debug level)
	// notice (moderately verbose, what you want in production probably)
	// warning (only very important / critical messages are logged)
	LogLevel RedisLogLevel `json:"logLevel,omitempty"`

	// Set the number of databases. The default database is DB 0, you can select
	// a different one on a per-connection basis using SELECT <dbid> where
	// dbid is a number between 0 and 'databases'-1
	Databases int `json:"databases,omitempty"`
}

// RedisStatus defines the observed state of Redis
type RedisStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// status of redis cluster
	Status Status `json:"status,omitempty"`

	// master pod name
	Master string `json:"master,omitempty"`
}

// Redis is the Schema for the redis API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSpec   `json:"spec,omitempty"`
	Status RedisStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RedisList contains a list of Redis
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}
