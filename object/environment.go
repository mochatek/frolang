package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

// Adds value to supplied identifier in the environment
func (environment *Environment) Set(name string, object Object) Object {
	environment.store[name] = object
	return object
}

// Updates value of supplied identifier in the environment in which it was declared
func (environment *Environment) Update(name string, object Object) Object {
	for env := environment; env != nil; env = env.outer {
		if _, ok := env.store[name]; ok {
			env.store[name] = object
			return object
		}
	}
	environment.store[name] = object
	return object
}

// Retrieves value of supplied identifier from environment
// If identifier is not present in current environment, look up in outer environment (Scope chain)
func (environment *Environment) Get(name string) (Object, bool) {
	object, ok := environment.store[name]
	if !ok && environment.outer != nil {
		return environment.outer.Get(name)
	}
	return object, ok
}

// Constructor function for global environment
// *outer points to null as this is the outermost environment
func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

// Constructor function for local environment
// *outer points to the outer environment thereby creating the scope chain
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
