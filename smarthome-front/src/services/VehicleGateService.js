
class VehicleGateService {
    getRequestOptions() {
        return  {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                };
    }  

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

    async toPrivate(id) {
        const url = `http://localhost:8081/api/vehicle-gate/private/${id}`;

        const response = await fetch(url, this.getRequestOptions());

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async toPublic(id) {
        const url = `http://localhost:8081/api/vehicle-gate/public/${id}`;

        const response = await fetch(url, this.getRequestOptions());

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async open(id) {
        const url = `http://localhost:8081/api/vehicle-gate/open/${id}`;

        const response = await fetch(url, this.getRequestOptions());

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async close(id) {
        const url = `http://localhost:8081/api/vehicle-gate/close/${id}`;

        const response = await fetch(url, this.getRequestOptions());

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }
}

export default new VehicleGateService();