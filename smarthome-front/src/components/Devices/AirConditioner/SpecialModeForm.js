import React, { Component } from 'react';
import {
    Button,
    TextField,
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
    DialogContentText,
    DialogActions,
    IconButton,
    Snackbar,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import './SpecialModeForm.css'


class SpecialModeForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            start: '',
            end: '',
            mode: '',
            temperature: 20,
            selectedDays: [],
            specialModes: [],
            showDialog: false,
            snackbarMessage: '',
            showSnackbar: false,
            openSnackbar: false,
        };
    }

    handleSelectedDays = (event) => {
        const selectedDays = event.target.value;

        this.setState((prevState) => ({
            selectedDays,
        }));
    };

    handleAdd = () => {
        const { start, end, mode, temperature, selectedDays, specialModes } = this.state;
        const { acModes, minTemp, maxTemp } = this.props;
        const selectedMode = mode === '' ? acModes[0] : this.state.mode;
        // check if a day is selected
        if(selectedDays.length === 0) {
            this.setState({ snackbarMessage: "Please select a day" });
            this.handleClick();
            return;
        }
        // check if a temp between min and max
        if (temperature < minTemp || temperature > maxTemp) {
            this.setState({ snackbarMessage: "Temperature is out of the range" });
            this.handleClick();
            return;
        }
        // check if is start < end
        // if (new Date(`2000-01-01 ${start}`) >= new Date(`2000-01-01 ${end}`)) {
        //     this.setState({ snackbarMessage: "Start time must be before end time" });
        //     this.handleClick();
        //     return;
        // }
        // check if something already exists for that day, if there is, check if the start or end is between those 2 times
        const existingModeForDay = specialModes.find(
            (item) =>
                item.selectedDays.some((day) => selectedDays.includes(day)) &&
                ((new Date(`2000-01-01 ${start}`) >= new Date(`2000-01-01 ${item.start}`) &&
                    new Date(`2000-01-01 ${start}`) <= new Date(`2000-01-01 ${item.end}`)) ||
                    (new Date(`2000-01-01 ${end}`) >= new Date(`2000-01-01 ${item.start}`) &&
                        new Date(`2000-01-01 ${end}`) <= new Date(`2000-01-01 ${item.end}`)))
        );
    
        if (existingModeForDay) {
            this.setState({ snackbarMessage: "There is already a mode for selected day and time" });
            this.handleClick();
            return;
        }

        // if everything alright add mode
        const specialMode = {
            start,
            end,
            selectedMode,
            temperature,
            selectedDays,
        };

        console.log(specialMode)

        // reset data
        this.setState({
            specialModes: [...specialModes, specialMode],
            start: '',
            end: '',
            mode: selectedMode,
            temperature: 20,
            selectedDays: [],
        });
    };

    openDialog = () => {
        this.setState({ showDialog: true });
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

    handleSave = () => {
        this.setState({ showDialog: false });
        this.props.onAdd(this.state.specialModes);
    }

    render() {

        const { acModes, minTemp, maxTemp } = this.props;

        return (
        <div>
            <Button variant="contained" color="primary" onClick={this.openDialog} disabled={acModes.length === 0 || minTemp >= maxTemp}>
                Add Special Mode
            </Button>

            <Dialog open={this.state.showDialog} onClose={this.closeDialog}>
                <DialogTitle style={{marginBottom: '25px'}}>Add Special Mode</DialogTitle>
                <DialogContent>
                    <div className='firstRow'>
                        <span className="a">Start:</span>
                        <FormControl>
                            <Input style={{cursor: 'pointer'}} type="time" value={this.state.start} onChange={(e) => this.setState({ start: e.target.value })} />
                        </FormControl>

                        <span style={{ marginLeft: '1.5em' }} className="a">End:</span>
                        <FormControl>
                            <Input style={{cursor: 'pointer'}} type="time" value={this.state.end} onChange={(e) => this.setState({ end: e.target.value })} />
                        </FormControl>
                    </div>

                    <div className='secondRow'>
                        <span className="a">Mode:</span>
                        <FormControl>
                            <Select value={this.state.mode === '' ? acModes[0] : this.state.mode } onChange={(e) => this.setState({ mode: e.target.value })}>
                                {acModes.map((mode) => (
                                    <MenuItem key={mode} value={mode}>
                                    {mode.charAt(0).toUpperCase() + mode.slice(1)}
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>

                        <span className="a">Temperature:</span>
                        <FormControl style={{ width: '80px' }}>
                            <Input
                                type="number"
                                value={this.state.temperature}
                                onChange={(e) => this.setState({ temperature: e.target.value })}
                            />
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
                                <TableCell>Mode</TableCell>
                                <TableCell>Temperature</TableCell>
                                <TableCell>Day</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {this.state.specialModes.map((item, index) => (
                                <TableRow key={index}>
                                    <TableCell>{item.start}</TableCell>
                                    <TableCell>{item.end}</TableCell>
                                    <TableCell>{item.selectedMode}</TableCell>
                                    <TableCell>{item.temperature}</TableCell>
                                    <TableCell>
                                        {item.selectedDays.map((day, dayIndex) => (
                                            <Chip key={dayIndex} label={day} />
                                        ))}
                                    </TableCell>
                                    <TableCell>
                                        <IconButton color="secondary" onClick={() => this.handleDelete(index)}>
                                            <CloseIcon />
                                        </IconButton>
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
    );
    }
}

export default SpecialModeForm;
