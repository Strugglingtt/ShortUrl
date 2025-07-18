package service

import (
	"context"

	pb "backend-shorturl/api/shorturl/public/v1"
)

type PublicService struct {
	pb.UnimplementedPublicServer
}

func NewPublicService() *PublicService {
	return &PublicService{}
}

func (s *PublicService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}
