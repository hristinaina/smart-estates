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

import './Registration.css';
import authService from '../../services/AuthService'

const Registration = () => {
    const [name, setName] = useState('');
    const [surname, setSurname] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [isButtonDisabled, setIsButtonDisabled] = useState(true);
    const [open, setOpen] = React.useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState('');

    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const [profileImage, setProfileImage] = useState(null);
    const [imagePreview, setImagePreview] = useState(null);

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    const passwordRegex = /^(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])(.{8,})$/;

    const navigate = useNavigate();

    const handleImageChange = (event) => {
      const file = event.target.files[0];
  
      // Kreirajte URL za prikaz slike
      const previewURL = URL.createObjectURL(file);
  
      setProfileImage(file);
      setImagePreview(previewURL);
    };

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

    const handleNameChange = (event) => {
      setName(event.target.value);
      event.target.value.trim() ===  '' || password.trim() === '' || confirmPassword.trim() === '' || confirmPassword !== password || surname.trim() === ''
      ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handleSurnameChange = (event) => {
      setSurname(event.target.value);
      event.target.value.trim() ===  '' || password.trim() === '' || confirmPassword.trim() === '' || confirmPassword !== password || name.trim() === ''
      ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handleEmailChange = (event) => {
      setEmail(event.target.value);
      event.target.value.trim() ===  '' || !emailRegex.test(event.target.value.trim()) || password.trim() === '' || confirmPassword.trim() === '' || confirmPassword !== password
      ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
        event.target.value.trim() ===  '' || !passwordRegex.test(event.target.value.trim()) || isValidEmail() || confirmPassword.trim() === '' || event.target.value.trim() !== confirmPassword
        ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const handleConfirmPasswordChange = (event) => {
      setConfirmPassword(event.target.value);
      event.target.value.trim() ===  '' || !passwordRegex.test(event.target.value.trim()) || isValidEmail() || password.trim() === '' || event.target.value.trim() !== password
      ? checkButtonDisabled(true) : checkButtonDisabled(false)
    };

    const isValidEmail = () => {
      return email.trim() === '' || !emailRegex.test(email.trim())
    }

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

    const resetFormFields = () => {
      setName('');
      setSurname('');
      setEmail('');
      setPassword('');
      setConfirmPassword('');
      // TODO refresh slike
    };

    // sign up
    const handleSignUp = async () => {
      // TODO promeni sliku
      const result = await authService.regUser(email, password, name, surname, "allaalal");
    
      if (result.success) {
          resetFormFields()
          setIsButtonDisabled(true)
          setSnackbarMessage("Successfully registered. Check your email!");
          handleClick()       
          // TODO posalji mejl
      } else {
        setSnackbarMessage(result.error);
        handleClick()
      }
    }

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
    <div className='ground'>

      <div className='left'>
        <p className='title-reg'>Sign up</p>

        <form>
        <div className='input-fields'>
            <div className='label'> Name:</div>
            <TextField
                value={name}
                onChange={handleNameChange}
                id="name"
                className='text-field'
                sx={{ m: 1, width: '30ch' }}
                placeholder="John"
                helperText="Required"
            />
        </div> 

        <div className='input-fields'>
            <div className='label'> Surname:</div>
            <TextField
                value={surname}
                onChange={handleSurnameChange}
                id="surname"
                className='text-field'
                sx={{ m: 1, width: '30ch' }}
                placeholder="Smith"
                helperText="Required"
            />
        </div> 

        <div className='input-fields'>
            <div className='label'> Email:</div>
            <TextField
                value={email}
                onChange={handleEmailChange}
                id="email"
                className='text-field'
                sx={{ m: 1, width: '30ch' }}
                placeholder="someone@example.com"
                helperText="Required"
                type='email'
            />
        </div>    
        <div className='input-fields'>
          <div className='label'>Password:</div>
          <TextField
            id="password"
            className='text-field'
            type={showPassword ? 'text' : 'password'}
            sx={{ m: 1, width: '30ch' }}
            placeholder='P@ssw0rd123'
            helperText="Required. Min 8 characters, special character, capital latter"
            value={password}
            onChange={handlePasswordChange}
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

        <div className='input-fields'>
          <div className='label'>Confirm password:</div>
          <TextField
            id="confirm-password"
            className='text-field'
            type={showConfirmPassword ? 'text' : 'password'}
            sx={{ m: 1, width: '30ch' }}
            placeholder='P@ssw0rd123'
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

        <div >
          <div htmlFor="profileImage" className='label' style={{marginBottom: "5px"}}>Profile Image:</div>
          <input
            type="file"
            id="profileImage"
            className='input-image'
            accept="image/*"
            onChange={handleImageChange} />
        </div>

        {/* Show choosen image */}
        {imagePreview && (
          <div>
            <img className='cropped-image' src={imagePreview} alt="Profile Preview"/>
          </div>
        )}

        <Button 
        id='signup'
          variant="contained" 
          color="secondary" 
          disabled={isButtonDisabled}
          onClick={handleSignUp}
          style={{marginTop: "50px", color: "#806894", textTransform: 'none'}} 
          sx={{ m: 1, width: '39ch' }}>
            Sign up
        </Button>

        <Snackbar
        open={open}
        autoHideDuration={3000}
        onClose={handleClose}
        message={snackbarMessage}
        action={action}/>
        </form>
      </div>

      <div className='right'>
        <p className='welcome'>Welcome to Smart Home!</p>
        <p className='desc'>One place to remotely manage all your devices!</p>
        <Link to="/">        
          <Button sx={theme.customStyles.myCustomButton} variant="contained" color="primary">Already have an account? Login</Button>
        </Link>
      </div>
    </div>
    </ThemeProvider>
  );
};

export default Registration;
