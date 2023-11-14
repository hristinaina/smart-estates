const devices = [
    {
      id: 1,
      image: '/images/real_estate_example.png',
      name: 'Luxury Apartment',
      type: 'Spacious apartment with stunning views',
      status: 'Online',
    },
    {
      id: 2,
      image: '/images/real_estate_example.png',
      name: 'Modern House',
      type: 'Contemporary house with a beautiful garden',
      status: 'Offline',
    },
    {
      id: 3,
      image: '/images/real_estate_example.png',
      name: 'Beachfront Villa',
      type: 'Villa with direct access to the beach',
      status: 'Online',
    },
    {
      id: 4,
      image: '/images/real_estate_example.png',
      name: 'Cozy Cottage',
      type: 'Charming cottage in a peaceful location',
      status: 'Online',
    },
    {
      id: 5,
      image: '/images/real_estate_example.png',
      name: 'City Penthouse',
      type: 'Penthouse with panoramic city views',
      status: 'Offline',
    },
    {
      id: 6,
      image: '/images/real_estate_example.png',
      name: 'Beachfront Villa',
      type: 'Villa with direct access to the beach',
      status: 'Online',
    },
    {
      id: 7,
      image: '/images/real_estate_example.png',
      name: 'Cozy Cottage',
      type: 'Charming cottage in a peaceful location',
      status: 'Online',
    },
    {
      id: 8,
      image: '/images/real_estate_example.png',
      name: 'City Penthouse',
      type: 'Penthouse with panoramic city views',
      status: 'Offline',
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
}

export default new DeviceService();

