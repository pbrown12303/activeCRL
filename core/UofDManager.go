package core

import (
	"github.com/pkg/errors"
)

// UofDInitializationFunction is a function that adds core elements to the uOfD during its initialization process. These
// functions are called by the UofDManager after a new UniverseOfDiscourse has been created
type UofDInitializationFunction func(uOFD *UniverseOfDiscourse, hl *Transaction) error

// UofDPostInitializationFunction is an application-specific function that is called after a UofD has been created and all
// UofDInitializationFunctions have been invoked.
type UofDPostInitializationFunction func(uOfD *UniverseOfDiscourse, hl *Transaction) error

// UofDManager manages a universe of discourse and the functions used to initialize it
type UofDManager struct {
	UofD                        *UniverseOfDiscourse
	initializationFunctions     []UofDInitializationFunction
	postInitializationFunctions []UofDPostInitializationFunction
}

// AddInitializationFunction adds a function that will be called during the UniverseOfDiscourse initialization. The function is intended
// to be used by applications to add core concepts to the uOfD before any application data is added
func (mgr *UofDManager) AddInitializationFunction(function UofDInitializationFunction) {
	mgr.initializationFunctions = append(mgr.initializationFunctions, function)
}

// AddPostInitializationFunction adds a function that will be called during the UniverseOfDiscourse initialization. The function is intended
// to be used by applications to perform activities after all core concepts have been added to the uOfD
func (mgr *UofDManager) AddPostInitializationFunction(function UofDInitializationFunction) {
	mgr.initializationFunctions = append(mgr.initializationFunctions, function)
}

// Initialize establishes an initialized UniverseOfDiscourse. It creates the uOfD, calls all of the initialization functions,
// and then calls all of the post-initialization functions.
func (mgr *UofDManager) Initialize() error {
	mgr.UofD = NewUniverseOfDiscourse()
	hl := mgr.UofD.NewTransaction()
	defer hl.ReleaseLocksAndWait()
	for _, function := range mgr.initializationFunctions {
		err := function(mgr.UofD, hl)
		if err != nil {
			errors.Wrap(err, "UofDManager.Initialize failed")
		}
	}
	for _, function := range mgr.postInitializationFunctions {
		err := function(mgr.UofD, hl)
		if err != nil {
			errors.Wrap(err, "UofDManager.Initialize failed")
		}
	}
	return nil
}
