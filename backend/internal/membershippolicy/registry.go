package membershippolicy

import "fmt"

type Registry struct {
	policies map[string]Policy
}

func NewRegistry() *Registry {
	general := NewGeneralPolicy()
	ssbm := NewSSBMPolicy()

	return &Registry{
		policies: map[string]Policy{
			general.ProgramName(): general,
			ssbm.ProgramName():    ssbm,
		},
	}
}

func (r *Registry) Get(programName string) (Policy, error) {
	policy, exists := r.policies[programName]
	if !exists {
		return nil, fmt.Errorf(
			"no membership policy registered for program %q",
			programName,
		)
	}

	return policy, nil
}
