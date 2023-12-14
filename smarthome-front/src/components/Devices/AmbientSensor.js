import { Component } from 'react';
import {Line} from 'react-chartjs-2';
import 'chartjs-adapter-date-fns'
import './Devices.css';
import 'chart.js/auto';
import { Navigation } from '../Navigation/Navigation';
import './AmbientSensor.css'
import authService from '../../services/AuthService'
import AmbientSensorService from '../../services/AmbientSensorService';


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
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());

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
            console.log("rezultat", values)
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
            // const formattedTimestamps = timestamps.map((timestamp) => {
            //     const date = new Date(timestamp);
            //     return date.toLocaleTimeString('en-US', { hour: 'numeric', minute: 'numeric' });
            // });

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
        console.log("uslooo")
        console.log(message)
        // console.log(this.values)
        // console.log(data.labels)

        const newValue = JSON.parse(message);

        // console.log(newValue.humidity)
        // console.log(newValue['temperature'])
        // console.log(timestamp)

        // const refreshData = 

        // const formattedTimestamp = new Date(newValue.timestamp).toLocaleTimeString('en-US', { hour: 'numeric', minute: 'numeric' });

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
                {/* <div className='top-bar'> */}
                    <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer", float: "left" }} onClick={this.handleBackArrow} />
                    <span className='buttons'>
                        <span onClick={() => this.setActiveGraph(1)} className={this.state.activeGraph === 1 ? 'active-button' : 'non-active-button'}>Real Time</span>
                        <span onClick={() => this.setActiveGraph(2)} className={this.state.activeGraph === 2 ? 'active-button' : 'non-active-button'}>History</span>
                    </span>
                {/* </div> */}

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
