import React, { Component, useState } from "react";
import './NewRealEstate.css';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import axios from 'axios';


function LocationMarker({ onMapClick }) {
    const [position, setPosition] = useState(null);
    const [address] = useState('');

    const map = useMapEvents({
        click(e) {
            const clickedPosition = e.latlng;
            setPosition(clickedPosition);
            if (onMapClick) {
                onMapClick(clickedPosition);
            }
        },
        locationfound(e) {
            setPosition(e.latlng);
            console.log("TESTTTT");
            console.log(e.latlng);
            map.flyTo(e.latlng, map.getZoom());
        },
    });

    return position === null ? null : (
        <Marker position={position}>
            <Popup>{address}</Popup>
        </Marker>
    );
}



export class NewRealEstate extends Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedType: 'apartment',
            selectedCity: 'Novi Sad',
            address: '',
            selectedImage: null,
        };

        this.position = [45.23598471651923, 19.83932472361301]; // Initial map position

    }

    handleTypeChange = (event) => {
        this.setState({ selectedType: event.target.value });
    }

    handleCityChange = (event) => {
        if (typeof event == "string") {
            this.setState({ selectedCity: event });
        }
        else {
            this.setState({ selectedCity: event.target.value });
        }
    }

    handleAddressChange = (newAddress) => {
        axios
                .get(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${newAddress.lat}&lon=${newAddress.lng}`)
                .then((response) => {
                    console.log(response.data.address);
                    const obj = response.data.address;
                    if (obj.house_number != undefined)
                        this.setState({'address': response.data.address.road + " " + response.data.address.house_number + 
                        ", " + response.data.address.city_district});
                    else 
                        this.setState({'address': response.data.address.road + 
                        ", " + response.data.address.city_district});
                    this.handleCityChange(response.data.address.city_district);
                    return response.data;

                })
                .catch((error) => {
                    console.error('Error fetching address:', error);
                });
    }

    handleImageChange = (e) => {
        const file = e.target.files[0];
        this.setState({selectedImage: file})
    
        // here image can be uploaded to server
        console.log('Selected Image:', file);
      };
    
    render() {
        return (
            <div>
                <div id="new-real-estate-container-parent">
                    <div id="new-real-estate-container">
                        <p id="new-real-estate-title">New Real Estate</p>
                        <p className="new-real-estate-label">Type</p>
                        <select 
                            className="new-real-estate-select"
                            value={this.state.selectedType}
                            onChange={this.handleTypeChange}>
                            <option value="apartment">APARTMENT</option>
                            <option value="house">HOUSE</option>
                            <option value="villa">VILLA</option>
                        </select>
                        <p className="new-real-estate-label">City</p>
                        <select 
                            className="new-real-estate-select"
                            value={this.state.selectedCity}
                            onChange={this.handleCityChange}>
                            <option
                                value="novi-sad">NOVI SAD, SERBIA</option>
                            <option value="belgrade">BELGRADE, SERBIA</option>
                        </select>
                        <p className="new-real-estate-label">Address</p>
                        <input 
                            className="new-real-estate-input" 
                            type="text" 
                            name="address" 
                            placeholder="Type address or choose on the map"
                            value={this.state.address} // Set the input value from the state
                            onChange={(e) => this.setState({ address: e.target.value })}
                            />
                        
                        <div id="maps">
                            <MapContainer 
                                center={this.position} 
                                zoom={15} 
                                style={{ height: '100%', width: '100%' }}
                                >
                                <TileLayer
                                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                                />
                               <LocationMarker
                                    onMapClick={(clickedPosition) => {
                                        // handle the clicked position and address here
                                        this.handleAddressChange(clickedPosition);
                                    }}
                                />
                            </MapContainer>
                        </div>
                        <p className="new-real-estate-label">Square Footage (m2)</p>
                        <input 
                            className="new-real-estate-input" 
                            type="number" 
                            name="footage" 
                            min="0"
                            placeholder="Type square footage of the real estate..."
                            />
                        <p className="new-real-estate-label">Number of Floors</p>
                        <input 
                            className="new-real-estate-input" 
                            type="number" 
                            name="floors" 
                            min="1"
                            placeholder="Type number of floors..."
                            />
                        <br/>
                        <div 
                            id="upload-image-container"
                            onClick={() => this.fileInput.click()}>
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={this.handleImageChange}
                                    style={{ display: 'none' }}
                                    ref={(fileInput) => (this.fileInput = fileInput)}
                                />
                            <img id="upload-image" src="/images/photo.png"/>
                            <p id="upload-image-p">Upload image</p>
                        </div>
                        <span>
                            <button
                                id="cancel-button" className="btn">
                                    CANCEL
                            </button>
                            <button
                            id="confirm-button" className="btn">
                                CONFIRM
                            </button>
                        </span>
                    </div>
                </div>
                
            </div>
        )
    }

}