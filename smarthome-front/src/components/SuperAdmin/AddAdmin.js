import React, {useState} from 'react';
import { Link, useNavigate } from 'react-router-dom';
import theme from '../../theme';
import { ThemeProvider } from '@emotion/react';

import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import TextField from '@mui/material/TextField'
import Snackbar from '@mui/material/Snackbar';
import CloseIcon from '@mui/icons-material/Close';

import '../User/Form.css'; 
import superAdminService from '../../services/SuperAdmin' 
import { Navigation } from '../Navigation/Navigation';


const AddAdmin = () => {
    const [nameAdmin, setNameAdmin] = useState('');
    const [surnameAdmin, setSurnameAdmin] = useState('');
    const [emailAdmin, setEmailAdmin] = useState('');
    const [isButtonAddDisabled, setIsButtonAddDisabled] = useState(true);

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    const navigate = useNavigate();

    const [open, setOpen] = React.useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState(''); 

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

    // add new admin
    const handleSignUpAdmin = async () => {
        const result = await superAdminService.AddAdmin(nameAdmin, surnameAdmin, emailAdmin)
    
        if (result.success) {
            setSnackbarMessage("New admin is added");
            handleClick()
            setEmailAdmin('')
            setNameAdmin('')
            setSurnameAdmin('')
            setIsButtonAddDisabled(true)
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
    <Navigation />
    <div className='container'>
        <p className='almost-done'>Add new admin</p>
        <p className='subtitle'>Join new members to our team!</p>
        <form>

        <div className='fields'>
            <div className='fields-name'>Name:</div>
            <TextField
                    value={nameAdmin}
                    onChange={handleAdminNameChange}
                    sx={{ m: 1, width: '34ch' }}
                    id="name"
                    placeholder="John" />

            <div className='input-fields'>
                <div className='fields-name'>Surname:</div>
                <TextField
                                value={surnameAdmin}
                                onChange={handleAdminSurnameChange}
                                sx={{ m: 1, width: '34ch' }}
                                id="surname"
                                className='text-field'
                                placeholder="Smith" />
            </div>

            <div className='input-fields'>
                <div className='fields-name'>Email:</div>
                <TextField
                            value={emailAdmin}
                            onChange={handleAdminEmailChange}
                            sx={{ m: 1, width: '34ch' }}
                            id="email"
                            className='text-field'
                            placeholder="someone@example.com"
                            type='email' />
            </div>
        </div>
            <Button 
                id='save'
                variant="contained" 
                color="primary" 
                disabled={isButtonAddDisabled}
                onClick={handleSignUpAdmin}
                style={{marginTop: "50px", textTransform: 'none'}} 
                >
                    Add Admin
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

export default AddAdmin;
