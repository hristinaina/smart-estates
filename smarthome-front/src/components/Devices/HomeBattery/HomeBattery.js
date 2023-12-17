
import { Component } from 'react';
import '../Devices.css';
import { Navigation } from '../../Navigation/Navigation';
import mqtt from 'mqtt';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import authService from '../../../services/AuthService';
import DeviceService from '../../../services/DeviceService';
import 'chart.js/auto';
import SPGraph from '../SolarPanel/SPGraph';
import { TextField } from '@mui/material';
import { Button } from 'reactstrap';
import { Snackbar } from "@mui/material";


export class HomeBattery extends Component {
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

        const device = await DeviceService.getHB(this.id);
        console.log(device);

        const user = authService.getCurrentUser();
        this.Name = device.Device.Name;
        const historyData = await DeviceService.getSPGraphData(this.id, user.Email, "2023-12-12", "2023-12-23");
    
        this.setState({
            device: device,
            data: historyData,
            email: user.Email,
            startDate: "2023-12-12",
            endDate: "2023-12-23",
        });
    }

    handleFormSubmit = async (e) => {
        e.preventDefault();

        const { email, startDate, endDate } = this.state;
        console.log(email, startDate, endDate);
        const historyData = await DeviceService.getSPGraphData(this.id, email, startDate, endDate);
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
        //TODO CRTEZ BATERIJE I PRIKAZATI KOLIKO JE TRENUTNO POPUNJENA
        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{this.Name}</span>
                <div className='sp-container'>
                    <div id="sp-left-card">
                        <p className='sp-card-title'>Device Data</p>
                        <p className='sp-data-text'>Maximum capacity (kWh):</p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.Size} InputProps={{
                            readOnly: true,
                        }} />
                    
                        <p className='sp-data-text'>Occupied capacity (kWh): </p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.CurrentValue} InputProps={{
                            readOnly: true,
                        }} />
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
