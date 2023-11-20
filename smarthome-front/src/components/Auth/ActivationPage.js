import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import theme from '../../theme';
import authService from '../../services/AuthService'
import Button from '@mui/material/Button';


export class ActivationPage extends Component {
  static displayName = ActivationPage.name;

  constructor(props) {
    super(props);
    this.state = {
      activationStatus: "",
      isMounted: false,
    };
  }

  componentDidMount() {
    this.setState({ isMounted: true }, () => {
      this.activateAccount();
    });
  }

  componentWillUnmount() {
    this.setState({ isMounted: false });
  }

  activateAccount = async () => {
    try {
      const result = await authService.activateAccount();

      if (this.state.isMounted) {
        if (result.success) {
          this.setState({ activationStatus: 'You have successfully activated your account!' });
        } else {
          this.setState({
            activationStatus:
              'An error occurred while activating the account.',
          });
        }
      }
    } catch (error) {
      if (this.state.isMounted) {
        this.setState({
          activationStatus: 'An error occurred while activating the account.',
        });
        console.error('Greska:', error);
      }
    }
  };

  render() {
    const { activationStatus } = this.state;

    return (
      <div>
        <h1>Account activation check</h1>
        {activationStatus !== null && (
          <>
            <p>{activationStatus}</p>
            {activationStatus === 'You have successfully activated your account!' && (
              <Link to="/">
                <Button
                  sx={theme.customStyles.myCustomButton}
                  variant="contained"
                  color="primary"
                >
                  Go to login
                </Button>
              </Link>
            )}
          </>
        )}
      </div>
    );
  }
}


