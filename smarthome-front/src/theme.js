// theme.js

import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#806894',
    },
    secondary: {
      main: '#CEC0D9',
    },
  },
  customStyles: {
    myCustomButton: {
      textTransform: 'none',
      fontFamily: 'Poppins',
      width: '37%',
    },
    // Add more custom styles as needed
  },
});

export default theme;
