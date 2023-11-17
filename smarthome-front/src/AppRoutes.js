import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { NewRealEstate } from './components/RealEstate/NewRealEstate';
import { Devices } from './components/Devices/Devices';
import Login from './components/Login/Login';
import Registration from './components/Registration/Registration';
import { ActivationPage } from './components/Auth/ActivationPage';

const AppRoutes = [
  {
    path: '/real-estates',
    element: <RealEstates />
  },
  {
    path: '/',
    element: <Login />
  },
  {
    path: '/reg',
    element: <Registration />
  },
  {
    path: '/new-real-estate',
    element: <NewRealEstate/>
  },
  {
    path: '/devices',
    element: <Devices />
  },
  {
    path: '/activate',
    element: <ActivationPage />
  }
];

export default AppRoutes;
