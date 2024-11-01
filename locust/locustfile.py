from locust import HttpUser, task, between

class StudentLoadTest(HttpUser):
    wait_time = between(1, 3)  # Wait time between tasks

    @task
    def add_student_agronomia(self):
        self.client.post("/agronomia/add_student", json={
            "student": "Alvaro Garcia",
            "age": 20,
            "faculty": "Agronomia",
            "discipline": 1
        })

    @task
    def add_student_ingenieria(self):
        self.client.post("/ingenieria/add_student", json={
            "student": "John Doe",
            "age": 22,
            "faculty": "Ingenieria",
            "discipline": 2
        })
