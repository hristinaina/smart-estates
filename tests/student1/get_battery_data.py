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
    def get_battery_data(self):
        response = self.client.get(
            "/api/hb/6",
            headers={"Authorization": f"Bearer {self.token}"},
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")

    @task
    def get_battery_last_hour(self):
        response = self.client.get(
            "/api/hb/last-hour/6",
            headers={"Authorization": f"Bearer {self.token}"},
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")

