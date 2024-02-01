
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

// todo prepraviti na tabelu umjesto grafa i uzimanje mejla kao sto je Tasija uradila
export class EVCharger extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            data: [],
            email: '',
            startDate: '',
            endDate: '',
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await SolarPanelService.getSPById(this.id); //todo use different service

        // todo update data (number of connections, maxValuePercentage)
        const updatedData =
        {
            ...device,
            //Value: lastValue,
        }

        const user = authService.getCurrentUser();
        this.Name = device.Device.Name;
        const historyData = await SolarPanelService.getSPGraphData(this.id, user.Email, "2023-12-12", "2023-12-23");
    
        this.setState({
            device: updatedData,
            data: historyData,
            email: user.Email,
            startDate: "2023-12-12",
            endDate: "2023-12-23",
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
    // todo change this with full data for a car
    handleMqttMessage(topic, message) {
        const { device } = this.state;
        const newValue = message.toString();
        const updatedData =
        {
            ...device,
            Value: newValue,
        }
        this.setState({
            device: updatedData,
        });
    }

    handleFormSubmit = async (e) => {
        e.preventDefault();

        const { email, startDate, endDate } = this.state;
        console.log(email, startDate, endDate);
        const historyData = await SolarPanelService.getSPGraphData(this.id, email, startDate, endDate);
        this.setState({
            data: historyData,
        });
    };

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
        const { device, data, email, startDate, endDate } = this.state;
        //todo na lijevoj karti dodati neki dinamicki ispis konekcija zavisno od broja i sakrikti prikaz
        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{this.Name}</span>
                <div className='sp-container'>
                    <div id="sp-left-card">
                        <p className='sp-card-title'>Device Data</p>
                        <p className='sp-data-text'>ChargingPower:</p>
                        <p><b>{device.NumberOfPanels} </b> </p>
                        <p className='sp-data-text'>Max Charging Percentage</p>
                        <p><b>{device.SurfaceArea}</b></p> 
                        
                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Actions History</p>
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
                        <SPGraph data={data} />
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
