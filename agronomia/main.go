package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	// Only allow GET method
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

func addStudent(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var student StudentInfo
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate faculty
	if student.Faculty != "Ingenieria" && student.Faculty != "Agronomia" {
		http.Error(w, "Invalid faculty type", http.StatusBadRequest)
		return
	}

	// Validate discipline
	if student.Discipline < 1 || student.Discipline > 3 {
		http.Error(w, "Invalid discipline type", http.StatusBadRequest)
		return
	}

	// Log the received data
	log.Printf("Received data: student=%s, age=%d, faculty=%s, discipline=%d",
		student.Student, student.Age, student.Faculty, student.Discipline)

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student information received successfully",
	})
}

func main() {
	// Register routes
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/add_student", addStudent)

	// Start server
	log.Println("Server starting on :8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}