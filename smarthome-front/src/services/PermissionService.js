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
}

export default new PermissionService();
