package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/db"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"

	"github.com/jjauzion/ws-backend/internal/logger"
	pb "github.com/jjauzion/ws-backend/proto"
)

type grpcServer struct {
	pb.UnimplementedApiServer
}

func (s *grpcServer) StartTask(ctx context.Context, req *pb.StartTaskReq) (*pb.StartTaskRep, error) {
	log := logger.ProvideLogger()
	log.Info("starting StartTask")
	start := time.Now()
	var err error
	var rep *pb.StartTaskRep
	if req.WithGPU {
		rep = &pb.StartTaskRep{
			Job:    &pb.Job{Dataset: "s3://test-dataset", DockerImage: "docker.io/jjauzion/ws-mock-container"},
			TaskId: uuid.New().String(),
		}
	} else {
		err = errNoTasksInQueue
	}
	log.Info("ended StartTask", zap.Duration("in", time.Since(start)))
	return rep, err
}

func (s *grpcServer) EndTask(context.Context, *pb.EndTaskReq) (*pb.EndTaskRep, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func RunGRPC(bootstrap bool) {
	ctx := context.Background()
	lg, cf, dbal, err := buildDependencies()
	if err != nil {
		return
	}

	if bootstrap {
		if err := db.Bootstrap(ctx, dbal); err != nil {
			lg.Error("bootstrap failed", zap.Error(err))
			return
		}
	}
	t, err := dbal.GetOldestTask(ctx)
	if err != nil {
		lg.Error("", zap.Error(err))
	}
	if t == nil {
		lg.Info("no task in queue")
	} else {
		lg.Info("oldest task is", zap.Any("task", t))
		dbal.UpdateTaskStatus(ctx, t.ID, db.StatusRunning)
	}

	port := ":" + cf.WS_GRPC_PORT
	lis, err := net.Listen("tcp", port)
	if err != nil {
		lg.Fatal("failed to listen", zap.Error(err))
	}
	lg.Info("grpc server listening on", zap.String("port", port))
	s := grpc.NewServer()
	pb.RegisterApiServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		lg.Fatal("failed to serve", zap.Error(err))
	}
}
