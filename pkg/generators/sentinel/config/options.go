package config

import (
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
)

// Options holds configuration for the pods
type Options struct{}

// NewOptions returns an Options struct for the given saasv1alpha1.SentinelSpec
func NewOptions(spec saasv1alpha1.SentinelSpec) Options {
	return Options{}
}
