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
    
    test_cases = [
        {
            "student": "John Doe",
            "age": 33,
            "faculty": "Ingenieria",
            "discipline": 2
        },
        {
            "student": "Macaco test de kafka 3",
            "age": 22,
            "faculty": "Agronomia",
            "discipline": 2
        },
        {
            "student": "Alberto Hernandez test de kafka 3",
            "age": 20,
            "faculty": "Ingenieria",
            "discipline": 2
        },
    ]
    
    # Test both services
    for service in ["agronomia", "ingenieria"]:
        test_health_endpoint(base_url, service)
        time.sleep(1)
        
        for test_case in test_cases:
            test_add_student_endpoint(base_url, service, test_case)
            time.sleep(1)

if __name__ == "__main__":
    main()