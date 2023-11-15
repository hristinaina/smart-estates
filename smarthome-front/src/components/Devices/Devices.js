
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'
import { Divider } from '@mui/material';
import { Link } from 'react-router-dom';
import { color } from '@mui/system';

export class Devices extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
        };
    }

    async componentDidMount() {
        try {
            const result = await DeviceService.getDevices(1);
            this.setState({ data: result });
        } catch (error) {
            // Handle error
        }
    }

    render() {
        const { data } = this.state;

        return (
            <div>
                <Navigation />
                <div id="tools">
                    <Link to="/"><img src='/images/arrow.png' id='arrow' /></Link>
                    <span className='estate-title'>Ta i ta nekretnina</span>
                    <p id="add-device">
                        <img alt="." src="/images/plus.png" id="plus" />
                        Add Device
                    </p>
                </div>
                <Divider style={{width: "87%", marginLeft: 'auto', marginRight: 'auto', marginBottom: '20px'}}/>
                <DevicesList devices={data} />
            </div>
        )
    }
}

const DevicesList = ({ devices }) => {
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
                                src={device.image}
                                className='device-img'
                            />
                            <div className='device-info'>
                                <p className='device-title'>{device.name}</p>
                                <p className='device-text'>{device.type}</p>
                                {device.status=="Online" && (<p className='device-text' style={{color: 'green'}}>{device.status}</p>)}
                                {device.status=="Offline" && (<p className='device-text' style={{color: 'red'}}>{device.status}</p>)}
                            </div>
                        </div>
                    ))}
                </div>
            ))}
        </div>
    );
};
