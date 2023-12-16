import { Component } from "react";
import authService from "../../../services/AuthService";
import { Navigation } from "../../Navigation/Navigation";
import VehicleGateService from "../../../services/VehicleGateService";
import './VehicleGate.css';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import { TextField } from '@mui/material';
import { Button } from 'reactstrap';
import mqtt from 'mqtt';

export class VehicleGate extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            licensePlate: '',
            startDate: '',
            endDate: '',
            enterLicensePlate: '',
            enter: false,
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = '';
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await VehicleGateService.get(this.id);

        const user = authService.getCurrentUser();
        this.Name = device.ConsumptionDevice.Device.Name;
        await this.setState({device: device});

        try {        
            // mqtt connection
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                        clientId: "react-front-nvt-2023-lamp",
                        clean: false,
                        keepalive: 60
                });
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('vg/open/' + this.id);
                });
    
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessage(topic, message);
                });
            }
        } catch (error) {
            console.error(error);
        }
        
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleModeChange = async(mode) => {
        let device = this.state.device;
        if (device.Mode != mode) {
            device.Mode = mode;
            if (mode === 0) {
                await VehicleGateService.toPrivate(this.id);
            } else {
                await VehicleGateService.toPublic(this.id);
            }
        }
        await this.setState({device: device})
    }

    async handleMqttMessage(topic, message) {
        const tokens = message.toString().split('+');
        let device = this.state.device;
        console.log(message.toString());
        if (tokens[0] == "open") {
            device.IsOpen = true;
            if (tokens[2] == "enter") {
                await this.setState({device: device, enterLicensePlate: tokens[1], enter: true});
            }
            else {
                await this.setState({device: device, enterLicensePlate: tokens[1], enter: false});
            }

        }
        else {
            device.IsOpen = false;
            await this.setState({device: device, enterLicensePlate: '', enter: false});
        }
    }

    render() {            
        const { licensePlate, startDate, endDate} = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow}/>
                <span className='estate-title'>{this.Name}</span>
                <div className="sp-container">
                    <div id="sp-left-card">
                        <p className="sp-card-title">Device Data</p>
                        <p className="sp-data-text">Mode</p>
                        <p className="vg-description">{this.state.device.Mode === 0 ? 'Private' : 'Public'}</p>
                        <img src='/images/private.png' className={`vg-icon vg-padlock ${this.state.device.Mode === 1 ? 'unlocked': ''}`} onClick={ () => this.handleModeChange(0)}/>
                        <img src='/images/public.png' className={`vg-icon vg-padlock ${this.state.device.Mode === 0 ? 'unlocked': ''}`} onClick={ () => this.handleModeChange(1)}/>
                        <p className="sp-data-text">State</p>
                        <p className="vg-description">{this.state.device.IsOpen === true ? 'Opened' : 'Closed'}</p>
                        <p className="vg-description">{this.state.enterLicensePlate} {this.state.enter === true ? ' is entering...' : ''}</p>
                        <img src='/images/closed-gate.png' className={`vg-icon ${this.state.device.IsOpen === true ? 'unlocked' : ''}`} />
                        <img src='/images/opened-gate.png' className={`vg-icon ${this.state.device.IsOpen === false ? 'unlocked' : ''}`} />
                        <div id="vg-box">
                        <p className="sp-data-text">Trusted License Plates</p>
                        <List id="vg-list">
                            <ListItem>    
                                <ListItemText primary="NS-123-45"/>
                            </ListItem>
                            <ListItem>
                                <ListItemText primary="NS-123-56"/>
                            </ListItem>
                        </List>
                        </div>
                        <span className='vg-description vg-add'>Add License Plate</span>
                    </div>

                    <div id="sp-right-card">
                            <p className="sp-card-title">Reports</p>
                            <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                License Plate:
                                <TextField style={{ backgroundColor: "white" }} type="text" value={licensePlate} onChange={(e) => this.setState({ licensePlate: e.target.value })} />
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
                            <Button id='sp-data-button'>Confirm</Button>
                        </form>
                    </div>
                </div>
            </div>
        )
    }
}