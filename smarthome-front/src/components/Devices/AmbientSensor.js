import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import mqtt from 'mqtt';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import authService from '../../services/AuthService'


export class AmbientSensor extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            switchOn: false,
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
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

                // Handle incoming MQTT messages
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessage(topic, message);
                });
            }
        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }

        let socket = new WebSocket("ws://localhost:8082/ws")
        console.log("Attempting Websocket Connection")

        socket.onopen = () => {
            console.log("Successfully Connected")
        }

        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event)
        }

        socket.onmessage = (msg) => {
            console.log(msg)
        }
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    handleSwitchToggle = () => {
        const topic = "lamp/switch/" + this.id;

        this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));
        const message = (!this.state.switchOn).toString();
        this.mqttClient.publish(topic, message);
    };

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        const { device } = this.state;
        const newValue = message.toString();
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
                    <p className='device-title'>Id: {this.id}</p>
                    {/* {switchOn ? (<p className='device-text'>Value: {device.Value}</p>) : null} */}
                    <p className='device-text'>Last Value: {device.Value}</p>
                    <Stack direction="row" spacing={1} alignItems="center">
                        {/* <Typography>Off</Typography> */}
                        <Switch
                            checked={switchOn}
                            onChange={this.handleSwitchToggle}
                        />
                        <Typography>On</Typography>
                    </Stack>
                </div>
            </div>
        )
    }
}
