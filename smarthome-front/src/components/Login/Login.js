import React, {useState} from 'react';
import { Link } from 'react-router-dom';

import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

import './Login.css'; 


const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    const passwordRegex = /^(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])(.{8,})$/;

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    const handleMouseDownPassword = (event) => {
        event.preventDefault();
    };

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        event.target.value.trim() ===  '' || !emailRegex.test(event.target.value.trim()) || password.trim() === '' 
        ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
        event.target.value.trim() ===  '' || !passwordRegex.test(event.target.value.trim()) || username.trim() === '' 
        ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const checkButtonDisabled = (value) => {
        value ? setIsButtonDisabled(true) : setIsButtonDisabled(false);
    };

    const handleLogin = () => {
        // send values form to server
    }


  return (
    <div className='background'>
      <div className='left-side'>
        <p className='title-login'>Login</p>
        <form>
        <div className='fields'>
            <div style={{marginRight: "250px"}}> Email:</div>
            <TextField
                value={username}
                onChange={handleUsernameChange}
                id="username"
                sx={{ m: 1, width: '30ch' }}
                placeholder="e.g. someone@example.com"
                helperText="Required"
                type='email'
            />
        </div>    
        <div className='fields'>
            <div className='label'>Password:</div>
            <TextField
                id="password"
                type={showPassword ? 'text' : 'password'}
                sx={{ m: 1, width: '30ch' }}
                placeholder='e.g. !mikaMIKIC'
                helperText="Required. Min 8 characters, special character, capital latter"
                value={password}
                onChange={handlePasswordChange}
                required
                InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <IconButton
                            aria-label="toggle password visibility"
                            onClick={handleClickShowPassword}
                            onMouseDown={handleMouseDownPassword}>
                                {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                    </InputAdornment>
                ),
                }}
            />
        </div>
            <Button 
                variant="contained" 
                color="primary" 
                disabled={isButtonDisabled}
                onClick={handleLogin}
                style={{marginTop: "50px"}} 
                sx={{ m: 1, width: '39ch' }}>
                    Login
            </Button>
        </form>
      </div>
      <div className='right-side'>
        <p className='title'>Welcome to Smart Home!</p>
        <p className='text'>One place to remotely manage all your devices!</p>
        <Link to="/reg">
            <Button className='reg' variant="contained" color="secondary">No account yet? Sign up</Button>
        </Link>
      </div>
    </div>
  );
};

export default Login;
