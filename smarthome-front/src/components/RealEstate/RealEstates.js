
import React,{ Component, useState } from 'react';
import './RealEstates.css';
import Dialog from '../Dialog/Dialog';
import RealEstateService from '../../services/RealEstateService';

import { Navigation } from '../Navigation/Navigation';

export class RealEstates extends Component {

    constructor(props) {
        super(props);

        this.state = {
            isAdmin: true,
            showApproveDialog: false,
            showDiscardDialog: false,
            realEstates: [],
        };
    }

    async componentDidMount() {
        try {
            if (!this.state.isAdmin) {
                const result = await RealEstateService.getRealEstates();
                this.setState({realEstates: result})
                console.log(result);
            } else {
                const result = await RealEstateService.getPendingRealEstates();
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

    handleConfirmDiscard = () => {
        this.setState({showDiscardDialog: false});
        console.log("Request discarted...")
    }

    handleConfirmApprove = () => {
        this.setState({showApproveDialog: false});
        console.log("Approved");
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
                    {this.state.realEstates.map((realEstate, index) => (
                        <div className='real-estate-card'>
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
                    ))}
                </div>

                {this.state.isAdmin && (
                    <div id="bottom-bar">
                    <button className='bottom-bar-btn' id='bottom-bar-approve' onClick={this.handleApprove}>APPROVE</button>
                    <button className='bottom-bar-btn' id='bottom-bar-discard' onClick={this.handleDiscard}>DISCARD</button>
                 </div>
                )}
                 
                {this.state.showApproveDialog && (
                <Dialog
                    title="Approve Real Estate Request"
                    message="Are you sure you want to approve selected real-estate request?"
                    onConfirm={this.handleConfirmApprove}
                    onCancel={this.handleCancel}
                />
                )}

                {this.state.showDiscardDialog && (
                <Dialog
                    title="Discard Real Estate Request"
                    message="Are you sure you want to discard selected real-estate request?"
                    onConfirm={this.handleConfirmDiscard}
                    onCancel={this.handleCancel}
                />
                )}
            </div>
        )
    }
}