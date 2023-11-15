import React from 'react';
import Login from "./components/Login/Login";
import Registration from "./components/Registration/Registration";


const AppRoutes = [
  {
    path: '/',
    element: <Login />
  },
  
  {
    path: '/reg',
    element: <Registration />
  }
];

export default AppRoutes;
