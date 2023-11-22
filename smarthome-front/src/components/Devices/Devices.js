
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'
import { Divider } from '@mui/material';
import { Link } from 'react-router-dom';
import ImageService from '../../services/ImageService';
import { DeviceUnknown } from '@mui/icons-material';

export class Devices extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
            deviceImages: {},
        };
    }

    async componentDidMount() {
        try {
            const result = await DeviceService.getDevices(2);
            await this.setState({ data: result });
            console.log("dataaaaaaaaaa");
            console.log(result)

            const deviceImages = {};
            for (const device of result) {
                console.log("nameeeeeeee", device.Name);
                const imageUrl = await ImageService.getImage("devices&" + device.Name);
                console.log("urlllllll");
                console.log(imageUrl);
                deviceImages[device.Id] = imageUrl;
            }
            await this.setState({deviceImages});
            console.log("devicessss");
            console.log(this.state.deviceImages);
        } catch (error) {
            // Handle error
            console.log("error");
            console.error(error);
        }
    }

    render() {
        const { data, deviceImages } = this.state;

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
                <Divider style={{width: "87%", marginLeft: 'auto', marginRight: 'auto', marginBottom: '20px'}}/>
                <DevicesList devices={data} deviceImages={deviceImages}/>
            </div>
        )
    }
}

const DevicesList = ({ devices, deviceImages }) => {
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
                        <div key={index} className='device-card'>
                            <img
                                alt='device'
                                src={deviceImages[device.Id]}
                                className='device-img'
                            />
                            <div className='device-info'>
                                <p className='device-title'>{device.Name}</p>
                                <p className='device-text'>{device.Type}</p>
                                {device.IsOnline && (<p className='device-text' style={{color: 'green'}}>Online</p>)}
                                {!device.IsOnline && (<p className='device-text' style={{color: 'red'}}>Offline</p>)}
                            </div>
                        </div>
                    ))}
                </div>
            ))}
        </div>
    );
};
