import { Component } from 'react';
import '../Devices/Devices.css';
import { Navigation } from '../Navigation/Navigation';
import authService from '../../services/AuthService';
import 'chart.js/auto';
import { Autocomplete, TextField, Button, Box, Snackbar } from '@mui/material';
import "./ElectricityOverview.css"
import SearchSelect from './SearchSelect';
import ElectricityService from '../../services/ElectricityService';
import AdminGraph from './AdminGraph';

export class ElectricityOverview extends Component {
    constructor(props) {
        super(props);
        this.state = {
            consumptionData: [],
            productionData: [],
            ratioData: [],
            edData: [],
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
            selectedOption: { label: '6h', value: '-6h' },
            selectedTypeOption: { label: 'By City', value: 'city' },
            startDate: '',
            endDate: '',
            options: [
                { label: '6h', value: '-6h' },
                { label: '12h', value: '-12h' },
                { label: '24h', value: '-24h' },
                { label: 'last week', value: '-7d' },
                { label: 'last month', value: '-30d' },
            ],
            typeOptions: [
                { label: 'By City', value: 'city' },
                { label: 'By Real Estate', value: 'rs' }
            ],
        };
        this.selectedOptions = [];
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");
    }

    isTimestampInLastHour = (timestamp) => {
        const currentTimestamp = new Date();
        const timestampDate = new Date(timestamp);

        const timeDifference = currentTimestamp - timestampDate;

        return timeDifference <= 3600000;
    };

    updateGraph = async (value) => {
        const resultC = await ElectricityService.getConsumptionGraphDataForDropdownSelect("consumption", this.state.selectedTypeOption.value, this.selectedOptions, value);
        const resultP = await ElectricityService.getConsumptionGraphDataForDropdownSelect("solar_panel", this.state.selectedTypeOption.value, this.selectedOptions, value);
        const ratio = await ElectricityService.getRatioGraphDataForDropdownSelect(this.state.selectedTypeOption.value, this.selectedOptions, value, "");
        const edratio = await ElectricityService.getRatioGraphDataForDropdownSelect(this.state.selectedTypeOption.value, this.selectedOptions, value, "electrical_distribution");

        let showMinutes = true;
        if (!["-6h", "-12h", "-24h"].includes(this.state.selectedOption.value))
            showMinutes = false
        const graphData = this.convertResultToGraphData(resultC.result.result, showMinutes, "Electricity Consumption")
        const graphProduction = this.convertResultToGraphData(resultP.result.result, showMinutes, "Electricity Production")
        const graphRatio = this.convertResultToGraphData(ratio.result.result, showMinutes, "Production & Consumption Ratio")
        const edRatio = this.convertResultToGraphData(edratio.result.result, showMinutes, "Electrical Distribution")

        this.setState({
            consumptionData: graphData,
            productionData: graphProduction,
            ratioData: graphRatio,
            edData: edRatio,
        });
    }

    handleTypeSelectChange = (event, selectedOption) => {
        this.setState({ selectedTypeOption: selectedOption });
    };

    handleOptionChange = async (event, value) => {
        if (this.selectedOptions.length === 0) {
            this.setState({ snackbarMessage: "You haven't selected any subjects" });
            this.handleClick();
            return;
        }
        this.setState({ selectedOption: value });
        await this.updateGraph(value.value)
    };

    handleDateChange = (fieldName, event) => {
        this.setState({ [fieldName]: event.target.value });
    };

    handleButtonClick = async () => {
        if (this.selectedOptions.length === 0) {
            this.setState({ snackbarMessage: "You haven't selected any subjects" });
            this.handleClick();
            return;
        }
        if (this.state.startDate === '' || this.state.endDate === '') {
            this.setState({ snackbarMessage: "Please enter dates" });
            this.handleClick();
            return;
        }

        if (new Date(this.state.startDate) > new Date(this.state.endDate)) {
            this.setState({ snackbarMessage: "Start date cannot be greater than end date" });
            this.handleClick();
            return;
        }
        const oneMonth = 30 * 24 * 60 * 60 * 1000;
        const twoDays = 2 * 24 * 60 * 60 * 1000
        const difference = new Date(this.state.endDate) - new Date(this.state.startDate);
        let showMinutes = true;
        if (difference > oneMonth) {
            this.setState({ snackbarMessage: 'The difference between start date and end date must not be more than one month' });
            this.handleClick();
            return;
        }
        if (difference > twoDays) {
            showMinutes = false;
        }
        const resultC = await ElectricityService.getConsumptionGraphDataForDates("consumption", this.state.selectedTypeOption.value, this.selectedOptions, this.state.startDate, this.state.endDate);
        const resultP = await ElectricityService.getConsumptionGraphDataForDates("solar_panel", this.state.selectedTypeOption.value, this.selectedOptions, this.state.startDate, this.state.endDate);
        const ratio = await ElectricityService.getRatioGraphDataForDates(this.state.selectedTypeOption.value, this.selectedOptions, this.state.startDate, this.state.endDate, "");
        const edratio = await ElectricityService.getRatioGraphDataForDates(this.state.selectedTypeOption.value, this.selectedOptions, this.state.startDate, this.state.endDate, "electrical_distribution");

        const graphData = this.convertResultToGraphData(resultC.result.result, showMinutes, "Electricity Consumption")
        const graphProduction = this.convertResultToGraphData(resultP.result.result, showMinutes, "Electricity Production")
        const graphRatio = this.convertResultToGraphData(ratio.result.result, showMinutes, "Production & Consumption Ratio")
        const edRatio = this.convertResultToGraphData(edratio.result.result, showMinutes, "Electrical Distribution")

        this.setState({
            consumptionData: graphData,
            productionData: graphProduction,
            ratioData: graphRatio,
            edData: edRatio,
        });
    }

