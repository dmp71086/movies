package api

import (
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RPCValidationError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		return status.Error(codes.InvalidArgument, valErr.Error())
	}

	return err
}
