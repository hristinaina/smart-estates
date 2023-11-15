CREATE DATABASE smart_home;
use smart_home;

drop table realestate;

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
    UserId INT
);

INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId)
VALUES
  (1, 'Villa B Dorm',  0, '123 Main St', 'Cityville', 150.5, 2, 'path/to/picture1.jpg', 0, 0),
  (2, 'Neka kuca nmp', 1, '456 Oak Ave', 'Townton', 200.75, 3, 'path/to/picture2.jpg', 1, 1),
  (3, 'Joj stvarno nzm', 0, '789 Pine Ln', 'Villageto wn', 30.25, 1, 'path/to/picture3.jpg', 0, 1),
  (4, 'Spavamise', 1, '101 Elm Blvd', 'Hamlet City', 700.0, 2, 'path/to/picture4.jpg', 2, 2);
  
  select * from realestate;

