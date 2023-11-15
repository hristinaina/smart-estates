import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { NewRealEstate } from './components/RealEstate/NewRealEstate';
import { Devices } from './components/Devices/Devices';

const AppRoutes = [
  {
    path: '/real-estates',
    element: <RealEstates />
  },
  {
    path: '/new-real-estate',
    element: <NewRealEstate/>
  },
  {
    path: '/devices',
    element: <Devices />
  }
];

export default AppRoutes;
