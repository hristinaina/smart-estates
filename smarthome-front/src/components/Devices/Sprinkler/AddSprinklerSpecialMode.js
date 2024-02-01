import React, { Component } from 'react';
import {
    Button,
    Select,
    MenuItem,
    InputLabel,
    FormControl,
    Input,
    Chip,
    Table,
    TableContainer,
    TableHead,
    TableBody,
    TableRow,
    TableCell,
    Paper,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    IconButton,
    Snackbar,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import '../AirConditioner/SpecialModeForm.css'
import SprinklerService from '../../../services/SprinklerService';


class AddSprinklerSpecialMode extends Component {
    constructor(props) {
        super(props);
        this.state = {
            start: '',
            end: '',
            mode: '',
            selectedDays: [],
            specialModes: [],
            newSpecialModes: [],
            showDialog: false,
            snackbarMessage: '',
            showSnackbar: false,
            openSnackbar: false,
        };
        this.id = parseInt(this.extractDeviceIdFromUrl());
    }

    async componentDidMount() {
    }

    getSelectedDays(selectedDays) {
        return selectedDays.split(',').filter(day => day !== "");
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleSelectedDays = (event) => {
        const selectedDays = event.target.value;

        this.setState((prevState) => ({
            selectedDays,
        }));
    };

    handleAdd = () => {
        const { start, end, mode, selectedDays, specialModes } = this.state;
        // check if a day is selected
        if(selectedDays.length === 0) {
            this.setState({ snackbarMessage: "Please select a day" });
            this.handleClick();
            return;
        }
        // check if something already exists for that day, if there is, check if the start or end is between those 2 times
        const startDate = new Date(`2000-01-01 ${start}`);
        const endDate = new Date(`2000-01-01 ${end}`);

        // Check if end is before start, then adjust the date to the next day
        if (endDate < startDate) {
            endDate.setDate(endDate.getDate() + 1);
        }
        const existingModeForDay = specialModes.find(
            (item) =>
            {
                const itemStartDate = new Date(`2000-01-01 ${item.start}`);
                const itemEndDate = new Date(`2000-01-01 ${item.end}`);

                // Check if item.end is before item.start, then adjust the date to the next day
                if (itemEndDate < itemStartDate) {
                    itemEndDate.setDate(itemEndDate.getDate() + 1);
                }
                return (
                item.selectedDays.some((day) => selectedDays.includes(day)) &&
                (((startDate >= itemStartDate && startDate <= itemEndDate)) ||
                ((endDate >= itemStartDate && endDate <= itemEndDate)) ||
                (startDate <= itemStartDate && endDate >= itemEndDate)))
            }
              
                    
        );
        if (existingModeForDay) {
            this.setState({ snackbarMessage: "There is already a mode for selected day and time" });
            this.handleClick();
            return;
        }

        // if everything ok add mode
        const specialMode = {
            start,
            end,
            selectedDays,
        };

        let newSpecialModes = this.state.newSpecialModes

        // reset data
        this.setState({
            specialModes: [...specialModes, specialMode],
            start: '',
            end: '',
            selectedDays: [],
            newSpecialModes: [...newSpecialModes, specialMode],
        });
    };

    openDialog = async() => {
        this.setState({ showDialog: true });
        await this.setState({specialModes: this.props.modes});
    };

    closeDialog = () => {
        this.setState({ showDialog: false });
    };

    handleDelete = (index) => {
        const specialModesCopy = [...this.state.specialModes];
        specialModesCopy.splice(index, 1);
        this.setState({ specialModes: specialModesCopy });
    };

    // snackbar
    handleClick = () => {
        this.setState({ openSnackbar: true });
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        this.setState({ openSnackbar: false });
    };

    handleSave = async() => {
        this.setState({ showDialog: false });
        this.state.newSpecialModes.forEach(mode => {
            SprinklerService.addMode(mode, this.id);
        });
        this.props.onAdd(this.state.newSpecialModes);
        this.setState({specialModes: [], newSpecialModes: []});
    }

    render() {

        return (
        <div>
            <p id="special-mode" onClick={this.openDialog}>
                Add Special Mode
            </p>

            <Dialog open={this.state.showDialog} onClose={this.closeDialog}>
                <DialogTitle style={{marginBottom: '25px'}}>Add Special Mode</DialogTitle>
                <DialogContent>
                    <div className='firstRow' style={{marginBottom: '2em'}}>
                        <span className="a">Start:</span>
                        <FormControl>
                            <Input style={{cursor: 'pointer'}} type="time" value={this.state.start} onChange={(e) => this.setState({ start: e.target.value })} />
                        </FormControl>

                        <span style={{ marginLeft: '1.5em' }} className="a">End:</span>
                        <FormControl>
                            <Input style={{cursor: 'pointer'}} type="time" value={this.state.end} onChange={(e) => this.setState({ end: e.target.value })} />
                        </FormControl>
                    </div>

                    <FormControl id="ac-dropdown">
                        <InputLabel id="multi-select-dropdown-label">Select Days</InputLabel>
                        <Select
                            labelId="multi-select-dropdown-label"
                            id="multi-select"
                            multiple
                            value={this.state.selectedDays}
                            onChange={this.handleSelectedDays}
                            renderValue={(selected) => selected.join(', ')}>
                        <MenuItem value="Monday">Monday</MenuItem>
                        <MenuItem value="Tuesday">Tuesday</MenuItem>
                        <MenuItem value="Wednesday">Wednesday</MenuItem>
                        <MenuItem value="Thursday">Thursday</MenuItem>
                        <MenuItem value="Friday">Friday</MenuItem>
                        <MenuItem value="Saturday">Saturday</MenuItem>
                        <MenuItem value="Sunday">Sunday</MenuItem>
                        </Select>
                    </FormControl>
                <Button variant="contained" color="primary" onClick={this.handleAdd}>ADD</Button>


                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>Start</TableCell>
                                <TableCell>End</TableCell>
                                <TableCell>Day</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {this.state.specialModes.map((item, index) => (
                                <TableRow key={index}>
                                    <TableCell>{item.start}</TableCell>
                                    <TableCell>{item.end}</TableCell>
                                    <TableCell>
                                        {item.selectedDays.map((day, dayIndex) => (
                                            <Chip key={dayIndex} label={day} />
                                        ))}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </DialogContent>

            <DialogActions style={{ display: 'flex', justifyContent: 'space-between' }}>
                <div>
                    <Button onClick={this.closeDialog} color="primary">
                        Close
                    </Button>
                </div>
                <div>
                    <Button variant="contained" onClick={this.handleSave} color="primary">
                        Save
                    </Button>
                </div>
            </DialogActions>
            </Dialog>

            <Snackbar
                    open={this.state.openSnackbar}
                    autoHideDuration={3000}
                    onClose={this.handleClose}
                    message={this.state.snackbarMessage}
                    action={this.action}
                />
        </div>
    );}
}

export default AddSprinklerSpecialMode;
