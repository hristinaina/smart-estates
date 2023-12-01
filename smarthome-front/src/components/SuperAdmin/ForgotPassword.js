import React, {useState} from 'react';
import { Link, useNavigate } from 'react-router-dom';
import theme from '../../theme';
import { ThemeProvider } from '@emotion/react';

import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import TextField from '@mui/material/TextField';
import Snackbar from '@mui/material/Snackbar';
import CloseIcon from '@mui/icons-material/Close';

import './ResetPassword.css'; 
import resetPasswordService from '../../services/ResetPassword' 


const ForgotPassword = () => {
    const [email, setEmail] = useState('');
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    const navigate = useNavigate();

    const [open, setOpen] = React.useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState(''); 

    const handleEmailChange = (event) => {
        setEmail(event.target.value);
        event.target.value.trim() ===  '' ||  (!emailRegex.test(event.target.value.trim()) && event.target.value.trim() !==  'admin') 
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

    // send mail
    const handleResetPassword = async () => {
        const result = await resetPasswordService.SendResetPasswordEmail(email);
        if (result.success) {
            setSnackbarMessage(result.message);
            handleClick()
            setEmail('')
            navigate('/')
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

    <div className='container' >
        <p className='almost-done'>Forgot Password ?</p>
        <p className='subtitle'>Don't worry! It happens. Please enter the address associated with your account.</p>
        <form>

        <div className='fields'>
            <div className='fields-name'>Email:</div>
            <TextField
                className='text-field'
                id="email"
                sx={{ m: 1, width: '27%' }}
                value={email}
                onChange={handleEmailChange}
                placeholder="someone@example.com"
                helperText="Required"
                type='email' />
        </div>
            <Button 
                id='save'
                variant="contained" 
                color="primary" 
                disabled={isButtonDisabled}
                onClick={handleResetPassword}
                style={{marginTop: "50px", textTransform: 'none'}} 
                >
                    Submit
            </Button>

            <div className="remember">
                <Link to={"/"} style={{ textDecoration: 'none'}}>
                    <span id="remember-password">Ahh.. Now I remember my password</span>
                </Link>
    
            </div>
            <Snackbar
        open={open}
        autoHideDuration={1000}
        onClose={handleClose}
        message={snackbarMessage}
        action={action}
        />
        </form>
    </div>      
    </ThemeProvider>
    );
};

export default ForgotPassword;
