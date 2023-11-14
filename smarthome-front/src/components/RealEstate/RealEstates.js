
import { Component } from 'react';
import './RealEstates.css';

export class RealEstates extends Component {

    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <p id="add-real-estate">
                    <img alt="." src="/images/plus.png" id="plus" /> 
                    Add Real-Estate
                </p>
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
        )
    }
 }