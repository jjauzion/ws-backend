package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/jjauzion/ws-backend/conf"
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
	conf conf.Configuration
	dbal db.Dbal
	log  logger.Logger
}

func (srv *grpcServer) RegisterServer(s *grpc.Server) {
	pb.RegisterApiServer(s, srv)
}

func (s *grpcServer) StartTask(ctx context.Context, req *pb.StartTaskReq) (*pb.StartTaskRep, error) {
	s.log.Info("starting StartTask")
	start := time.Now()

	t, err := s.dbal.GetNextTask(ctx)
	if err != nil {
		s.log.Error("", zap.Error(err))
	}
	if t == nil {
		s.log.Info("no task in queue")
		return nil, errNoTasksInQueue
	}
	s.log.Info("oldest task is", zap.Any("task", t))
	err = s.dbal.UpdateTaskStatus(ctx, t.ID, db.StatusRunning)
	if err != nil {
		return nil, err
	}
	var rep *pb.StartTaskRep
	rep = &pb.StartTaskRep{
		Job:    &pb.Job{Dataset: t.Job.Dataset, DockerImage: t.Job.DockerImage},
		TaskId: uuid.New().String(),
	}
	s.log.Info("ended StartTask", zap.Duration("in", time.Since(start)))
	return rep, err
}

func (s *grpcServer) EndTask(context.Context, *pb.EndTaskReq) (*pb.EndTaskRep, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func RunGRPC(bootstrap bool) {
	ctx := context.Background()
	app, _, err := buildApplication()
	if err != nil {
		return
	}

	if bootstrap {
		if err := db.Bootstrap(ctx, app.dbal); err != nil {
			app.log.Error("bootstrap failed", zap.Error(err))
			return
		}
	}

	var srv = grpcServer{conf: app.conf, dbal: app.dbal}

	port := ":" + app.conf.WS_GRPC_PORT
	lis, err := net.Listen("tcp", port)
	if err != nil {
		app.log.Fatal("failed to listen", zap.Error(err))
	}
	app.log.Info("grpc server listening on", zap.String("port", port))
	s := grpc.NewServer()
	defer s.Stop()
	srv.RegisterServer(s)
	if err := s.Serve(lis); err != nil {
		app.log.Fatal("failed to serve", zap.Error(err))
	}
}
