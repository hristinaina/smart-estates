class LampService {

    async getGraphData(from, to){
        try {
            const response = await fetch('http://localhost:8081/api/lamp/graph/' 
            + from + '/' + to);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data: ', error);
            throw error;
        }
    } 
}

export default new LampService();