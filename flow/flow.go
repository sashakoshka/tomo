package flow

// Flow represents any multi-stage process that relies on callbacks to advance
// to other stages. It allows for such a process to be laid flat instead of
// nested.
type Flow struct {
	// Transition specifies a function to call before moving on to another
	// stage, such as clearing out a container so that new elements may be
	// added to it.
	Transition func ()

	// Stages is a map that pairs stage names with stage functions.
	Stages map [string] func ()
	
	stage string
}

// Switch transitions the flow to a different stage, running the specified
// transition callback first.
func (flow Flow) Switch (stage string) {
	stageCallback := flow.Stages[stage]
	if stageCallback == nil { return }
	if flow.Transition != nil { flow.Transition() }
	flow.stage = stage
	stageCallback()
}

// SwitchFunc returns a function that calles Switch with the specfied stage
// name. This is useful for creating callbacks.
func (flow Flow) SwitchFunc (stage string) (callback func ()) {
	return func () {
		flow.Switch(stage)
	}
}

// Stage returns the name of the current stage the flow is on.
func (flow Flow) Stage () (name string) {
	return flow.stage
}
