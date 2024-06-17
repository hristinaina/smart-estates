import random
from datetime import datetime, timedelta
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS

INFLUX_TOKEN = "Ws3mz5dBgA-fOQF1F3WJnpdMaRmbCtEXoMlSCur-gL7oteVeFzknIuyonFmX5hMA9GgjmBs0Ahidfv6orkCy-w=="
ORG = "Smart Home"
BUCKET = "bucket"
URL = "http://localhost:8086"

# Initialize the InfluxDB client
influxdb_client = InfluxDBClient(url=URL, token=INFLUX_TOKEN, org=ORG)

def generate_data_gates(device_id, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=SYNCHRONOUS)
    current_time = start_time

    while current_time <= end_time:
        action = "exit" if random.randint(0, 1) == 0 else "enter"
        success = random.choice([True, False])
        license_plate = random.choice(["NS-123-45", "NS-456-22", "NS-222-34"])

        point = (
            Point("gates")
            .tag("Id", str(device_id))
            .tag("Action", action)
            .tag("Success", str(success))
            .field("LicensePlate", license_plate)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 10 seconds
        current_time += timedelta(hours=12)

def generate_data_sprinkler(device_id, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=SYNCHRONOUS)
    current_time = start_time

    modes = ["on", "off"]
    users = ["kvucic6@gmail.com", "auto"]

    while current_time <= end_time:
        mode = random.choice(modes)
        user = random.choice(users)

        point = (
            Point("sprinkler")
            .tag("device_id", str(device_id))
            .field("action", mode)
            .field("user_id", user)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 6 hours
        current_time += timedelta(hours=6)

def generate_lamp_data(device_id, start_time, end_time):
    write_api = influxdb_client.write_api(write_options=SYNCHRONOUS)
    current_time = start_time

    while current_time <= end_time:
        percentage = random.uniform(0, 100)
        point = (
            Point("lamp")
            .tag("Id", str(device_id))
            .tag("DeviceName", "Lampica u sobici")
            .field("LicensePlate", percentage)
            .time(current_time, WritePrecision.NS)
        )

        write_api.write(bucket=BUCKET, org=ORG, record=point)
        print(f"Written point: {point} at {current_time}")

        # Increment time by 1 hour
        current_time += timedelta(seconds=10)

def populate_vg():
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    generate_data_gates(10, start_time, end_time)

    # Close the client
    influxdb_client.close()


def delete_data(start_time, stop_time, measurement, tag_key, tag_value):
    delete_api = influxdb_client.delete_api()
    predicate = f'_measurement="{measurement}" AND device_id="7"'
    delete_api.delete(
        start=start_time,
        stop=stop_time,
        bucket=BUCKET,
        org=ORG,
        predicate=predicate
    )


def populate_sprinkler():
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    generate_data_sprinkler(11, start_time, end_time)

    influxdb_client.close()


def populate_lamp():
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    generate_lamp_data(7, start_time, end_time)

    influxdb_client.close()


if __name__ == '__main__':
    # # Define the time range for the data to be deleted
    # stop_time = "2024-06-16T00:00:00Z"
    # start_time = "1970-01-01T00:00:00Z"  # From the beginning of time

    # # Call the function to delete data
    # delete_data(start_time=start_time, stop_time=stop_time, measurement="sprinkler", tag_key="device_id", tag_value="10")

    # # Close the client
    # influxdb_client.close()
    populate_lamp()
