import React, { useState, useEffect } from 'react';
import authService from '../../services/AuthService';

const ActivationPage = () => {
  const [activationStatus, setActivationStatus] = useState(null);

  useEffect(() => {
    const activateAccount = async () => {
      try {
        const result = await authService.activateAccount();

        if (result.success) {
          setActivationStatus('Nalog je uspešno aktiviran!');
        } else {
          setActivationStatus('Došlo je do greške prilikom aktivacije naloga.');
        }
      } catch (error) {
        setActivationStatus('Došlo je do greške prilikom slanja zahteva.');
        console.error('Greška prilikom slanja zahteva:', error);
      }
    };

    activateAccount();
  }, []);

  return (
    <div>
      <h1>Provera aktivacije naloga</h1>
      {activationStatus && <p>{activationStatus}</p>}
    </div>
  );
};

export default ActivationPage;
