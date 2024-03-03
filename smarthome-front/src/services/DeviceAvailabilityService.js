
class AvailabilityService {

    async get(id, startDate, endDate) {
        const url = `http://localhost:8081/api/devices/availability`;

        const requestOptions = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"DeviceId": id, "UserEmail": "", "StartDate": startDate, "EndDate": endDate}),
        };

        const response = await fetch(url, requestOptions);

        if (!response.ok) {
            throw new Error(`HTTP Error! Status: ${response.status}`);
        }

        return await response.json();
    }
}

export default new AvailabilityService();