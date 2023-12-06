class RealEstateService {

    async get() {
        try {
            const response = await fetch('http://localhost:8081/api/real-estates/');
            const data = await response.json();
            return data;
          } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
          }
    }

    async getById(id) {
      try {
          const response = await fetch('http://localhost:8081/api/real-estates/' + id);
          const data = await response.json();
          return data;
        } catch (error) {
          console.error('Error fetching data:', error);
          throw error;
        }
    }

    async getPending() {
      try {
        const response = await fetch('http://localhost:8081/api/real-estates/pending');
        const data = await response.json();
        return data;
      } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
      }
    }

    async getAllByUserId(id) {
      console.log("usaoo");
      try {
        const response = await fetch('http://localhost:8081/api/real-estates/user/' + id);
        const data = await response.json();
        console.log("data: ", data);
        return data;
      } catch(error) {
        console.error("Error fetching data:", error);
        throw error;
      }
    }

    async changeState(state, id, reason) {
      var url  = '';
      if (state === 0) url = `http://localhost:8081/api/real-estates/${id}/0`;
      else {url = `http://localhost:8081/api/real-estates/${id}/1`};

      const requestOptions = {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          'DiscardReason': reason,
        }),
 
      };

      const response = await fetch(url, requestOptions);

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const responseData = await response.json();
      return responseData;
    }

    async add(estate) {
      console.log("usaooooooo");
      console.log(estate);
      const url = 'http://localhost:8081/api/real-estates/';

      const requestOptions = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(estate),
      };
      const response = await fetch(url, requestOptions);

      if (!response.ok){
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const responseData = await response.json();
      return responseData;
      
    }

}

export default new RealEstateService();