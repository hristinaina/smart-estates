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


def generate_consumption(start_time, end_time, battery_id):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time

    while current_time <= end_time:
        temperature = 25 + random.uniform(-5, 5)  # Generate temperature between 10 and 20 degrees
        device_id = random.randint(1,3)
        if device_id == 2:
            estate_id = 1
        else:
            estate_id = 2
        point = (
            Point("consumption")
            .tag("device_id", device_id)
            .tag("battery_id", battery_id)
            .tag("estate_id", estate_id)
            .field("electricity", temperature)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 10 seconds
        current_time += timedelta(minutes=20)


def generate_data_home_battery(deviceId, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time

    # Use a fixed period for the sine wave, e.g., 24 hours (86400 seconds)
    period = 60400
    amplitude = 1  # Amplitude of the sine wave
    offset = 1  # Offset to shift the sine wave into the positive range

    while current_time <= end_time:
        value = amplitude * math.sin(2 * math.pi * (current_time.timestamp() % period) / period) + offset

        point = (
            Point("home_battery")
            .tag("device_id", deviceId)
            .field("currentValue", value)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 12 hours
        current_time += timedelta(minutes=2)


if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    #generate_data_home_battery(6, start_time, end_time)
    generate_consumption(start_time, end_time, 6)
    generate_consumption(start_time, end_time, 14)
    generate_data_home_battery(14, start_time, end_time)
    # Close the client
    influxdb_client.close()
