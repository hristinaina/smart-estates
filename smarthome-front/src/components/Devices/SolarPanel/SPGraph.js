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
        beginAtZero: true, 
        min: 0, 
        max: 1, 
        ticks: {
          stepSize: 1,
          callback: function(value) {
            // Display only 0 and 1 on the y-axis
            if (value === 0) {
              return "off";
            } else if (value === 1){
              return "on";
            } else return null; // Skip other values
          },
        },
      },
    },
  };

  return <Line style={{marginTop: "30px", marginBottom: "5px"}} data={chartData} options={options} />;
};

export default SPGraph;
