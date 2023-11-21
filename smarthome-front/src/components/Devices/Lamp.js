
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import mqtt from 'mqtt';
import { useParams } from 'react-router-dom';
import { width } from '@mui/system';
import { Link } from 'react-router-dom';


export class Lamp extends Component {
    constructor(props) {
        super(props);
        this.state = {
            device: {},
        };
        this.mqttClient = null;
        this.id = this.extractDeviceIdFromUrl();
    }

    async componentDidMount() {
        const { device } = this.state;  // todo instead of this get device data from back having device id
        const updatedData = 
        {
            ...device,
            Value: 0,
        }
        this.setState({
            device: updatedData,
        });

        try {
            this.mqttClient = mqtt.connect('ws://broker.emqx.io:8083/mqtt');

            // Subscribe to the MQTT topic for device status
            this.mqttClient.on('connect', () => {
                this.mqttClient.subscribe('device/data/' + this.id);
            });

            // Handle incoming MQTT messages
            this.mqttClient.on('message', (topic, message) => {
                this.handleMqttMessage(topic, message);
            });
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

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    render() {
        const { device } = this.state;

        return (
            <div>
                <Navigation />
                <Link to="/devices"><img src='/images/arrow.png' id='arrow' style={{margin: "55px 0 0 90px"}} /></Link>
                <div style={{width: "fit-content", marginLeft: "auto", marginRight: "auto", marginTop: "10%"}}>
                    <p className='device-title'>Id: {this.id}</p>
                    <p className='device-text'>Value: {device.Value}%</p>
                </div>
            </div>
        )
    }
}
