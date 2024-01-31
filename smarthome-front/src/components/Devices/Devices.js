
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'
import { Divider } from '@mui/material';
import { Link } from 'react-router-dom';
import RealEstateService from '../../services/RealEstateService'
import mqtt from 'mqtt';
import authService from '../../services/AuthService'

import ImageService from '../../services/ImageService';
import { DeviceUnknown } from '@mui/icons-material';
import PermissionService from '../../services/PermissionService';

export class Devices extends Component {
    connected = false;
    navigationToNewDevice = false;

    constructor(props) {
        super(props);
        this.state = {
            data: [],
            deviceImages: {},
            name: '',
            owner: false,
        };
        this.mqttClient = null;
        this.connecting = false; //change to true if you want to use this
        this.id = parseInt(localStorage.getItem('real-estate'));
    }

    async componentDidMount() {
        
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        try {
            const result = await DeviceService.getDevices(this.id);
            console.log(result)
            this.setState({ data: result });

            // todo ako nije pozvati api koja vraca sve uredjaje koje moze da vidi ulogovani korisnik za izabranu nekretninu
            const currentUser = authService.getCurrentUser()
            console.log(localStorage.getItem('owner'))
            if (currentUser.Id != localStorage.getItem('owner')) {
                // todo ne moze da dodaje uredjaje
                console.log("nije vlasnik!!")
                console.log(this.id)
                console.log(currentUser.Id)
                const devices = await PermissionService.getDevices(this.id, currentUser.Id)
                console.log(devices)
                this.setState({ data: devices });
            }
            else {
                this.setState({ owner: true });
            }

            const deviceImages = {};
            for (const device of result) {
                const imageUrl = await ImageService.getImage("devices&" + device.Name);
                deviceImages[device.Id] = imageUrl;
            }
            await this.setState({ deviceImages });
        } catch (error) {
            console.log("Error fetching data from the server");
            console.log(error);
            window.location.assign("/");
        }

        const result = await RealEstateService.getById(this.id);
        this.setState({ name: result.Name });

        try {
            if (!this.connected) {  // to avoid reconnecting because this renders 2 times !!!
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-devices",
                    clean: false,
                    keepalive: 60
                });
                console.log("Connected to mqtt broker");
                // Subscribe to the MQTT topic for device status
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('device/status/+');
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
        if (!this.navigationToNewDevice && this.connected) {
            localStorage.removeItem('real-estate');
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        console.log("handle message");
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
            window.location.assign("/ambient-sensor/" + device.Id)
        else if (device.Type === 'Air conditioner')
            window.location.assign("/air-conditioner/" + device.Id)
        else if (device.Type === 'Washing machine')
            window.location.assign("/washing-machine/" + device.Id)
        else if (device.Type === 'Lamp')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Vehicle gate')
            window.location.assign("/vehicle-gate/" + device.Id)
        else if (device.Type === 'Sprinkler')
            window.location.assign("/lamp/" + device.Id)
        else if (device.Type === 'Solar panel')
            window.location.assign("/sp/" + device.Id)
        else if (device.Type === 'Battery storage')
            window.location.assign("/hb/" + device.Id)
        else if (device.Type === 'Electric vehicle charger')
            window.location.assign("/lamp/" + device.Id)
    }

    extractDeviceIdFromTopic(topic) {
        const parts = topic.split('/');
        return parts[parts.length - 1];
    }

    handleNavigationToNewDevice = () => {
        this.navigationToNewDevice = true;
    };

    render() {
        const { data, deviceImages, name } = this.state;
        const connecting = this.connecting;
        return (
            <div>
                <Navigation />
                <div id="tools">
                    <Link to="/real-estates"><img src='/images/arrow.png' id='arrow' /></Link>
                    <span className='estate-title'>{name}</span>
                    {this.state.owner && (
                        <p id="add-device">
                            <Link to="/new-device" onClick={this.handleNavigationToNewDevice}>
                                <img alt="." src="/images/plus.png" id="plus" />
                                Add Device
                            </Link>
                        </p>
                    )}
                </div>
                <Divider style={{ width: "87%", marginLeft: 'auto', marginRight: 'auto', marginBottom: '20px' }} />
                <DevicesList devices={data} deviceImages={deviceImages} onClick={this.handleClick} connecting={connecting} />
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
