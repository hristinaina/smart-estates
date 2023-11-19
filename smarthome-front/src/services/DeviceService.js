class DeviceService {

  async getDevices(realEstateId) {
    try {
        const response = await fetch('http://localhost:8081/api/devices/estate/' + realEstateId);
        const data = await response.json();
        return this.replaceTypeWithString(data);
      } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
      }
  }

  replaceTypeWithString(data){
    for (let i = 0; i < data.length; i++){
      let d = data[i];
      if (d.Type === 0) d.Type = 'Ambient Sensor'
      else if(d.Type === 1) d.Type = 'Air conditioner' 
      else if(d.Type === 2) d.Type = 'Washing machine' 
      else if(d.Type === 3) d.Type = 'Lamp' 
      else if(d.Type === 4) d.Type = 'Vehicle gate' 
      else if(d.Type === 5) d.Type = 'Sprinkler' 
      else if(d.Type === 6) d.Type = 'Solar panel' 
      else if(d.Type === 7) d.Type = 'Battery storage' 
      else if(d.Type === 8) d.Type = 'Electric vehicle charger' 
    }
    return data
  }

  async createDevice(device) {
    try {
      const response = await fetch('http://localhost:8081/api/devices/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(device),
      });
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }
}

export default new DeviceService();

