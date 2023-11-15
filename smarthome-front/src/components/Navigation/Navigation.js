import React, { Component } from 'react';
import { Navbar, NavItem, NavLink, Collapse, NavbarToggler } from 'reactstrap';
import { Link } from 'react-router-dom';
import './Navigation.css';

export class Navigation extends Component {
    static displayName = Navigation.name;

    constructor(props) {
        super(props);
    }

    render() {
        const isAdmin = true; //todo call function to get user role

        return (
            <header>
                <Navbar className="navbar">
                    {isAdmin && (
                        <ul>
                            <span className="logo">Smart Home</span>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/real-estates">Home</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/database-manager">Nekaj</NavLink>
                            </NavItem>
                            <NavItem>
                                <NavLink tag={Link} className="text-light" to="/reports">Nekaj</NavLink>
                            </NavItem>
                            <NavItem className="logout">
                                <NavLink tag={Link} className="text-light" to="/">Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                    {!isAdmin && (
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
                                <NavLink tag={Link} className="text-light" to="/">Log out</NavLink>
                            </NavItem>
                        </ul>
                    )}
                </Navbar>
            </header>
        );
    }
}
