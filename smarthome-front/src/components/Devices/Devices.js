
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'
import { Divider } from '@mui/material';
import { Link } from 'react-router-dom';
import mqtt from 'mqtt';


export class Devices extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
        };
        this.mqttClient = null;
        this.connecting = false; //change to true if you want to use this
        this.id = parseInt(this.extractEstateFromUrl());
    }

    async componentDidMount() {
        try {
            const result = await DeviceService.getDevices(this.id);
            this.setState({ data: result });
        } catch (error) {
            console.log("Error fetching data from the server");
            console.log(error);
        }

        try {
            this.mqttClient = mqtt.connect('ws://broker.emqx.io:8083/mqtt');

            // Subscribe to the MQTT topic for device status
            this.mqttClient.on('connect', () => {
                this.mqttClient.subscribe('device/status/+');
            });

            // Handle incoming MQTT messages
            this.mqttClient.on('message', (topic, message) => {
                this.handleMqttMessage(topic, message);
            });
        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }

        // setTimeout(() => {
        //     const { data } = this.state;
        //     this.connecting = false;
        //     this.setState({
        //         data: data,
        //     });
        // }, 5000);
        }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        const { data } = this.state;
        const deviceId = this.extractDeviceIdFromTopic(topic);
        const status = message.toString();
        console.log(deviceId, status);
        // Update the IsOnline status based on the received MQTT message
        const updatedData = data.map((device) =>
            device.Id == deviceId
                ? {
                    ...device,
                    IsOnline: status == 'online',
                }
                : device
        );
        console.log(updatedData)

        this.setState({
            data: updatedData,
        });
    }

    //todo navigate to appropriate page
    handleClick(device) {
        if (device.Type === 'Ambient Sensor')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Air conditioner')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Washing machine')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Lamp')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Vehicle gate')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Sprinkler')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Solar panel')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Battery storage')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Electric vehicle charger')
            window.location.assign("/lamp/" + device.Id)
    }

    extractDeviceIdFromTopic(topic) {
        const parts = topic.split('/');
        return parts[parts.length - 1];
    }

    extractEstateFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    render() {
        const { data } = this.state;
        const connecting = this.connecting;
        const newDeviceNav = "/new-device/" + this.id;

        return (
            <div>
                <Navigation />
                <div id="tools">
                    <Link to="/real-estates"><img src='/images/arrow.png' id='arrow' /></Link>
                    <span className='estate-title'>Ta i ta nekretnina</span>
                    <p id="add-device">
                        <Link to={newDeviceNav}>
                            <img alt="." src="/images/plus.png" id="plus" />
                            Add Device
                        </Link>
                    </p>
                </div>
                <Divider style={{ width: "87%", marginLeft: 'auto', marginRight: 'auto', marginBottom: '20px' }} />
                <DevicesList devices={data} onClick={this.handleClick} connecting={connecting}/>
            </div>
        )
    }
}

const DevicesList = ({ devices, onClick, connecting }) => {
    const chunkSize = 5; // Number of items per row

    const chunkArray = (arr, size) => {
        return Array.from({ length: Math.ceil(arr.length / size) }, (v, i) =>
            arr.slice(i * size, i * size + size)
        );
    };

    const rows = chunkArray(devices, chunkSize);

    return (
        <div id='devices-container'>
            {rows.map((row, rowIndex) => (
                <div key={rowIndex} className='device-row'>
                    {row.map((device, index) => (
                        <div key={index} className='device-card' onClick={() => onClick(device)}>
                            <img
                                alt='device'
                                src={device.Picture}
                                className='device-img'
                            />
                            <div className='device-info'>
                                <p className='device-title'>{device.Name}</p>
                                <p className='device-text'>{device.Type}</p>
                                {device.IsOnline && (<p className='device-text' style={{ color: 'green' }}>Online</p>)}
                                {!device.IsOnline && (<p className='device-text' style={{ color: 'red' }}>Offline</p>)}
                                {/* {!device.IsOnline && !connecting && (<p className='device-text' style={{ color: 'red' }}>Offline</p>)}
                                {connecting && !device.IsOnline && (<p className='device-text'>Connecting</p>)} */}
                            </div>
                        </div>
                    ))}
                </div>
            ))}
        </div>
    );
};
