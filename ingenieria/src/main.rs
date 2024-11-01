use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use tonic::Request;

pub mod proto {
    tonic::include_proto!("student");
}
use proto::{student_service_client::StudentServiceClient, Student};

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

async fn send_to_grpc(student: &StudentInfo) -> Result<String, Box<dyn std::error::Error>> {
    let mut client = StudentServiceClient::connect("http://discipline2-service:50051").await?;

    let request = Request::new(Student {
        student: student.student.clone(),
        age: student.age as u32,
        faculty: student.faculty.clone(),
        discipline: student.discipline as u32,
    });

    let response = client.process_discipline2_student(request).await?;
    Ok(response.into_inner().message)
}

#[post("/add_student")]
async fn add_student(student_info: web::Json<StudentInfo>) -> impl Responder {
    if student_info.faculty != "Ingenieria" && student_info.faculty != "Agronomia" {
        return HttpResponse::BadRequest().body("Invalid faculty type");
    }

    if ![1, 2, 3].contains(&student_info.discipline) {
        return HttpResponse::BadRequest().body("Invalid discipline type");
    }

    println!(
        "Received data: student={}, age={}, faculty={}, discipline={}",
        student_info.student, student_info.age, student_info.faculty, student_info.discipline
    );

    // If discipline is 2, send to gRPC service
    if student_info.discipline == 2 {
        match send_to_grpc(&student_info).await {
            Ok(response) => {
                println!("gRPC Service Response: {}", response);
            }
            Err(e) => {
                println!("Error sending to gRPC service: {}", e);
                return HttpResponse::InternalServerError().body("Error processing discipline 2 student");
            }
        }
    }

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