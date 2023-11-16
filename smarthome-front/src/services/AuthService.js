const loginUser = async (email, password) => {
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
  
  export default { loginUser };
  