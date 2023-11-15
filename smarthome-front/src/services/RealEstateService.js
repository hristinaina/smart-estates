class RealEstateService {

    static async getRealEstates() {
        try {
            const response = await fetch('http://localhost:8081/api/real-estates/');
            const data = await response.json();
            return data;
          } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
          }
    }

    static async getPendingRealEstates() {
      try {
        const response = await fetch('http://localhost:8081/api/real-estates/pending');
        const data = await response.json();
        return data;
      } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
      }
    }

}

export default RealEstateService;