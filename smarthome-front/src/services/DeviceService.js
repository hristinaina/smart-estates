class DeviceService {

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

  async getSPLastValue(deviceId) {
    try {
      const response = await fetch('http://localhost:8081/api/sp/lastValue/' + deviceId, {
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

  async getSPById(id) {
    try {
      const response = await fetch('http://localhost:8081/api/sp/' + id, {
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

  async getHB(id) {
    try {
      const response = await fetch('http://localhost:8081/api/hb/' + id, {
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

  async getHBGraphDataByRS(id) {
    try {
      const response = await fetch('http://localhost:8081/api/hb/last-hour/' + id, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching data: ', error);
      throw error;
    }
  }

  async getGraphDataForDropdownSelect(estateId, time) {
    try {
      const response = await fetch('http://localhost:8081/api/hb/selected-time/' + estateId, {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({ time })
      })
      // console.log(response)

      if (response.ok) {
        const data = await response.json();
        // console.log(data)
        return { result: data };
      } else {
        const data = await response.json();
        return { result: false, error: data.error };
      }
    } catch (error) {
      console.error('Greška :', error);
      return { result: false, error: 'Network error' };
    }
  }

  async getGraphDataForDates(id, start, end) {
    try {
      const response = await fetch('http://localhost:8081/api/hb/selected-date/' + id, {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({ start, end })
      })
      // console.log(response)

      if (response.ok) {
        const data = await response.json();
        console.log(data)
        return { result: data };
      } else {
        const data = await response.json();
        return { result: false, error: data.error };
      }
    } catch (error) {
      console.error('Greška :', error);
      return { result: false, error: 'Network error' };
    }
  }
}

export default new DeviceService();

