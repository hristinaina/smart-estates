import React,{ Component, useState } from 'react';
import './RealEstates.css';
import { NewRealEstate } from './NewRealEstate';
import { Navigation } from '../Navigation/Navigation';
import authService from '../../services/AuthService'

export class RealEstates extends Component {
    constructor(props) {
        super(props);

        this.state = {
            showNewRealEstate: false,
            user : null
        };
    }

    toggleNewRealEstate = () => {
        this.setState((prevState) => ({
            showNewRealEstate: !prevState.showNewRealEstate,
        }));
    }

    componentDidMount() {
        const currentUser = authService.getCurrentUser()
        this.setState({user: currentUser})
      }

    render() {
        const { user } = this.state;

        if (!user) return null;
        
        return (
            <div>
                <Navigation />
                {!this.state.showNewRealEstate && (
                <p id="add-real-estate" onClick={this.toggleNewRealEstate}>
                    <img alt="." src="/images/plus.png" id="plus" />
                    Add Real-Estate
                </p>
                )}

                {this.state.showNewRealEstate ? (
                <NewRealEstate />
                ) : (
                <div id='real-estates-container'>
                    <div className='real-estate-card'>
                    <img alt='real-estate' src='/images/real_estate_example.png' className='real-estate-img' />
                    <div className='real-estate-info'>
                        <p className='real-estate-title'>Villa B dorm</p>
                        <p className='real-estate-text'>Location: Maldives</p>
                        <p className='real-estate-text'>Square Footage: 102 m2</p>
                        <p className='real-estate-text'>Number of Floors: 2</p>
                        <p className='real-estate-text state-color'>Accepted</p>
                    </div>
                    </div>
                </div>
                )}
            </div>
        )
    }
}