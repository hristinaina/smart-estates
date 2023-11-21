import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { NewRealEstate } from './components/RealEstate/NewRealEstate';
import { Devices } from './components/Devices/Devices';
import Login from './components/Login/Login';
import Registration from './components/Registration/Registration';
import { NewDevice } from './components/Devices/NewDevice';
import { Lamp } from './components/Devices/Lamp';

const AppRoutes = [
  {
    path: '/',
    element: <Login />
  },
  {
    path: '/reg',
    element: <Registration />
  },
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
  },
  {
    path: '/new-device',
    element: <NewDevice />
  },
  {
    path: "/lamp/:id",
    element: <Lamp />
  }
];

export default AppRoutes;
