package util

import "k8s.io/apimachinery/pkg/util/intstr"

func IntStrPtr(value intstr.IntOrString) *intstr.IntOrString {
	return &value
}
