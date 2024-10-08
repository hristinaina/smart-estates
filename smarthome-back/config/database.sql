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
    StatusTimeStamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    LastValue DOUBLE DEFAULT -1
);

CREATE TABLE consumptionDevice (
    DeviceId INT PRIMARY KEY,
    PowerSupply INT NOT NULL,
    PowerConsumption DOUBLE,
    FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE airConditioner (
    DeviceId INT PRIMARY KEY,
    MinTemperature FLOAT(32) NOT NULL,
    MaxTemperature FLOAT(32) NOT NULL,
    Mode VARCHAR(255) NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES consumptionDevice(DeviceId)
);

CREATE TABLE specialModes (
    DeviceId INT PRIMARY KEY,
    StartTime TIME NOT NULL,
    EndTime TIME NOT NULL,
    Mode VARCHAR(255) NOT NULL,
    Temperature FLOAT(32) NOT NULL,
    SelectedDays VARCHAR(255) NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES airConditioner(DeviceId)
);

CREATE TABLE washingMachine (
    DeviceId INT PRIMARY KEY,
    Mode VARCHAR(255) NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES consumptionDevice(DeviceId)
);

CREATE TABLE machineMode (
    Id INT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Duration FLOAT(32) NOT NULL,
    Temp VARCHAR(255) NOT NULL
);

CREATE TABLE machineScheduledMode (
	Id INT PRIMARY KEY,
    DeviceId INT NOT NULL,
    StartTime TIMESTAMP NOT NULL,
    ModeId INT NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES washingMachine(DeviceId),
    FOREIGN KEY (ModeId) REFERENCES machineMode(Id)
);

CREATE TABLE sprinkler (
    DeviceId INT PRIMARY KEY,
    IsOn BOOL,
    FOREIGN KEY (DeviceId) REFERENCES consumptiondevice(DeviceId)
);

CREATE TABLE sprinklerSpecialMode(
    Id INT PRIMARY KEY AUTO_INCREMENT,
    DeviceId INT NOT NULL,
    StartTime TIME NOT NULL,
    EndTiME TIME NOT NULL,
    SelectedDays VARCHAR(255) NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES sprinkler(DeviceId)
);

CREATE TABLE evCharger (
    DeviceId INT PRIMARY KEY,
    ChargingPower DOUBLE NOT NULL,
    Connections INT NOT NULL,
    FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE lamp (
    DeviceId INT PRIMARY KEY,
    IsOn bool,
    LightningLevel int,
    FOREIGN KEY (DeviceId) REFERENCES consumptionDevice(DeviceId)
);

CREATE TABLE vehicleGate (
    DeviceId INT PRIMARY KEY,
    IsOpen bool,
    Mode int,
    FOREIGN KEY (DeviceId) REFERENCES consumptionDevice(DeviceId)
);

CREATE TABLE homeBattery (
    DeviceId INT PRIMARY KEY,
    Size DOUBLE NOT NULL,
    CurrentValue DOUBLE  DEFAULT 0.0,
    FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE solarPanel (
    DeviceId INT PRIMARY KEY,
    SurfaceArea DOUBLE NOT NULL,
    Efficiency DOUBLE NOT NULL,
    NumberOfPanels INT NOT NULL,
    IsOn BOOLEAN,
    FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

CREATE TABLE licensePlate (
    DeviceId INT,
    PlateNumber VARCHAR(255),
    FOREIGN KEY (DeviceId) REFERENCES vehicleGate(DeviceId)
);

CREATE TABLE permission (
    Id INT PRIMARY KEY AUTO_INCREMENT,
    RealEstateId INT NOT NULL,
    DeviceId INT NOT NULL,
    UserEmail VARCHAR(255) NOT NULL,
    isActive BOOLEAN NOT NULL,
    isDeleted BOOLEAN NOT NULL,
    FOREIGN KEY (RealEstateId) REFERENCES realestate(Id),
    FOREIGN KEY (DeviceId) REFERENCES device(Id)
);

INSERT INTO user (Id, Email, Password, Name, Surname, Role)
VALUES
    (1, 'nesa@gmail.com', '$2a$10$OwhZfsWej4wtXgey1DDWduCIGw4mI65wYdFcRY1ICno71gJKqNrJq', 'Nenad', 'Peric', 0),
    (2, 'nata@gmail.com', '$2a$10$OwhZfsWej4wtXgey1DDWduCIGw4mI65wYdFcRY1ICno71gJKqNrJq', 'Natasa', 'Maric', 1),
    (3, 'anastasijas557@gmail.com', '$2a$10$rs45oZDdYuLSOmzOdsJGS..HJ.9zmguT0r4cUt131XKqkac4P/7iu', 'Anastasija', 'Savic', 1),
    (4, 'kvucic6@gmail.com', '$2a$10$rs45oZDdYuLSOmzOdsJGS..HJ.9zmguT0r4cUt131XKqkac4P/7iu', 'Katarina', 'Vucic', 1);
    (5, 'hristinacina@gmail.com', '$2a$10$OwhZfsWej4wtXgey1DDWduCIGw4mI65wYdFcRY1ICno71gJKqNrJq', 'Hristina', 'Adamovic', 1),

INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId, DiscardReason)
VALUES
    (1, 'Villa B Dorm',  0, '123 Main St', 'Novi Sad', 150.5, 2, 'path/to/picture1.jpg', 0, 1, ''),
    (2, 'Neka kuca nmp', 1, '456 Oak Ave', 'Beograd', 200.75, 3, 'path/to/picture2.jpg', 1, 2, ''),
    (3, 'Joj stvarno nzm', 0, '789 Pine Ln', 'Zrenjanin', 30.25, 1, 'path/to/picture3.jpg', 1, 3, ''),
    (4, 'Spavamise', 1, '101 Elm Blvd', 'Novi Sad', 700.0, 2, 'path/to/picture4.jpg', 2, 2, 'jer mi se moze'),
    (5, 'ma ne znam', 0, '102 Elm Blvd', 'Beograd', 65.0, 2, 'path/to/picture5.jpg', 0, 2, ''),
    (6, 'Spavamise2', 1, '103 Elm Blvd', 'Zrenjanin', 70.0, 2, 'path/to/picture6.jpg', 0, 2, '');

INSERT INTO device (Id, Name, Type, RealEstate, IsOnline, StatusTimeStamp, LastValue)
VALUES
    (1, 'Masina Sladja', 2,  2, false, '2023-12-06 15:30:00', -1),
    (2, 'Prsk prsk', 5, 1, false, '2023-12-06 15:30:00', -1),
    (3, 'Neka klima', 1, 2, false, '2023-12-06 15:30:00', -1),
    (4, 'Panelcic', 6, 2, false, '2023-12-06 15:30:00', -1),
    (5, 'Punjac1', 8, 2, false, '2023-12-06 15:30:00', -1),
    (6, 'Baterija1', 7, 2, false, '2023-12-06 15:30:00', -1),
    (7, 'Lampica u sobici', 3, 1, false, '2023-12-06 15:30:00', -1),
    (8, 'Lampetina', 3, 1, false, '2023-12-06 15:30:00', -1),
    (9, 'Kapijica', 4, 1, false, '2023-12-06 15:30:00', -1),
    (10, 'Kapijetina', 4, 1, false, '2023-12-06 15:30:00', -1),
    (11, 'Prskalica', 5, 1, false, '2024-01-29 15:30:00', -1),
    (12, 'Panel 2', 6, 3, false, '2023-12-06 15:30:00', -1),
    (13, 'Punjac2', 8, 3, false, '2023-12-06 15:30:00', -1),
    (14, 'Baterija2', 7, 1, false, '2023-12-06 15:30:00', -1),
    (15, 'Masina 2', 2,  3, false, '2023-12-06 15:30:00', -1),
    (16, 'Klima 2', 1, 1, false, '2023-12-06 15:30:00', -1),
    (17, 'Senzorcic', 0,  2, false, '2023-12-06 15:30:00', -1),
    (18, 'Senzor amb', 0, 1, false, '2023-12-06 15:30:00', -1);
    
INSERT INTO consumptionDevice (DeviceId, PowerSupply, PowerConsumption)
VALUES
    (1, 1, 200),
    (2, 0, 0),
    (3, 1, 300),
	(7, 0, 50),
    (8, 1, 75),
    (9, 1, 223),
    (10, 0, 0),
    (11, 0, 0),
    (15, 1, 200),
    (16, 1, 300),
    (17, 0, 0),
    (18, 1, 50);

INSERT INTO airConditioner (DeviceId, MinTemperature, MaxTemperature, Mode)
VALUES
    (3, 10, 30, 'Cooling,Automatic,Ventilation');

INSERT INTO specialModes (DeviceId, StartTime, EndTime, Mode, Temperature, SelectedDays)
VALUES
    (3, '08:00:00', '16:00:00', 'Ventilation', 22, 'Monday,Tuesday,Wednesday');

INSERT INTO airConditioner (DeviceId, MinTemperature, MaxTemperature, Mode)
VALUES
    (16, 10, 30, 'Cooling,Automatic,Ventilation');

INSERT INTO specialModes (DeviceId, StartTime, EndTime, Mode, Temperature, SelectedDays)
VALUES
    (16, '08:00:00', '16:00:00', 'Ventilation', 22, 'Monday,Tuesday,Wednesday');

INSERT INTO washingMachine (DeviceId, Mode)
VALUES
    (1, 'Cotton,Synthetic,Quick,Delicate');

INSERT INTO machineMode (Id, Name, Duration, Temp)
VALUES
    (1, 'Cotton', 120, '60°C'),
    (2, 'Synthetic', 60, '40°C'),
    (3, 'Quick', 30, '30°C'),
    (4, 'Delicate', 90, '30°C');

INSERT INTO machineScheduledMode (Id, DeviceId, StartTime, ModeId)
VALUES
    (1, 1, '2023-12-06 15:30:00', 1),
    (2, 1, '2024-01-28 23:00:00', 1);

INSERT INTO washingMachine (DeviceId, Mode)
VALUES
    (15, 'Cotton,Synthetic,Quick,Delicate');

INSERT INTO sprinkler (DeviceId, IsOn)
VALUES
    (11, false);

INSERT INTO sprinklerSpecialMode (Id, DeviceId, StartTime, EndTime, SelectedDays)
VALUES
    (1, 11, '08:00:00', '10:00:00', 'Thursday,Sunday');

INSERT INTO evCharger (DeviceId, ChargingPower, Connections)
VALUES
    (5, 10, 2),
    (13, 10, 2);
    
INSERT INTO lamp(DeviceId, IsOn, LightningLevel)
VALUES
	(7, false, 0),
    (8, true, 2);

INSERT INTO vehicleGate(DeviceId, IsOpen, Mode)
VALUES
    (9, false, 0),
    (10, false, 1);

INSERT INTO solarPanel (DeviceId, SurfaceArea, Efficiency, NumberOfPanels, IsOn)
VALUES
    (4, 2.3, 30, 3, true),
    (12, 2.3, 30, 3, true);

INSERT INTO homeBattery (DeviceId, Size, CurrentValue)
VALUES
    (6, 10, 2),
    (14, 10, 2);

INSERT INTO licensePlate(DeviceId, PlateNumber)
VALUES
    (9, 'NS-123-45'),
    (9, 'NS-456-22'),
    (9, 'NS-222-34'),
    (10, 'NS-123-45');