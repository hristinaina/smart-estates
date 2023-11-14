import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { Devices } from './components/Devices/Devices';

const AppRoutes = [
  {
    path: '/',
    element: <RealEstates />
  },
  {
    path: '/devices',
    element: <Devices />
  }
];

export default AppRoutes;
