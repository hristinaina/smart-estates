import React, { Component } from "react";
import './NewRealEstate.css';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';


export class NewRealEstate extends Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedType: 'apartment',
            selectedCity: 'Novi Sad',
        };
    }

    handleTypeChange = (event) => {
        this.setState({ selectedType: event.target.value });
    }

    handleCityChange = (event) => {
        this.setState({ selectedCity: event.target.value });
    }

    render() {

        return (
            <div>
                <div id="new-real-estate-container-parent">
                    <div id="new-real-estate-container">
                        <p id="new-real-estate-title">New Real Estate</p>
                        <p className="new-real-estate-label">Type</p>
                        <select 
                            id="new-real-estate-select"
                            value={this.state.selectedType}
                            onChange={this.handleTypeChange}>
                            <option value="apartment">APARTMENT</option>
                            <option value="house">HOUSE</option>
                            <option value="villa">VILLA</option>
                        </select>
                        <p className="new-real-estate-label">City</p>
                        <select 
                            id="new-real-estate-select"
                            value={this.state.selectedCity}
                            onChange={this.handleCityChange}>
                            <option
                                value="novi-sad">NOVI SAD, SERBIA</option>
                            <option value="belgrade">BELGRADE, SERBIA</option>
                        </select>
                        <p className="new-real-estate-label">Address</p>
                        <input 
                            className="new-real-estate-input" 
                            type="text" name="address" 
                            placeholder="Type address or choose on the map"/>
                        


                        {/* or use material design select */}
                        {/* <Select id="new-real-estate-select"
                            value={this.state.selectedValue}
                            onChange={this.handleChange}
                            label="Select an option"
                        >
                            <MenuItem value="">
                            </MenuItem>
                            <MenuItem value="apartment">Apartment</MenuItem>
                            <MenuItem value="house">House</MenuItem>
                            <MenuItem value="villa">Villa</MenuItem>
                        </Select> */}
                    </div>
                </div>
                
            </div>
        )
    }

}