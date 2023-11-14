import React, { Component } from "react";
import './NewRealEstate.css';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';


export class NewRealEstate extends Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedType: 'apartment',
            selectedCity: 'Novi Sad',
        };

        this.position = [45.23598471651923, 19.83932472361301]; // Initial map position

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
                        
                        <div id="maps">
                            <MapContainer center={this.position} zoom={13} style={{ height: '100%', width: '100%' }}>
                                <TileLayer
                                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                                />
                                <Marker position={this.position}>
                                    <Popup>
                                    A pretty CSS3 popup. <br /> Easily customizable.
                                    </Popup>
                                </Marker>
                            </MapContainer>
                        </div>

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