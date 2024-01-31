import React, { Component } from 'react';
import { Table, TableBody, TableCell, TableHead, TableRow } from '@mui/material';

class SensorDataTable extends Component {
    constructor(props) {
        super(props);
        this.state = {
            minTemp: 0.0,
            maxTemp: 0.0,
            avgTemp: 0.0,
            minHmd: 0.0,
            maxHmd: 0.0,
            avgHmd: 0.0,
        };
    }

    componentDidmounnt() {
        this.calculateAverage()
    }

    componentDidUpdate(prevProps) {
        if (prevProps.data !== this.props.data) {
            this.calculateAverage()
        }
    }

    calculateAverage() {
        const { data } = this.props;

        const temperatures = Object.values(data).map(entry => entry.temperature);
        const humidities = Object.values(data).map(entry => entry.humidity);
        
        const minTemperature = Math.min(...temperatures);
        const maxTemperature = Math.max(...temperatures);
        const avgTemperature = temperatures.reduce((acc, val) => acc + val, 0) / temperatures.length;
        
        const minHumidity = Math.min(...humidities);
        const maxHumidity = Math.max(...humidities);
        const avgHumidity = humidities.reduce((acc, val) => acc + val, 0) / humidities.length;

        this.setState({minTemp: minTemperature, maxTemp: maxTemperature, avgTemp: avgTemperature})
        this.setState({minHmd: minHumidity, maxHmd: maxHumidity, avgHmd: avgHumidity})
    }


    render() {
        const {minTemp, maxTemp, avgTemp, minHmd, maxHmd, avgHmd} = this.state

        return (
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell style={{fontWeight: "900"}}>Statistics</TableCell>
                        <TableCell>Temperature (°C)</TableCell>
                        <TableCell>Humidity (%)</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    <TableRow>
                        <TableCell>Maximum value</TableCell>
                        <TableCell style={{color: "#8B0000", fontWeight: "bold"}}>{maxTemp.toFixed(2)}</TableCell>
                        <TableCell style={{color: "#8B0000", fontWeight: "bold"}}>{maxHmd.toFixed(2)}</TableCell>
                    </TableRow>
                    <TableRow>
                        <TableCell>Minimum value</TableCell>
                        <TableCell style={{color: "#00008B", fontWeight: "bold"}}>{minTemp.toFixed(2)}</TableCell>
                        <TableCell style={{color: "#00008B", fontWeight: "bold"}}>{minHmd.toFixed(2)}</TableCell>
                    </TableRow>
                    <TableRow>
                        <TableCell>Average value</TableCell>
                        <TableCell style={{color: "#8B8000", fontWeight: "bold"}}>{avgTemp.toFixed(2)}</TableCell>
                        <TableCell style={{color: "#8B8000", fontWeight: "bold"}}>{avgHmd.toFixed(2)}</TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        );
    };
}    

export default SensorDataTable;
