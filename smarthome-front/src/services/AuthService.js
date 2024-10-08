class AuthService {
  
  async loginUser(email, password) {
    try {
      const response = await fetch('http://localhost:8081/api/users/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
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


  async validateUser() {
    const response = await fetch('http://localhost:8081/api/users/validate', {
      method: 'GET',
      credentials: 'include',
    });

    if (response.status === 200) {
      const data = await response.json();
      if (data["message"].Role !== 2 || data["message"]["IsLogin"])
      {
        this.setUser(data["message"])
        return true
      } 
      return false

    } else if (response.status === 401) {
      return false

    } else {
      console.error('Greška prilikom provere korisnika:', response.status);
      return false;
    }
  };

  async logoutUser() {
    try {
      const response = await fetch('http://localhost:8081/api/users/logout', {
        method: 'POST',
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

  async regUser(email, password, name, surname, role=1) {
    try{
      const response = await fetch('http://localhost:8081/api/users/verificationMail', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password, name, surname, role })
      });

      if (response.ok) {
        return { success: true };
      } else {
        const data = await response.json();
        return { success: false, error: data.error };
      }
    } catch (error) {
      console.error('Greška :', error);
      return { success: false, error: 'Network error' };
    }
  }; 

  async activateAccount() {
    try{
      const queryParams = new URLSearchParams(window.location.search);
      const token = queryParams.get('token');

      console.log(token)
      const response = await fetch('http://localhost:8081/api/users/activate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token })
      });

      if (response.ok) {
        return { success: true };
      } else {
        const data = await response.json();
        return { success: false, error: data.error };
      }
    } catch (error) {
      console.error('Greška :', error);
      return { success: false, error: 'Network error' };
    }
  }; 

  async setUser(user) {
    localStorage.setItem('user', JSON.stringify(user));
  }
  
  getCurrentUser() {
    const user = JSON.parse(localStorage.getItem('user'));
    if (!user) {
      window.location.href = '/';
    }
    return user;
  }
}

const authService = new AuthService();

export default authService;