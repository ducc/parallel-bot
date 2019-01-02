package framework

import "sync"

type Registry struct {
	commands map[string]Command
	mutex    *sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
		mutex:    &sync.RWMutex{},
	}
}

func (r *Registry) Register(cmd Command) {
	r.mutex.Lock()
	r.commands[cmd.Name()] = cmd
	r.mutex.Unlock()
}

func (r *Registry) Get(name string) Command {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if cmd, ok := r.commands[name]; ok {
		return cmd
	}
	return nil
}
