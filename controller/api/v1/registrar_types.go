/*
Copyright 2024.

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

// RegistrarSpec defines the desired state of Registrar
type RegistrarSpec struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
	Region     string `json:"region,omitempty"`
}

// RegistrarStatus defines the observed state of Registrar
type RegistrarStatus struct {
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Registrar is the Schema for the registrars API
type Registrar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrarSpec   `json:"spec,omitempty"`
	Status RegistrarStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RegistrarList contains a list of Registrar
type RegistrarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Registrar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Registrar{}, &RegistrarList{})
}
