package main

// IoTEvent defines AWS IoT event message.
type IoTEvent struct {
    Sensor   string `json:"sensor"`
}

func handler(event IoTEvent) error {
	fmt.Print(event.sensor + "\n")
	return nil
}

