package mapping

import (
	"context"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
)

// Codec creates a codec that maps method names.
type Codec struct {
	methods map[string]string
}

// New returns a new Mapping codec.
func New(ctx context.Context) *Codec {
	return &Codec{
		methods: make(map[string]string),
	}
}

// Add adds a mapping.
func (c *Codec) Add(from string, to string) {
	c.methods[from] = to
}

// NewRequest returns a new CodecRequest of type MappingRequest.
func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	outerCR := &MappingRequest{
		methods: c.methods,
	}
	jsonC := json.NewCodec()
	innerCR := jsonC.NewRequest(r)
	outerCR.CodecRequest = innerCR.(*json.CodecRequest)

	return outerCR
}

// MappingRequest decodes and encodes a single request. MappingCodecRequest
// implements gorilla/rpc.CodecRequest interface primarily by embedding
// the CodecRequest from gorilla/rpc/json. By selectively adding
// CodecRequest methods to TRCodecRequest, we can modify that behaviour
// while maintaining all the other remaining CodecRequest methods from
// gorilla's rpc/json implementation
type MappingRequest struct {
	*json.CodecRequest
	methods map[string]string
}

// Method returns the decoded method as a string of the form "Service.Method"
// after checking for, and correcting a lowercase method name
// By being of lower depth in the struct , Method will replace the implementation
// of Method() on the embedded CodecRequest. Because the request data is part
// of the embedded json.CodecRequest, and unexported, we have to get the
// requested method name via the embedded CR's own method Method().
// Essentially, this just intercepts the return value from the embedded
// gorilla/rpc/json.CodecRequest.Method(), checks/modifies it, and passes it
// on to the calling rpc server.
func (r *MappingRequest) Method() (string, error) {
	m, err := r.CodecRequest.Method()
	if err != nil {
		return "", err
	}

	if _, exists := r.methods[m]; exists {
		return r.methods[m], nil
	}

	return m, nil
}
