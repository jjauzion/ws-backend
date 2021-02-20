package pkg

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	pebkac = func(s string, errCode codes.Code) error { return status.Error(errCode, s) }

	errAnonUser                 = pebkac("anon is forbidden from doing this", codes.PermissionDenied)
	errBadReq                   = pebkac("bad request data", codes.InvalidArgument)
	errBadUserIDHeader          = pebkac("missing or bad '"+grpcUserID+"' head key", codes.InvalidArgument)
	errExpectedNonEmptyID       = pebkac("expected non-empty id", codes.InvalidArgument)
	errForbidden                = pebkac("forbidden", codes.PermissionDenied)
	errNonUniqueEmail           = pebkac("this email is already registered", codes.AlreadyExists)
	errUserNotFound             = pebkac("User not found", codes.NotFound)
	errNoTasksInQueue           = pebkac("no tasks in queue", codes.InvalidArgument)
	errUnsetCropSelectorVersion = pebkac("unset crop selector version", codes.Internal)
)

func exactlyOneOf(fields string) error {
	return pebkac("exactly one of "+fields+" must be set", codes.InvalidArgument)
}
