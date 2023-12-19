
import { Component } from 'react';
import '../Devices.css';
import { Navigation } from '../../Navigation/Navigation';
import mqtt from 'mqtt';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import authService from '../../../services/AuthService';
import 'chart.js/auto';
import SPGraph from './SPGraph';
import { TextField } from '@mui/material';
import { Button } from 'reactstrap';
import './SolarPanel.css'
import { Snackbar } from "@mui/material";
import SolarPanelService from '../../../services/SolarPanelService';


export class SolarPanel extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            switchOn: false,
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

        const device = await SolarPanelService.getSPById(this.id);
        let lastValue = await SolarPanelService.getSPLastValue(this.id);
        if (lastValue == null) lastValue = 0.0;
        const updatedData =
        {
            ...device,
            Value: lastValue,
        }
        console.log(device);

        const user = authService.getCurrentUser();
        this.Name = device.Device.Name;
        const historyData = await SolarPanelService.getSPGraphData(this.id, user.Email, "2023-12-12", "2023-12-23");
    
        this.setState({
            device: updatedData,
            switchOn: device.IsOn,
            data: historyData,
            email: user.Email,
            startDate: "2023-12-12",
            endDate: "2023-12-23",
        });

        try {
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-sp",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('sp/data/' + this.id);
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
        const topic = "sp/switch/" + this.id;

        this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));
        var message = {
            "IsOn": (!this.state.switchOn),
            "UserEmail": authService.getCurrentUser().Email,
        }
        this.mqttClient.publish(topic, JSON.stringify(message));

        this.setState({ snackbarMessage: "Successfully changed switch state!" });
        this.handleClick();
    };

    // Handle incoming MQTT messages
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
        const { device, switchOn, data, email, startDate, endDate } = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{this.Name}</span>
                <div className='sp-container'>
                    <div id="sp-left-card">
                        <p className='sp-card-title'>Device Data</p>
                        <p className='sp-data-text'>Number of panels:</p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.NumberOfPanels} InputProps={{
                            readOnly: true,
                        }} />
                        <p className='sp-data-text'>Surface area per panel (m<sup>2</sup>):</p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.SurfaceArea} InputProps={{
                            readOnly: true,
                        }} />
                        <p className='sp-data-text'>Efficiency per panel (%):</p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.Efficiency} InputProps={{
                            readOnly: true,
                        }} />
                        {/* {switchOn ? (<p className='device-text'>Value: {device.Value}</p>) : null} */}
                        <p className='sp-data-text'>Status: </p>
                        <Stack direction="row" className="status-alingment" spacing={1} alignItems="center">
                            <Typography style={{ display: "inline", fontSize: "1.1em" }}>Off</Typography>
                            <Switch
                                checked={switchOn}
                                onChange={this.handleSwitchToggle}
                            />
                            <Typography style={{ display: "inline", fontSize: "1.1em" }}>On</Typography>
                        </Stack>
                        <p className='sp-data-text'>Produced electricity in previous minute (kW/m<sup>2</sup>): </p>
                        <p><b>{device.Value}</b></p>
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
