class PermissionService {

    async sendGrantValues(request) {
        try {
            const response = await fetch('http://localhost:8081/api/permission', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify(request)
            })
            
            if (response.ok) {
                const data = await response.json();
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

    async verifyAccount() {
        try {
            const queryParams = new URLSearchParams(window.location.search);
            const token = queryParams.get('token');

            const response = await fetch('http://localhost:8081/api/permission/verify', {
                method: 'POST',
                body: JSON.stringify({token})
            })
            
            if (response.ok) {
                return { success: true };
            } else {
                const data = await response.json();
                return { success: false, error: data.error };
            }
        } 
        catch (error) {
            console.error('Greška :', error);
            return { success: false, error: 'Network error' };
        }
    }

    async getPermissionsByRealEstateId(realEstateId) {
        try {
            const response = await fetch('http://localhost:8081/api/permission/' + realEstateId, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            const data = await response.json();
            if (data != null)
                return data;
            else 
                return [];
        } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }
}

export default new PermissionService();
