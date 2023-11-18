class DeviceService {

  async getDevices(realEstateId) {
    try {
        const response = await fetch('http://localhost:8081/api/devices/estate/' + realEstateId);
        const data = await response.json();
        return data;
      } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
      }
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

