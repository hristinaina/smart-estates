class ResetPasswordService {

    async SendResetPasswordEmail(email) {
        try {
            const response = await fetch('http://localhost:8081/api/users/verify-email', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email }),
            credentials: 'include',
            });
        
            if (response.ok) {
                const data = await response.json();
                return { success: true, message: data.message };
            } else {
                const data = await response.json();
                return { success: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška prilikom slanja zahteva:', error);
            return { success: false, error: 'Network error' };
        }
    };


    async ResetPassword(password) {
        try {
            const queryParams = new URLSearchParams(window.location.search);
            const token = queryParams.get('token');   

            const response = await fetch('http://localhost:8081/api/users/reset-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password, token }),
            credentials: 'include',
            });

            if (response.ok) {
                return { success: true };
            } else {
                const data = await response.json();
                return { success: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška prilikom slanja zahteva:', error);
            return { success: false, error: 'Network error' };
        }
    };

    async TokenExist() {
        const queryParams = new URLSearchParams(window.location.search);
        const token = queryParams.get('token'); 
        
        console.log(token)

        if (token !== null) {
            return token;
        } else {
            return '';
        }  
    }
}

const resetPasswordService = new ResetPasswordService();

export default resetPasswordService;