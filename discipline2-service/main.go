package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "discipline2-service/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStudentServiceServer
}

func (s *server) ProcessDiscipline2Student(ctx context.Context, student *pb.Student) (*pb.Response, error) {
	log.Printf("Received Discipline 2 Student: Name=%s, Age=%d, Faculty=%s",
		student.Student, student.Age, student.Faculty)
	
	return &pb.Response{
		Message: fmt.Sprintf("Successfully processed discipline 2 student: %s", student.Student),
	}, nil
}

func main() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStudentServiceServer(s, &server{})
	
	log.Printf("gRPC server listening on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}