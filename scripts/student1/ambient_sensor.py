import math
import random
from datetime import datetime, timedelta
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import ASYNCHRONOUS

INFLUX_TOKEN = "UsAetsOrZuhhTCu7lHiVkF46qbD2h55Tkv0eNyCKx9H5IxF0IyNEkNTMHigzIiZVWaCgSoCoMwoQnk38lnwbXw=="
ORG = "Smart Home"
BUCKET = "bucket"

url = "http://localhost:8086"
influxdb_client = InfluxDBClient(url=url, token=INFLUX_TOKEN, org=ORG)

def generate_data_ambient_sensor(deviceId, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=ASYNCHRONOUS)
    current_time = start_time
    
    while current_time <= end_time:
        temperature = 25 + random.uniform(-5, 5)  # Generate temperature between 10 and 20 degrees
        humidity = 35 + random.uniform(-10, 10)  # Generate humidity between 20 and 40 percent

        point = (
            Point("ambient_sensor")
            .tag("device-id", deviceId)
            .field("temperature", temperature)
            .field("humidity", humidity)
            .time(current_time, WritePrecision.NS)
        )
        
        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")
        
        # Increment time by 10 seconds
        current_time += timedelta(seconds=10)

if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    generate_data_ambient_sensor(17, start_time, end_time)
    # generate_data_ambient_sensor(18, start_time, end_time)

    # Close the client
    influxdb_client.close()
