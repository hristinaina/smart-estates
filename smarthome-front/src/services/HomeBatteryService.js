class HomeBatteryService {

    async getHB(id) {
        try {
            const response = await fetch('http://localhost:8081/api/hb/' + id, {
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

    async getHBGraphDataByRS(id) {
        try {
            const response = await fetch('http://localhost:8081/api/hb/last-hour/' + id, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching data: ', error);
            throw error;
        }
    }

    async getGraphDataForDropdownSelect(estateId, time) {
        try {
            const response = await fetch('http://localhost:8081/api/hb/selected-time/' + estateId, {
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

    async getGraphDataForDates(id, start, end) {
        try {
            const response = await fetch('http://localhost:8081/api/hb/selected-date/' + id, {
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

export default new HomeBatteryService();