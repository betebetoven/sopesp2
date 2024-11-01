package main

import (
	"context"
	"log"
	"time"

	pb "discipline2-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewStudentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	student := &pb.Student{
		Student:    "Test Student",
		Age:       20,
		Faculty:   "Ingenieria",
		Discipline: 2,
	}

	r, err := c.ProcessDiscipline2Student(ctx, student)
	if err != nil {
		log.Fatalf("could not process student: %v", err)
	}
	log.Printf("Response: %s", r.Message)
}