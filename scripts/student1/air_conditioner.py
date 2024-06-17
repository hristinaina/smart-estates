import random
from datetime import datetime, timedelta
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS

INFLUX_TOKEN = "UsAetsOrZuhhTCu7lHiVkF46qbD2h55Tkv0eNyCKx9H5IxF0IyNEkNTMHigzIiZVWaCgSoCoMwoQnk38lnwbXw=="
ORG = "Smart Home"
BUCKET = "bucket"

url = "http://localhost:8086"
influxdb_client = InfluxDBClient(url=url, token=INFLUX_TOKEN, org=ORG)

def save_ac_to_influxdb(device_id, mode, previous_mode, user, switch_ac, timestamp):
    write_api = influxdb_client.write_api(write_options=SYNCHRONOUS)
    action = 0 if not switch_ac else 1
    
    point = (
        Point("air_conditioner1")
        .tag("device_id", str(device_id))
        .field("action", action)
        .field("mode", mode)
        .field("user_id", user)
        .time(timestamp, WritePrecision.NS)
    )
    
    write_api.write(bucket=BUCKET, org=ORG, record=point)
    write_api.flush()

    if previous_mode:
        # Use the same timestamp for the previous mode
        point = (
            Point("air_conditioner1")
            .tag("device_id", str(device_id))
            .field("action", 0)
            .field("mode", previous_mode)
            .field("user_id", "auto")
            .time(timestamp, WritePrecision.NS)
        )
        
        write_api.write(bucket=BUCKET, org=ORG, record=point)
        write_api.flush()
    
    # print(f"Air Conditioner data written to InfluxDB at {timestamp}")

def generate_ac_data(device_id, start_time, end_time):
    modes = ['Cooling', 'Automatic', 'Ventilation']
    users = ["Natasa Maric", "auto"]

    current_time = start_time
    while current_time <= end_time:
        for _ in range(2):  # Generate two records per day
            mode = random.choice(modes)
            previous_mode = random.choice(modes) if random.choice([True, False]) else ""
            user = random.choice(users)
            switch_ac = random.choice([True, False])

            save_ac_to_influxdb(device_id, mode, previous_mode, user, switch_ac, current_time)
            
            # Increment time randomly within the day
            current_time += timedelta(hours=random.randint(1, 12))
        current_time = datetime.combine(current_time.date() + timedelta(days=1), datetime.min.time())  # Move to the next day

if __name__ == "__main__":
    # Define the start and end time for the data
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(days=90)  # 3 months ago

    # Generate and save random data for a specific device
    generate_ac_data(3, start_time, end_time)
    # generate_ac_data(16, start_time, end_time)

    # Close the client
    influxdb_client.close()