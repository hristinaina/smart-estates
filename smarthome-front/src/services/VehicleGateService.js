
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

    async getLicensePlates(id) {
        try {
            const url = `http://localhost:8081/api/vehicle-gate/license-plate/${id}`;
            const response = await fetch(url);
            const data = await response.json();
            return data;
        }
        catch (error) {
            console.error(error);
            throw error;
        }
    }

    async addLicensePlate(id, licensePlate) {
        const url = `http://localhost:8081/api/vehicle-gate/license-plate`

        const requestOptions = {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({"DeviceId": id, "LicensePlate": licensePlate}),
          };

        const response = await fetch(url, requestOptions);

        if (!response.ok) 
        {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const responseData = await response.json();
        return responseData;
    }

    async getCountGraphData(id, from, to, licensePlate) {
        try {
            let url = '';
            licensePlate = licensePlate.trim()
            console.log(licensePlate);
            if (licensePlate === '') {
                console.log('ovdje');
                url = `http://localhost:8081/api/vehicle-gate/count/${id}/${from}/${to}/-1`;
            } else {
                url = `http://localhost:8081/api/vehicle-gate/count/${id}/${from}/${to}/${licensePlate}`;
            }
            console.log(url);
            const response = await fetch(url);
            const data = await response.json();
            return data;
        } catch (error) {
            console.error(error);
            throw error;
        }
    }

    addNewGraphData(graphType, oldData, newData) {
            // license plate entries count graph
            console.log("dataaaa");
            console.log(oldData);
            if (!this.checkDates(newData["start_date"], newData["end_date"])) {
                return oldData;
            }
            if (newData["current_license_plate"] != '' & newData["current_license_plate"] != newData["license_plate"]) {
                return oldData;
            }
            
            if (graphType == 0) {
                if (oldData.labels.includes(newData["license_plate"])) {
                    const index = oldData.labels.indexOf(newData["license_plate"]);
                    let data = oldData.datasets[0].data;
                    let count = data[index] + 1;
                    data[index] = count;
                    oldData.datasets.data = data;
                }
                else {
                    const index  = oldData.labels.length;
                    let labels = oldData.labels;
                    labels.push(newData["license_plate"]);
                    console.log(index);
                    console.log(oldData.datasets[0]);
                    oldData.datasets[0].data.push(1);
                    oldData.labels = labels;
                  

                }
            }
            // TODO: handle other types of graphs
        
        
        // console.log(oldData);
        return oldData;
    }

    checkDates(startDate2, endDate2) {
        if (startDate2 != '' || endDate2 != '') {
            const startDate = new Date(startDate2);
            const endDate = new Date(endDate2);
    
            const currentDate = new Date();
    
            if (currentDate >= startDate && currentDate <= endDate) {
                return true;
            }
            return false;
        }
        return true;
    }

    checkDateOrder(startDate2, endDate2) {
        const startDate = new Date(startDate2);
        const endDate = new Date(endDate2);

        const currentDate = new Date();

        if (endDate <= startDate) {
            return false;
        }
        return true;
        
    }
    
}

export default new VehicleGateService();