    getRandomColor() {
        const random = () => Math.floor(Math.random() * 256);
        const r = random();
        const g = random();
        const b = random();
        const a = 1; // You can adjust the alpha (transparency) value if needed
        return `rgba(${r},${g},${b},${a})`;
    }

    convertResultToGraphData(values, showMinutes, type) {
        if (values == null) {
            return {
                timestamps: [],
                datasets: [],
                x: {},
            }
        }
        //values is a map["estateId"]map[timestamp]float
        // Step 1: Combine timestamps from all inner maps
        const allTimestamps = Array.from(new Set(Object.values(values).flatMap(innerMap => Object.keys(innerMap)))).sort();
        // Step 2: Create arrays with values for each inner map
        const keyValuesArrays = Object.entries(values).map(([key, innerMap]) => ({
            label: key,
            data: allTimestamps.map(timestamp => innerMap[timestamp] || 0),
            borderColor: this.getRandomColor(),
            borderWidth: 2,
            fill: false,
        }));
        const graphData = {
            timestamps: allTimestamps,
            datasets: keyValuesArrays,
            x: {
                type: 'time',
                time: {
                    displayFormats: {
                        quarter: 'HH:MM'
                    }
                }
            },
            title: type,
        };
        if (!showMinutes) {
            graphData.x.time = {
                unit: 'day',
                displayFormats: {
                    day: 'MMM d',
                },
            };
        }
        return graphData;
    }

    handleSearchSelectChange = (selectedOptions) => {
        this.selectedOptions = selectedOptions;
    };

    // snackbar
    handleClick = () => {
        this.setState({ open: true });
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        this.setState({ open: false });
    };

    render() {
        const { consumptionData, productionData, ratioData, edData, startDate, endDate, selectedOption, options, selectedTypeOption, typeOptions } = this.state;

        return (
            <div>
                <Navigation />
                <div className='c-card'>
                    <p className='sp-card-title'>Electricity consumption overview</p>
                    <div className='c-tools-container'>
                        <Autocomplete
                            value={selectedTypeOption}
                            onChange={this.handleTypeSelectChange}
                            options={typeOptions}
                            getOptionLabel={(option) => option.label}
                            style={{ width: '260px', marginLeft: "auto" }}
                            renderInput={(params) => (
                                <TextField
                                    style={{ backgroundColor: "white" }}
                                    {...params}
                                    label="Select Category"
                                    InputLabelProps={{
                                        shrink: true,
                                    }}
                                />
                            )}
                            isOptionEqualToValue={(option, value) => option.value === value.value}
                            renderOption={(props, option, { selected }) => (
                                <li {...props}>
                                    <span>{option.label}</span>
                                </li>
                            )}
                            disableClearable />
                        <SearchSelect
                            options={selectedTypeOption.value}
                            onOptionsChange={this.handleSearchSelectChange}
                        />
                    </div>
                    <div className='c-tools-container'>
                        <Autocomplete
                            value={selectedOption}
                            onChange={this.handleOptionChange}
                            options={options}
                            getOptionLabel={(option) => option.label}
                            style={{ width: '260px', marginLeft: "auto" }}
                            renderInput={(params) => (
                                <TextField
                                    style={{ backgroundColor: "white" }}
                                    {...params}
                                    label="Select Time Range"
                                    InputLabelProps={{
                                        shrink: true,
                                    }}
                                />
                            )}
                            isOptionEqualToValue={(option, value) => option.value === value.value}
                            renderOption={(props, option, { selected }) => (
                                <li {...props}>
                                    <span>{option.label}</span>
                                </li>
                            )}
                            disableClearable />
                        <p style={{ marginLeft: "20px", marginRight: "20px" }}><b>or</b></p>
                        <Box display="flex" alignItems="center" style={{ marginRight: "auto" }}>
                            <TextField
                                style={{ backgroundColor: "white", marginRight: "5px" }}
                                label="Start Date"
                                type="date"
                                value={startDate}
                                onChange={(e) => this.handleDateChange('startDate', e)}
                                InputLabelProps={{
                                    shrink: true,
                                }}
                                inputProps={{
                                    max: new Date(new Date().setDate(new Date().getDate() + 1)).toISOString().split('T')[0],
                                }}
                            />
                            <TextField
                                style={{ backgroundColor: "white", marginRight: "7px" }}
                                label="End Date"
                                type="date"
                                value={endDate}
                                onChange={(e) => this.handleDateChange('endDate', e)}
                                InputLabelProps={{
                                    shrink: true,
                                }}
                                inputProps={{
                                    max: new Date(new Date().setDate(new Date().getDate() + 1)).toISOString().split('T')[0],
                                }}
                            />
                            <Button variant="contained" color="primary" onClick={this.handleButtonClick}>
                                Apply
                            </Button>
                        </Box>
                    </div>
                </div>
                <div className='sp-container'>
                    <div className='c-left-card'><AdminGraph data={consumptionData} /></div>
                    <div className='c-right-card'><AdminGraph data={productionData} /></div>
                </div>
                <div className='sp-container'>
                    <div className='c-left-card'><AdminGraph data={ratioData} /></div>
                    <div className='c-right-card'><AdminGraph data={edData} /></div>
                </div>
                <Snackbar
                    open={this.state.open}
                    autoHideDuration={3000}
                    onClose={this.handleClose}
                    message={this.state.snackbarMessage}
                    action={this.action}
                />
            </div>
        )
    }
}

