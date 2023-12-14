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

    async getDataForSelectedTime(id, time) {
        try {
            const response = await fetch('http://localhost:8081/api/ambient/selected-time/' + id, {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ time })
            })
            // console.log(response)
            
            if (response.ok) {
                const data = await response.json();
                // console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }

    async getDataForSelectedDate(id, start, end) {
        try {
            const response = await fetch('http://localhost:8081/api/ambient/selected-date/' + id, {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ start, end })
            })
            // console.log(response)
            
            if (response.ok) {
                const data = await response.json();
                console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }
}

export default new AmbientSensorService();
