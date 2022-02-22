/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeepBlueSpec defines the desired state of DeepBlue
type DeepBlueSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DeepBlue. Edit deepblue_types.go to remove/update
	Foo  string `json:"foo,omitempty"`
	Name string `json:"name,omitempty"`
	Book string `json:"book,omitempty"`
}

// DeepBlueStatus defines the observed state of DeepBlue
type DeepBlueStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,Status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeepBlue is the Schema for the deepblues API
type DeepBlue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeepBlueSpec   `json:"spec,omitempty"`
	Status DeepBlueStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeepBlueList contains a list of DeepBlue
type DeepBlueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeepBlue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeepBlue{}, &DeepBlueList{})
}
