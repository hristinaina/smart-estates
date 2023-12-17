import React, { Component } from "react";
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
import Dialog from "../../Dialog/Dialog";
import { Snackbar } from "@mui/material";
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';


export class VehicleGate extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {LicensePlates: []},
            licensePlate: '',
            startDate: '',
            endDate: '',
            enterLicensePlate: '',
            enter: false,
            exit: false,
            showAddLicensePlateDialog: false,
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = '';
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await VehicleGateService.get(this.id);
        console.log(device);

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
                await this.setState({device: device, enterLicensePlate: tokens[1], enter: true, exit: false});
            }
            else {
                await this.setState({device: device, enterLicensePlate: tokens[1], enter: false, exit: true});
            }

        }
        else {
            device.IsOpen = false;
            await this.setState({device: device, enterLicensePlate: '', enter: false, exit: false});
        }
    }

    openAddLicensePlateDialog = async () => {
        await this.setState({showAddLicensePlateDialog: true})
    }

    handleCancel = async () => {
        await this.setState({showAddLicensePlateDialog: false})
    }

    handleAddLicensePlate = async(licensePlate) => {
        licensePlate = licensePlate.trim();
        const pattern = /^[A-Z]{2}-\d{3}-\d{2}$/;
        if (licensePlate == "") {
            await this.setState({snackbarMessage: "Can't add empty license plate"});
            this.handleClick();
            return;
        } else if (!pattern.test(licensePlate)) {
            await this.setState({snackbarMessage: "Please check inputted license plate"});
            this.handleClick();
            return;
        }
        await VehicleGateService.addLicensePlate(this.state.device.ConsumptionDevice.Device.Id, licensePlate);
        let device = this.state.device;
        let licensePlates = device.LicensePlates;
        licensePlates.push(licensePlate);
        device.LicensePlates = licensePlates;
        await this.setState({device: device, showAddLicensePlateDialog: false, snackbarMessage: "Trusted license plate successfully added!"});
        this.handleClick();
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
                        <p className="vg-description">{this.state.exit === false ? this.state.enterLicensePlate : ''} {this.state.enter === true ? ' is entering...' : ''}</p>
                        <p className="vg-description">{this.state.enter === false ? this.state.enterLicensePlate : ''} {this.state.exit === true ? ' is exiting...' : ''}</p>
                        <img src='/images/closed-gate.png' className={`vg-icon ${this.state.device.IsOpen === true ? 'unlocked' : ''}`} />
                        <img src='/images/opened-gate.png' className={`vg-icon ${this.state.device.IsOpen === false ? 'unlocked' : ''}`} />
                        <div id="vg-box">
                            <p className="sp-data-text">Trusted License Plates</p>
                            <List id="vg-list">
                                {this.state.device.LicensePlates.map((licensePlate, index) => (
                                    <ListItem key={index}>
                                        <ListItemText primary={licensePlate} />
                                    </ListItem>
                                ))}
                            </List>
                        </div>
                        <span className='vg-description vg-add'><p onClick={this.openAddLicensePlateDialog}>Add License Plate</p></span>
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
                {this.state.showAddLicensePlateDialog && (
                <Dialog
                    title="Add Trusted License Plate"
                    message="Note that this vehicle will be able to enter property even when the mode is set to private."
                    onConfirm={this.handleAddLicensePlate}
                    onCancel={this.handleCancel}
                    isDiscard={true}
                    inputPlaceholder="Write license plate number here..."
                />
                )}
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