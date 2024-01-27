// HistoryGraph.js
import React from 'react';
import { Line } from 'react-chartjs-2';

const SPGraph = ({ data }) => {
  const chartData = {
    labels: data.Labels,
    datasets: [
      {
        label: 'Switch status',
        data: data.Values,
        borderColor: 'rgba(75,192,192,1)',
        borderWidth: 2,
        fill: false,
        stepped: 'before', // Set the stepped property to 'before'
      },
    ],
  };

    const options = {
      scales: {
        y: {
          beginAtZero: false,
        },
      },
  };

  return <Line style={{marginTop: "30px", marginBottom: "5px"}} data={chartData} options={options} />;
};

export default SPGraph;
