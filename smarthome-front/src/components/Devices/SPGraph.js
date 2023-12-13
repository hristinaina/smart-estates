// HistoryGraph.js
import React from 'react';
import { Line } from 'react-chartjs-2';

const SPGraph = ({ data }) => {
  const chartData = {
    labels: data.labels,
    datasets: [
      {
        label: 'Switch status',
        data: data.values,
        borderColor: 'rgba(75,192,192,1)',
        borderWidth: 2,
        fill: false,
        stepped: 'before', // Set the stepped property to 'before'
      },
    ],
  };

  const options = {
    scales: {
      x: [
        {
          type: 'time',
          time: {
            unit: 'hour',
          },
        },
      ],
      y: [
        {
          ticks: {
            beginAtZero: true,
            max: 1, // Set max to 1 for binary true/false values
          },
        },
      ],
    },
    plugins: {
      legend: {
        display: true,
        position: 'top',
      },
    },
  };

  return <Line data={chartData} options={options} />;
};

export default SPGraph;
