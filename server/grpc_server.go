package server

import (
	"context"
	"github.com/jjauzion/ws-backend/db"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/jjauzion/ws-backend/proto"
)

type grpcServer struct {
	pb.UnimplementedApiServer
}

func (s *grpcServer) StartTask(context.Context, *pb.StartTaskReq) (*pb.StartTaskRep, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (s *grpcServer) EndTask(context.Context, *pb.EndTaskReq) (*pb.EndTaskRep, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func RunGRPC(bootstrap bool) {
	lg, cf, dbh, err := dependencies()
	if err != nil {
		return
	}

	if bootstrap {
		if err := db.Bootstrap(dbh); err != nil {
			lg.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	port := ":" + cf.WS_GRPC_PORT
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lg.Info("grpc server listening on", zap.String("port", port))
	s := grpc.NewServer()
	pb.RegisterApiServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
