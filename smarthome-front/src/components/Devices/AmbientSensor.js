import { Component } from 'react';
import {Line} from 'react-chartjs-2';
import './Devices.css';
import 'chart.js/auto';
import { Navigation } from '../Navigation/Navigation';
import mqtt from 'mqtt';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import authService from '../../services/AuthService'
import AmbientSensorService from '../../services/AmbientSensorService';


export class AmbientSensor extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            switchOn: false,
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
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());

        this.options = {
            scales: {
                y: {
                    beginAtZero: true,
                },
            },
        };
        

        this.values = null

    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const { device } = this.state;  // todo instead of this get device from back by deviceId
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
            this.values = new Map(Object.entries(values))
            // console.log("rezultat", values)
            // console.log(typeof(values))

            const timestamps = Object.keys(values);
            const humidityData = timestamps.map((timestamp) => values[timestamp].humidity);
            const temperatureData = timestamps.map((timestamp) => values[timestamp].temperature);

            // console.log("vreme: ", timestamps)
            // console.log("humidity: ", humidityData)
            // console.log("temperature: ", temperatureData)
        
            // if (this.chartInstance) {
            //     this.chartInstance.destroy();
            // }
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

            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-AmbientSensor",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic for device status
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('device/data/' + this.id);
                });

                // Handle incoming MQTT valuess
                this.mqttClient.on('values', (topic, values) => {
                    this.handleMqttvalues(topic, values);
                });
            }
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

    populateGraph = (message) => {
        const { data } = this.state;
        // console.log("uslooo")
        // console.log(message)
        // console.log(this.values)

        const newValue = JSON.parse(message);

        // console.log(newValue.humidity)
        // console.log(newValue['temperature'])
        // console.log(timestamp)

        const updatedChartData = {
            // labels: [...data.labels, newValue.timestamp], 
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

    handleSwitchToggle = () => {
        const topic = "lamp/switch/" + this.id;

        this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));
        const values = (!this.state.switchOn).toString();
        this.mqttClient.publish(topic, values);
    };

    // Handle incoming MQTT valuess
    handleMqttvalues(topic, values) {
        const { device } = this.state;
        const newValue = values.toString();
        const updatedData =
        {
            ...device,
            Value: newValue + "%",
        }
        this.setState({
            device: updatedData,
        });
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    render() {
        const { device, switchOn } = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <div style={{ width: "fit-content", marginLeft: "auto", marginRight: "auto", marginTop: "10%" }}>
                {/* /    <p className='device-title'>Id: {this.id}</p> */}
                    {/* {switchOn ? (<p className='device-text'>Value: {device.Value}</p>) : null} */}
                    {/* <p className='device-text'>Last Value: {device.Value}</p> */}
                    {/* <Stack direction="row" spacing={1} alignItems="center"> */}
                        {/* <Typography>Off</Typography> */}
                        {/* <Switch */}
                            {/* checked={switchOn} */}
                            {/* onChange={this.handleSwitchToggle} */}
                        {/* /> */}
                        {/* <Typography>On</Typography> */}
                    {/* </Stack> */}
                </div>

                <Line ref={(ref) => (this.chartInstance = ref)} id='graph' data={this.state.data} options={this.options} />
            </div>
        )
    }
}
