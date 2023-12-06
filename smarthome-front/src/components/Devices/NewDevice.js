import React, { Component } from "react";
import './Devices.css';
import { Link } from 'react-router-dom';
import DeviceService from "../../services/DeviceService";
import ImageService from "../../services/ImageService";
import authService from '../../services/AuthService'
import { Snackbar } from "@mui/material";


export class NewDevice extends Component {

    constructor(props) {
        super(props);
        this.state = {
            selectedType: 0,
            selectedImage: null,
            imagePreview: null,
            name: "",
            powerConsumption: 200,
            minTemp: 16,
            maxTemp: 31,
            batterySize: 13,
            chargingPower: 2.3,
            connections: 1,
            selectedPowerSupply: 0,
            showPowerSupply: true,
            showPowerConsumption: false,
            showAirConditioner: false,
            showBatterySize: false,
            showCharger: false,
            isButtonDisabled: true,
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
        };
        this.id = parseInt(localStorage.getItem("real-estate"));
    }

    types = [
        { value: 0, label: 'Ambient Sensor' },
        { value: 1, label: 'Air conditioner' },
        { value: 2, label: 'Washing machine' },
        { value: 3, label: 'Lamp' },
        { value: 4, label: 'Vehicle gate' },
        { value: 5, label: 'Sprinkler' },
        { value: 6, label: 'Solar panel' },
        { value: 7, label: 'Battery storage' },
        { value: 8, label: 'Electric vehicle charger' },
    ];

    handleTypeChange = (event) => {
        const selectedType = event.target.value;

        this.setState({ selectedType }, () => {
            // Callback function after state is updated
            this.setState({
                showAirConditioner: false,
                showBatterySize: false,
                showCharger: false,
            });

            if (selectedType == 6 || selectedType == 7 || selectedType == 8) {
                this.setState({
                    showPowerSupply: false,
                    showPowerConsumption: false,
                });
            } else {
                this.setState({
                    showPowerSupply: true,
                });
                if (this.state.selectedPowerSupply == 1) {
                    this.setState({
                        showPowerConsumption: true,
                    });
                }
            }

            if (selectedType == 7) {
                this.setState({
                    showBatterySize: true,
                });
            } else if (selectedType == 1) {
                this.setState({
                    showAirConditioner: true,
                });
            } else if (selectedType == 8) {
                this.setState({
                    showCharger: true,
                });
            }
            this.checkButton();
        });
    };


    handlePowerSupplyChange = (event) => {
        const selectedPowerSupply = event.target.value;

        this.setState({ selectedPowerSupply }, () => {
            if (selectedPowerSupply == 1) {
                this.setState({
                    showPowerConsumption: true,
                });
            } else {
                this.setState({
                    showPowerConsumption: false,
                });
            }
            this.checkButton();
        });
    };

    handleNameChange = (event) => {
        const name = event.target.value;

        this.setState({ name }, () => {
            this.checkButton();
        });
    }

    handleChargingPower = (event) => {
        const chargingPower = event.target.value;

        this.setState({ chargingPower }, () => {
            this.checkButton();
        });
    }

    handleConnections = (event) => {
        const connections = event.target.value;

        this.setState({ connections }, () => {
            this.checkButton();
        });
    }

    handleMinTemp = (event) => {
        const minTemp = event.target.value;

        this.setState((prevState) => ({
            ...prevState,
            minTemp,
        }), () => {
            this.checkButton();
        });
    }

    handleMaxTemp = (event) => {
        const maxTemp = event.target.value;

        this.setState((prevState) => ({
            ...prevState,
            maxTemp,
        }), () => {
            this.checkButton();
        });
    }

    handleBatterySize = (event) => {
        const batterySize = event.target.value;

        this.setState((prevState) => ({
            ...prevState,
            batterySize,
        }), () => {
            this.checkButton();
        });
    }

    handlePowerConsumption = (event) => {
        const powerConsumption = event.target.value;

        this.setState((prevState) => ({
            powerConsumption,
        }), () => {
            this.checkButton();
        });
    }


    handleImageChange = (event) => {
        const file = event.target.files[0];

        this.setState(
            {
                selectedImage: file,
                imagePreview: URL.createObjectURL(file),
            },
            () => {
                this.checkButton();
            }
        );
    };

    checkButton = () => {
        if (this.state.name.trim() == '' || this.state.selectedImage == null) {
            this.setState({ isButtonDisabled: true })
        }
        else {
            if (this.state.selectedType == 1 && (this.state.minTemp >= this.state.maxTemp || this.state.minTemp < -40 || this.state.maxTemp > 60)) {
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedType == 7 && (this.state.batterySize > 100000 || this.state.batterySize < 1)) {
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedPowerSupply == 1 && (this.state.powerConsumption > 60000 || this.state.powerConsumption <= 0)
                && (this.state.selectedType != 6 && this.state.selectedType != 7 && this.state.selectedType != 8)) {
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedType == 8 && (this.state.connections < 1 || this.state.connections > 20
                || this.state.chargingPower < 1 || this.state.chargingPower > 360)) {
                this.setState({ isButtonDisabled: true })
            }
            else {
                this.setState({ isButtonDisabled: false })
            }
        }
    }

