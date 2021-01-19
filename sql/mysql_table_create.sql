/* Simple script to create the MySQL DB */

CREATE TABLE Device (
	Id INT NOT NULL AUTO_INCREMENT,
	Name NVARCHAR(100) NOT NULL,
	Description NVARCHAR(4086) NULL,
	DeviceEUI NVARCHAR(64) NULL,
	DevAddr NVARCHAR(32) NULL,
	CONSTRAINT PK_Device PRIMARY KEY (Id),
	CONSTRAINT UX_Device_Name UNIQUE KEY (Name)
);

CREATE TABLE Location (
	DeviceId INT NOT NULL,
	LoggedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Location POINT NOT NULL,
	Id BIGINT NOT NULL AUTO_INCREMENT,
	CONSTRAINT PK_Location PRIMARY KEY (Id),
	CONSTRAINT UX_Location_LoggedOnDeviceId UNIQUE KEY (LoggedOn, DeviceId), /* Device can't be at different places at the same time */
	CONSTRAINT FK_Location_Device FOREIGN KEY (DeviceId) REFERENCES Device(Id) ON DELETE NO ACTION ON UPDATE CASCADE,
	INDEX IX_Location_LoggedOn (LoggedOn)
);

CREATE TABLE BatteryStatus (
	Id BININT NOT NULL AUTO_INCREMENT,
	DeviceId INT NOT NULL,
	LoggedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	RawValue INT NOT NULL,
	BatteryLevel TINYINT NULL,
	CONSTRAINT PK_BatteryStatus PRIMARY KEY (Id),
	CONSTRAINT UX_BatteryStatus_LoggedOnDeviceId UNIQUE KEY (LoggedOn, DeviceId), /* A device can't reports it's battery value multiple times at the same time */
	CONSTRAINT FK_BatteryStatus_Device FOREIGN KEY (DeviceId) REFERENCES (Device(Id ON DELETE NO ACTION ON UPDATE CASCADE,
	INDEX IX_BatteryStatus_LoggedOn (LoggedOn)
);
