import React, {useEffect, useState} from 'react';
import { Link, useNavigate } from 'react-router-dom';
import theme from '../../theme';
import { ThemeProvider } from '@emotion/react';

import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import Snackbar from '@mui/material/Snackbar';
import CloseIcon from '@mui/icons-material/Close';

import './Login.css'; 
import authService from '../../services/AuthService'


const Login = () => {
    useEffect(() => {
        const fetchData = async () => {
            try {
                const result = await authService.validateUser();
                console.log(result)
                !result ? navigate('/'): navigate('/real-estates');
            } catch (error) {
            console.error('Error:', error);
            }
        };
        
        fetchData();
    }, []);

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    const passwordRegex = /^(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])(.{8,})$/;

    const navigate = useNavigate();

    const [open, setOpen] = React.useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState(''); 

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    const handleMouseDownPassword = (event) => {
        event.preventDefault();
    };

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        event.target.value.trim() ===  '' ||  (!emailRegex.test(event.target.value.trim()) && event.target.value.trim() !==  'admin') || password.trim() === '' 
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

    // snackbar
    const handleClick = () => {
        setOpen(true);
    };

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        setOpen(false);
    };

    // login
    const handleLogin = async () => {
        const result = await authService.loginUser(username, password);
    
        if (result.success) {
            const result = await authService.validateUser()
            !result ? navigate('/reset-password?token=superadmin'): navigate('/real-estates');
        } else {
            setSnackbarMessage(result.error);
            handleClick()
        }
    };

    const action = (
    <React.Fragment>
        <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleClose}>
        <CloseIcon fontSize="small" />
        </IconButton>
    </React.Fragment>
    );


    return (
    <ThemeProvider theme={theme}>
    <div className='background'>
        <div className='left-side'>
        <p className='title-login'>Login</p>
        <form>
        <div className='fields'>
            <div className='label'> Email:</div>
            <TextField
                value={username}
                onChange={handleUsernameChange}
                id="username"
                sx={{ m: 1, width: '30ch' }}
                placeholder="someone@example.com"
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
        <Link style={{ textDecoration: 'none'}} to="/forgot-password">
            <div className='forgot-password'>Forgot password ?</div>
        </Link>
            <Button 
                id='login'
                variant="contained" 
                color="primary" 
                disabled={isButtonDisabled}
                onClick={handleLogin}
                style={{marginTop: "50px", textTransform: 'none'}} 
                sx={{ m: 1, width: '39ch' }}>
                    Login
            </Button>
            <Snackbar
        open={open}
        autoHideDuration={1000}
        onClose={handleClose}
        message={snackbarMessage}
        action={action}
        />
        </form>
        </div>
        <div className='right-side'>
            <p className='title'>Welcome to Smart Home!</p>
            <p className='text'>One place to remotely manage all your devices!</p>
            <Link to="/reg">
                <Button className="reg" sx={theme.customStyles.myCustomButton} variant="contained" color="secondary">No account yet? Sign up</Button>
            </Link>
        </div>
    </div>
    </ThemeProvider>
    );
};

export default Login;
