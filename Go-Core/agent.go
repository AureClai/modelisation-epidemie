package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var agentsCreated uint = 0
var randomSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	Healthy   uint = 0
	Sick      uint = 1
	Recovered uint = 2
	Dead      uint = 3
)

type Agent struct {
	ID       uint
	Radius   float64
	State    uint
	Position Vect2
	Speed    Vect2
	Movable  bool
	WillDie  bool
	TimeSick float64
}

type StartAgentsParameters struct {
	Position Vect2 `json:"position"`
	Speed    Vect2 `json:"speed"`
	State    uint  `json:"state"`
	Movable  bool  `json:"movable"`
}

func (agent *Agent) GetInfo() string {
	return fmt.Sprintf("%v;%v;%v;%v", agent.ID, agent.State, agent.Position.X, agent.Position.Y)
}

type AgentList [](*Agent)

func (alist *AgentList) RemoveAgent(agent *Agent) {
	i := 0
	list := ([](*Agent))(*alist)
	for _, elem := range list {
		if elem == agent {
			break
		}
		i++
	}
	list[len(list)-1], list[i] = list[i], list[len(list)-1]
	*alist = list[:len(list)-1]
}

func (alist *AgentList) RandomChoice() *Agent {
	list := ([](*Agent))(*alist)
	i := rand.Intn(len(list))
	return list[i]
}

func CopyList(list AgentList) AgentList {
	newList := make(AgentList, 0)
	for _, elem := range list {
		newList = append(newList, elem)
	}
	return newList
}

func NewAgent(settings *SimulationSettings) *Agent {
	agentsCreated++
	windowSizeX := settings.WindowSizeX
	windowSizeY := settings.WindowSizeY
	agentRadius := settings.AgentRadius
	agentSpeed := settings.AgentStartSpeed
	pDeath := settings.PDeath

	position := Vect2{
		X: randomSeed.Float64()*(windowSizeX-2*agentRadius) + agentRadius,
		Y: randomSeed.Float64()*(windowSizeY-2*agentRadius) + agentRadius,
	}
	angle := randomSeed.Float64() * 2 * math.Pi
	speed := Vect2{
		X: agentSpeed * math.Cos(angle),
		Y: agentSpeed * math.Sin(angle),
	}
	willDie := false
	if randomSeed.Float64() < pDeath {
		willDie = true
	}
	return &Agent{
		ID:       agentsCreated,
		Radius:   agentRadius,
		State:    Healthy,
		Position: position,
		Speed:    speed,
		Movable:  true,
		WillDie:  willDie,
		TimeSick: 0,
	}
}

func (agent *Agent) GetSick(time float64) {
	agent.State = Sick
	agent.TimeSick = time
}

func (agent *Agent) testContact(other *Agent) bool {
	distance := norm(substract(other.Position, agent.Position))
	minDistance := agent.Radius + other.Radius
	return distance < minDistance
}

func (agent *Agent) testContactWithWall(wall *Wall) (bool, float64) {
	segmentFactor := dot(substract(agent.Position, wall.Start), wall.Direction())
	var d float64
	if segmentFactor <= 0 {
		d = norm(substract(agent.Position, wall.Start))
	} else if segmentFactor >= wall.Length() {
		d = norm(substract(agent.Position, wall.End))
	} else {
		wallStartToAgentCenter := substract(agent.Position, wall.Start)
		d = norm(substract(wallStartToAgentCenter, scalar_times(wall.Direction(), dot(wallStartToAgentCenter, wall.Direction()))))
	}
	return d <= agent.Radius+wall.Radius, segmentFactor
}

