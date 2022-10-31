package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"github.com/phayes/freeport"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func NewTestPortForwarder(cfg *rest.Config, key types.NamespacedName, localPort, podPort uint32,
	out io.Writer, stopCh chan struct{}, readyCh chan struct{}) (*portforward.PortForwarder, error) {

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		key.Namespace, key.Name)
	hostIP := strings.TrimPrefix(cfg.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(cfg)
	if err != nil {
		return nil, err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, podPort)}, stopCh, readyCh, out, out)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func PortForward(cfg *rest.Config, key types.NamespacedName, podPort uint32) (uint32, chan struct{}, error) {

	localPort, err := freeport.GetFreePort()
	if err != nil {
		return 0, nil, err
	}

	errCh := make(chan error)
	stopCh := make(chan struct{})
	readyCh := make(chan struct{})

	go func() {
		defer GinkgoRecover()
		fw, err := NewTestPortForwarder(cfg, key, uint32(localPort), podPort, GinkgoWriter, stopCh, readyCh)
		if err != nil {
			errCh <- err
		}
		err = fw.ForwardPorts()
		if err != nil {
			errCh <- err
		}
	}()

	for {
		select {
		case err := <-errCh:
			return 0, nil, err
		case <-readyCh:
			return uint32(localPort), stopCh, nil
		}
	}
}
