from locust import HttpUser, task, between
import json
import random
import time

class StudentUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        # Separate test cases by faculty
        self.ingenieria_cases = [
            {
                "student": "Alice Even Discipline 1",
                "age": 22,
                "faculty": "Ingenieria",
                "discipline": 1
            },
            {
                "student": "Bob Odd Discipline 1",
                "age": 33,
                "faculty": "Ingenieria",
                "discipline": 1
            },
            {
                "student": "Charlie Even Discipline 2",
                "age": 20,
                "faculty": "Ingenieria",
                "discipline": 2
            },
            {
                "student": "Daisy Odd Discipline 2",
                "age": 25,
                "faculty": "Ingenieria",
                "discipline": 2
            },
            {
                "student": "Evan Even Discipline 3",
                "age": 28,
                "faculty": "Ingenieria",
                "discipline": 3
            },
            {
                "student": "Fay Odd Discipline 3",
                "age": 35,
                "faculty": "Ingenieria",
                "discipline": 3
            }
        ]
        
        self.agronomia_cases = [
            {
                "student": "Ana Even Discipline 1",
                "age": 20,
                "faculty": "Agronomia",
                "discipline": 1
            },
            {
                "student": "Bruno Odd Discipline 1",
                "age": 21,
                "faculty": "Agronomia",
                "discipline": 1
            },
            {
                "student": "Carla Even Discipline 2",
                "age": 22,
                "faculty": "Agronomia",
                "discipline": 2
            },
            {
                "student": "David Odd Discipline 2",
                "age": 23,
                "faculty": "Agronomia",
                "discipline": 2
            },
            {
                "student": "Elena Even Discipline 3",
                "age": 24,
                "faculty": "Agronomia",
                "discipline": 3
            },
            {
                "student": "Felix Odd Discipline 3",
                "age": 25,
                "faculty": "Agronomia",
                "discipline": 3
            }
        ]

    @task(1)
    def health_check_ingenieria(self):
        self.client.get("/ingenieria/health")

    @task(1)
    def health_check_agronomia(self):
        self.client.get("/agronomia/health")

    @task(3)
    def add_ingenieria_student(self):
        test_case = random.choice(self.ingenieria_cases).copy()
        test_case["student"] = f"{test_case['student']}_{int(time.time())}"
        
        self.client.post(
            "/ingenieria/add_student",
            json=test_case,
            headers={"Content-Type": "application/json"}
        )

    @task(3)
    def add_agronomia_student(self):
        test_case = random.choice(self.agronomia_cases).copy()
        test_case["student"] = f"{test_case['student']}_{int(time.time())}"
        
        self.client.post(
            "/agronomia/add_student",
            json=test_case,
            headers={"Content-Type": "application/json"}
        )