const devices = [
  {
    id: 2,
    image: '/images/lamp.png',
    name: 'Bakina lampa',
    type: 'Lamp',
    status: 'Offline',
  },
  {
    id: 3,
    image: '/images/washing_machine.png',
    name: 'Masina Sladja',
    type: 'Washing Machine',
    status: 'Online',
  },
  {
    id: 4,
    image: '/images/sprinkler.png',
    name: 'Prsk prsk',
    type: 'Sprinkler',
    status: 'Online',
  },
  {
    id: 5,
    image: '/images/solar_panel.png',
    name: 'Samo kes',
    type: 'Solar Panel',
    status: 'Offline',
  },
  {
    id: 6,
    image: '/images/washing_machine.png',
    name: 'Masina Masa',
    type: 'Washing machine',
    status: 'Online',
  },
  {
    id: 7,
    image: '/images/sprinkler.png',
    name: 'Curaljka',
    type: 'Sprinkler',
    status: 'Online',
  },
  {
    id: 8,
    image: '/images/lamp.png',
    name: 'Zvezda Severnjaca',
    type: 'Lamp',
    status: 'Offline',
  },
  {
    id: 1,
    image: '/images/solar_panel.png',
    name: 'Panelcici',
    type: 'Solar Panel',
    status: 'Online',
  },
  // ... more real estates
];

class DeviceService {

  async getDevices(realEstateId) {
    // try {
    //     const response = await fetch('http://localhost:8081/api/devices/' + realEstateId);
    //     const data = await response.json();
    //     return data;
    //   } catch (error) {
    //     console.error('Error fetching data:', error);
    //     throw error;
    //   }
    return devices;
  }

  async createDevice(device) {
    try {
      console.log(device);
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

