class SuperAdminService {
    async ResetPassword(password) {
        try {
          const response = await fetch('http://localhost:8081/api/users/reset-password', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password }),
            credentials: 'include',
          });
      
          if (response.ok) {
            return { success: true };
          } else {
            const data = await response.json();
            return { success: false, error: data.error };
          }
        } catch (error) {
          console.error('Greška prilikom slanja zahtjeva:', error);
          return { success: false, error: 'Network error' };
        }
    };

    async AddAdmin(name, surname, email) {
      try {
        const response = await fetch('http://localhost:8081/api/users/add-admin', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ name, surname, email }),
          credentials: 'include',
        });
    
        if (response.ok) {
          return { success: true };
        } else {
          const data = await response.json();
          return { success: false, error: data.error };
        }
      } catch (error) {
        console.error('Greška prilikom slanja zahtjeva:', error);
        return { success: false, error: 'Network error' };
      }
  };
}

const superAdminService = new SuperAdminService();

export default superAdminService;