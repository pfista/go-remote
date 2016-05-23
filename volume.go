package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/volume", volume)
	http.ListenAndServe(":8082", nil)
}

type Volume struct {
	OutputVolume string `json:"outputVolume"`
	InputVolume  string `json:"inputVolume"`
	AlertVolume  string `json:"alertVolume"`
	OutputMuted  bool   `json:"outputMuted"`
}

func volume(w http.ResponseWriter, r *http.Request) {
	vol := &Volume{}

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	vol.OutputVolume = r.FormValue("outputVolume")

	if vol.OutputVolume != "" {
		setVolume(vol.OutputVolume)
	}

	vol.update()

	json.NewEncoder(w).Encode(vol)
}

// Updates vol to the system's current volume state
func (vol *Volume) update() {
	result := execApplescript("get volume settings")

	cleaned := strings.Split(result, ",")
	vol.OutputVolume = (strings.Split(cleaned[0], ":"))[1]
	vol.InputVolume = (strings.Split(cleaned[1], ":"))[1]
	vol.AlertVolume = (strings.Split(cleaned[2], ":"))[1]
	muted, err := strconv.ParseBool((strings.Split(cleaned[3], ":")[1]))
	if err != nil {
		fmt.Print(err)
	}

	vol.OutputMuted = muted
}

func setVolume(level string) {
	result := execApplescript("set volume output volume " + level)
	fmt.Println("Result: " + result)
}

func execApplescript(command string) string {
	cmd := exec.Command("osascript", "-e", command)
	output, err := cmd.CombinedOutput()
	prettyOutput := strings.Replace(string(output), "\n", "", -1)
	if err != nil {
		fmt.Println(err)
	}
	return prettyOutput
}
