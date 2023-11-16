class RealEstateService {

    static async get() {
        try {
            const response = await fetch('http://localhost:8081/api/real-estates/');
            const data = await response.json();
            return data;
          } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
          }
    }

    static async getPending() {
      try {
        const response = await fetch('http://localhost:8081/api/real-estates/pending');
        const data = await response.json();
        return data;
      } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
      }
    }

    // static async changeState(state, id) {
    //   console.log("usaoo " + state);
    //   console.log(id.toString());
    //   try {
    //     if (state === 0) {
    //       console.log("............................");
    //       console.log('http://localhost:8081/api/real-estates/' + id.toString() + "/0");
    //       const response = await fetch('http://localhost:8081/api/real-estates/' + id.toString() + "/0");
    //       const data = await response.json();
    //       return data;
    //     } else {
    //       const response = await fetch('http://localhost:8081/api/real-estates/' + toString(id) + "/1");
    //       const data = await response.json();
    //       return data;
    //     }
    //   } catch (error) {
    //     console.error("Error fetching data: ", error);
    //     throw error;
    //   }
    // }

    static async changeState(state, id) {
      var url  = '';
      if (state === 0) url = `http://localhost:8081/api/real-estates/${id}/0`;
      else {url = `http://localhost:8081/api/real-estates/${id}/1`};

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

}

export default RealEstateService;