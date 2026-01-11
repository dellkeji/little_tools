package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceSpec defines the desired state of Service
type ServiceSpec struct {
	// Image is the container image to deploy
	Image string `json:"image"`

	// Replicas is the number of desired replicas
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`

	// Port is the port the service listens on
	// +kubebuilder:default=8080
	Port int32 `json:"port,omitempty"`

	// ServiceType is the type of Kubernetes service
	// +kubebuilder:default="ClusterIP"
	ServiceType string `json:"serviceType,omitempty"`

	// ConfigData contains configuration data for the service
	ConfigData map[string]string `json:"configData,omitempty"`

	// Ingress configuration
	Ingress *IngressSpec `json:"ingress,omitempty"`

	// Environment variables
	Env []EnvVar `json:"env,omitempty"`

	// Resource requirements
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

// IngressSpec defines ingress configuration
type IngressSpec struct {
	// Enabled indicates if ingress should be created
	Enabled bool `json:"enabled"`

	// Host is the hostname for the ingress
	Host string `json:"host,omitempty"`

	// Path is the path for the ingress rule
	// +kubebuilder:default="/"
	Path string `json:"path,omitempty"`

	// TLS configuration
	TLS *TLSSpec `json:"tls,omitempty"`

	// Annotations for the ingress
	Annotations map[string]string `json:"annotations,omitempty"`
}

// TLSSpec defines TLS configuration for ingress
type TLSSpec struct {
	// Enabled indicates if TLS should be enabled
	Enabled bool `json:"enabled"`

	// SecretName is the name of the TLS secret
	SecretName string `json:"secretName,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ResourceRequirements defines resource requirements
type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// ServiceStatus defines the observed state of Service
type ServiceStatus struct {
	// Phase represents the current phase of the service
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ReadyReplicas is the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// URL is the external URL if ingress is enabled
	URL string `json:"url,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
//+kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Service is the Schema for the services API
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceList contains a list of Service
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Service{}, &ServiceList{})
}
