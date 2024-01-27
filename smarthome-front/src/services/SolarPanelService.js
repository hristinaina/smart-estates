class SolarPanelService {

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
}

export default new SolarPanelService();