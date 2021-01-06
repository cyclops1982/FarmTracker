package loramsgs

type SodaqUniversalTracker struct {
	Unixtime uint32
	RawVoltage uint8
	BoardTemperature int8
	Latitude int32
	Longitude int32
	Altitude uint16
	Speed uint8
	Course uint8
	Satelites uint8
}


