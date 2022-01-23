package greeter

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"protobuf-http/proto/pb/greeterpb"
	"sync"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func AddRoutes(r *mux.Router, svc Service) {
	r.Methods("POST").Path("/greeter").Handler(kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(*greeterpb.GreeterRequest)
			return svc.SayGreeting(ctx, req)
		},
		decodeSayGreetingRequest,
		encodeProtobufResponse,
	))
}

func decodeSayGreetingRequest(_ context.Context, r *http.Request) (interface{}, error) {

	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufferPool.Put(b)

	if _, err := io.Copy(b, r.Body); err != nil {
		return nil, errors.Wrap(err, "cannot read request body")
	}

	req := &greeterpb.GreeterRequest{}
	if err := proto.Unmarshal(b.Bytes(), req); err != nil {
		return nil, errors.Wrap(err, "request is not a valid protobuf text")
	}

	return req, nil
}

func encodeProtobufResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {

	m, ok := resp.(protoreflect.ProtoMessage)
	if !ok {
		return errors.New("response type is not a protobuf message")
	}

	b, err := proto.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "cannot marshal response protobuf message")
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return errors.Wrap(err, "cannot write protobuf text to response")
	}

	return nil
}
