import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../../services/AuthService'


const AuthenticationGuard = ({ children }) => {
  const history = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      const isAuthenticated = await authService.validateUser();

      if (!isAuthenticated) {
        console.log("usao je ovdee")
        // Ako korisnik nije autentikovan, preusmeri ga na login
        history.push('/');
      }
    };

    checkAuthentication();
  }, [history]);

  return children;
};

export default AuthenticationGuard;
