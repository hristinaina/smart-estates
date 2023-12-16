
class VehicleGateService {

    async get(id) {
        try {
            const response = await fetch(`http://localhost:8081/api/vehicle-gate/${id}`);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data: ', error);
            throw error;
        }
    }
}

export default new VehicleGateService();