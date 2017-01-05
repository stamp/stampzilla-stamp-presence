package main

type State struct {
	Sensor1 bool
	Sensor2 bool
	Sensor3 bool
	Sensor4 bool
	Door    bool
}

func NewState() *State {
	return &State{}
}
