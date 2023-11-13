// App.js
import React from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AppRoutes from './AppRoutes';
import { ThemeProvider } from '@mui/material/styles';  // Uvezite ThemeProvider
import theme from './theme';  // Uvezite vašu temu (prilagodite putanju)


function App() {
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
