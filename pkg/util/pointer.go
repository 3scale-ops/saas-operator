package util

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func IntStrPtr(value intstr.IntOrString) *intstr.IntOrString {
	return &value
}

func Metav1DurationPtr(value time.Duration) *metav1.Duration {
	d := &metav1.Duration{Duration: value}
	return d
}

func Pointer[T any](t T) *T {
	return &t
}
