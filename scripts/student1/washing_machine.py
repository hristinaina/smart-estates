import random
from datetime import datetime, timedelta
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS

INFLUX_TOKEN = "UsAetsOrZuhhTCu7lHiVkF46qbD2h55Tkv0eNyCKx9H5IxF0IyNEkNTMHigzIiZVWaCgSoCoMwoQnk38lnwbXw=="
ORG = "Smart Home"
BUCKET = "bucket"

url = "http://localhost:8086"
influxdb_client = InfluxDBClient(url=url, token=INFLUX_TOKEN, org=ORG)

def generate_and_save_wm_data(device_id, start_time, end_time):
    modes = ['Cotton', 'Synthetic', 'Quick', 'Delicate']
    users = ["Natasa Maric", "auto"]
    write_api = influxdb_client.write_api(write_options=SYNCHRONOUS)

    current_time = start_time
    while current_time <= end_time:
        for _ in range(2):  # Generate two records per day
            mode = random.choice(modes)
            previous_mode = random.choice(modes) if random.choice([True, False]) else ""
            user = random.choice(users)
            switch_wm = random.choice([True, False])
            action = "Turn off" if not switch_wm else "Turn on"

            point = (
                Point("washing_machine")
                .tag("device_id", str(device_id))
                .field("action", action)
                .field("mode", mode)
                .field("user_id", user)
                .time(current_time, WritePrecision.NS)
            )
            write_api.write(bucket=BUCKET, org=ORG, record=point)

            if previous_mode:
                # Use the same timestamp for the previous mode
                point = (
                    Point("washing_machine")
                    .tag("device_id", str(device_id))
                    .field("action", "Turn off")
                    .field("mode", previous_mode)
                    .field("user_id", "auto")
                    .time(current_time, WritePrecision.NS)
                )
                write_api.write(bucket=BUCKET, org=ORG, record=point)
        
        # Increment time randomly within the day
        current_time += timedelta(hours=random.randint(1, 12))
        # Move to the next day if we reach or exceed the end of the current day
        if current_time.time() >= datetime.max.time():
            current_time = datetime.combine(current_time.date() + timedelta(days=1), datetime.min.time())
        
        write_api.flush()
    
    print("Washing Machine data written to InfluxDB")

if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    # Generate and save random data for a specific device
    generate_and_save_wm_data(1, start_time, end_time)
    # generate_and_save_wm_data(15, start_time, end_time)

    # Close the client
    influxdb_client.close()
