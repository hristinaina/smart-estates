
class LampService {

    async turnOn(id) {
        const url = `http://localhost:8081/api/lamp/on/${id}`;

        const requestOptions = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
        };

        const response = await fetch(url, requestOptions);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async turnOff(id) {
        const url = `http://localhost:8081/api/lamp/off/${id}`;

        const requestOptions = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
        };

        const response = await fetch(url, requestOptions);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async getGraphData(id, from, to){
        try {
            const url = 'http://localhost:8081/api/lamp/graph/' + id + "/" + from + '/' + to;
            const response = await fetch('http://localhost:8081/api/lamp/graph/' + id + "/" + from + '/' + to);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data: ', error);
            throw error;
        }
    }

    async getAllGraphData(id) {
        let ls = new LampService();
        const result24h = await ls.getGraphData(id, '-24h', '-1');
        const result12h = await ls.getGraphData(id, '-12h', '-1');
        const result6h = await ls.getGraphData(id, '-6h', '-1');
        const result1h = await ls.getGraphData(id, '-1h', '-1');
        const result30m = await ls.getGraphData(id, '-30m', '-1');
        const result7d = await ls.getGraphData(id, '-7d', '-1');
        const result30d = await ls.getGraphData(id, '-30d', '-1');

        let allResults = [result30m, result1h, result6h, result12h, result24h, result7d, result30d];
        let allLabels = ["last 30min", "last 1h", "last 6h", "last 12h", "last 24h", "last 7 days", "last 30 days"];
        
        let data = new Map();

        let values = [];
        for (let i = 0; i < 101; i++){
            values.push(0);
        }
        allResults.forEach((res, i) => {
            if (res.data != null)
            {
                res.data.forEach(element => {
                    values[element.Value] = element.Count;
                });
                data.set(allLabels[i], values);
                values = []
            }

        });

        return data;
    }

    createGraphData(graphData) {
        let values = []
        for (let i = 0; i < 101; i++) {
            values.push(0);
        }
        if (graphData != null) {
            graphData.forEach((res, i) => {
            values[res.Value] = res.Count; 
            });
        }

        return values;
    }

    getRandomColor() {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        const alpha = 1;
        return `rgba(${r},${g},${b},${alpha})`;
    }
}

export default new LampService();