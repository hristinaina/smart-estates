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


def generate_data_ev_charger(deviceId, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time
    actions = ['start', 'end', 'percentageChange']

    while current_time <= end_time:
        action = random.choice(actions)
        plug_id = random.randint(0, 1)
        percentage = random.uniform(0.2, 0.9)
        if action != 'percentageChange':
            user = 'auto'
        else:
            user = 'nata@gmail.com'
        point = (
            Point("ev_charger")
            .tag("device_id", deviceId)
            .tag("user_id", user)
            .tag("plug_id", plug_id)
            .tag("action", action)
            .field("value", percentage)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 12 hours
        current_time += timedelta(hours=12)


if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    generate_data_ev_charger(13, start_time, end_time)

    # Close the client
    influxdb_client.close()
