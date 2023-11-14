import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { NewRealEstate } from './components/RealEstate/NewRealEstate';

const AppRoutes = [
  {
    path: '/real-estates',
    element: <RealEstates />
  },
  {
    path: '/new-real-estate',
    element: <NewRealEstate/>
  }
];

export default AppRoutes;
