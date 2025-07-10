# Go Backend for Parking Spot Finder

## Setup
1. Install Go and MySQL
2. Create database and table:

CREATE DATABASE parkingdb;
USE parkingdb;

CREATE TABLE parking_spots (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255),
  lat DOUBLE,
  lng DOUBLE,
  available_spots INT
);

3. Configure `.env` file
4. Run:

```
go mod tidy
go run main.go
```
