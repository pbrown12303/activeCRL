package core

import (
	"os"

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
func (graphPtr *CrlGraph) AddConceptRecursively(concept Concept, trans *Transaction) error {
	err := graphPtr.addConcept(concept, trans)
	if err != nil {
		return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
	}
	ownedConcepts := concept.GetOwnedConcepts(trans)
	for _, el := range ownedConcepts {
		switch el.GetConceptType() {
		case Refinement:
			ref := el
			abstractConcept := ref.GetAbstractConcept(trans)
			refinedConcept := ref.GetRefinedConcept(trans)
			if abstractConcept != nil && refinedConcept != nil {
				abstractConceptID := abstractConcept.GetConceptID(trans)
				if !graphPtr.gvgraph.IsNode(abstractConceptID) {
					err := graphPtr.addConcept(abstractConcept, trans)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
				refinedConceptID := refinedConcept.GetConceptID(trans)
				if !graphPtr.gvgraph.IsNode(refinedConceptID) {
					err := graphPtr.addConcept(refinedConcept, trans)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
				err := graphPtr.addRefinementEdge(abstractConcept, refinedConcept, trans)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
		case Element, Reference, Literal:
			err := graphPtr.AddConceptRecursively(el, trans)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
			}
			err = graphPtr.addOwnerEdge(concept, el, trans)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.AddConceptRecursively failed")
			}
		}
	}
	return nil
}

func (graphPtr *CrlGraph) addConcept(concept Concept, trans *Transaction) error {
	switch concept.GetConceptType() {
	case Refinement:
		abstractConcept := concept.GetAbstractConcept(trans)
		refinedConcept := concept.GetRefinedConcept(trans)
		if abstractConcept != nil && refinedConcept != nil {
			abstractConceptID := abstractConcept.GetConceptID(trans)
			if !graphPtr.gvgraph.IsNode(abstractConceptID) {
				err := graphPtr.addConcept(abstractConcept, trans)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
			refinedConceptID := refinedConcept.GetConceptID(trans)
			if !graphPtr.gvgraph.IsNode(refinedConceptID) {
				err := graphPtr.addConcept(refinedConcept, trans)
				if err != nil {
					return errors.Wrap(err, "CrlGraph.addConcept failed")
				}
			}
			err := graphPtr.addRefinementEdge(abstractConcept, refinedConcept, trans)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.addConcept failed")
			}
		}
	case Element, Reference, Literal:
		id := "\"" + concept.GetConceptID(trans) + "\""
		label := concept.GetLabel(trans)
		typeName := ConceptTypeToString(concept.GetConceptType())
		if !graphPtr.gvgraph.IsNode(id) {
			nodeAttrs := make(map[string]string)
			nodeAttrs["shape"] = "none"
			referencedAttributeNameExt := ""
			switch concept.GetConceptType() {
			case Reference:
				referencedAttributeNameExt = "<TR><TD>" + concept.GetReferencedAttributeName(trans).String() + "</TD></TR>"
			}
			nodeAttrs["label"] = "<<TABLE><TR><TD>" + typeName + "</TD></TR><TR><TD>" + label + "</TD></TR><TR><TD>" + id + "</TD></TR>" + referencedAttributeNameExt + " </TABLE>>"
			err := graphPtr.gvgraph.AddNode("", id, nodeAttrs)
			if err != nil {
				return errors.Wrap(err, "CrlGraph.addConcept failed")
			}
			// Make sure the owner is displayed
			owner := concept.GetOwningConcept(trans)
			if owner != nil {
				ownerID := owner.GetConceptID(trans)
				if !graphPtr.gvgraph.IsNode(ownerID) {
					err := graphPtr.addConcept(owner, trans)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
					err = graphPtr.addOwnerEdge(owner, concept, trans)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
			}
			switch concept.GetConceptType() {
			case Reference:
				referencedConceptID := concept.GetReferencedConceptID(trans)
				if referencedConceptID != "" {
					referencedConcept := concept.GetReferencedConcept(trans)
					if !graphPtr.gvgraph.IsNode(referencedConceptID) {
						err := graphPtr.addConcept(referencedConcept, trans)
						if err != nil {
							return errors.Wrap(err, "CrlGraph.addConcept failed")
						}
					}
					err := graphPtr.addReferencedElementEdge(concept, referencedConcept, trans)
					if err != nil {
						return errors.Wrap(err, "CrlGraph.addConcept failed")
					}
				}
			}
		}
	}
	return nil
}

func (graphPtr *CrlGraph) addOwnerEdge(parent Concept, child Concept, trans *Transaction) error {
	parentID := "\"" + parent.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(parentID) {
		return errors.New("CrlGraph.addOwnerEdge called with parent node not present")
	}
	childID := "\"" + child.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(childID) {
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

func (graphPtr *CrlGraph) addReferencedElementEdge(reference Concept, referencedElement Concept, trans *Transaction) error {
	referenceID := "\"" + reference.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(referenceID) {
		return errors.New("CrlGraph.addReferencedElementEdge called with reference node not present")
	}
	referencedElementID := "\"" + referencedElement.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(referencedElementID) {
		return nil
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

func (graphPtr *CrlGraph) addRefinementEdge(abstractConcept Concept, refinedConcept Concept, trans *Transaction) error {
	abstractConceptID := "\"" + abstractConcept.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(abstractConceptID) {
		return errors.New("CrlGraph.addRefinementEdge called with abstractConcept node not present")
	}
	refinedConceptID := "\"" + refinedConcept.GetConceptID(trans) + "\""
	if !graphPtr.gvgraph.IsNode(refinedConceptID) {
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
