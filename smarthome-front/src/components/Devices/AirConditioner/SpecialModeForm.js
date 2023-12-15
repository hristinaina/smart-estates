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
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import './SpecialModeForm.css'


class SpecialModeForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            start: '',
            end: '',
            mode: 'cooling',
            temperature: 20,
            selectedDays: [],
            specialModes: [],
            showDialog: false,
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

        const specialMode = {
            start,
            end,
            mode,
            temperature,
            selectedDays,
        };

        this.props.onAdd(specialMode);

        // Resetovanje polja nakon dodavanja
        this.setState({
            specialModes: [...specialModes, specialMode],
            start: '',
            end: '',
            mode: 'cooling',
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

    render() {
        return (
        <div>
            <Button variant="contained" color="primary" onClick={this.openDialog}>
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
                            {/* <InputLabel>Mode</InputLabel> */}
                            <Select value={this.state.mode} onChange={(e) => this.setState({ mode: e.target.value })}>
                                <MenuItem value="cooling">Cooling</MenuItem>
                                <MenuItem value="heating">Heating</MenuItem>
                                <MenuItem value="automatic">Automatic</MenuItem>
                                <MenuItem value="ventilation">Ventilation</MenuItem>
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
                        <MenuItem value="Monday">Ponedeljak</MenuItem>
                        <MenuItem value="Tuesday">Utorak</MenuItem>
                        <MenuItem value="Wednesday">Sreda</MenuItem>
                        <MenuItem value="Thursday">Cetvrtak</MenuItem>
                        <MenuItem value="Friday">Petak</MenuItem>
                        <MenuItem value="Saturday">Subota</MenuItem>
                        <MenuItem value="Sunday">Nedelja</MenuItem>
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
                                    <TableCell>{item.mode}</TableCell>
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
        </div>
    );
    }
}

export default SpecialModeForm;
