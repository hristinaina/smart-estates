import React, { Component, useState } from "react";
import './NewRealEstate.css';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import axios from 'axios';
import { RealEstates } from "./RealEstates";
import { Navigation } from "../Navigation/Navigation";
import RealEstateService from "../../services/RealEstateService";
import { Snackbar } from "@mui/material";
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';
import ImageService from "../../services/ImageService";


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
            userId: 2,
            selectedType: '0',
            selectedCity: 'Novi Sad',
            searchCity: '',
            address: '',
            selectedImage: null,
            showRealEstates: false,
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
            imagePreview: null,
        };

        this.position = [45.23598471651923, 19.83932472361301]; // Initial map position

    }

    cities = [
        { vlaue: 'novi-sad', label: 'NOVI SAD, SERBIA'},
        { value: 'belgrade', label: 'BELGRADE, SERBIA'},
        { value: 'zrenjanin', label: 'ZRENJANIN, SREBIA'}
    ];

    getFilteredCities = () => {
        return this.cities.filter((city) =>
            city.label.toLowerCase().includes(this.state.searchCity.toLowerCase())
        );
    };

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
                const obj = response.data.address;
                if (obj.house_number !== undefined)
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
        this.setState({selectedImage: file, imagePreview: URL.createObjectURL(file)})
    
        // here image can be uploaded to server
        console.log('Selected Image:', file);
    }

    confirm = async () => {
        console.log("opet je usao");
        var estate = {
            "Name": document.getElementsByName("name")[0].value.trim(),
            "Type": Number(this.state.selectedType),
            "Address": this.state.address.trim(),
            "City": this.state.selectedCity.trim(),
            "SquareFootage": Number(document.getElementsByName("footage")[0].value),
            "NumberOfFloors": Number(document.getElementsByName("floors")[0].value),
            "Picture": "blabla",
            "State": 0,
            "User": this.state.userId,
        }

        try {
            const result = await RealEstateService.add(estate);
            console.log(result);
            this.handleUpload()
        } catch (error) {
            console.log("Error")
            console.error(error);
            this.setState({snackbarMessage: "Please check input fields!"})
            this.handleClick()
        }
    }

    cancel = () => {
        window.location.href = '/real-estates';
    }

     // snackbar
    handleClick = () => {
        this.setState({open: true});
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
          return;
        }
        this.setState({open: false});
      };

    action = (
        <React.Fragment>
            <IconButton
            size="small"
            aria-label="close"
            color="inherit"
            onClick={this.handleClose}>
            <CloseIcon fontSize="small" />
            </IconButton>
        </React.Fragment>
        );
    

    
    handleUpload = async () => {
        if (!this.state.selectedImage) {
            this.setState({snackbarMessage: "Please check input fields!"});
            this.handleClick();
            return;
        }

        const formData = new FormData();
        formData.append('image', this.state.selectedImage);
        
        try {
            var name = String(document.getElementsByName('name')[0].value).trim();
            const substr = this.state.selectedImage.name.split(".")[1].trim();
            name += "." + substr;
            await ImageService.uploadImage(formData, name);
            window.location.href = '/real-estates';
        } catch (error) {
            console.log("Error");
            console.error(error);
            this.setState({snackbarMessage: "Error uploading image!"});
            this.handleClick();
        }
    };
    
    render() {
        const filteredCities = this.getFilteredCities();
        return (
            <div>
                <Navigation/>
                {this.state.showRealEstates ? (
                <RealEstates />
                ) : (
                <div id="new-real-estate-container-parent">
                    <div id="new-real-estate-container">
                        <p id="new-real-estate-title">New Real Estate</p>
                        <p className="new-real-estate-label">Name</p>
                        <input 
                            className="new-real-estate-input" 
                            type="text" 
                            name="name" 
                            maxLength="50"
                            placeholder="Type the name of your real estate"
                            />
                        <p className="new-real-estate-label">Type</p>
                        <select 
                            className="new-real-estate-select"
                            value={this.state.selectedType}
                            onChange={this.handleTypeChange}>
                            <option value='0'>APARTMENT</option>
                            <option value='1'>HOUSE</option>
                            <option value="2">VILLA</option>
                        </select>
                        <p className="new-real-estate-label">City and country</p>
                        <select 
                            className="new-real-estate-select"
                            value={this.state.selectedCity}
                            onChange={this.handleCityChange}>
                            {filteredCities.map((city) => (
                                <option key={city.value} value={city.value}>
                                    {city.label}
                                </option>
                            ))}
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
                            <img alt="Real Estate" id="upload-image" src="/images/photo.png"/>
                            <p id="upload-image-p">Upload image</p>
                            {this.state.imagePreview && (
                            <div>
                                <img className='cropped-image' src={this.state.imagePreview} alt="Uploaded Image Preview" />
                            </div>
                        )}
                        </div>
                        <span>
                            <button
                                id="cancel-button" className="btn" onClick={this.cancel}>
                                    CANCEL
                            </button>
                            <button
                            id="confirm-button" className="btn" onClick={this.confirm}>
                                CONFIRM
                            </button>
                        </span>
                    </div>
                </div> )}

                <Snackbar
                    open={this.state.open}
                    autoHideDuration={3000}
                    onClose={this.handleClose}
                    message={this.state.snackbarMessage}
                    action={this.action}
                />
            </div>
        )
    }

}