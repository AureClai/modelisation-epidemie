package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var simulationDuration float64 = 30 // seconds
var dt float64 = 1.0 / 60.0

type Simulation struct {
	Agents   AgentList
	Walls    WallList
	Infos    []string
	Settings *SimulationSettings
}

type SimulationSettings struct {
	Walls               WallList                   `json:"walls"`
	WindowSizeX         float64                    `json:"window_size_x"`
	WindowSizeY         float64                    `json:"window_size_y"`
	Duration            float64                    `json:"duration"`
	Dt                  float64                    `json:"dt"`
	TimeToRecover       float64                    `json:"time_to_recover"`
	FracRandomUnmovable float64                    `json:"frac_unmovable"`
	NbRandomAgents      uint                       `json:"nb_random_agents"`
	NbRandomSicks       uint                       `json:"nb_random_sick"`
	PDeath              float64                    `json:"death_proportion"`
	AgentStartSpeed     float64                    `json:"agents_start_speed"`
	AgentRadius         float64                    `json:"agents_radius"`
	StartAgParam        [](*StartAgentsParameters) `json="start_agents"`
}

func NewSimulation(settings *SimulationSettings) *Simulation {
	walls := settings.Walls
	agents := instanciate_agents(walls, settings)
	return &Simulation{
		Agents:   agents,
		Walls:    walls,
		Infos:    make([]string, 0),
		Settings: settings,
	}
}

func (sim *Simulation) Run() {
	simu_time := 0.0
	aliveAgents := CopyList(sim.Agents)
	for simu_time < simulationDuration {
		//fmt.Println(simu_time)
		simu_time += dt
		// Collision agent/walls
		bouceWithWalls(aliveAgents, sim.Walls, simu_time)

		// Collision between agents
		bounce(aliveAgents, simu_time)
		hasDied := make(AgentList, 0)
		// Move alive Agents
		for _, agent := range aliveAgents {
			isDead := agent.updatePos(simu_time, sim.Settings)
			if isDead {
				hasDied = append(hasDied, agent)
			}
		}
		// Update dead agent
		for _, deadAgent := range hasDied {
			aliveAgents.RemoveAgent(deadAgent)
			fmt.Printf("%v has been removed from alives \n", deadAgent.ID)
		}

		//Get  Info
		for _, agent := range sim.Agents {
			sim.Infos = append(sim.Infos, fmt.Sprintf("%v;%v", simu_time, agent.GetInfo()))
		}
	}

}

func (sim *Simulation) SaveResults() {
	dirName := "Results_" + time.Now().Format("20060201_150405")
	os.Mkdir("."+string(filepath.Separator)+dirName, 0777)
	// Positions
	file, err := os.Create("." + string(filepath.Separator) + dirName + string(filepath.Separator) + "positions.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"time;id;state;x;y"})
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range sim.Infos {
		err = writer.Write([]string{value})
		if err != nil {
			fmt.Println(err)
		}
	}

	// Settings
	filepath := "." + string(filepath.Separator) + dirName + string(filepath.Separator) + "settings.json"
	jsonFile, _ := json.MarshalIndent(sim.Settings, "", " ")
	_ = ioutil.WriteFile(filepath, jsonFile, 0644)
}
