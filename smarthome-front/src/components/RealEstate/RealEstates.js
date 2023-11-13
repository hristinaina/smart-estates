
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
                    <img alt="." src="/images/plus.png" id="plus" /> Add Real-Estate</p>
            </div>
        )
    }
 }