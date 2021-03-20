package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	pebkac = func(s string, errCode codes.Code) error { return status.Error(errCode, s) }

	errBadReq         = pebkac("bad request data", codes.InvalidArgument)
	errForbidden      = pebkac("forbidden", codes.PermissionDenied)
	errNoTasksInQueue = pebkac("no tasks in queue", codes.Code(600))
)

func exactlyOneOf(fields string) error {
	return pebkac("exactly one of "+fields+" must be set", codes.InvalidArgument)
}
