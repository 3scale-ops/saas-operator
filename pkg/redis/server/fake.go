package server

import "github.com/3scale/saas-operator/pkg/redis/client"

// NewFakeServerWithFakeClient returns a fake server with a fake client that will return the
// provided responses when called. This is only intended for testing.
func NewFakeServerWithFakeClient(host, port string, responses ...client.FakeResponse) *Server {
	rsp := []client.FakeResponse{}
	return &Server{
		host:   host,
		port:   port,
		client: &client.FakeClient{Responses: append(rsp, responses...)},
	}
}
