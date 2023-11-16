// App.js
import React, { useEffect, useState } from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AppRoutes from './AppRoutes';
import { ThemeProvider } from '@mui/material/styles';  // Uvezite ThemeProvider
import theme from './theme';  // Uvezite vašu temu (prilagodite putanju)
import authService from './services/AuthService'
import Login from './components/Login/Login';


function App() {
  const [isUserValid, setIsUserValid] = useState(false);

  useEffect(() => {
    const checkUserValidity = async () => {
      const validationResult = await authService.validateUser();

      if (validationResult.success) {
        console.log('Korisnik je validan:', validationResult.user);
        setIsUserValid(true);
      } else {
        console.error('Greška prilikom provere korisnika:', validationResult.error);
        setIsUserValid(false);
      }
    };

    checkUserValidity();
  }, []); 

  return (
    <ThemeProvider theme={theme}>
      <Router>
        {isUserValid ? (
          <Routes>
            {AppRoutes.map((route, index) => {
              const { element, ...rest } = route;
              return <Route key={index} {...rest} element={element} />;
            })}
          </Routes>
        ) : (
          <Routes>
            <Route path="/login" element={<Login />} />
          </Routes>
        )}
      </Router>
    </ThemeProvider>
  );
}

export default App;
