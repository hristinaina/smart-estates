
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
import PermissionService from '../../../services/PermissionService';
import ProductionGraph from './ProductionGraph';
import DeviceHeader from '../DeviceHeader/DeviceHeader';


export class SolarPanel extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            switchOn: false,
            data: [],
            email: '',  // email of the selected option from user emails
            userEmails: [],
            startDate: '',
            productionData: [],
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

        //const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/sp/');
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
        const historyData = await SolarPanelService.getSPGraphData(this.id, user.Email, "2024-06-13", "2024-06-20");
        const resultP = await SolarPanelService.getProduction(this.id, "2024-06-13", "2024-06-20");
        const graphProduction = this.convertResultToGraphData(resultP)
        let users = await PermissionService.getPermissions(this.id, user.EstateId);
        users.push(user.Email);
        users.push("all");
        users.push("none");
        this.setState({
            device: updatedData,
            switchOn: device.IsOn,
            data: historyData,
            productionData: graphProduction,
            email: user.Email,
            userEmails: users,
            startDate: '2024-06-13',
            endDate: '2024-06-20',
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
        const resultP = await SolarPanelService.getProduction(this.id, startDate, endDate);
        const graphProduction = this.convertResultToGraphData(resultP)
        this.setState({
            data: historyData,
            productionData: graphProduction,
        });
    };

    
    convertResultToGraphData(values) {
        if (values == null) {
            return {
                timestamps: [],
                consumptionData: []
            }
        }
        const timestamps = Object.keys(values);
        const consumptionData = timestamps.map((timestamp) => values[timestamp]);
        const graphData = {
            timestamps: timestamps,
            consumptionData: consumptionData
        }
        return graphData
    }

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
        const { productionData, device, switchOn, data, email, startDate, endDate, userEmails } = this.state;

        return (
            <div>
                <Navigation />
                <DeviceHeader handleBackArrow={this.handleBackArrow} name={this.Name} />
                <div className='sp-container'>
                    <div id="sp-left-card">
                        <p className='sp-card-title'>Device Data</p>
                        <p className='sp-data-text'>Number of panels:</p>
                        <p><b>{device.NumberOfPanels} </b> </p>
                        <p className='sp-data-text'>Surface area per panel (m<sup>2</sup>):</p>
                        <p><b>{device.SurfaceArea}</b></p> 
                        <p className='sp-data-text'>Efficiency per panel (%):</p>
                        <p><b>{device.Efficiency}</b></p>
                        {/* {switchOn ? (<p className='device-text'>Value: {device.Value}</p>) : null} */}
                        <p className='sp-data-text'>Produced electricity in previous minute (kW/m<sup>2</sup>): </p>
                        <p><b>{device.Value}</b></p>
                        <p className='sp-data-text'>Status: </p>
                        <Stack direction="row" className="status-alingment" spacing={1} alignItems="center">
                            <Typography style={{ display: "inline", fontSize: "1.1em" }}>Off</Typography>
                            <Switch
                                checked={switchOn}
                                onChange={this.handleSwitchToggle}
                            />
                            <Typography style={{ display: "inline", fontSize: "1.1em" }}>On</Typography>
                        </Stack>
                    </div>
                    <div id='sp-right-card'>
                        <form onSubmit={this.handleFormSubmit} className='sp-container-input'>
                            <label>
                                User:
                                <select style={{width: "200px", cursor: "pointer"}}
                                    className="new-real-estate-select"
                                    value={email}
                                    onChange={(e) => this.setState({ email: e.target.value })}>
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
                            <Button type="submit" id='sp-data-button'>Fetch Data</Button>
                        </form>

                        <div className='card'>
                            <p className='sp-card-title'>Switch History</p>
                            <SPGraph data={data}/>
                        </div>
                    </div>
                </div>
                <div className='center-graph card'>
                    <p className='sp-card-title'>Electricity produced</p>
                    <ProductionGraph data={productionData} />
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
