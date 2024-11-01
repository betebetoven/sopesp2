import requests
import json
from typing import Dict
import time

def test_health_endpoint(base_url: str) -> None:
    """Test the health check endpoint"""
    try:
        response = requests.get(f"{base_url}/health")
        print("\n=== Health Check Test ===")
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json()}")
        assert response.status_code == 200, "Health check failed!"
        assert response.json()["status"] == "healthy", "Unexpected health check response!"
        print("✅ Health check test passed!")
    except Exception as e:
        print(f"❌ Health check test failed: {str(e)}")

def test_add_student_endpoint(base_url: str, student_data: Dict) -> None:
    """Test the add_student endpoint with provided data"""
    try:
        response = requests.post(
            f"{base_url}/add_student",
            headers={"Content-Type": "application/json"},
            json=student_data
        )
        print(f"\n=== Add Student Test ({student_data['student']}) ===")
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
    base_url = "http://localhost:8080"
    
    test_cases = [
        {
            "student": "John Doe",
            "age": 20,
            "faculty": "Ingenieria",
            "discipline": 1
        },
        {
            "student": "Jane Smith de rust dockerizado",
            "age": 22,
            "faculty": "Agronomia",
            "discipline": 2
        },
        # Invalid faculty test case
        {
            "student": "Invalid Faculty",
            "age": 21,
            "faculty": "Medicine",
            "discipline": 1
        },
        # Invalid discipline test case
        {
            "student": "Invalid Discipline",
            "age": 23,
            "faculty": "Ingenieria",
            "discipline": 5
        }
    ]
    
    test_health_endpoint(base_url)
    time.sleep(1)
    
    for test_case in test_cases:
        test_add_student_endpoint(base_url, test_case)
        time.sleep(1)

if __name__ == "__main__":
    main()