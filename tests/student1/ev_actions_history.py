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
    def get_ev_actions(self):
        data = {
            "DeviceId": 5,
            "UserEmail": "all",
            "StartDate": "2024-01-16T12:34:56Z",
            "EndDate": "2024-02-15T14:34:56Z"
        }
        response = self.client.put(
            "/api/ev/actions",
            headers={"Authorization": f"Bearer {self.token}"},
            json=data
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")
