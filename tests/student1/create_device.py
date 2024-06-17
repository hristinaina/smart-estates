from locust import HttpUser, task, between


class MyUser(HttpUser):
    wait_time = between(1, 3)
    host = "http://localhost:8080"

    def on_start(self):
        response = self.client.post("/api/users/login", json={"email": "nata@gmail.com", "password": "Pass!123"})
        if response.status_code == 200:
            self.token = response.json()["token"]
        else:
            print("Failed to log in")
            self.token = None

    @task
    def create_ambient_sensor(self):
        ambient_sensor = {
            "Name": "AmbSenzorcic",
            "Type": 0,
            "RealEstate": 1,
            "PowerSupply": 0
        }
        response = self.client.post(
            "/api/devices/",
            headers={"Authorization": f"Bearer {self.token}"},
            json=ambient_sensor
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")
