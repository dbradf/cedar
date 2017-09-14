package mock

import (
	"github.com/mongodb/amboy"
	"github.com/mongodb/amboy/dependency"
	"github.com/pkg/errors"
	"github.com/mongodb/anser/db"
	"github.com/mongodb/anser/model"
)

type Environment struct {
	Queue             amboy.Queue
	Session           *Session
	Network           *DependencyNetwork
	MetaNS            model.Namespace
	MigrationRegistry map[string]db.MigrationOperation
	ProcessorRegistry map[string]db.Processor
	DependecyManagers map[string]*DependencyManager
	IsSetup           bool
	SetupError        error
	SessionError      error
	QueueError        error
	NetworkError      error
}

func NewEnvironment() *Environment {
	return &Environment{
		Session:           NewSession(),
		Network:           NewDependencyNetwork(),
		DependecyManagers: make(map[string]*DependencyManager),
		MigrationRegistry: make(map[string]db.MigrationOperation),
		ProcessorRegistry: make(map[string]db.Processor),
	}
}

func (e *Environment) Setup(q amboy.Queue, uri string) error {
	e.Queue = q
	e.Session.URI = uri
	e.IsSetup = true

	return e.SetupError
}

func (e *Environment) GetSession() (db.Session, error) {
	if e.SessionError != nil {
		return nil, e.SessionError
	}
	return e.Session, nil
}

func (e *Environment) GetQueue() (amboy.Queue, error) {
	if e.QueueError != nil {
		return nil, e.QueueError
	}

	return e.Queue, nil
}

func (e *Environment) GetDependencyNetwork() (model.DependencyNetworker, error) {
	if e.NetworkError != nil {
		return nil, e.NetworkError
	}

	return e.Network, nil
}

func (e *Environment) RegisterManualMigrationOperation(name string, op db.MigrationOperation) error {
	if _, ok := e.MigrationRegistry[name]; ok {
		return errors.Errorf("migration operation %s already exists", name)
	}

	e.MigrationRegistry[name] = op
	return nil
}

func (e *Environment) GetManualMigrationOperation(name string) (db.MigrationOperation, bool) {
	op, ok := e.MigrationRegistry[name]
	return op, ok
}

func (e *Environment) RegisterDocumentProcessor(name string, docp db.Processor) error {
	if _, ok := e.ProcessorRegistry[name]; ok {
		return errors.Errorf("document processor named %s already registered", name)
	}

	e.ProcessorRegistry[name] = docp
	return nil
}

func (e *Environment) GetDocumentProcessor(name string) (db.Processor, bool) {
	docp, ok := e.ProcessorRegistry[name]
	return docp, ok
}

func (e *Environment) MetadataNamespace() model.Namespace { return e.MetaNS }

func (e *Environment) NewDependencyManager(n string, q map[string]interface{}, ns model.Namespace) dependency.Manager {
	e.DependecyManagers[n] = &DependencyManager{
		Name:     n,
		Query:    q,
		NS:       ns,
		JobEdges: dependency.NewJobEdges(),
	}

	return e.DependecyManagers[n]
}