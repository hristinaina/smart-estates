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
    def get_rs_consumption(self):
        rs_consumption = {
            "type": "rs",
            "selectedOptions": ["Neka kuca nmp"],
            "start": "2024-02-15T12:34:56Z",
            "end": "2024-01-16T12:34:56Z",
            "queryType": "consumption",
            "batteryId": "6"
        }
        response = self.client.post(
            "/api/consumption/selected-date",
            headers={"Authorization": f"Bearer {self.token}"},
            json=rs_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")


    @task
    def get_rs_production(self):
        rs_consumption = {
            "type": "rs",
            "selectedOptions": ["Neka kuca nmp"],
            "start": "2024-02-15T12:34:56Z",
            "end": "2024-01-16T12:34:56Z",
            "queryType": "solar_panel",
            "batteryId": "6"
        }
        response = self.client.post(
            "/api/consumption/selected-date",
            headers={"Authorization": f"Bearer {self.token}"},
            json=rs_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")
