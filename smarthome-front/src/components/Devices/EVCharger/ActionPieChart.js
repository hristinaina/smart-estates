import React, { Component } from 'react';
import { Pie } from 'react-chartjs-2';

class ActionPieChart extends Component {
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
        this.calculateActionPercentages();
    }

    componentDidUpdate(prevProps) {
        if (prevProps.data !== this.props.data) {
            if(this.props.graph === 1)
                this.calculateActionPercentages();
            else if (this.props.graph === 2)
                this.calculateUserPercentages();
            else if (this.props.graph === 3)
                this.calculatePlugPercentages();
        }
    }

    calculateActionPercentages = () => {
        const { data } = this.props;
        const modeCounts = {};
        console.log(data);
        Object.values(data).forEach(entry => {
            const mode = entry.Action;
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

    calculateUserPercentages = () => {
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

    calculatePlugPercentages = () => {
        const { data } = this.props;
        const userActivity = {};

        Object.values(data).forEach(entry => {
            let user = entry.Plug; 
            if (entry.Plug != -1)
                user = entry.Plug + 1;
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

export default ActionPieChart;
