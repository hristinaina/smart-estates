from datetime import datetime, timedelta
import random
import pytz
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS

url = f"http://localhost:8086"
token = "N7Ndk5eXmrtrxxj_hvMYp9ZOtGPPC7kMn-EQMoVq9Ogc15ZcZ_yamIwyye9tm8W1ESlA_NwZ2ktUh9XhREuTTw=="
org = "Smart Home"
bucket = "bucket"

# client = InfluxDBClient(url=url, token=token, org=org)
# write_api = client.write_api(write_options=SYNCHRONOUS)

utc = pytz.utc

def add_data_for_ambient_sensor(device_id):
    start_time = datetime.now(utc) - timedelta(days=90)  # Start time 90 days ago in UTC
    end_time = datetime.now(utc)  # End time is current UTC time

    current_time = start_time
    while current_time <= end_time:
        for i in range(6):  # Simulira 6 očitavanja u minuti (svakih 10 sekundi)
            timestamp = current_time.isoformat()
            temperature = round(random.uniform(15, 25), 2) 
            humidity = round(random.uniform(40, 80), 2)
            
            point = Point("measurement1") \
                .tag("device_id", str(device_id)) \
                .field("temperature", temperature) \
                .field("humidity", humidity) \
                .time(timestamp)
            
            write_api.write(bucket=bucket, org=org, record=point)
            print(f"Upisana temperatura: {temperature} °C @ {timestamp}")
            print(f"Upisana vlaznost vazduha: {humidity} °C @ {timestamp}")

            current_time += timedelta(seconds=10)  # Povećava vreme za 10 sekundi za sledeće očitavanje
        
        current_time += timedelta(seconds=50)


def write_data_to_influxdb():
    # Povezivanje na InfluxDB
    client = InfluxDBClient(url=url, token="N7Ndk5eXmrtrxxj_hvMYp9ZOtGPPC7kMn-EQMoVq9Ogc15ZcZ_yamIwyye9tm8W1ESlA_NwZ2ktUh9XhREuTTw==", org=org)

    # API za pisanje
    write_api = client.write_api(write_options=SYNCHRONOUS)

    try:
        # Generisanje i upis podataka
        for i in range(10):
            point = Point("measurement") \
                .tag("location", "office") \
                .field("temperature", random.uniform(20, 25)) \
                .field("humidity", random.uniform(50, 60)) \
                .time(datetime.utcnow(), WritePrecision.NS)

            write_api.write(bucket=bucket, org=org, record=point)

        print("Podaci su uspešno upisani u InfluxDB!")

    except Exception as e:
        print(f"Greska pri upisu podataka: {e}")

    finally:
        # Zatvaranje veze
        client.close()


if __name__ == '__main__':
    # add_data_for_ambient_sensor(11)
    write_data_to_influxdb()
    