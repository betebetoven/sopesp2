import requests
import json
from typing import Dict
import time

def test_health_endpoint(base_url: str, service: str) -> None:
    """Test the health check endpoint"""
    try:
        response = requests.get(f"{base_url}/{service}/health")
        print(f"\n=== Health Check Test ({service}) ===")
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json()}")
        assert response.status_code == 200, "Health check failed!"
        assert response.json()["status"] == "healthy", "Unexpected health check response!"
        print("✅ Health check test passed!")
    except Exception as e:
        print(f"❌ Health check test failed: {str(e)}")

def test_add_student_endpoint(base_url: str, service: str, student_data: Dict) -> None:
    """Test the add_student endpoint with provided data"""
    try:
        response = requests.post(
            f"{base_url}/{service}/add_student",
            headers={"Content-Type": "application/json"},
            json=student_data
        )
        print(f"\n=== Add Student Test ({service} - {student_data['student']}) ===")
        print(f"Request Data: {json.dumps(student_data, indent=2)}")
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json()}")
        
        if response.status_code == 200:
            print("✅ Add student test passed!")
        else:
            print("❌ Add student test failed!")
    except Exception as e:
        print(f"❌ Add student test failed: {str(e)}")

def main():
    base_url = "http://34.57.143.146"  # Your ingress IP
    service = "agronomia"  # Specify the target service
    
    test_cases = [
        # Discipline 1
        {
            "student": "Alice Even Discipline 1",
            "age": 22,
            "faculty": "Agronomía",
            "discipline": 1
        },
        {
            "student": "Bob Odd Discipline 1",
            "age": 33,
            "faculty": "Agronomía",
            "discipline": 1
        },
        # Discipline 2
        {
            "student": "Charlie Even Discipline 2",
            "age": 20,
            "faculty": "Agronomía",
            "discipline": 2
        },
        {
            "student": "Daisy Odd Discipline 2",
            "age": 25,
            "faculty": "Agronomía",
            "discipline": 2
        },
        # Discipline 3
        {
            "student": "Evan Even Discipline 3",
            "age": 28,
            "faculty": "Agronomía",
            "discipline": 3
        },
        {
            "student": "Fay Odd Discipline 3",
            "age": 35,
            "faculty": "Agronomía",
            "discipline": 3
        },
    ]
    
    # Test Agronomia service only
    test_health_endpoint(base_url, service)
    time.sleep(1)
    
    for test_case in test_cases:
        test_add_student_endpoint(base_url, service, test_case)
        time.sleep(1)

if __name__ == "__main__":
    main()
