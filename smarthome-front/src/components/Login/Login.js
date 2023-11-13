// Login.js
import React, {useState} from 'react';
import './Login.css'; // Uvozite CSS datoteku
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';


const Login = () => {
    const [showPassword, setShowPassword] = useState(false);

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    const handleMouseDownPassword = (event) => {
        event.preventDefault();
    };


  return (
    <div className='background'>
      <div className='left-side'>
        <p className='title-login'>Login</p>
        <div className='fields'>
            <div className='label'> Username:</div>
            <TextField
                id="username"
                sx={{ m: 1, width: '30ch' }}
                placeholder="Type here"
            />
        </div>    
        <div className='fields'>
            <div className='label'>Password:</div>
            <TextField
                id="password"
                type={showPassword ? 'text' : 'password'}
                sx={{ m: 1, width: '30ch' }}
                placeholder='Type here'
                InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                    <IconButton
                        aria-label="toggle password visibility"
                        onClick={handleClickShowPassword}
                        onMouseDown={handleMouseDownPassword}
                    >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                    </InputAdornment>
                ),
                }}
            />
        </div>

        <Button variant="contained" color="primary" style={{marginTop: "50px"}} sx={{ m: 1, width: '39ch' }}>Login</Button>
      </div>
      <div className='right-side'>
        <p className='title'>Welcome to Smart Home!</p>
        <p className='text'>One place to remotely manage all your devices!</p>
        <Button className='reg' variant="contained" color="secondary">No account yet? Sign up</Button>
      </div>
    </div>
  );
};

export default Login;
