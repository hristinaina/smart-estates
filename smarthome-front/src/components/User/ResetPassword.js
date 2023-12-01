import React, {useState, useEffect} from 'react';
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

import './Form.css'; 
import resetPasswordService from '../../services/ResetPassword' 
import superAdminService from '../../services/SuperAdmin';
import authService from '../../services/AuthService';


const ResetPassword = () => {

  const navigate = useNavigate()
  const [isSuperadmin, setIsSuperadmin] = useState(false);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = await resetPasswordService.TokenExist();
        console.log('token', token)
        if (token === '') {
            navigate('/')
        }
        else if (token === 'superadmin') {
          setIsSuperadmin(true)
        }
      } catch(error) {
        console.error('Greška prilikom provere autentičnosti:', error);
      }
    };
    checkAuth();
  }, []);

    const [confirmPassword, setConfirmPassword] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);

    const passwordRegex = /^(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])(.{8,})$/;;

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
      if(isSuperadmin) {
        const result = await superAdminService.ResetPassword(password);
        if (result.success) {
            await authService.validateUser()
            navigate('/real-estates');
        } else {
            setSnackbarMessage(result.error);
            handleClick()
        }
      }
      else {
        const result = await resetPasswordService.ResetPassword(password);
        if (result.success) {
          navigate('/')
        } else {
            setSnackbarMessage(result.error);
            handleClick()
        }
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
      {isSuperadmin && (
          <div>
            <p className='almost-done'>Almost done...</p>
            <p className='subtitle'>For security, you must reset your password</p>
          </div>
        )}

      {!isSuperadmin && (
        <div>
          <p className='almost-done'>Reset Password</p>
        </div>
      )}


        <form>

        <div className='fields'>
            <div className='fields-name'>Password:</div>
            <TextField
                id="password"
                type={showPassword ? 'text' : 'password'}
                sx={{ m: 1, width: '27%' }}
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
            sx={{ m: 1, width: '27%' }}
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

            {!isSuperadmin && (
              <div className="remember">
                <Link to={"/"} style={{ textDecoration: 'none'}}>
                    <span id="remember-password">Ahh.. Now I remember my password</span>
                </Link>
              </div>
            )}

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
