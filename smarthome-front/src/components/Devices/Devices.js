
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'
import { Divider } from '@mui/material';
import { Link } from 'react-router-dom';
import mqtt from 'mqtt';

import ImageService from '../../services/ImageService';
import { DeviceUnknown } from '@mui/icons-material';

export class Devices extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
            deviceImages: {},
        };
        this.mqttClient = null;
        this.connecting = false; //change to true if you want to use this
        this.id = parseInt(localStorage.getItem('real-estate'));
    }

    async componentDidMount() {
        try {
            const result = await DeviceService.getDevices(this.id);
            this.setState({ data: result });

            const deviceImages = {};
            for (const device of result) {
                const imageUrl = await ImageService.getImage("devices&" + device.Name);
                deviceImages[device.Id] = imageUrl;
            }
            await this.setState({deviceImages});
        } catch (error) {
            console.log("Error fetching data from the server");
            console.log(error);
        }

        try {
            this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                clientId: "react-front-nvt-2023-devices",
            });

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
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        this.setState((prevState) => {
            const { data } = prevState;
            const deviceId = parseInt(this.extractDeviceIdFromTopic(topic));
            const status = message.toString();

            // Update the IsOnline status based on the received MQTT message
            const updatedData = data.map((device) =>
                device.Id == deviceId
                    ? {
                        ...device,
                        IsOnline: status === 'online',
                    }
                    : device
            );

            return {
                data: updatedData,
            };
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

    render() {
        const { data, deviceImages } = this.state;
        const connecting = this.connecting;
        return (
            <div>
                <Navigation />
                <div id="tools">
                    <Link to="/real-estates"><img src='/images/arrow.png' id='arrow' /></Link>
                    <span className='estate-title'>Ta i ta nekretnina</span>
                    <p id="add-device">
                        <Link to="/new-device">
                            <img alt="." src="/images/plus.png" id="plus" />
                            Add Device
                        </Link>
                    </p>
                </div>
                <Divider style={{ width: "87%", marginLeft: 'auto', marginRight: 'auto', marginBottom: '20px' }} />
                <DevicesList devices={data} deviceImages={deviceImages} onClick={this.handleClick} connecting={connecting}/>
            </div>
        )
    }
}

const DevicesList = ({ devices, deviceImages, onClick, connecting }) => {
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
                                src={deviceImages[device.Id]}
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
