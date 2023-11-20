import React, { useState, useEffect } from 'react';
import './Account.css';
import { Navigation } from '../Navigation/Navigation';
import authService from '../../services/AuthService';
import theme from '../../theme';

import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';

const Account = () => {
  const [selectedOption, setSelectedOption] = useState('PROFILE');
  const [name, setName] = useState('');
  const [surname, setSurname] = useState('');
  const [email, setEmail] = useState('');

  const [nameAdmin, setNameAdmin] = useState('');
  const [surnameAdmin, setSurnameAdmin] = useState('');
  const [emailAdmin, setEmailAdmin] = useState('');

  const [isButtonUpdateDisabled, setIsButtonUpdateDisabled] = useState(true);
  const [isButtonAddDisabled, setIsButtonAddDisabled] = useState(true);

  const [user, setUser] = useState({});

  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

  useEffect(() => {
    const user = authService.getCurrentUser();
    setName(user.Name)
    setSurname(user.Surname)
    setEmail(user.Email)
    console.log("user ", user)
    setUser(user);
  }, [setUser]);

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

  const handleAdminNameChange = (event) => {
    setNameAdmin(event.target.value);
    event.target.value.trim() ===  '' || surnameAdmin.trim() === '' || emailAdmin.trim() === ''
        ? checkButtonAddDisabled(true) : checkButtonAddDisabled(false)
  };

  const handleAdminSurnameChange = (event) => {
    setSurnameAdmin(event.target.value);
    event.target.value.trim() ===  '' || nameAdmin.trim() === '' || emailAdmin.trim() === ''
        ? checkButtonAddDisabled(true) : checkButtonAddDisabled(false)
  };

  const handleAdminEmailChange = (event) => {
    setEmailAdmin(event.target.value);
    event.target.value.trim() ===  '' ||  !emailRegex.test(event.target.value.trim()) || nameAdmin.trim() === '' || surnameAdmin.trim() === '' 
        ? checkButtonAddDisabled(true) : checkButtonAddDisabled(false)
  };

  const checkButtonAddDisabled = (value) => {
    value ? setIsButtonAddDisabled(true) : setIsButtonAddDisabled(false);
  };

  const handleOptionChange = (option) => {
    setSelectedOption(option);
  };

  const handleUpdateProfile = () => {
    // Logika za ažuriranje profila
    console.log('Profil je ažuriran!');
  };

  const handleSignUpAdmin = () => {
    // Logika za registraciju admina
    console.log('Admin je registrovan!');
  };

  return (
    <div>
    <Navigation />
    <div className="user-profile-container">
      <div className="side-menu">
        <div className='container-image'>
            <img id='profile-image' src="/images/user.png" alt="User" />
            <img id='add-image' src="/images/plus_purple.png" alt="Add Image"/>
        </div>
        <p></p>
        <div
          className={`menu-option ${selectedOption === 'PROFILE' ? 'selected' : ''}`}
          onClick={() => handleOptionChange('PROFILE')}
        >
          PROFILE
        </div>
        <div
          className={`menu-option ${selectedOption === 'ADD_ADMIN' ? 'selected' : ''}`}
          onClick={() => handleOptionChange('ADD_ADMIN')}
        >
          ADD ADMIN
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
                        readOnly: user.Name !== '' // postavlja readOnly na true ako je surname različito od praznog stringa
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
                    InputProps={{
                        readOnly: user.Surname !== '' // postavlja readOnly na true ako je surname različito od praznog stringa
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
                style={{ display: (user.Name === '' || user.Surname === '') ? 'block' : 'none' }}>UPDATE</Button> 
          </form>
          </>
        )}

        {selectedOption === 'ADD_ADMIN' && (
          <>
            <div className="admin-form">
            <form className='update-form'> 
                <p className='about-you'>Add new admin</p>
                    <div className='user-data'>
                        <div className='field-name'> Name:</div>
                        <TextField
                            value={nameAdmin}
                            onChange={handleAdminNameChange}
                            sx={{ m: 1, width: '34ch' }}
                            id="name"
                            placeholder="John" />
                    </div> 

                    <div className='user-data'>
                        <div className='field-name'> Surname:</div>
                        <TextField
                            value={surnameAdmin}
                            onChange={handleAdminSurnameChange}
                            sx={{ m: 1, width: '34ch' }}
                            id="surname"
                            className='text-field'
                            placeholder="Smith" />
                    </div> 

                    <div className='user-data'>
                        <div className='field-name'> Email:</div>
                        <TextField
                            value={emailAdmin}
                            onChange={handleAdminEmailChange}
                            sx={{ m: 1, width: '34ch' }}
                            id="email"
                            className='text-field'
                            placeholder="someone@example.com"
                            type='email' />
                    </div>

                    <Button 
                        id='add-admin-btn'
                        variant="contained" 
                        color="primary" 
                        disabled={isButtonAddDisabled}
                        onClick={handleUpdateProfile}
                        sx={theme.customStyles.myCustomButton}>ADD ADMIN</Button> 
            </form>
            </div>
          </>
        )}
      </div>
    </div>
    </div>
  );
};

export default Account;
