class AmbientSensorService {

    async getGraphData(id){
        try {
            const response = await fetch('http://localhost:8081/api/ambient/last-hour/' + id);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data: ', error);
            throw error;
        }
    } 
}

export default new AmbientSensorService();
