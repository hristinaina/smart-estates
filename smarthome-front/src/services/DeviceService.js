class DeviceService {

  async get(id) {
    try {
      const response = await fetch(`http://localhost:8081/api/devices/${id}`);
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }

  async getDevices(realEstateId) {
    try {
      const response = await fetch('http://localhost:8081/api/devices/estate/' + realEstateId, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });
      const data = await response.json();
      if (data != null)
        return this.replaceTypeWithString(data);
      else return [];
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }

  replaceTypeWithString(data) {
    for (let i = 0; i < data.length; i++) {
      let d = data[i];
      if (d.Type === 0) d.Type = 'Ambient Sensor'
      else if (d.Type === 1) d.Type = 'Air conditioner'
      else if (d.Type === 2) d.Type = 'Washing machine'
      else if (d.Type === 3) d.Type = 'Lamp'
      else if (d.Type === 4) d.Type = 'Vehicle gate'
      else if (d.Type === 5) d.Type = 'Sprinkler'
      else if (d.Type === 6) d.Type = 'Solar panel'
      else if (d.Type === 7) d.Type = 'Battery storage'
      else if (d.Type === 8) d.Type = 'Electric vehicle charger'
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
        credentials: 'include',
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Unknown error occurred");
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data:', error.message);
      if (error.message == "Unexpected end of JSON input") return null;
      throw error;
    }
  }

  async getSPGraphData(deviceId, email, startDate, endDate) {
    const gdata = {
      "DeviceId": deviceId,
      "UserEmail": email,
      "StartDate": startDate,
      "EndDate": endDate
    }
    try {
      const response = await fetch('http://localhost:8081/api/sp/graphData', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(gdata),
        credentials: 'include',
      });
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }

  async getACHistoryData(deviceId, email, startDate, endDate) {
    const hdata = {
      "DeviceId": deviceId,
      "UserEmail": email,
      "StartDate": startDate,
      "EndDate": endDate
    }
    try {
      const response = await fetch('http://localhost:8081/api/ac/history', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(hdata),
        credentials: 'include',
      });
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }

  async getDeviceById(id, path) {
    try {
        const response = await fetch(path + id, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
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

