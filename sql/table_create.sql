/* Simple script to create the MySQL DB */

CREATE TABLE Location(
	DeviceId	INT NOT NULL,
	LoggedOn	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Location	POINT NOT NULL,
	Id		BIGINT NOT NULL AUTO_INCREMENT,
	CONSTRAINT PK_Location PRIMARY KEY (Id),
	CONSTRAINT UX_Location_LoggedOnDeviceId UNIQUE KEY (LoggedOn, DeviceId), /* Device can't be at different places at the same time */
	INDEX IX_Location_When (LoggedOn) /* Make it fast to order by When. */
)


