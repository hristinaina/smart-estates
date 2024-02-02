class EVChargerService {

    async get(id) {
        try {
            const response = await fetch('http://localhost:8081/api/ev/' + id, {
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

    async getLastPercentage(id) {
        try {
            const response = await fetch('http://localhost:8081/api/ev/lastPercentage/' + id, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            let data = await response.json();
            if (data.error) data = 0.9;
            return data;
        } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }

    
  async getTableActions(deviceId, email, startDate, endDate) {
    const hdata = {
      "DeviceId": deviceId,
      "UserEmail": email,
      "StartDate": startDate,
      "EndDate": endDate
    }
    try {
      const response = await fetch('http://localhost:8081/api/ev/actions', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(hdata),
        credentials: 'include',
      });
      const data = await response.json();
      console.log(data);
      return data.result;
    } catch (error) {
      console.error('Error fetching data:', error);
      throw error;
    }
  }
}

export default new EVChargerService();