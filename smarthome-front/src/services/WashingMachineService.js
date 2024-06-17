class WashingMachineService {

    async scheduledMode(requestData) {
        try {
            const response = await fetch('http://localhost:8081/api/wm/schedule', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
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

    async getScheduledMode(deviceId) {
        try {
            const response = await fetch('http://localhost:8081/api/wm/schedule/' + deviceId, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            });
            const data = await response.json();

            if (data != null)
                return data;
            else return [];
        } 
        
        catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }

    async getWMHistoryData(deviceId, email, startDate, endDate) {
        const hdata = {
            "DeviceId": deviceId,
            "UserEmail": email,
            "StartDate": startDate,
            "EndDate": endDate
        }
        try {
                const response = await fetch('http://localhost:8081/api/wm/history', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(hdata),
                credentials: 'include',
            });
            const data = await response.json();
            return data;
        } 
        catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }
}

export default new WashingMachineService();
