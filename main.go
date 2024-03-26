package main

func main() {
	/*
		devices := createDevices(10) // initialize 10 devices including secret key sharing and initialized communication channels

		genericDerivation := GenericDerivation{devices: devices}
		benchDerivation(genericDerivation)

		tvrfDerivation := TVRFDerivation{devices: devices}
		benchDerivation(tvrfDerivation)
	*/
	InitDevices(4, 6)
}
