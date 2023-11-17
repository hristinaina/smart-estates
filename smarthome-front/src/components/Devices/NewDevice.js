import React, { Component } from "react";
import './Devices.css';
import { Link } from 'react-router-dom';


export class NewDevice extends Component {

    constructor(props) {
        super(props);
        this.state = {
            selectedType: 'ambient-sensor',
            selectedImage: null,
            imagePreview: null,
            name: "",
            powerConsumption: 200,
            minTemp: 16,
            maxTemp: 31,
            batterySize: 13,
            chargingPower: 2.3,
            connections: 1,
            selectedPowerSupply: 'autonomous',
            showPowerSupply: true,
            showPowerConsumption: false,
            showAirConditioner: false,
            showBatterySize: false,
            showCharger: false,
            isButtonDisabled: true,
        };
    }

    types = [
        { value: 'ambient-sensor', label: 'Ambient Sensor' },
        { value: 'air', label: 'Air conditioner' },
        { value: 'washing-machine', label: 'Washing machine' },
        { value: 'lamp', label: 'Lamp' },
        { value: 'vehicle-gate', label: 'Vehicle gate' },
        { value: 'sprinkler', label: 'Sprinkler' },
        { value: 'solar-panel', label: 'Solar panel' },
        { value: 'battery-storage', label: 'Battery storage' },
        { value: 'electric-vehicle-charger', label: 'Electric vehicle charger' },
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

            if (selectedType === 'solar-panel' || selectedType === 'battery-storage' || selectedType === 'electric-vehicle-charger') {
                this.setState({
                    showPowerSupply: false,
                    showPowerConsumption: false,
                });
            } else {
                this.setState({
                    showPowerSupply: true,
                });
                if(this.state.selectedPowerSupply==="home"){
                    this.setState({
                        showPowerConsumption: true,
                    });
                }
            }

            if (selectedType === 'battery-storage') {
                this.setState({
                    showBatterySize: true,
                });
            } else if (selectedType === 'air') {
                this.setState({
                    showAirConditioner: true,
                });
            } else if (selectedType === 'electric-vehicle-charger') {
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
            if (selectedPowerSupply === 'home') {
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
        if (this.state.name.trim() === '' || this.state.selectedImage === null) {
            this.setState({ isButtonDisabled: true })
        }
        else {
            if (this.state.selectedType === 'air' && (this.state.minTemp>= this.state.maxTemp ||this.state.minTemp< -40 || this.state.maxTemp > 60)){
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedType === 'battery-storage' && (this.state.batterySize> 100000 ||this.state.batterySize< 1)){
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedPowerSupply === 'home' && (this.state.powerConsumption> 60000 ||this.state.powerConsumption<= 0) 
            && (this.state.selectedType != 'solar-panel' && this.state.selectedType != 'battery-storage' && this.state.selectedType != 'electric-vehicle-charger')){
                this.setState({ isButtonDisabled: true })
            }
            else if (this.state.selectedType === 'electric-vehicle-charger' && (this.state.connections < 1 ||this.state.connections> 20
                ||this.state.chargingPower< 1 ||this.state.chargingPower> 360)){
                this.setState({ isButtonDisabled: true })
            }
            else{
                this.setState({ isButtonDisabled: false })
            }
        }
    }

    createDevice = () => {
        console.log("api sent");
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
                                    <option value="autonomous">Autonomous</option>
                                    <option value="home">Home (home battery/network)</option>
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
                                <p className="new-real-estate-label">Minimal temperature (celsius):</p>
                                <input
                                    className="new-real-estate-input"
                                    type="number"
                                    name="min-temp"
                                    placeholder="Enter the min temp of your air conditioner (in celsius)"
                                    value={this.state.minTemp}
                                    onChange={this.handleMinTemp}
                                />
                                <p className="new-real-estate-label">Maximal temperature: (celsius)</p>
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
                                id="confirm-button" className={`btn ${this.state.isButtonDisabled ? 'disabled' : ''}`} disabled={this.state.isButtonDisabled}
                                 onClick={this.createDevice}>
                                CONFIRM
                            </button>
                        </span>
                    </div>
                </div>
            </div>
        )
    }
}