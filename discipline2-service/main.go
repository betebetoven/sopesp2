package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "time"

    pb "discipline2-service/proto"
    "google.golang.org/grpc"
    "cloud.google.com/go/pubsub"
)

type server struct {
    pb.UnimplementedStudentServiceServer
    pubsubClient *pubsub.Client
    topic        *pubsub.Topic
}

type StudentResult struct {
    Student    string    `json:"student"`
    Age       uint32    `json:"age"`
    Faculty   string    `json:"faculty"`
    Discipline uint32    `json:"discipline"`
    Result    string    `json:"result"`
    Timestamp time.Time `json:"timestamp"`
}

func determineWinner(student *pb.Student) string {
    // Simple algorithm: based on student's age
    // If age is even -> winner, if odd -> loser
    if student.Age%2 == 0 {
        return "winner"
    }
    return "loser"
}

func (s *server) ProcessDiscipline2Student(ctx context.Context, student *pb.Student) (*pb.Response, error) {
    result := determineWinner(student)
    
    // Prepare student result
    studentResult := StudentResult{
        Student:    student.Student,
        Age:       student.Age,
        Faculty:   student.Faculty,
        Discipline: student.Discipline,
        Result:    result,
        Timestamp: time.Now(),
    }

    // Convert to JSON
    jsonData, err := json.Marshal(studentResult)
    if err != nil {
        log.Printf("Error marshaling student result: %v", err)
        return nil, err
    }

    // Publish to Pub/Sub
    msg := &pubsub.Message{
        Data: jsonData,
    }

    id, err := s.topic.Publish(ctx, msg).Get(ctx)
    if err != nil {
        log.Printf("Failed to publish message: %v", err)
        return nil, err
    }

    log.Printf("Published message with ID: %s", id)
    log.Printf("Processed student: %s (Result: %s)", student.Student, result)

    return &pb.Response{
        Message: fmt.Sprintf("Successfully processed discipline 2 student: %s (Result: %s)", student.Student, result),
    }, nil
}

func main() {
    ctx := context.Background()

    // Create Pub/Sub client
    pubsubClient, err := pubsub.NewClient(ctx, "servicio-440415")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer pubsubClient.Close()

    // Get topic
    topic := pubsubClient.Topic("student-results")

    // Start gRPC server
    port := 50051
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterStudentServiceServer(s, &server{
        pubsubClient: pubsubClient,
        topic:       topic,
    })
    
    log.Printf("gRPC server listening on port %d", port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}