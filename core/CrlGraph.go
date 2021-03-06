package core

import (
	"os"
	"reflect"

	"github.com/awalterschulze/gographviz"
	"github.com/pkg/errors"
)

// CrlGraph is a graphviz representation of a specified set of CRL data structures.
// At present, Refinements cannot be the referencedConcept of a Reference.
type CrlGraph struct {
	gvgraph *gographviz.Graph
}

// NewCrlGraph returns an initialized CrlGraph
func NewCrlGraph(graphName string) *CrlGraph {
	graph := &CrlGraph{}
	graph.gvgraph = gographviz.NewGraph()
	graph.gvgraph.SetDir(true)
	graph.gvgraph.SetStrict(true)
	graph.gvgraph.SetName(graphName)
	return graph
}

// AddConceptRecursively will add the given concept and all its child descendants to the graph.
// it will also add any referenced concepts, but not recursively. Existing concepts will not be duplicated.
func (graphPtr *CrlGraph) AddConceptRecursively(concept Element, hl *HeldLocks) error {
	err := graphPtr.addConcept(concept, hl)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
	}
	ownedConcepts := concept.GetOwnedConcepts(hl)
	for _, el := range ownedConcepts {
		switch el.(type) {
		case Refinement:
			ref := el.(Refinement)
			abstractConcept := ref.GetAbstractConcept(hl)
			refinedConcept := ref.GetRefinedConcept(hl)
			if abstractConcept != nil && refinedConcept != nil {
				abstractConceptID := abstractConcept.GetConceptID(hl)
				if graphPtr.gvgraph.IsNode(abstractConceptID) == false {
					err := graphPtr.addConcept(abstractConcept, hl)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
				refinedConceptID := refinedConcept.GetConceptID(hl)
				if graphPtr.gvgraph.IsNode(refinedConceptID) == false {
					err := graphPtr.addConcept(refinedConcept, hl)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
				err := graphPtr.addRefinementEdge(abstractConcept, refinedConcept, hl)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
		case Element, Literal, Reference:
			err := graphPtr.AddConceptRecursively(el, hl)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
			}
			err = graphPtr.addOwnerEdge(concept, el, hl)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
			}
		}
	}
	return nil
}