func (agent *Agent) updatePos(simu_time float64, settings *SimulationSettings) bool {
	windowSizeX := settings.WindowSizeX
	windowSizeY := settings.WindowSizeY
	timeToRecover := settings.TimeToRecover
	is_dead := false
	if agent.Movable {
		agent.Position = add(agent.Position, scalar_times(agent.Speed, dt))

		// Wall bounce
		var tRollback float64
		if agent.Position.X+agent.Radius >= windowSizeX {
			// Ball bounce right wall
			tRollback = (agent.Radius - windowSizeX + agent.Position.X) / (agent.Speed.X)
			agent.Position.X = agent.Position.X - tRollback*agent.Speed.X
			agent.Speed.X = -agent.Speed.X
			agent.Position.X = agent.Position.X + tRollback*agent.Speed.X
		}
		if agent.Position.X-agent.Radius < 0 {
			// Ball bounce right wall
			tRollback = -(agent.Radius - agent.Position.X) / (agent.Speed.X)
			agent.Position.X = agent.Position.X - tRollback*agent.Speed.X
			agent.Speed.X = -agent.Speed.X
			agent.Position.X = agent.Position.X + tRollback*agent.Speed.X
		}
		if agent.Position.Y+agent.Radius >= windowSizeY {
			// Ball bounce right wall
			tRollback = (agent.Radius - windowSizeY + agent.Position.Y) / (agent.Speed.Y)
			agent.Position.Y = agent.Position.Y - tRollback*agent.Speed.Y
			agent.Speed.Y = -agent.Speed.Y
			agent.Position.Y = agent.Position.Y + tRollback*agent.Speed.Y
		}
		if agent.Position.Y-agent.Radius < 0 {
			// Ball bounce right wall
			tRollback = -(agent.Radius - agent.Position.Y) / (agent.Speed.Y)
			agent.Position.Y = agent.Position.Y - tRollback*agent.Speed.Y
			agent.Speed.Y = -agent.Speed.Y
			agent.Position.Y = agent.Position.Y + tRollback*agent.Speed.Y
		}
	}

	// Change of state
	if agent.State == Sick {
		if simu_time > agent.TimeSick+timeToRecover {
			if agent.WillDie {
				agent.State = Dead
				is_dead = true
				fmt.Printf("%v has died.\n", agent.ID)
			} else {
				agent.State = Recovered
				fmt.Printf("%v has recovered.\n", agent.ID)
			}
		}
	}

	return is_dead
}

func (agent *Agent) bounce(other *Agent, simu_time float64) {
	// Collision detection
	if agent.testContact(other) {
		// Collision
		var tRollback float64
		// Calculation of point of collision
		x1 := agent.Position.X
		x2 := other.Position.X
		xV1 := agent.Speed.X
		xV2 := other.Speed.X
		y1 := agent.Position.Y
		y2 := other.Position.Y
		yV1 := agent.Speed.Y
		yV2 := other.Speed.Y
		a := (xV2-xV1)*(xV2-xV1) + (yV2-yV1)*(yV2-yV1)
		b := (x2-x1)*(xV1-xV2) + (y2-y1)*(yV1-yV2)
		c := (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) - (agent.Radius+other.Radius)*(agent.Radius+other.Radius)
		delta := b*b - 4*a*c
		root1 := (-b - math.Sqrt(delta)) / (2 * a)
		root2 := (-b + math.Sqrt(delta)) / (2 * a)
		if root1 >= 0 {
			tRollback = root1
		} else {
			tRollback = root2
		}
		old_pos1 := substract(agent.Position, scalar_times(agent.Speed, tRollback))
		old_pos2 := substract(other.Position, scalar_times(other.Speed, tRollback))
		// Calculate new speed
		normal := normalize(substract(old_pos2, old_pos1))
		old_v1 := copy(agent.Speed)
		old_v2 := copy(other.Speed)

		// The two are movables
		if agent.Movable && other.Movable {
			agent.Speed = add(scalar_times(normal, dot(old_v2, normal)), substract(old_v1, scalar_times(normal, dot(old_v1, normal))))
			other.Speed = add(scalar_times(normal, dot(old_v1, normal)), substract(old_v2, scalar_times(normal, dot(old_v2, normal))))

			agent.Position = add(old_pos1, scalar_times(agent.Speed, tRollback))
			other.Position = add(old_pos2, scalar_times(other.Speed, tRollback))
		}
		// Only the agent is movable
		if agent.Movable && !other.Movable {
			agent.Speed = add(scalar_times(normal, -dot(old_v1, normal)), substract(old_v1, scalar_times(normal, dot(old_v1, normal))))
			agent.Position = add(old_pos1, scalar_times(agent.Speed, tRollback))
		} else {
			// Only the other is movable
			other.Speed = add(scalar_times(normal, -dot(old_v2, normal)), substract(old_v2, scalar_times(normal, dot(old_v2, normal))))
			other.Position = add(old_pos2, scalar_times(other.Speed, tRollback))
		}

		// Contamination :
		if agent.State == Healthy && other.State == Sick {
			// The agent is contamined by the other
			agent.GetSick(simu_time)
			fmt.Printf("%v has contamined %v\n", other.ID, agent.ID)
		}
		if agent.State == Sick && other.State == Healthy {
			// The other is contamined by the agent
			other.GetSick(simu_time)
			fmt.Printf("%v has contamined %v\n", agent.ID, other.ID)
		}
	}
}

