package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	current_value := 0.0
	for {
		url := fmt.Sprintf("http://localhost:31270/set/CurrentDrivableActor/IndependentBrake_F.InputValue?Value=%f", current_value)
		req, _ := http.NewRequest(http.MethodPatch, url, nil)
		req.Header.Set("DTGCommKey", "RIgk+NV9JtMaoKko2sjidHERmaDI2S0AjYUn711mS9A=")
		client.Do(req)

		time.Sleep(10 * time.Millisecond)

		current_value += 0.05
		if current_value >= 1.0 {
			break
		}
	}
}
