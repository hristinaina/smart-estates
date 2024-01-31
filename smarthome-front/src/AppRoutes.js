import React from 'react';
import { RealEstates } from './components/RealEstate/RealEstates';
import { NewRealEstate } from './components/RealEstate/NewRealEstate';
import { Devices } from './components/Devices/Devices';
import Login from './components/Login/Login';
import Registration from './components/Registration/Registration';
import { ActivationPage } from './components/Auth/ActivationPage';
import { NewDevice } from './components/Devices/NewDevice';
import ResetPassword from './components/User/ResetPassword';
import Account from './components/User/Account';
import { Lamp } from './components/Devices/Lamp';
import AddAdmin from './components/SuperAdmin/AddAdmin';
import ForgotPassword from './components/User/ForgotPassword';
import { AmbientSensor } from './components/Devices/AmbientSensor/AmbientSensor';
import { SolarPanel } from './components/Devices/SolarPanel/SolarPanel';
import { AirConditioner } from './components/Devices/AirConditioner/AirConditioner';
import { HomeBattery } from './components/Devices/HomeBattery/HomeBattery';
import { VehicleGate } from './components/Devices/VehicleGate/VehicleGate';
import { ElectricityOverview } from './components/Admin/ElectricityOverview';

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
  },
  {
    path: '/new-device',
    element: <NewDevice />
  },
  {
    path: '/reset-password',
    element: <ResetPassword />
  },
  {
    path: '/account',
    element: <Account/>
  },
  {
    path: '/add-admin',
    element: <AddAdmin />
  },
  {
    path: '/forgot-password',
    element: <ForgotPassword/>
  },
  {
    path: "/lamp/:id",
    element: <Lamp />
  },
  {
    path: "/ambient-sensor/:id",
    element: <AmbientSensor />
  },
  {
    path: "/air-conditioner/:id",
    element: <AirConditioner />
  },
  {
    path: "/vehicle-gate/:id",
    element: <VehicleGate/>
  },
  {
    path: "/sp/:id",
    element: <SolarPanel />
  },
  {
    path: "/hb/:id",
    element: <HomeBattery />
  },
  {
    path: "/consumption",
    element: <ElectricityOverview />
  },
];

export default AppRoutes;
