import React, { Component } from 'react';
import { Pie } from 'react-chartjs-2';

class PieChart extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: {
                labels: [],
                datasets: [{
                    data: [],
                    backgroundColor: [],
                }],
            },
        };
    }

    componentDidMount() {
        this.calculatePercentages();
    }

    componentDidUpdate(prevProps) {
        if (prevProps.data !== this.props.data) {
            if(this.props.graph === 1)
                this.calculatePercentages();
            else if (this.props.graph === 2)
                this.calculatePercentagesOffOn();
            else 
                this.calculateUserActivity();
        }
    }

    calculatePercentages = () => {
        const { data } = this.props;
        const modeCounts = {};

        Object.values(data).forEach(entry => {
            const mode = entry.Mode;
            modeCounts[mode] = (modeCounts[mode] || 0) + 1;
        });

        const totalEntries = Object.values(data).length;

        const percentages = {};
        Object.entries(modeCounts).forEach(([mode, count]) => {
            percentages[mode] = (count / totalEntries) * 100;
        });

        const labels = Object.keys(percentages);
        const values = Object.values(percentages);
        const backgroundColors = this.generateRandomColors(labels.length); 

        this.setState({
            data: {
                labels: labels,
                datasets: [{
                    data: values,
                    backgroundColor: backgroundColors, 
                }],
            },
        });
    };

    calculatePercentagesOffOn = () => {
        const { data } = this.props;
        console.log(data)
        let turnOnCount = 0;
        let turnOffCount = 0;

        Object.values(data).forEach(entry => {
            if (entry.Action === "Turn on") {
                turnOnCount++;
            } else if (entry.Action === "Turn off") {
                turnOffCount++;
            }
        });

        const totalEntries = Object.values(data).length;

        const turnOnPercentage = (turnOnCount / totalEntries) * 100;
        const turnOffPercentage = (turnOffCount / totalEntries) * 100;

        const backgroundColors = this.generateRandomColors(2); 

        this.setState({
            data: {
                labels: ["Turn On", "Turn Off"],
                datasets: [{
                    data: [turnOnPercentage, turnOffPercentage],
                    backgroundColor: backgroundColors, 
                }],
            },
        });
    };

    calculateUserActivity = () => {
        const { data } = this.props;
        const userActivity = {};

        Object.values(data).forEach(entry => {
            const user = entry.User;
            userActivity[user] = (userActivity[user] || 0) + 1;
        });

        const totalEntries = Object.values(data).length;

        const userPercentages = {};
        Object.entries(userActivity).forEach(([user, count]) => {
            userPercentages[user] = (count / totalEntries) * 100;
        });

        const labels = Object.keys(userPercentages);
        const values = Object.values(userPercentages);
        const backgroundColors = this.generateRandomColors(labels.length); 

        this.setState({
            data: {
                labels: labels,
                datasets: [{
                    data: values,
                    backgroundColor: backgroundColors,
                }],
            },
        });
    };

    generateRandomColors = (numColors) => {
        const colors = [];
        for (let i = 0; i < numColors; i++) {
            const color = '#' + Math.floor(Math.random() * 16777215).toString(16);
            colors.push(color);
        }
        return colors;
    };

    render() {
        return (
            <div>
                <Pie
                    data={this.state.data}
                    options={{
                        maintainAspectRatio: false,
                    }}
                />
            </div>
        );
    }
}

export default PieChart;
