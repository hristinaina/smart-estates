import React, {useState} from 'react';
import { Link } from 'react-router-dom';

import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

import './Registration.css';

const Registration = () => {
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [profileImage, setProfileImage] = useState(null);
    const [imagePreview, setImagePreview] = useState(null);

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


  return (
    <div className='ground'>
      <div className='left'>
        <p className='title-reg'>Sign up</p>
        <div className='input-fields'>
            <div className='label'> Username:</div>
            <TextField
                id="username"
                className='text-field'
                sx={{ m: 1, width: '30ch' }}
                placeholder="Type here"
            />
        </div>    
        <div className='input-fields'>
          <div className='label'>Password:</div>
          <TextField
            id="password"
            className='text-field'
            type={showPassword ? 'text' : 'password'}
            sx={{ m: 1, width: '30ch' }}
            placeholder='Type here'
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
          <div style={{marginRight: "150px"}}>Confirm password:</div>
          <TextField
            id="confirm-password"
            className='text-field'
            type={showConfirmPassword ? 'text' : 'password'}
            sx={{ m: 1, width: '30ch' }}
            placeholder='Type here'
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
          <div htmlFor="profileImage" className='label' style={{marginBottom: "5px", marginRight: "200px"}}>Profile Image:</div>
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

        <Button variant="contained" color="secondary" style={{marginTop: "50px", color: "#806894"}} sx={{ m: 1, width: '39ch' }}>Sign up</Button>
      </div>
      <div className='right'>
        <p className='welcome'>Welcome to Smart Home!</p>
        <p className='desc'>One place to remotely manage all your devices!</p>
        <Link to="/">        
          <Button variant="contained" color="primary">Already have an account? Login</Button>
        </Link>
      </div>
    </div>
  );
};

export default Registration;
