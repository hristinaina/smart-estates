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
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }
}

export default new EVChargerService();