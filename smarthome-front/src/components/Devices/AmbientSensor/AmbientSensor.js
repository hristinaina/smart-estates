import React, { Component } from 'react';
import {Line} from 'react-chartjs-2';
import 'chartjs-adapter-date-fns'
import '../Devices.css';
import 'chart.js/auto';
import { Navigation } from '../../Navigation/Navigation';
import './AmbientSensor.css'
import authService from '../../../services/AuthService'
import AmbientSensorService from '../../../services/AmbientSensorService';
import { Autocomplete, TextField, Button, Box, Grid, IconButton, Snackbar } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import DeviceService from '../../../services/DeviceService';


export class AmbientSensor extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            switchOn: false,
            activeGraph: 1,
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'Humidity',
                        data: [],
                        borderColor: 'rgba(128,104,148,1)',
                        borderWidth: 2,
                        fill: false,
                    },
                    {
                        label: 'Temperature',
                        data: [],
                        borderColor: 'rgba(255, 99, 132, 1)', 
                        borderWidth: 2,
                        fill: false,
                    }, 
                ],
            },
            latestData: null,
            open: false,
            snackbarMessage: '',
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
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";

        this.options = {
            scales: {
                x: {
                    type: 'time',
                    time: {
                        displayFormats: {
                            quarter: 'HH:MM'
                        }
                    }
                },
                y: {
                    beginAtZero: true,
                },
            },
        };
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/ambient/');
        console.log(device)
        this.Name = device.Device.Name;
        // const { device } = this.state;  // todo instead of this get device from back by deviceId
        const updatedData =
        {
            ...device,
            Value: "Loading...",
        }
        this.setState({
            device: updatedData,
        });

        try {
            const result = await AmbientSensorService.getGraphData(this.id);
            const values = result.result
            console.log("rezultat", values)
            // console.log(typeof(values))

            const timestamps = Object.keys(values);
            const humidityData = timestamps.map((timestamp) => values[timestamp].humidity);
            const temperatureData = timestamps.map((timestamp) => values[timestamp].temperature);

            await this.setState({
                data: {
                    labels: timestamps,
                    datasets: [
                        {
                            label: 'Humidity',
                            data: humidityData,
                            borderColor: 'rgba(128,104,148,1)',
                            borderWidth: 2,
                            fill: false,
                        },
                        {
                            label: 'Temperature',
                            data: temperatureData,
                            borderColor: 'rgba(255, 99, 132, 1)', 
                            borderWidth: 2,
                            fill: false,
                        },
                    ],
                },
            });

        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }

        let socket = new WebSocket("ws://localhost:8082/ambient")
        console.log("Attempting Websocket Connection")

        socket.onopen = () => {
            console.log("Successfully Connected")
            socket.send(this.id)
        }

        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event)
        }

        socket.onmessage = (msg) => {
            if(this.state.activeGraph === 1) 
                this.populateGraph(msg.data)
        }
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    isTimestampInLastHour = (timestamp) => {
        const currentTimestamp = new Date();
        const timestampDate = new Date(timestamp);
    
        const timeDifference = currentTimestamp - timestampDate;

        return timeDifference <= 3600000;
    };

    historyGraph = async (values) => {
        const timestamps = Object.keys(values);
        const humidityData = timestamps.map((timestamp) => values[timestamp].humidity);
        const temperatureData = timestamps.map((timestamp) => values[timestamp].temperature);

        await this.setState({
            data: {
                labels: timestamps,
                datasets: [
                    {
                        label: 'Humidity',
                        data: humidityData,
                        borderColor: 'rgba(128,104,148,1)',
                        borderWidth: 2,
                        fill: false,
                    },
                    {
                        label: 'Temperature',
                        data: temperatureData,
                        borderColor: 'rgba(255, 99, 132, 1)', 
                        borderWidth: 2,
                        fill: false,
                    },
                ],
            },
        });
    }

    populateGraph = (message) => {
        const { data } = this.state;

        const newValue = JSON.parse(message);

        const updatedChartData = {
            labels: data.labels.filter((label) => this.isTimestampInLastHour(label)).concat(newValue.timestamp),
            datasets: [
                {
                    label: 'Humidity',
                    data: [...data.datasets[0].data, newValue.humidity],
                    borderColor: 'rgba(128,104,148,1)',
                    borderWidth: 2,
                    fill: false,
                },
                {
                    label: 'Temperature',
                    data: [...data.datasets[1].data, newValue.temperature], 
                    borderColor: 'rgba(255, 99, 132, 1)',
                    borderWidth: 2,
                    fill: false,
                },
            ],
        };

        this.setState({
            data: updatedChartData,
        });
    };

    updateGraph = async (value) => {
        const result = await AmbientSensorService.getDataForSelectedTime(this.id, value);
        await this.historyGraph(result.result.result)
    }

    handleOptionChange = async (event, value) => {
        this.setState({ selectedOption: value });
        await this.updateGraph(value.value)
    };
    
    handleDateChange = (fieldName, event) => {
        this.setState({ [fieldName]: event.target.value });
    };
    
    handleButtonClick = async () => {
        if(this.state.startDate === '' || this.state.endDate === '') {
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
            this.setState({ snackbarMessage: 'The difference between start date and end date must not be more than one month'});
            this.handleClick();
            return;
        }
        const result = await AmbientSensorService.getDataForSelectedDate(this.id, this.state.startDate, this.state.endDate);
        console.log("datum graf ", result.result.result)
        result.result.result != null ? await this.historyGraph(result.result.result) : await this.historyGraph([]) 
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
        const { selectedOption, startDate, endDate, options } = this.state;

        const action = (
            <React.Fragment>
                <IconButton
                size="small"
                aria-label="close"
                color="inherit"
                onClick={this.handleClose}>
                <CloseIcon fontSize="small" />
                </IconButton>
            </React.Fragment>);

        return (
            <div>
                <Navigation />
                <span style={{float: "left", marginRight: "250px"}}>
                <img src='/images/arrow.png' id='arrow' alt='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer"}} onClick={this.handleBackArrow} />
                <span className='ambient-sensor-title'>{this.Name}</span>
                </span>
                
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
        <Snackbar
        open={this.state.open}
        autoHideDuration={3000}
        onClose={this.handleClose}
        message={this.state.snackbarMessage}
        action={action}/>

                <div className='canvas'>
                    {this.state.activeGraph === 1 && <Line ref={(ref) => (this.chartInstance = ref)} id='graph' data={this.state.data} options={this.options} />}
                    {this.state.activeGraph === 2 && <Line ref={(ref) => (this.chartInstance = ref)} id='graph' data={this.state.data} options={this.options} />}
                </div>
            </div>
        )
    }

    setActiveGraph = (graphNumber) => {
        this.setState({ activeGraph: graphNumber });
    }
}
