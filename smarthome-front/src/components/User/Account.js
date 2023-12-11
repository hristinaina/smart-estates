import React, { useState, useEffect } from 'react';
import './Account.css';
import { Navigation } from '../Navigation/Navigation';
import authService from '../../services/AuthService';
import superAdminService from '../../services/SuperAdmin';
import theme from '../../theme';

import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Snackbar from '@mui/material/Snackbar';
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';
import ImageService from '../../services/ImageService';


const Account = () => {
  const [selectedOption, setSelectedOption] = useState('PROFILE');
  const [profileImage, setProfileImage] = useState('');
  const [name, setName] = useState('');
  const [surname, setSurname] = useState('');
  const [email, setEmail] = useState('');
  const [isButtonUpdateDisabled, setIsButtonUpdateDisabled] = useState(true);

  const [user, setUser] = useState({});

  const [open, setOpen] = React.useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState(''); 

  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

  useEffect(() => {
    const valid = authService.validateUser();
    if (!valid) window.location.assign("/");

    const user = authService.getCurrentUser();
    setImage(user.Email.replace('@', ''));
    setName(user.Name)
    setSurname(user.Surname)
    setEmail(user.Email)
    setUser(user);
  }, [setUser]);

  const setImage = async(img) => {
    const url = await ImageService.getImage(img);
    setProfileImage(url);
  };

  const handleNameChange = (event) => {
    setName(event.target.value);
    event.target.value.trim() ===  '' || surname.trim() === '' 
        ? checkButtonUpdateDisabled(true) : checkButtonUpdateDisabled(false)
  };

  const handleSurnameChange = (event) => {
    setSurname(event.target.value);
    event.target.value.trim() ===  '' || name.trim() === '' 
        ? checkButtonUpdateDisabled(true) : checkButtonUpdateDisabled(false)
  };

  const checkButtonUpdateDisabled = (value) => {
    value ? setIsButtonUpdateDisabled(true) : setIsButtonUpdateDisabled(false);
  };

  const handleOptionChange = (option) => {
    setSelectedOption(option);
  };

  const handleUpdateProfile = async () => {
    const result = await superAdminService.EditSuperAdmin(name, surname, email)

    if (result.success) {
        setSnackbarMessage("Successful account change");
        handleClick();
        await authService.validateUser()
        user.Name = name;
        user.Surname = surname;
    } else {
        setSnackbarMessage(result.error);
        handleClick()
    }
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
    <div>
    <Navigation />
    <div className="user-profile-container">
      <div className="side-menu">
        <div className='container-image'>
          <img
            id="profile-image"
            src={user['Role'] === 1 ? profileImage : "/images/user.png"}
            alt="User"/>
          <img id='add-image' src="/images/plus_purple.png" alt="Add Image"/>
        </div>

        <div className='name-surname'>
          <p style={{display: 'inline'}}>{user.Name}</p>
          <p style={{display: 'inline', marginLeft: '15px'}}>{user.Surname}</p>
        </div>

        <div
          className={`menu-option ${selectedOption === 'PROFILE' ? 'selected' : ''}`}
          onClick={() => handleOptionChange('PROFILE')}>
            PROFILE
        </div>
      </div>

      <div className="content">
        {selectedOption === 'PROFILE' && (
          <>
          <form className='update-form'> 
          <p className='about-you'>About you</p>
            <div className='user-data'>
                <div className='field-name'> Name:</div>
                <TextField
                    value={name}
                    onChange={handleNameChange}
                    sx={{ m: 1, width: '34ch' }}
                    id="name"
                    placeholder="Add your name"
                    InputProps={{
                        readOnly: user.Name !== ''
                      }} />
            </div> 

            <div className='user-data'>
                <div className='field-name'> Surname:</div>
                <TextField
                    value={surname}
                    onChange={handleSurnameChange}
                    sx={{ m: 1, width: '34ch' }}
                    id="surname"
                    className='text-field'
                    placeholder="Add your surname"
                    helperText={user.Surname !== '' ? '' : 'Required'}
                    InputProps={{
                        readOnly: user.Surname !== ''
                      }}  />
            </div> 

            <div className='user-data'>
                <div className='field-name'> Email:</div>
                <TextField
                    value={email}
                    sx={{ m: 1, width: '34ch' }}
                    id="email"
                    className='text-field'
                    placeholder="someone@example.com"
                    type='email' 
                    InputProps={{
                        readOnly: true,
                      }}/>
            </div>
            <Button 
                id='update'
                variant="contained" 
                color="primary" 
                disabled={isButtonUpdateDisabled}
                onClick={handleUpdateProfile}
                sx={theme.customStyles.myCustomButton}
                style={{ display: (user.Name === '' || user.Surname === '') ? 'block' : 'none' }}>UPDATE
            </Button> 

            <Snackbar
              open={open}
              autoHideDuration={1000}
              onClose={handleClose}
              message={snackbarMessage}
              action={action}/>
          </form>
          </>
        )}
      </div>
    </div>
    </div>
  );
};

export default Account;
