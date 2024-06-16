// HistoryGraph.js
import React from 'react';
import { Line } from 'react-chartjs-2';

const AdminGraph = ({ data }) => {
    if (data.datasets == undefined) {
        data.datasets = [];
    }
    const chartData = {
        labels: data.timestamps,
        datasets: data.datasets,
    };

    const options = {
        scales: {
            x: data.x,
            y: {
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'kW/m2',
                },
            },
        },
        plugins: {
            title: {
                display: true,
                padding: {
                    top: 10,
                    bottom: 10,
                },
            },
        },
    };

    return <Line id='graph' style={{ marginTop: "30px", marginBottom: "5px" }} data={chartData} options={options} />;
};

export default AdminGraph;
