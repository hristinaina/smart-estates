// App.js
import React, { useEffect, useState } from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Routes, Navigate, useNavigate} from 'react-router-dom';
import AppRoutes from './AppRoutes';
import { ThemeProvider } from '@mui/material/styles';  // Uvezite ThemeProvider
import theme from './theme';  // Uvezite vašu temu (prilagodite putanju)


function App() {

  // useEffect(() => {
  //   const checkAuthentication = async () => {
  //     try {
  //       const response = await fetch('http://localhost:8081/api/users/validate', {
  //         method: 'GET',
  //         credentials: 'include',
  //       });

  //       if (response.status === 200) {
  //         // Ako je korisnik autentikovan, postavi stanje na autentikovan
  //         setAuthenticated(true);
  //       }
  //     } catch (error) {
  //       console.error('Greška prilikom provere korisnika:', error);
  //     }
  //   };

  //   checkAuthentication();
  // }, []);

  return (
    <ThemeProvider theme={theme}>
      <Router>
        <Routes>
          {AppRoutes.map((route, index) => {
            const { element, ...rest } = route;
            return <Route key={index} {...rest} element={element} />;
          })}
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
