import React, { Component } from 'react';
import { Navbar, NavItem, NavLink, Collapse, NavbarToggler } from 'reactstrap';
import { Link } from 'react-router-dom';
import './Navigation.css';
import authService from '../../services/AuthService'


export class Navigation extends Component {
    static displayName = Navigation.name;  
    state = {
        role: null,
      };

    handleLogout = async () => {
        const result = await authService.logoutUser();
    
        if (result.success) {
            localStorage.removeItem('user')
          console.log('Uspešno ste se odjavili!');
        } else {
          console.error('Greška prilikom odjavljivanja:', result.error);
        }
      };

    constructor(props) {
        super(props);
    }

    componentDidMount() {
        // Dohvati podatke iz lokalnog skladišta
        const user = JSON.parse(localStorage.getItem('user'));
        if (user) {
            this.setState({ role: user['Role'] });
          }
      }

    render() {
        const { role } = this.state;

        return (
            <header>
                <Navbar className="navbar">
                    {role==1 && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/home">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/database-manager">Nekaj</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/reports">Nekaj</NavLink>
                            </NavItem>
                            <NavItem className="logout">
                                <NavLink tag={Link} className="text-light" to="/" onClick={this.handleLogout}>Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                    {role==0 && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/home">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/database-manager">Admin1</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/reports">ADmin2</NavLink>
                            </NavItem>
                            <NavItem className="logout">
                                <NavLink tag={Link} className="text-light" to="/" onClick={this.handleLogout}>Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                </Navbar>
            </header>
        );
    }
}
