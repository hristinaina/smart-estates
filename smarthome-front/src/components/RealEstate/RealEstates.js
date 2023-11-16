
import React,{ Component, setState } from 'react';
import './RealEstates.css';
import Dialog from '../Dialog/Dialog';
import RealEstateService from '../../services/RealEstateService';

import { Navigation } from '../Navigation/Navigation';

export class RealEstates extends Component {

    constructor(props) {
        super(props);

        this.state = {
            isAdmin: true,
            isDisabled: true,
            showApproveDialog: false,
            showDiscardDialog: false,
            selectedRealEstate: -1,
            realEstates: [],
        };
    }

    async componentDidMount() {
        try {
            if (!this.state.isAdmin) {
                const result = await RealEstateService.get();
                this.setState({realEstates: result})
                console.log(result);
            } else {
                const result = await RealEstateService.getPending();
                this.setState({realEstates: result})
                console.log(result);
            }
           
        } catch (error) {
            console.log("error");
            console.error(error);
        }
    }

    handleApprove = () => {
        this.setState({showApproveDialog: true, 
                       showDiscardDialog: false,
                      })
    }

    handleDiscard = () => {
        this.setState({showDiscardDialog: true,
                        showApproveDialog: false,
                      })
    }

    handleConfirmDiscard = (reason) => {
        this.setState({showDiscardDialog: false, isDisabled: true, selectedRealEstate: -1});
        console.log("Request discarted...");
        console.log("Reason:", reason);
        this.changeState(1, reason);
    }

    handleConfirmApprove = () => {
        this.setState({showApproveDialog: false, isDisabled: true, selectedRealEstate: -1});
        console.log("Approved");
        this.changeState(0, '');
    }

    handleCancel = () => {
        this.setState({showApproveDialog: false,
                       showDiscardDialog: false,
                      });
        console.log("Cancelled...");
    }

    handleAddRealEstateClick = () => {
        window.location.href = '/new-real-estate';
    }

    handleCardClick = (id) => {
        this.setState({selectedRealEstate: id, isDisabled: false});
    }

    changeState = async (state, reason) => {
        try {
            const result = await RealEstateService.changeState(state, this.state.selectedRealEstate, reason);
            console.log("Success");
            console.log(result);
            this.componentDidMount();
        } catch (error) {
            console.log("Error");
            console.error(error);
        }
    }
 
    render() {

        return (
            <div id="real-estates-parent-container">
                <Navigation />
                {!this.state.isAdmin && (
                <p id="add-real-estate" onClick={this.handleAddRealEstateClick}>
                    <img alt="." src="/images/plus.png" id="plus" />
                    Add Real-Estate
                </p>
                )}
                
                <div id='real-estates-container'>
                    {this.state.realEstates !== null  ? (
                    this.state.realEstates.map((realEstate) => (
                        <div 
                            key={realEstate.Id}
                            className={`real-estate-card ${(realEstate.Id !== this.state.selectedRealEstate && this.state.isAdmin === true) ? 'not-selected-card' : 'selected-card'}`} 
                            onClick={() => this.handleCardClick(realEstate.Id)}>
                            <img alt='real-estate' src='/images/real_estate_example.png' className='real-estate-img' />
                            <div className='real-estate-info'>
                                <p className='real-estate-title'>{realEstate.Name}</p>
                                <p className='real-estate-text'>Address: {realEstate.Address}</p>
                                <p className='real-estate-text'>Square Footage: {realEstate.SquareFootage}</p>
                                <p className='real-estate-text'>Number of Floors: {realEstate.NumberOfFloors}</p>
                                <p className={`real-estate-text ${realEstate.State === 1 ? 'accepted' : realEstate.State === 0 ? 'pending' : 'declined'}`}>
                                    {realEstate.State === 1 ? 'Accepted' : realEstate.State === 0 ? 'Pending' : 'Declined'}
                                </p>
                            </div>
                        </div>
                    ))): (<p id="nothing-available">No real estates available.</p>)}
                </div>

                {this.state.isAdmin && (
                    <div id="bottom-bar">
                    <button className='bottom-bar-btn' id='bottom-bar-approve' onClick={this.handleApprove} disabled={this.state.isDisabled}>APPROVE</button>
                    <button className='bottom-bar-btn' id='bottom-bar-discard' onClick={this.handleDiscard} disabled={this.state.isDisabled}>DISCARD</button>
                 </div>
                )}
                 
                {this.state.showApproveDialog && (
                <Dialog
                    title="Approve Real Estate Request"
                    message="Are you sure you want to approve selected real-estate request?"
                    onConfirm={this.handleConfirmApprove}
                    onCancel={this.handleCancel}
                    isDiscard={false}
                />
                )}

                {this.state.showDiscardDialog && (
                <Dialog
                    title="Discard Real Estate Request"
                    message="Are you sure you want to discard selected real-estate request?"
                    onConfirm={this.handleConfirmDiscard}
                    onCancel={this.handleCancel}
                    isDiscard={true}
                />
                )}
            </div>
        )
    }
}