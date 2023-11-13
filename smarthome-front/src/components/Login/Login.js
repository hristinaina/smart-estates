// Login.js
import React from 'react';
import './Login.css'; // Uvozite CSS datoteku
import Button from '@mui/material/Button';

const Login = () => {
  return (
    <div className='background'>
      <div className='left-side'>
        {/* Vaša forma za login ide ovde */}
        <form>
          {/* Polja forme, dugme za prijavu itd. */}
        </form>
      </div>
      <div className='right-side'>
        {/* Tekst i dugme za registraciju idu ovde */}
        <p className='title'>Welcome to Smart Home!</p>
        <p className='text'>One place to remotely manage all your devices!</p>
        <Button className='reg' variant="contained" color="secondary">No account yet? Sign up</Button>
      </div>
    </div>
  );
};

export default Login;
