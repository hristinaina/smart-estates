import React,{ Component, useState } from 'react';
import './RealEstates.css';
import { Navigation } from '../Navigation/Navigation';
import authService from '../../services/AuthService'
import Dialog from '../Dialog/Dialog';
import RealEstateService from '../../services/RealEstateService';
import ImageService from '../../services/ImageService';

export class RealEstates extends Component {

    constructor(props) {
        super(props);

        this.state = {
            showNewRealEstate: false,
            user : null,
            isAdmin: false,
            userId: -1,
            isDisabled: true,
            showApproveDialog: false,
            showDiscardDialog: false,
            selectedRealEstate: -1,
            realEstates: [],
            realEstateImages: {},
        };
    }

    async componentDidMount() {
        const currentUser = authService.getCurrentUser();
        if (currentUser['Role'] === 0 || currentUser['Role'] === 2) {
            await this.setState({isAdmin: true});
        } else {
            await this.setState({isAdmin: false});
        }
        
        await this.setState({user: currentUser, userId: currentUser.Id, });
        console.log("comp did mount");
        try {
            if (!this.state.isAdmin) {
                const result = await RealEstateService.getAllByUserId(currentUser.Id);
                await this.setState({realEstates: result})

            } else {
                const result = await RealEstateService.getPending();
                await this.setState({realEstates: result})
            }
            const realEstateImages = {};
            for (const realEstate of this.state.realEstates) {
                const imageUrl = await ImageService.getImage(realEstate.Name);
                realEstateImages[realEstate.Id] = imageUrl;
            }
            this.setState({realEstateImages});
           
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
        this.changeState(1, reason);
    }

    handleConfirmApprove = () => {
        this.setState({showApproveDialog: false, isDisabled: true, selectedRealEstate: -1});
        this.changeState(0, '');
    }

    handleCancel = () => {
        this.setState({showApproveDialog: false,
                       showDiscardDialog: false,
                      });
    }

    handleAddRealEstateClick = () => {
        window.location.href = '/new-real-estate';
    }

    handleCardClick = (id) => {
        this.setState({selectedRealEstate: id, isDisabled: false});
    }

    changeState = async (state, reason) => {
        try {
            await RealEstateService.changeState(state, this.state.selectedRealEstate, reason);
            this.componentDidMount();
        } catch (error) {
            console.log("Error");
            console.error(error);
        }
    }
 
    render() {
        const { user } = this.state;

        if (!user) return null;
        
        return (
            <div id="real-estates-parent-container">
                <Navigation />
                {!this.state.isAdmin && (
                <p id="add-real-estate" onClick={this.handleAddRealEstateClick}>
                    <img alt="" src="/images/plus.png" id="plus" />
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
                            <img alt='real-estate' src={this.state.realEstateImages[realEstate.Id]} className='real-estate-img' />
                            <div className='real-estate-info'>
                                <p className='real-estate-title'>{realEstate.Name}</p>
                                <p className='real-estate-text'>
                                Type: {realEstate.Type === 0 ? 'HOME' : realEstate.Type === 1 ? 'APARTMENT' : 'VILLA'} </p>
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