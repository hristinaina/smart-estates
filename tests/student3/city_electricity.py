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
    def get_city_consumption(self):
        city_consumption = {
            "type": "city",
            "selectedOptions": ["Beograd"],
            "time": "-24h",
            "queryType": "consumption",
            "batteryId": "6"
        }
        response = self.client.post(
            "/api/consumption/selected-time",
            headers={"Authorization": f"Bearer {self.token}"},
            json=city_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")

    @task
    def get_city_production(self):
        city_consumption = {
            "type": "city",
            "selectedOptions": ["Beograd"],
            "time": "-24h",
            "queryType": "solar_panel",
            "batteryId": "6"
        }
        response = self.client.post(
            "/api/consumption/selected-time",
            headers={"Authorization": f"Bearer {self.token}"},
            json=city_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")

    @task
    def get_city_ratio(self):
        city_consumption = {
            "type": "city",
            "selectedOptions": ["Beograd"],
            "time": "-24h",
            "batteryId": ""
        }
        response = self.client.post(
            "/api/consumption/ratio/selected-time",
            headers={"Authorization": f"Bearer {self.token}"},
            json=city_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")

    @task
    def get_city_ed(self):
        city_consumption = {
            "type": "city",
            "selectedOptions": ["Beograd"],
            "time": "-24h",
            "batteryId": "electrical_distribution"
        }
        response = self.client.post(
            "/api/consumption/ratio/selected-time",
            headers={"Authorization": f"Bearer {self.token}"},
            json=city_consumption
        )
        if response.status_code != 200:
            self.environment.runner.quit()
        else:
            print("success")
