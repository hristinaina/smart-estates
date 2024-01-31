class SprinklerService {

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
            const response = await fetch(`http://localhost:8081/api/sprinkler/${id}`);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error("Error fetching data: ", error);
            throw error;
        }
    }

    async getSpecialModes(id) {
        try {
            const response = await fetch(`http://localhost:8081/api/sprinkler/mode/${id}`);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error("Error fetching data: ", error);
            throw error;
        }
    }

    async changeState(id, isOn) {
        let url = '';
        if (isOn) url = `http://localhost:8081/api/sprinkler/${id}/on`;
        else  url = `http://localhost:8081/api/sprinkler/${id}/off`;
        const response = await fetch(url, this.getRequestOptions());

        if (!response.ok) {
            throw new Error(`HTTP error! Status ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async addMode(mode, id) {
        let selectedDays = "";
        mode.selectedDays.forEach(day => {
            selectedDays += day + ","
        });
        const url =  `http://localhost:8081/api/sprinkler/mode/${id}`;
        const requestOptions = {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({"StartTime": mode.start + ":00", "EndTime": mode.end + ":00", "SelectedDays": selectedDays}),
        };
        
        const response = await fetch(url, requestOptions);

        if (!response.ok) {
            throw new Error(`HTTP error! Status ${response.status}`);
        }

        const responseData = await response.json();
        return responseData;
    }

    async deleteMode(modeId) {
        const url =  `http://localhost:8081/api/sprinkler/mode/${modeId}`;

        const requestOptions = {
            method: 'DELETE',
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
}

export default new SprinklerService();