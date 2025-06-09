/*
Copyright 2025.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Deployment struct {
	Image   string            `json:"image"`
	Ports   []int32           `json:"ports"`
	EnvVars map[string]string `json:"envVars,omitempty"`
}

// AutomatedAppDeploymentSpec defines the desired state of AutomatedAppDeployment.
type AutomatedAppDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Replicas    int32        `json:"replicas,omitempty"`
	Deployments []Deployment `json:"deployments"`
}

// AutomatedAppDeploymentStatus defines the observed state of AutomatedAppDeployment.
type AutomatedAppDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AutomatedAppDeployment is the Schema for the automatedappdeployments API.
type AutomatedAppDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutomatedAppDeploymentSpec   `json:"spec,omitempty"`
	Status AutomatedAppDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutomatedAppDeploymentList contains a list of AutomatedAppDeployment.
type AutomatedAppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutomatedAppDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutomatedAppDeployment{}, &AutomatedAppDeploymentList{})
}
