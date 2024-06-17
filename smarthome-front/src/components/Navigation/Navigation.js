import React, { Component } from 'react';
import { Navbar, NavItem, NavLink } from 'reactstrap';
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
        const user = authService.getCurrentUser()
        this.setState({ role: user['Role'] })
      }

    render() {
        const { role } = this.state;

        return (
            <header>
                <Navbar className="navbar">
                    {/* admin */}
                    {role===0 && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/real-estates">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/consumption">Consumption</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/account">Profile</NavLink>
                            </NavItem>
                            <NavItem className="logout">
                                <NavLink tag={Link} className="text-light" to="/" onClick={this.handleLogout}>Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                    {/* user */}
                    {role===1 && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/real-estates">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/account">Profile</NavLink>
                            </NavItem>
                            <NavItem className="logout">
                                <NavLink tag={Link} className="text-light" to="/" onClick={this.handleLogout}>Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                    {/* superadmin */}
                    {role===2 && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/real-estates">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/consumption">Consumption</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/account">Profile</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/add-admin">Add Admin</NavLink>
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
