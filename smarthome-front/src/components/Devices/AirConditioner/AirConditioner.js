import React, { Component } from 'react';
import '../Devices.css';
import'./AirConditioner.css'
import { Navigation } from '../../Navigation/Navigation';
import authService from '../../../services/AuthService'
import mqtt from 'mqtt';
import { Autocomplete, TextField, Button, Box, Grid, IconButton, Snackbar, Stack, Typography, Switch, Input, FormControl } from '@mui/material';
import DeviceService from '../../../services/DeviceService';


export class AirConditioner extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            mode: [],
            switchOn: false,
            data: [],
            email: '',
            startDate: '',
            endDate: '',
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
            temp: 20,
            currentTemp: 20.0
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/ac/');
        this.setState({mode: device.Mode.split(',')})
        this.setState({device: device})
    
        const updatedData =
        {
            ...device,
            Value: "Loading...",
        }
        console.log(device);

        const user = authService.getCurrentUser();
        // this.Name = device.Device.Name;
        // const historyData = await DeviceService.getSPGraphData(this.id, user.Email, "2023-12-12", "2023-12-23");
        // this.setState({
        //     device: updatedData,
        //     switchOn: device.IsOn,
        //     data: historyData,
        //     email: user.Email,
        //     startDate: "2023-12-12",
        //     endDate: "2023-12-23",
        // });

        try {
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-ac",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('ac/temp');
                    // this.mqttClient.subscribe('ac/data/' + this.id);
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

    handleSwitchToggle = () => {
        // const topic = "sp/switch/" + this.id;

        // this.setState((prevState) => ({
        //     switchOn: !prevState.switchOn,
        // }));
        // var message = {
        //     "IsOn": (!this.state.switchOn),
        //     "UserEmail": authService.getCurrentUser().Email,
        // }
        // this.mqttClient.publish(topic, JSON.stringify(message));

        // this.setState({ snackbarMessage: "Successfully changed switch state!" });
        // this.handleClick();
    };

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        this.setState({
            currentTemp: JSON.parse(message.toString()).temp
        });
    }

    // handleFormSubmit = async (e) => {
    //     e.preventDefault();

    //     const { email, startDate, endDate } = this.state;
    //     console.log(email, startDate, endDate);
    //     // const historyData = await DeviceService.getSPGraphData(this.id, email, startDate, endDate);
    //     // this.setState({
    //     //     data: historyData,
    //     // });
    // };

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
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
        const { device, mode, switchOn, data, email, startDate, endDate, temp, currentTemp} = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{this.Name}</span>
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Supported Modes</p>
                        <div style={{marginBottom: "25px"}}>
                            <span className='ac-current-temp'>Current temp:  </span>
                            {/* <p><b>{device.Value}</b></p> */}
                            <span><b>{ currentTemp }</b></span>
                        </div>                                                 
                        {mode.map((item) => (
                        <div key={item} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                            <Typography style={{ fontSize: '1.1em' }}>Off</Typography>
                            <Switch
                                checked={switchOn}
                                onChange={this.handleSwitchToggle}
                            />
                            <Typography style={{ fontSize: '1.1em' }}>On</Typography>
                            <span style={{ flex: 1 }}>{item}</span>
                            {item !== 'Ventilation' && (
                                <FormControl style={{ width: '80px' }}>
                                    <Input
                                        type="number"
                                        value={temp}
                                        onChange={this.handleTemperatureChange}
                                    />
                                </FormControl>
                            )}
                        </div>
                    ))}

                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Switch History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                Email:
                                <TextField style={{ backgroundColor: "white" }} type="text" value={email} onChange={(e) => this.setState({ email: e.target.value })} />
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
                            <Button type="submit" id='sp-data-button'>Fetch Data</Button>
                        </form>
                        {/* <SPGraph data={data} /> */}
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
