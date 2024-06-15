
import { Component } from 'react';
import '../Devices.css';
import { Navigation } from '../../Navigation/Navigation';
import mqtt from 'mqtt';
import authService from '../../../services/AuthService';
import 'chart.js/auto';
import SPGraph from '../SolarPanel/SPGraph';
import { TextField } from '@mui/material';
import { Button } from 'reactstrap';
import './EVCharger.css'
import { Snackbar } from "@mui/material";
import SolarPanelService from '../../../services/SolarPanelService';
import EVChargerService from '../../../services/EVChargerService';
import TableOfActions from './TableOfActions';
import PermissionService from '../../../services/PermissionService';

// todo prepraviti na tabelu umjesto grafa i uzimanje mejla kao sto je Tasija uradila
export class EVCharger extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            data: [],  //za tabelu
            connectionData: [],  // lista objekata car-simulation
            emailInput: '',
            userEmails: [],
            startDate: '',
            endDate: '',
            inputPercentage: 90,
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";
        this.email = "";
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await EVChargerService.get(this.id);

        let connectionData = [];
        for (let i = 0; i < device.Connections; i++) {
            connectionData.push({active: false})
        }

        const percentage = await EVChargerService.getLastPercentage(this.id);
        const user = authService.getCurrentUser();
        this.Name = device.Device.Name;
        this.email =  user.Email;
        const historyData = await EVChargerService.getTableActions(this.id, user.Email, "2023-12-12", "2024-02-07");
        //todo change this upper method
        let users = await PermissionService.getPermissions(this.id, user.EstateId);
        users.push(user.Email);
        users.push("all");
        users.push("none");
        this.setState({
            device: device,
            connectionData: connectionData,
            data: historyData,
            emailInput: user.Email,
            userEmails: users,
            startDate: "2023-12-12",
            endDate: "2024-02-07",
            inputPercentage: parseInt(percentage *100),
        });

        try {
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-evc",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('ev/data/' + this.id);
                    this.mqttClient.subscribe('ev/start/' + this.id);
                    this.mqttClient.subscribe('ev/end/' + this.id);
                });

                // Handle incoming MQTT messages
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessage(topic, message);
                });
            }
        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        const result = JSON.parse(message.toString())
        const { connectionData } = this.state;
        connectionData[result.PlugId] = {
            active: result.Active,
            maxCapacity: parseInt(result.MaxCapacity),
            currentCapacity: result.CurrentCapacity.toFixed(1),
        }
        this.setState({
            connectionData: connectionData,
        });
    }

    handleFormSubmit = async (e) => {
        e.preventDefault();
        const { emailInput, startDate, endDate } = this.state;
        console.log(emailInput, startDate, endDate);
        this.setState({
            data: {},
        });
        if(new Date(startDate) > new Date(endDate)) {
            this.setState({ snackbarMessage: "Start date must be before end date" });
            this.handleClick();
            return 
        }

        const historyData = await EVChargerService.getTableActions(this.id, emailInput, startDate, endDate);
        this.setState({
            data: historyData,
        });
    };

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleDateChange = (fieldName, event) => {
        this.setState({ [fieldName]: event.target.value });
    };

    handleButtonPercentageClick= () => {
        const topic = "ev/percentage/" + this.id;
        const { inputPercentage } = this.state;
        var message = {
            "CurrentCapacity": inputPercentage/100,
            "Email": authService.getCurrentUser().Email,
        }
        this.mqttClient.publish(topic, JSON.stringify(message));

        this.setState({ snackbarMessage: "Successfully changed maximum charging percentage!" });
        this.handleClick();
        console.log(inputPercentage);
    }

    handleBackArrow() {
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
        const { device, data, emailInput, startDate, endDate, inputPercentage, connectionData, userEmails } = this.state;
        const connectionsArray = Array.from({ length: device.Connections }, (_, index) => index + 1);
        console.log(connectionData);
        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{this.Name}</span>
                <div className='sp-container'>
                    <div id="sp-left-card">
                        <p className='sp-card-title'>Electrical Vehicle Charger</p>
                        <div className='box-container'>
                            <div className='ev-box'>
                                <p className='sp-data-text'>ChargingPower:</p>
                                <p ><b>{device.ChargingPower} kw</b></p>
                            </div>
                            <div className='ev-box'>
                                <p className='sp-data-text'>Max Charging Percentage:</p>
                                <div className='box-container'>
                                    <input
                                        className="new-real-estate-input"
                                        type="number"
                                        name="charging"
                                        maxLength="3"
                                        value={inputPercentage}
                                        onChange={(e) => this.setState({ inputPercentage: e.target.value })}
                                        style={{ display: "inline", width: "70px", marginLeft: "20px" }}
                                    />
                                    <Button className="ev-button" style={{ width: "80px", marginLeft: "15px" }} onClick={this.handleButtonPercentageClick}>Update</Button>
                                </div>
                            </div>
                        </div>
                        <div id="connections-container">
                            {connectionsArray.map((index) => (
                                <div key={index} className="connection-box">
                                    <p className="mark"><b>Plug {index}: </b></p>
                                    {connectionData[index-1].active &&
                                        (<div className='box-container'>
                                            <img src="/images/car.png" alt={`Car ${index}`} className="car-image" />
                                            <p style={{marginLeft: "20px" }}>{connectionData[index-1].currentCapacity}/{connectionData[index-1].maxCapacity}kwh</p>
                                            <p className='ev-right-data'><b>{parseInt(connectionData[index-1].currentCapacity/connectionData[index-1].maxCapacity*100)}%</b></p>
                                            <img src="/images/charging.png" className='charger-image' alt={`Charging ${index}`}/>
                                        </div>)}
                                    {!connectionData[index-1].active && (<p style={{marginLeft: "45px" }}>Free</p>)}
                                </div>
                            ))}
                        </div>
                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Actions History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                Email:
                                <select style={{width: "200px", cursor: "pointer"}}
                                    className="new-real-estate-select"
                                    value={emailInput}
                                    onChange={(e) => this.setState({ emailInput: e.target.value })}>
                                    {userEmails.map(email => (
                                        <option key={email} value={email}>
                                        {email}
                                        </option>
                                    ))}
                                </select>
                            </label>
                            <br />
                            <label>
                                Start Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={startDate} onChange={(e) => this.setState({ startDate: e.target.value })} />
                            </label>
                            <br />
                            <label>
                                End Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={endDate} onChange={(e) => this.setState({ endDate: e.target.value })} />
                            </label>
                            <br />
                            <Button type="submit" id='sp-data-button' className='button-height'>Fetch Data</Button>
                        </form>
                        <TableOfActions logData={data} />
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
