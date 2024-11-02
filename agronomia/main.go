package main

import (
    "context"
    "encoding/json"
    "log"
    "fmt"
    "net/http"
    "time"

    pb "agronomia/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// StudentInfo represents the student data structure
type StudentInfo struct {
    Student    string `json:"student"`
    Age       uint8  `json:"age"`
    Faculty   string `json:"faculty"`
    Discipline uint8  `json:"discipline"`
}

// healthCheck handles the health check requests
func healthCheck(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
}

func sendToGrpc(student StudentInfo) error {
    // Select the correct gRPC service based on discipline
    var serviceURL string
    switch student.Discipline {
    case 1:
        serviceURL = "discipline1-service:50051"
    case 2:
        serviceURL = "discipline2-service:50051"
    case 3:
        serviceURL = "discipline3-service:50051"
    default:
        return fmt.Errorf("invalid discipline: %d", student.Discipline)
    }

    // Create gRPC connection
    conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return err
    }
    defer conn.Close()

    // Create gRPC client
    c := pb.NewStudentServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Convert StudentInfo to protobuf Student
    pbStudent := &pb.Student{
        Student:    student.Student,
        Age:        uint32(student.Age),
        Faculty:    student.Faculty,
        Discipline: uint32(student.Discipline),
    }

    // Call gRPC service
    response, err := c.ProcessDiscipline2Student(ctx, pbStudent) // Assuming this method name is used across services
    if err != nil {
        return err
    }

    log.Printf("gRPC Service Response: %s", response.Message)
    return nil
}


func addStudent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var student StudentInfo
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if student.Faculty != "Ingenieria" && student.Faculty != "Agronomia" {
        http.Error(w, "Invalid faculty type", http.StatusBadRequest)
        return
    }

    if student.Discipline < 1 || student.Discipline > 3 {
        http.Error(w, "Invalid discipline type", http.StatusBadRequest)
        return
    }

    log.Printf("Received data: student=%s, age=%d, faculty=%s, discipline=%d",
        student.Student, student.Age, student.Faculty, student.Discipline)

    // If discipline is 2, send to gRPC service
    if student.Discipline == 2 || student.Discipline == 3|| student.Discipline == 1 {
        if err := sendToGrpc(student); err != nil {
            log.Printf("Error sending to gRPC service: %v", err)
            http.Error(w, "Error processing discipline 2 student", http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Student information received successfully",
    })
}

func main() {
    http.HandleFunc("/health", healthCheck)
    http.HandleFunc("/add_student", addStudent)

    log.Println("Server starting on added the validations :8081...")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}