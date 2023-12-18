import { Component } from 'react';
import '../Devices.css';
import { Navigation } from '../../Navigation/Navigation';
import authService from '../../../services/AuthService';
import DeviceService from '../../../services/DeviceService';
import 'chart.js/auto';
import SPGraph from '../SolarPanel/SPGraph';
import HBGraph from './HBGraph';
import AmbientSensorService from '../../../services/AmbientSensorService';
import { Autocomplete, TextField, Button, Box, Grid, IconButton, Snackbar } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';


export class HomeBattery extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            data: [],
            activeGraph: 1,
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
            selectedOption: { label: '6h', value: '-6h' },
            startDate: '',
            endDate: '',
            options: [
                { label: '6h', value: '-6h' },
                { label: '12h', value: '-12h' },
                { label: '24h', value: '-24h' },
                { label: 'last week', value: '-7d' },
                { label: 'last month', value: '-30d' },
            ],
        };
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getHB(this.id);
        console.log(device);
        const updatedData =
        {
            ...device,
            CurrentValue: device.CurrentValue.toFixed(3),
        }

        const user = authService.getCurrentUser();
        this.Name = device.Device.Name;

        const historyData = await AmbientSensorService.getGraphData(this.id);
        const graphData = this.convertResultToGraphData(historyData.result);
        this.setState({
            device: updatedData,
            data: graphData,
        });

        // Set up interval to fetch device data every minute
        this.apiRequestInterval = setInterval(() => {
            this.fetchDeviceData();
        }, 60000); // 60000 milliseconds = 1 minute

        let socket = new WebSocket("ws://localhost:8082/consumption")
        console.log("Attempting Websocket Connection")

        socket.onopen = () => {
            console.log("Successfully Connected")
            socket.send(this.id)
        }

        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event)
        }

        socket.onmessage = (msg) => {
            if (this.state.activeGraph === 1)
                this.populateGraph(msg.data)
        }
    }

    async fetchDeviceData() {
        const device = await DeviceService.getHB(this.id);
        const updatedData = {
            ...device,
            CurrentValue: device.CurrentValue.toFixed(3),
        };
        this.setState({
            device: updatedData,
        });
    }

    componentWillUnmount() {
        clearInterval(this.apiRequestInterval);
    }


    isTimestampInLastHour = (timestamp) => {
        const currentTimestamp = new Date();
        const timestampDate = new Date(timestamp);

        const timeDifference = currentTimestamp - timestampDate;

        return timeDifference <= 3600000;
    };

    populateGraph = (message) => {
        const { data } = this.state;

        const newValue = JSON.parse(message);
        if (newValue.estateId != this.state.device.Device.RealEstate) return;
        console.log("HEEEEEJ");
        console.log(newValue);
        console.log(data);
        const timestamps = data.timestamps.filter((label) => this.isTimestampInLastHour(label)).concat(newValue.timestamp);
        const consumptionData = data.consumptionData.concat(newValue.consumed);
        const updatedChartData = {
            timestamps: timestamps,
            consumptionData: consumptionData
        }
        this.setState({
            data: updatedChartData,
        });
    };

    updateGraph = async (value) => {
        const result = await AmbientSensorService.getDataForSelectedTime(this.id, value);
        const graphData = this.convertResultToGraphData(result.result.result)
        this.setState({
            data: graphData,
        });
    }

    setActiveGraph = (graphNumber) => {
        this.setState({ activeGraph: graphNumber });
    }

    handleOptionChange = async (event, value) => {
        this.setState({ selectedOption: value });
        await this.updateGraph(value.value)
    };

    handleDateChange = (fieldName, event) => {
        this.setState({ [fieldName]: event.target.value });
    };

    handleButtonClick = async () => {
        if (this.state.startDate === '' || this.state.endDate === '') {
            this.setState({ snackbarMessage: "Please enter dates" });
            this.handleClick();
            return;
        }

        if (new Date(this.state.startDate) > new Date(this.state.endDate)) {
            this.setState({ snackbarMessage: "Start date cannot be greater than end date" });
            this.handleClick();
            return;
        }
        const oneMonth = 30 * 24 * 60 * 60 * 1000;
        const difference = new Date(this.state.endDate) - new Date(this.state.startDate);

        if (difference > oneMonth) {
            this.setState({ snackbarMessage: 'The difference between start date and end date must not be more than one month' });
            this.handleClick();
            return;
        }
        const result = await AmbientSensorService.getDataForSelectedDate(this.id, this.state.startDate, this.state.endDate);
        console.log("datum graf ", result.result.result)
        const graphData = this.convertResultToGraphData(result.result.result)
        result.result.result != null ? await HBGraph(graphData) : await HBGraph([])
    };

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    convertResultToGraphData(values) {
        const timestamps = Object.keys(values);
        const consumptionData = timestamps.map((timestamp) => values[timestamp].consumed);
        const graphData = {
            timestamps: timestamps,
            consumptionData: consumptionData
        }
        return graphData
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
        const { device, data, startDate, endDate, selectedOption, options } = this.state;

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
                        <p className='sp-data-text'>Occupied size (kWh): </p>
                        <TextField style={{ backgroundColor: "white", width: "300px" }} type="number" value={device.CurrentValue} InputProps={{
                            readOnly: true,
                        }} />
                        {this.renderBatteryIcon(device.CurrentValue, device.Size)}
                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Estate consumption</p>
                        <span className='buttons'>
                            <span onClick={() => this.setActiveGraph(1)} className={this.state.activeGraph === 1 ? 'active-button' : 'non-active-button'}>Real Time</span>
                            <span onClick={() => { this.setActiveGraph(2); this.updateGraph(this.state.selectedOption.value) }} className={this.state.activeGraph === 2 ? 'active-button' : 'non-active-button'}>History</span>
                        </span>
                        {this.state.activeGraph === 2 &&
                            <div>
                                <Grid container spacing={2}>
                                    <Grid item xs={2}></Grid>
                                    <Grid item xs={3}>
                                        <Autocomplete
                                            value={selectedOption}
                                            onChange={this.handleOptionChange}
                                            options={options}
                                            getOptionLabel={(option) => option.label}
                                            style={{ width: '100%' }}
                                            renderInput={(params) => (
                                                <TextField
                                                    {...params}
                                                    label="Select Time Range"
                                                    InputLabelProps={{
                                                        shrink: true,
                                                    }}
                                                />
                                            )}
                                            isOptionEqualToValue={(option, value) => option.value === value.value}
                                            renderOption={(props, option, { selected }) => (
                                                <li {...props}>
                                                    <span>{option.label}</span>
                                                </li>
                                            )}
                                            disableClearable />
                                    </Grid>
                                    <Grid item xs={6}>
                                        <Box display="flex" alignItems="center" justifyContent="flex-end">
                                            <TextField
                                                label="Start Date"
                                                type="date"
                                                value={startDate}
                                                onChange={(e) => this.handleDateChange('startDate', e)}
                                                InputLabelProps={{
                                                    shrink: true,
                                                }}
                                                inputProps={{
                                                    max: new Date().toISOString().split('T')[0],
                                                }}
                                            />
                                            <TextField
                                                label="End Date"
                                                type="date"
                                                value={endDate}
                                                onChange={(e) => this.handleDateChange('endDate', e)}
                                                InputLabelProps={{
                                                    shrink: true,
                                                }}
                                                inputProps={{
                                                    max: new Date().toISOString().split('T')[0],
                                                }}
                                            />
                                            <Button variant="contained" color="primary" onClick={this.handleButtonClick}>
                                                Apply
                                            </Button>
                                        </Box>
                                    </Grid>
                                </Grid>

                            </div>}
                            <div className='canvas'>
                    {this.state.activeGraph === 1 && <HBGraph data={data} />}
                    {this.state.activeGraph === 2 && <HBGraph data={data}/>}
                </div>
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


    renderBatteryIcon(occupiedCapacity, maxCapacity) {
        // Calculate the percentage of occupied capacity
        const percentageOccupied = ((occupiedCapacity / maxCapacity) * 100).toFixed(0);

        // Set a minimum width for the filled battery part
        const minWidth = 1;

        // Calculate the width of the filled battery part based on the percentage
        const filledWidth = Math.max(minWidth, percentageOccupied);

        // Styles for the battery icon and the filled part
        const batteryStyle = {
            width: '100px', // Adjust the size of the battery icon
            height: '50px',
            background: '#ddd',
            position: 'relative',
            borderRadius: '5px',
            display: 'inline-block'
        };

        const filledStyle = {
            height: '100%',
            width: `${filledWidth}%`,
            background: 'green', // Adjust the color of the filled part
            position: 'absolute',
            borderRadius: '5px',
        };

        return (
            <div style={{ marginTop: "50px", display: "flex", alignItems: "center", justifyContent: "center" }}>
                <div style={{ fontSize: "20px", display: "inline", marginRight: "10px" }}>{percentageOccupied}%</div>
                <div style={batteryStyle}>
                    <div style={filledStyle}></div>
                </div>
            </div>
        );
    }
}

