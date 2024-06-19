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
    def get_solar_panel_actions(self):
        ambient_sensor = {
            "DeviceId": 4,
            "UserEmail": "nata@gmail.com",
            "StartDate": "2024-01-16T12:34:56Z",
            "EndDate": "2024-02-15T14:34:56Z"
        }
        response = self.client.put(
            "/api/sp/graphData",
            headers={"Authorization": f"Bearer {self.token}"},
            json=ambient_sensor
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")