    createDevice = async () => {
        console.log("api for new device sent");
        try {
            const data = {
                Name: this.state.name,
                Type: parseInt(this.state.selectedType),
                Picture: "/images/" + this.state.selectedImage.name,
                RealEstate: this.id,
                PowerSupply: parseInt(this.state.selectedPowerSupply),
                PowerConsumption: parseFloat(this.state.powerConsumption),
                MinTemperature: parseInt(this.state.minTemp),
                MaxTemperature: parseInt(this.state.maxTemp),
                ChargingPower: parseFloat(this.state.chargingPower),
                Connections: parseInt(this.state.connections),
                Size: parseFloat(this.state.batterySize),
                UserId: authService.getCurrentUser().Id,
            };
            const result = await DeviceService.createDevice(data);
            console.log(result);
            // uploading image
            const formData = new FormData();
            formData.append('image', this.state.selectedImage);
            var name = String(document.getElementsByName('name')[0].value).trim();
            const substr = this.state.selectedImage.name.split(".")[1].trim();
            name += "." + substr;
            await ImageService.uploadImage(formData, "devices&" + name);
            window.location.assign("/devices")
        } catch (error) {
            this.setState({ snackbarMessage: error.message });
            this.handleClick();
        }
    };

    cancel() {
        window.location.assign("/devices")
    }

    // snackbar
    handleClick = () => {
        this.setState({ open: true });
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        this.setState({ open: false });
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
                            type="text"
                            name="name"
                            placeholder="Type the name of the device"
                            value={this.state.name} // Set the input value from the state
                            onChange={this.handleNameChange}
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
                        {this.state.showPowerSupply && (
                            <div>
                                <p className="new-real-estate-label">Power supply type:</p>
                                <select
                                    className="new-real-estate-select"
                                    value={this.state.selectedPowerSupply}
                                    onChange={this.handlePowerSupplyChange}>
                                    <option value="0">Autonomous</option>
                                    <option value="1">Home (home battery/network)</option>
                                </select>
                            </div>
                        )}
                        {this.state.showPowerConsumption && (
                            <div>
                                <p className="new-real-estate-label">Power consumption (watts):</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="power-consumption"
                                    placeholder="Enter the power consumption of your device (in watts)"
                                    value={this.state.powerConsumption}
                                    onChange={this.handlePowerConsumption}
                                />
                            </div>
                        )}
                        {this.state.showAirConditioner && (
                            <div>
                                <p className="new-real-estate-label">Minimum temperature (celsius):</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="min-temp"
                                    placeholder="Enter the min temp of your air conditioner (in celsius)"
                                    value={this.state.minTemp}
                                    onChange={this.handleMinTemp}
                                />
                                <p className="new-real-estate-label">Maximum temperature: (celsius)</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="max-temp"
                                    placeholder="Enter the min temp of your air conditioner (in celsius)"
                                    value={this.state.maxTemp}
                                    onChange={this.handleMaxTemp}
                                />
                            </div>
                        )}
                        {this.state.showBatterySize && (
                            <div>
                                <p className="new-real-estate-label">Battery size (kWh):</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="battery-size"
                                    placeholder="Enter the battery size (in kWh)"
                                    value={this.state.batterySize}
                                    onChange={this.handleBatterySize}
                                />
                            </div>
                        )}
                        {this.state.showCharger && (
                            <div>
                                <p className="new-real-estate-label">Charging power (kWatts):</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="charging-power"
                                    placeholder="Enter the charging power (in kwatts)"
                                    value={this.state.chargingPower}
                                    onChange={this.handleChargingPower}
                                />
                                <p className="new-real-estate-label">Number of connections:</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="connections"
                                    placeholder="Enter the number of connections"
                                    value={this.state.connections}
                                    onChange={this.handleConnections}
                                />
                            </div>
                        )}
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
                            <img id="upload-image" src="/images/photo.png" alt="icon" />
                            <p id="upload-image-p">Upload image</p>
                        </div>
                        {/* Show choosen image */}
                        {this.state.imagePreview && (
                            <div>
                                <img className='cropped-image' src={this.state.imagePreview} alt="Device Preview" />
                            </div>
                        )}
                        <span>
                            <button
                                id="cancel-button" className="btn" onClick={this.cancel}>
                                CANCEL
                            </button>
                            <button
                                id="confirm-button" className={`btn ${this.state.isButtonDisabled ? 'disabled' : ''}`} disabled={this.state.isButtonDisabled}
                                onClick={this.createDevice}>
                                CONFIRM
                            </button>
                        </span>
                    </div>
                </div>
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