import React, {useState} from 'react';
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

import './ResetPassword.css'; 
import authService from '../../services/AuthService'
import superAdminService from '../../services/SuperAdmin' 


const ResetPassword = () => {
    const [confirmPassword, setConfirmPassword] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);

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

    const handleClickShowConfirmPassword = () => {
      setShowConfirmPassword(!showConfirmPassword);
    };

    const handleMouseDownConfirmPassword = (event) => {
        event.preventDefault();
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
        event.target.value.trim() ===  '' || !passwordRegex.test(event.target.value.trim()) || confirmPassword.trim() === '' || event.target.value.trim() !== confirmPassword
        ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handleConfirmPasswordChange = (event) => {
      setConfirmPassword(event.target.value);
      event.target.value.trim() ===  '' || !passwordRegex.test(event.target.value.trim()) || password.trim() === '' || event.target.value.trim() !== password
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

    // save reset password
    const handleResetPassword = async () => {
        const result = await superAdminService.ResetPassword(password);
        if (result.success) {
            await authService.validateUser()
            navigate('/real-estates');
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

      <div className='container'>
        <p className='almost-done'>Almost done...</p>
        <p className='subtitle'>For security you must reset your password</p>
        <form>

        <div className='fields'>
            <div className='fields-name'>Password:</div>
            <TextField
                id="password"
                type={showPassword ? 'text' : 'password'}
                sx={{ m: 1, width: '34ch' }}
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

        <div className='input-fields'>
          <div className='fields-name'>Confirm password:</div>
          <TextField
            id="confirm-password"
            className='text-field'
            type={showConfirmPassword ? 'text' : 'password'}
            sx={{ m: 1, width: '34ch' }}
            helperText="Required. Min 8 characters, special character, capital latter"
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={handleClickShowConfirmPassword}
                    onMouseDown={handleMouseDownConfirmPassword}>
                      {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
          </div>
        </div>
            <Button 
                id='save'
                variant="contained" 
                color="primary" 
                disabled={isButtonDisabled}
                onClick={handleResetPassword}
                style={{marginTop: "50px", textTransform: 'none'}} 
                >
                    Save
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
    </ThemeProvider>
  );
};

export default ResetPassword;
