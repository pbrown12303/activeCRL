package crleditor

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pkg/errors"
)

type workspaceFile struct {
	filename      string
	File          *os.File
	LoadedVersion int
	Info          os.FileInfo
	Domain        core.Concept
}

// CrlWorkspaceManager manages Crl Workspaces
type CrlWorkspaceManager struct {
	editor         *Editor
	workspaceFiles map[string]*workspaceFile
}

// NewCrlWorkspaceManager returns a configured CrlWorkspaceManager
func NewCrlWorkspaceManager(editor *Editor) *CrlWorkspaceManager {
	mgr := &CrlWorkspaceManager{}
	mgr.editor = editor
	return mgr
}

// Initialize initializes or re-initializes the workspace editor
func (mgr *CrlWorkspaceManager) Initialize() {
	mgr.workspaceFiles = make(map[string]*workspaceFile)

}

// ClearWorkspace deletes all of the files in the workspace that correspond to uOfD root elements
// and removes the corresponding entry in workspaceFiles
func (mgr *CrlWorkspaceManager) ClearWorkspace(workspacePath string, trans *core.Transaction) error {
	var err error
	rootElements := mgr.editor.uOfDManager.UofD.GetRootElements(trans)
	for id, wf := range mgr.workspaceFiles {
		if rootElements[id] == nil {
			err = mgr.deleteFile(wf)
			if err != nil {
				return errors.Wrap(err, "CrlEditor.ClearWorkspace failed")
			}
			delete(mgr.workspaceFiles, id)
		}
	}
	return nil
}

// CloseWorkspace saves and closes all workspace files
func (mgr *CrlWorkspaceManager) CloseWorkspace(trans *core.Transaction) error {
	err := mgr.SaveWorkspace(trans)
	if err != nil {
		return errors.Wrap(err, "CrlWorkspaceManager.CloseWorkspace failed")
	}
	for _, wsf := range mgr.workspaceFiles {
		err = wsf.File.Close()
		if err != nil {
			return errors.Wrap(err, "CrlWorkspaceManager.CloseWorkspace failed")
		}
	}
	return nil
}

// deleteFile deletes the file from the os
func (mgr *CrlWorkspaceManager) deleteFile(wf *workspaceFile) error {
	err := wf.File.Close()
	if err != nil {
		return errors.Wrap(err, "CrlEditor.delete file failed")
	}
	err = os.Remove(wf.filename)
	if err != nil {
		return errors.Wrap(err, "CrlEditor.delete file failed")
	}
	return nil
}

func (mgr *CrlWorkspaceManager) generateFilename(el core.Concept, trans *core.Transaction) string {
	return mgr.editor.userPreferences.WorkspacePath + "/" + el.GetLabel(trans) + "--" + el.GetConceptID(trans) + ".acrl"
}

// GetUofD returns the current UniverseOfDiscourse
func (mgr *CrlWorkspaceManager) GetUofD() *core.UniverseOfDiscourse {
	return mgr.editor.GetUofD()
}

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
func (mgr *CrlWorkspaceManager) newFile(el core.Concept, trans *core.Transaction) (*workspaceFile, error) {
	if mgr.editor.userPreferences.WorkspacePath == "" {
		return nil, errors.New("CrlBrowserEditor.NewFile called with no settings.WorkspacePath defined")
	}
	filename := mgr.generateFilename(el, trans)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	fileInfo, err2 := os.Stat(filename)
	if err2 != nil {
		return nil, err2
	}
	var wf workspaceFile
	wf.filename = filename
	wf.Domain = el
	wf.File = file
	wf.LoadedVersion = el.GetVersion(trans)
	wf.Info = fileInfo
	return &wf, nil
}

// openFile opens the file and returns a workspaceFile struct
func (mgr *CrlWorkspaceManager) openFile(fileInfo os.FileInfo, trans *core.Transaction) (*workspaceFile, error) {
	writable := (fileInfo.Mode().Perm() & 0200) > 0
	mode := os.O_RDONLY
	if writable {
		mode = os.O_RDWR
	}
	filename := mgr.editor.userPreferences.WorkspacePath + "/" + fileInfo.Name()
	file, err := os.OpenFile(filename, mode, fileInfo.Mode())
	if err != nil {
		return nil, err
	}
	fileContent := make([]byte, fileInfo.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}
	element, err2 := mgr.GetUofD().RecoverDomain(fileContent, trans)
	if err2 != nil {
		return nil, err2
	}
	if !writable {
		element.SetReadOnlyRecursively(true, trans)
	}
	var wf workspaceFile
	wf.filename = filename
	wf.Domain = element
	wf.Info = fileInfo
	wf.LoadedVersion = element.GetVersion(trans)
	wf.File = file
	return &wf, nil
}