func (agent *Agent) bouceWall(wall *Wall) {
	// The agents bounce with wall is like bouncing with a virtual ball with radius wall.Radius
	// and placed at the nearest point

	// test contact with wall
	hasContact, segmentFactor := agent.testContactWithWall(wall)
	if hasContact {
		// There is a contact with a wall
		if segmentFactor <= 0 {
			// If the nearest point is the start of the wall segment
			segmentFactor = 0
		} else if segmentFactor >= wall.Length() {
			// If the nearest point is the end of the wall segment
			segmentFactor = wall.Length()
		}
		// Assign the virtual center of the virtual ball
		virtualCenter := add(wall.Start, scalar_times(wall.Direction(), segmentFactor))

		// RollBack to avoid stability problems
		// (see equations...)
		x1 := agent.Position.X
		x2 := virtualCenter.X
		xV1 := agent.Speed.X
		xV2 := 0.0
		y1 := agent.Position.Y
		y2 := virtualCenter.Y
		yV1 := agent.Speed.Y
		yV2 := 0.0
		// 2nd degree polynom
		a := (xV2-xV1)*(xV2-xV1) + (yV2-yV1)*(yV2-yV1)
		b := (x2-x1)*(xV1-xV2) + (y2-y1)*(yV1-yV2)
		c := (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) -
			(agent.Radius+wall.Radius)*(agent.Radius+wall.Radius)
		delta := b*b - 4*a*c
		root1 := (-b - math.Sqrt(delta)) / (2 * a)
		root2 := (-b + math.Sqrt(delta)) / (2 * a)
		var tRollback float64
		if root1 >= 0 {
			tRollback = root1
		} else {
			tRollback = root2
		}

		// Roll back the agent position to avoid stability problem
		old_pos1 := substract(agent.Position, scalar_times(agent.Speed, tRollback))

		// Calculate the new speed (cf. Agent collision)
		normal := substract(virtualCenter, agent.Position)
		normal = scalar_times(normal, 1/norm(normal))

		old_v1 := copy(agent.Speed)
		agent.Speed = add(scalar_times(
			normal, -dot(old_v1, normal)), substract(old_v1, scalar_times(normal, dot(old_v1, normal))))

		// Calculate the new position (reprocess rollback)
		agent.Position = add(old_pos1, scalar_times(agent.Speed, tRollback))
	}
}

func place_agent_not_another(agents AgentList, walls WallList, settings *SimulationSettings) *Agent {
	correctly_placed := false
	var iteration uint = 0
	var agent *Agent = nil
	for !correctly_placed {
		iteration++
		agent = NewAgent(settings)
		correctly_placed = true

		// Test placement between agents
		for _, agent2 := range agents {
			if agent.testContact(agent2) {
				correctly_placed = false
				fmt.Println("Not well placed : AGENT")
				break
			}
		}

		// Test placement between agent and walls
		for _, wall := range walls {
			if cond, _ := agent.testContactWithWall(wall); cond {
				correctly_placed = false
				fmt.Println("Not well placed : WALL")
				break
			}
		}
		if iteration > 100 {
			fmt.Println("Max iterations (100) during placement, reduce radius or number or increase size of window")
			os.Exit(1)
		}
	}
	return agent
}

func instanciate_agents(walls WallList, settings *SimulationSettings) AgentList {
	// Start agents from settings
	startAgents := make(AgentList, 0)
	for _, params := range settings.StartAgParam {
		newAgent := NewAgent(settings)
		newAgent.Position = params.Position
		newAgent.Speed = params.Speed
		if params.State == Sick {
			newAgent.GetSick(0)
		} else {
			newAgent.State = params.State
		}
		newAgent.Movable = params.Movable
		startAgents = append(startAgents, newAgent)
	}

	// Instantiate the random agents
	agents := make(AgentList, 0)
	for i := 0; i < int(settings.NbRandomAgents); i++ {
		agents = append(agents, place_agent_not_another(agents, walls, settings))
	}

	// Unmovables
	remaining := CopyList(agents)
	for i := 0; i < int(float64(settings.NbRandomAgents)*settings.FracRandomUnmovable); i++ {
		agent := remaining.RandomChoice()
		agent.Movable = false
		agent.Speed = Vect2{0, 0}
		remaining.RemoveAgent(agent)
	}
	// Start sick
	for i := 0; i < int(settings.NbRandomSicks); i++ {
		agent := remaining.RandomChoice()
		agent.GetSick(0)
		remaining.RemoveAgent(agent)
	}

	// Append not random agents to random agents
	for _, agent := range startAgents {
		agents = append(agents, agent)
	}

	return agents
}

func bounce(agents AgentList, simu_time float64) {
	for i, firstAgent := range agents[:len(agents)-1] {
		for _, secondAgent := range agents[i+1:] {
			firstAgent.bounce(secondAgent, simu_time)
		}
	}
}

func bouceWithWalls(agents AgentList, walls WallList, simu_time float64) {
	for _, wall := range walls {
		for _, agent := range agents {
			agent.bouceWall(wall)
		}
	}
}
