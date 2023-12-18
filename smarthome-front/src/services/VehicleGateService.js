
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

    addNewGraphData(graphType, oldData, newData, id) {
        console.log("web socket radi");
        if (id == newData["vehicle_gate_id"]) {
            // license plate entries count graph
            if (graphType == 0) {
                if (oldData.labels.includes(newData["license_plate"])) {
                    const index = oldData.labels.indexOf(newData["license_plate"]);
                    let data = oldData.datasets[index].data;
                    let count = data[index] + 1;
                    data[index] = count;
                    oldData.datasets.data = data;
                }
                else {
                    const index  = oldData.labels.length;
                    let labels = oldData.labels;
                    labels.push(newData["license_plate"]);
                    console.log(index);
                    console.log(oldData.datasets[index]);
                    oldData.datasets[index].data.push(1);
                    oldData.labels = labels;
                  

                }
            }
            // TODO: handle other types of graphs
        }
        
        // console.log(oldData);
        return oldData;
    }
}

export default new VehicleGateService();