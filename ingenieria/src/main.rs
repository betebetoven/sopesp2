use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
struct StudentInfo {
    student: String,
    age: u8,
    faculty: String,
    discipline: u8,
}

#[derive(Serialize)]
struct HealthResponse {
    status: String,
}

#[derive(Serialize)]
struct SuccessResponse {
    message: String,
}

#[get("/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(HealthResponse {
        status: "healthy".to_string(),
    })
}

#[post("/add_student")]
async fn add_student(student_info: web::Json<StudentInfo>) -> impl Responder {
    // Check if the faculty is valid
    if student_info.faculty != "Ingenieria" && student_info.faculty != "Agronomia" {
        return HttpResponse::BadRequest().body("Invalid faculty type");
    }

    // Check if the discipline is valid
    if ![1, 2, 3].contains(&student_info.discipline) {
        return HttpResponse::BadRequest().body("Invalid discipline type");
    }

    // Process data (e.g., store it or print to console)
    println!(
        "Received data: student={}, age={}, faculty={}, discipline={}",
        student_info.student, student_info.age, student_info.faculty, student_info.discipline
    );

    HttpResponse::Ok().json(SuccessResponse {
        message: "Student information received successfully".to_string(),
    })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Server starting on :8080...");
    HttpServer::new(|| {
        App::new()
            .service(health_check)
            .service(add_student)
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}