// HistoryGraph.js
import React from 'react';
import { Line } from 'react-chartjs-2';

const HBGraph = ({ data, name }) => {
    const chartData = {
        labels: data.timestamps,
        datasets: [
            {
                label: name,
                data: data.consumptionData,
                borderColor: 'rgba(75,192,192,1)',
                borderWidth: 2,
                fill: false,
            },
        ],
    };

    const options = {
        scales: {
            x: data.x,
            y: {
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'kWh',
                },
            },
        },
    };

    return <Line id='graph' style={{ marginTop: "30px", marginBottom: "5px" }} data={chartData} options={options} />;
};

export default HBGraph;
