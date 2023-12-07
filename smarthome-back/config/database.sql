CREATE DATABASE smart_home;
use smart_home;

CREATE TABLE user (
                    Id INT PRIMARY KEY,
                    Email VARCHAR(255) UNIQUE,
                    Password VARCHAR(255),
                    Name VARCHAR(255),
                    Surname VARCHAR(255),
                    Role INT,
                    isLogin BOOLEAN DEFAULT false
);

CREATE TABLE realestate (
                            Id INT PRIMARY KEY,
                            Name VARCHAR(255),
                            Type INT,
                            Address VARCHAR(255),
                            City VARCHAR(255),
                            SquareFootage FLOAT(32),
                            NumberOfFloors INT,
                            Picture VARCHAR(255),
                            State INT,
                            UserId INT,
                            DiscardReason VARCHAR(255)
);

CREATE TABLE device (
                        Id INT PRIMARY KEY AUTO_INCREMENT,
                        Name VARCHAR(255) NOT NULL,
                        Type INT NOT NULL,
                        RealEstate INT NOT NULL,
                        IsOnline BOOLEAN,
                        StatusTimeStamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE consumptionDevice (
                                DeviceId INT PRIMARY KEY,
                                PowerSupply INT NOT NULL,
                                PowerConsumption DOUBLE,
                                FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE airConditioner (
                                DeviceId INT PRIMARY KEY,
                                MinTemperature INT NOT NULL,
                                MaxTemperature INT NOT NULL,
                                FOREIGN KEY (DeviceId) REFERENCES consumptionDevice(DeviceId)
);

CREATE TABLE evCharger (
                        DeviceId INT PRIMARY KEY,
                        ChargingPower DOUBLE NOT NULL,
                        Connections INT NOT NULL,
                        FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE homeBattery (
                            DeviceId INT PRIMARY KEY,
                            Size DOUBLE NOT NULL,
                            FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

INSERT INTO user (Id, Email, Password, Name, Surname, Role)
VALUES
    (1, 'nesa@gmail.com', '$2a$10$/fbbLLHt7hEZpMq3rQiWz.oF6cJRRNdO5Vek/NIGzAIyJp99jebrO', 'Nenad', 'Peric', 0),
    (2, 'nata@gmail.com', '$2a$10$rs45oZDdYuLSOmzOdsJGS..HJ.9zmguT0r4cUt131XKqkac4P/7iu', 'Natasa', 'Maric', 1);

INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId, DiscardReason)
VALUES
    (1, 'Villa B Dorm',  0, '123 Main St', 'Cityville', 150.5, 2, 'path/to/picture1.jpg', 0, 1, ''),
    (2, 'Neka kuca nmp', 1, '456 Oak Ave', 'Townton', 200.75, 3, 'path/to/picture2.jpg', 1, 1, ''),
    (3, 'Joj stvarno nzm', 0, '789 Pine Ln', 'Villageto wn', 30.25, 1, 'path/to/picture3.jpg', 0, 1, ''),
    (4, 'Spavamise', 1, '101 Elm Blvd', 'Hamlet City', 700.0, 2, 'path/to/picture4.jpg', 2, 2, 'jer mi se moze'),
    (5, 'ma ne znam', 0, '102 Elm Blvd', 'Hamlet City', 65.0, 2, 'path/to/picture5.jpg', 0, 2, ''),
    (6, 'Spavamise2', 1, '103 Elm Blvd', 'Hamlet City', 70.0, 2, 'path/to/picture6.jpg', 0, 2, '');

INSERT INTO device (Id, Name, Type, RealEstate, IsOnline, StatusTimeStamp)
VALUES
    (1, 'Masina Sladja', 2,  1, false, '2023-12-06 15:30:00'),
    (2, 'Prsk prsk', 5, 1, false, '2023-12-06 15:30:00'),
    (3, 'Neka klima', 1, 2, false, '2023-12-06 15:30:00'),
    (4, 'Panelcic', 6, 2, false, '2023-12-06 15:30:00'),
    (5, 'Punjac1', 8, 2, false, '2023-12-06 15:30:00'),
    (6, 'Baterija1', 7, 2, false, '2023-12-06 15:30:00');

INSERT INTO consumptionDevice (DeviceId, PowerSupply, PowerConsumption)
VALUES
    (1, 1, 200),
    (2, 0, 0),
    (3, 1, 300);

INSERT INTO airConditioner (DeviceId, MinTemperature, MaxTemperature)
VALUES
    (3, 10, 30);

INSERT INTO evCharger (DeviceId, ChargingPower, Connections)
VALUES
    (5, 10, 2);

INSERT INTO homeBattery (DeviceId, Size)
VALUES
    (6, 10);