// LoadSettings loads the settings saved in the workspace
func (mgr *CrlWorkspaceManager) LoadSettings(trans *core.Transaction) error {
	path := mgr.editor.getSettingsPath()
	_, err := os.Stat(path)
	if err != nil {
		// it is OK to not find the file
		mgr.editor.settings = &Settings{}
		return nil
	}
	fileSettings, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err
	}
	err = json.Unmarshal(fileSettings, mgr.editor.settings)
	jsonOpenDiagrams, _ := json.Marshal(mgr.editor.settings.OpenDiagrams)
	mgr.editor.transientDisplayedDiagrams.SetLiteralValue(string(jsonOpenDiagrams), trans)
	if err != nil {
		return err
	}
	return nil
}

// LoadUserPreferences loads the user preferences saved in the user's home directory
func (mgr *CrlWorkspaceManager) LoadUserPreferences(workspaceArg string) error {
	path := mgr.editor.getUserPreferencesPath()
	_, err := os.Stat(path)
	if err != nil {
		// it is OK to not find the file
		mgr.editor.userPreferences.WorkspacePath = workspaceArg
		return nil
	}
	fileSettings, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err
	}
	err = json.Unmarshal(fileSettings, mgr.editor.userPreferences)
	if err != nil {
		return err
	}
	return nil
}

// LoadWorkspace loads the workspace currently designated by the userPreferences.WorkspacePath. If the path is empty, it is a no-op.
func (mgr *CrlWorkspaceManager) LoadWorkspace(trans *core.Transaction) error {
	files, err := ioutil.ReadDir(mgr.editor.userPreferences.WorkspacePath)
	if err != nil {
		return errors.Wrap(err, "CrlWorkspaceManager.LoadWorkspace failed")
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".acrl") {
			workspaceFile, err := mgr.openFile(f, trans)
			if err != nil {
				return errors.Wrap(err, "CrlWorkspaceManager.LoadWorkspace failed loading "+f.Name())
			}
			mgr.workspaceFiles[workspaceFile.Domain.GetConceptID(trans)] = workspaceFile
		}
	}
	mgr.LoadSettings(trans)
	mgr.editor.SelectElementUsingIDString(mgr.editor.settings.Selection, trans)
	mgr.editor.diagramManager.DisplayDiagram(mgr.editor.settings.CurrentDiagram, trans)
	return nil
}

// saveFile saves the file and updates the fileInfo
func (mgr *CrlWorkspaceManager) saveFile(wf *workspaceFile, trans *core.Transaction) error {
	trans.ReadLockElement(wf.Domain)
	if wf.File == nil {
		return errors.New("CrlBrowserEditor.SaveFile called with nil file")
	}
	byteArray, err := mgr.GetUofD().MarshalDomain(wf.Domain, trans)
	if err != nil {
		return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
	}
	var length int
	length, err = wf.File.WriteAt(byteArray, 0)
	if err != nil {
		return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
	}
	err = wf.File.Truncate(int64(length))
	if err != nil {
		return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
	}
	err = wf.File.Sync()
	if err != nil {
		return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
	}
	oldFilename := wf.filename
	newFilename := mgr.generateFilename(wf.Domain, trans)
	if oldFilename != newFilename {
		err = wf.File.Close()
		if err != nil {
			return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
		}
		err = os.Rename(oldFilename, newFilename)
		if err != nil {
			return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
		}
		wf.filename = newFilename
		wf.File, err = os.OpenFile(newFilename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
		}
		wf.Info, err = os.Stat(newFilename)
		if err != nil {
			return errors.Wrap(err, "CrlBrowserEditor.saveFile failed")
		}
	}
	return nil
}

// SaveWorkspace saves all top-level concepts whose versions are different than the last retrieved version.
func (mgr *CrlWorkspaceManager) SaveWorkspace(trans *core.Transaction) error {
	rootElements := mgr.editor.uOfDManager.UofD.GetRootElements(trans)
	var err error
	for id, el := range rootElements {
		noSaveDomains := mgr.editor.getNoSaveDomains(trans)
		if !el.GetIsCore(trans) && noSaveDomains[el.GetConceptID(trans)] == nil {
			workspaceFile := mgr.workspaceFiles[id]
			if workspaceFile != nil {
				err = mgr.saveFile(workspaceFile, trans)
				if err != nil {
					return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
				}
			} else {
				workspaceFile, err = mgr.newFile(el, trans)
				if err != nil {
					return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
				}
				mgr.workspaceFiles[id] = workspaceFile
				err = mgr.saveFile(workspaceFile, trans)
				if err != nil {
					return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
				}
			}
		}
	}
	for id, wf := range mgr.workspaceFiles {
		if rootElements[id] == nil {
			err = mgr.deleteFile(wf)
			if err != nil {
				return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
			}
			delete(mgr.workspaceFiles, id)
		}
	}
	return nil
}
