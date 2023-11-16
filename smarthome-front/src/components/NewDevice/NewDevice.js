import React, { Component, useState } from "react";
import './NewDevice.css';
import { Link } from 'react-router-dom';


export class NewDevice extends Component {

    constructor(props) {
        super(props);
        this.state = {
            selectedType: 'buba',
            selectedImage: null,
            imagePreview: null,
            showNewDevice: true,
            showNewNekiTip: false,
            showNewNekiTip2: false,
        };
    }

    toggleTip1 = () => {
        this.setState((prevState) => ({
            showNewDevice: !prevState.showNewDevice,
            showNewNekiTip: !prevState.showNewNekiTip
        }));
    }

    toggleTip2 = () => {
        this.setState((prevState) => ({
            showNewDevice: !prevState.showNewDevice,
            showNewNekiTip2: !prevState.showNewNekiTip2
        }));
    }

    types = [
        { vlaue: 'buba', label: 'prikazz' },
        { value: 'bebica', label: 'berbicaaa' },
    ];

    handleTypeChange = (event) => {
        this.setState({ selectedType: event.target.value });
    }
    
    handleImageChange = (event) => {
        const file = event.target.files[0];
    
        this.setState({
            selectedImage: file,
            imagePreview: URL.createObjectURL(file),
          });
      };

    render() {
        const types = this.types;

        return (
            <div>
                <div id="new-real-estate-container-parent">
                    <div id="new-real-estate-container">
                        <p id="new-real-estate-title">New Device</p>
                        <p className="new-real-estate-label">Name</p>
                        <input
                            className="new-real-estate-input"
                            type="number"
                            name="footage"
                            min="0"
                            placeholder="Type the name of the device"
                        />
                        <p className="new-real-estate-label">Type</p>
                        <select
                            className="new-real-estate-select"
                            value={this.state.selectedType}
                            onChange={this.handleTypeChange}>
                            {types.map((type) => (
                                <option key={type.value} value={type.value}>
                                    {type.label}
                                </option>
                            ))}
                        </select>
                        <br />
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
                            <img id="upload-image" src="/images/photo.png" />
                            <p id="upload-image-p">Upload image</p>
                        </div>
                        {/* Show choosen image */}
                        {this.state.imagePreview && (
                            <div>
                                <img className='cropped-image' src={this.state.imagePreview} alt="Device Preview" />
                            </div>
                        )}
                        <span>
                            <Link to='/devices'>
                                <button
                                    id="cancel-button" className="btn">
                                    CANCEL
                                </button>
                            </Link>
                            <button
                                id="confirm-button" className="btn">
                                NEXT
                            </button>
                        </span>
                    </div>
                </div>

            </div>
        )
    }

}