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
}

export default new WashingMachineService();