func (graphPtr *CrlGraph) addConcept(concept Element, hl *HeldLocks) error {
	switch concept.(type) {
	case Refinement:
		ref := concept.(Refinement)
		abstractConcept := ref.GetAbstractConcept(hl)
		refinedConcept := ref.GetRefinedConcept(hl)
		if abstractConcept != nil && refinedConcept != nil {
			abstractConceptID := abstractConcept.GetConceptID(hl)
			if graphPtr.gvgraph.IsNode(abstractConceptID) == false {
				err := graphPtr.addConcept(abstractConcept, hl)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
			refinedConceptID := refinedConcept.GetConceptID(hl)
			if graphPtr.gvgraph.IsNode(refinedConceptID) == false {
				err := graphPtr.addConcept(refinedConcept, hl)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
			err := graphPtr.addRefinementEdge(abstractConcept, refinedConcept, hl)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.addConcept failed")
			}
		}
	case Element, Literal, Reference:
		id := "\"" + concept.GetConceptID(hl) + "\""
		label := concept.GetLabel(hl)
		typeName := reflect.TypeOf(concept).String()
		if graphPtr.gvgraph.IsNode(id) != true {
			nodeAttrs := make(map[string]string)
			nodeAttrs["shape"] = "none"
			nodeAttrs["label"] = "<<TABLE><TR><TD>" + typeName + "</TD></TR><TR><TD>" + label + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
			err := graphPtr.gvgraph.AddNode("", id, nodeAttrs)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.addConcept failed")
			}
			// Make sure the owner is displayed
			owner := concept.GetOwningConcept(hl)
			if owner != nil {
				ownerID := owner.GetConceptID(hl)
				if graphPtr.gvgraph.IsNode(ownerID) == false {
					err := graphPtr.addConcept(owner, hl)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
					err = graphPtr.addOwnerEdge(owner, concept, hl)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
			}
			switch concept.(type) {
			case Reference:
				ref := concept.(Reference)
				referencedConceptID := ref.GetReferencedConceptID(hl)
				if referencedConceptID != "" {
					referencedConcept := ref.GetReferencedConcept(hl)
					if graphPtr.gvgraph.IsNode(referencedConceptID) == false {
						err := graphPtr.addConcept(referencedConcept, hl)
						if err != nil {
							return errors.Wrap(err, "CrlGraph.addConcept failed")
						}
					}
					err := graphPtr.addReferencedElementEdge(ref, referencedConcept, hl)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
				// case Refinement:
				// 	ref := concept.(Refinement)
				// 	abstractConcept := ref.GetAbstractConcept(hl)
				// 	if abstractConcept != nil {
				// 		abstractConceptID := abstractConcept.GetConceptID(hl)
				// 		if graphPtr.gvgraph.IsNode(abstractConceptID) == false {
				// 			err := graphPtr.addConcept(abstractConcept, hl)
				// 			if err != nil {
				// 				return errors.Wrap(err, "CrlGraph.addConcept failed")
				// 			}
				// 		}
				// 		err := graphPtr.addAbstractEdge(ref, ref.GetAbstractConcept(hl), hl)
				// 		if err != nil {
				// 			return errors.Wrap(err, "CrlGraph.addConcept failed")
				// 		}
				// 	}
				// 	refinedConcept := ref.GetRefinedConcept(hl)
				// 	if refinedConcept != nil {
				// 		refinedConceptID := refinedConcept.GetConceptID(hl)
				// 		if graphPtr.gvgraph.IsNode(refinedConceptID) == false {
				// 			err := graphPtr.addConcept(refinedConcept, hl)
				// 			if err != nil {
				// 				return errors.Wrap(err, "CrlGraph.addConcept failed")
				// 			}
				// 		}
				// 		err := graphPtr.addRefinedEdge(ref, ref.GetRefinedConcept(hl), hl)
				// 		if err != nil {
				// 			return errors.Wrap(err, "CrlGraph.addConcept failed")
				// 		}
				// 	}
			}
		}
	}
	return nil
}

// func (graphPtr *CrlGraph) addAbstractEdge(refinement Element, abstractElement Element, hl *HeldLocks) error {
// 	refinementID := "\"" + refinement.GetConceptID(hl) + "\""
// 	abstractElementID := "\"" + abstractElement.GetConceptID(hl) + "\""
// 	refEdgeAttrs := make(map[string]string)
// 	refEdgeAttrs["arrowhead"] = "none"
// 	refEdgeAttrs["arrowtail"] = "onormal"
// 	refEdgeAttrs["dir"] = "both"
// 	refEdgeAttrs["constraint"] = "false"
// 	refEdgeAttrs["color"] = "green"
// 	err := graphPtr.gvgraph.AddEdge(refinementID, abstractElementID, true, refEdgeAttrs)
// 	if err != nil {
// 		return errors.Wrap(err, "CrlGraph.addReferencedElementEdge failed")
// 	}
// 	return nil
// }

func (graphPtr *CrlGraph) addOwnerEdge(parent Element, child Element, hl *HeldLocks) error {
	parentID := "\"" + parent.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(parentID) == false {
		return errors.New("CrlGraph.addOwnerEdge called with parent node not present")
	}
	childID := "\"" + child.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(childID) == false {
		return errors.New("CrlGraph.addOwnerEdge called with child node not present")
	}
	ownerEdgeAttrs := make(map[string]string)
	ownerEdgeAttrs["arrowhead"] = "none"
	ownerEdgeAttrs["arrowtail"] = "diamond"
	ownerEdgeAttrs["dir"] = "both"
	// ownerEdgeAttrs["weight"] = "100"
	err := graphPtr.gvgraph.AddEdge(parentID, childID, true, ownerEdgeAttrs)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.addOwnerEdge failed")
	}
	return nil
}

func (graphPtr *CrlGraph) addReferencedElementEdge(reference Element, referencedElement Element, hl *HeldLocks) error {
	referenceID := "\"" + reference.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(referenceID) == false {
		return errors.New("CrlGraph.addReferencedElementEdge called with reference node not present")
	}
	referencedElementID := "\"" + referencedElement.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(referencedElementID) == false {
		return errors.New("CrlGraph.addReferencedElementEdge called with referencedElement node not present")
	}
	refEdgeAttrs := make(map[string]string)
	refEdgeAttrs["arrowhead"] = "open"
	refEdgeAttrs["arrowtail"] = "none"
	refEdgeAttrs["dir"] = "both"
	// refEdgeAttrs["weight"] = "10"
	refEdgeAttrs["constraint"] = "false"
	refEdgeAttrs["color"] = "red"
	err := graphPtr.gvgraph.AddEdge(referenceID, referencedElementID, true, refEdgeAttrs)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.addReferencedElementEdge failed")
	}
	return nil
}

func (graphPtr *CrlGraph) addRefinementEdge(abstractConcept Element, refinedConcept Element, hl *HeldLocks) error {
	abstractConceptID := "\"" + abstractConcept.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(abstractConceptID) == false {
		return errors.New("CrlGraph.addRefinementEdge called with abstractConcept node not present")
	}
	refinedConceptID := "\"" + refinedConcept.GetConceptID(hl) + "\""
	if graphPtr.gvgraph.IsNode(refinedConceptID) == false {
		return errors.New("CrlGraph.addRefinementEdge called with refinedConcept node not present")
	}
	refEdgeAttrs := make(map[string]string)
	refEdgeAttrs["arrowhead"] = "none"
	refEdgeAttrs["arrowtail"] = "onormal"
	refEdgeAttrs["dir"] = "both"
	// refEdgeAttrs["weight"] = "1"
	refEdgeAttrs["constraint"] = "false"
	refEdgeAttrs["color"] = "turquoise"
	err := graphPtr.gvgraph.AddEdge(abstractConceptID, refinedConceptID, true, refEdgeAttrs)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.addRefinementEdge failed")
	}
	return nil
}

// ExportDOT writes a file containing the DOT representation of the graph
func (graphPtr *CrlGraph) ExportDOT(pathname string, filename string) error {
	file, err := graphPtr.newFile(pathname, filename)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.ExportDOT failed")
	}
	graphString := graphPtr.gvgraph.String()
	graphBytes := []byte(graphString)
	_, err2 := file.Write(graphBytes)
	if err2 != nil {
		return errors.Wrap(err, "CrlGraph.ExportDOT failed")
	}
	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "CrlGraph.ExportDOT failed")
	}
	return nil
}

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
func (graphPtr *CrlGraph) newFile(path string, filename string) (*os.File, error) {
	fullPath := path + "/" + filename + ".dot"
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "CrlGraph.newFile failed")
	}
	return file, nil
}
