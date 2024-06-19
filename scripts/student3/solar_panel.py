import math
import random
from datetime import datetime, timedelta
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import ASYNCHRONOUS

INFLUX_TOKEN = "5JrdbOYRvz9k03QXYOEO5LvHkB9pBVDRimwkr26IDRd9aDDvsHxBno4H3e412w7DZ--Rk8Ltgny2MURH4OUd5A=="
ORG = "Smart Home"
BUCKET = "bucket"

url = "http://localhost:8086"
influxdb_client = InfluxDBClient(url=url, token=INFLUX_TOKEN, org=ORG)


def generate_data_solar_panel(deviceId, start_time, end_time, battery_id, estate_id):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time

    while current_time <= end_time:
        temperature = 25 + random.uniform(-5, 5)  # Generate temperature between 10 and 20 degrees

        point = (
            Point("solar_panel")
            .tag("device_id", deviceId)
            .tag("battery_id", battery_id)
            .tag("estate_id", estate_id)
            .field("electricity", temperature)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 10 seconds
        current_time += timedelta(hours=1)

def generate_data_solar_panel_switch(deviceId, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time
    user = 'anastasijas557@gmail.com'
    while current_time <= end_time:
        isOn = random.randint(0, 1)

        point = (
            Point("solar_panel")
            .tag("device_id", deviceId)
            .tag("user_id", user)
            .field("isOn", isOn)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 10 seconds
        current_time += timedelta(hours=12)


if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    #generate_data_solar_panel(4, start_time, end_time, 6, 2)
    #generate_data_solar_panel(12, start_time, end_time, -1, 3)
    generate_data_solar_panel_switch(12, start_time, end_time)

    # Close the client
    influxdb_client.close()
