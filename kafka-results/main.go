package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/IBM/sarama"
)

type StudentResult struct {
    Student    string `json:"student"`
    Age       uint32 `json:"age"`
    Faculty   string `json:"faculty"`
    Discipline uint32 `json:"discipline"`
    Result    string `json:"result"`
    Timestamp string `json:"timestamp"`
}

var producer sarama.SyncProducer

func handleResult(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var result StudentResult
    if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Generate Kafka message
    msg := &sarama.ProducerMessage{
        Topic: "student-results",
        Value: sarama.StringEncoder(r.Body),
        Timestamp: time.Now(),
    }

    // Send to Kafka
    _, _, err := producer.SendMessage(msg)
    if err != nil {
        log.Printf("Failed to send to kafka: %v", err)
        http.Error(w, "Failed to process result", http.StatusInternalServerError)
        return
    }

    log.Printf("Processed result for student: %s (Discipline: %d, Result: %s)",
        result.Student, result.Discipline, result.Result)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Result processed successfully",
    })
}

func main() {
    // Configure Kafka
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    // Connect to Kafka
    var err error
    producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
    if err != nil {
        log.Fatalf("Failed to start Kafka producer: %v", err)
    }
    defer producer.Close()

    // Set up HTTP server
    http.HandleFunc("/result", handleResult)

    log.Println("Kafka results service starting on :8082...")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        log.Fatal(err)
    }
}