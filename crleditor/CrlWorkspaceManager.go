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
	Domain        core.Element
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
func (mgr *CrlWorkspaceManager) ClearWorkspace(workspacePath string, hl *core.Transaction) error {
	var err error
	rootElements := mgr.editor.uOfDManager.UofD.GetRootElements(hl)
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
func (mgr *CrlWorkspaceManager) CloseWorkspace(hl *core.Transaction) error {
	err := mgr.SaveWorkspace(hl)
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

func (mgr *CrlWorkspaceManager) generateFilename(el core.Element, hl *core.Transaction) string {
	return mgr.editor.userPreferences.WorkspacePath + "/" + el.GetLabel(hl) + "--" + el.GetConceptID(hl) + ".acrl"
}

// GetUofD returns the current UniverseOfDiscourse
func (mgr *CrlWorkspaceManager) GetUofD() *core.UniverseOfDiscourse {
	return mgr.editor.GetUofD()
}

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
func (mgr *CrlWorkspaceManager) newFile(el core.Element, hl *core.Transaction) (*workspaceFile, error) {
	if mgr.editor.userPreferences.WorkspacePath == "" {
		return nil, errors.New("CrlBrowserEditor.NewFile called with no settings.WorkspacePath defined")
	}
	filename := mgr.generateFilename(el, hl)
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
	wf.LoadedVersion = el.GetVersion(hl)
	wf.Info = fileInfo
	return &wf, nil
}

// openFile opens the file and returns a workspaceFile struct
func (mgr *CrlWorkspaceManager) openFile(fileInfo os.FileInfo, hl *core.Transaction) (*workspaceFile, error) {
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
	element, err2 := mgr.GetUofD().RecoverDomain(fileContent, hl)
	if err2 != nil {
		return nil, err2
	}
	if !writable {
		element.SetReadOnlyRecursively(true, hl)
	}
	var wf workspaceFile
	wf.filename = filename
	wf.Domain = element
	wf.Info = fileInfo
	wf.LoadedVersion = element.GetVersion(hl)
	wf.File = file
	return &wf, nil
}

// LoadSettings loads the settings saved in the workspace
func (mgr *CrlWorkspaceManager) LoadSettings() error {
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
func (mgr *CrlWorkspaceManager) LoadWorkspace(hl *core.Transaction) error {
	files, err := ioutil.ReadDir(mgr.editor.userPreferences.WorkspacePath)
	if err != nil {
		return errors.Wrap(err, "CrlWorkspaceManager.LoadWorkspace failed")
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".acrl") {
			workspaceFile, err := mgr.openFile(f, hl)
			if err != nil {
				return errors.Wrap(err, "CrlWorkspaceManager.LoadWorkspace failed")
			}
			mgr.workspaceFiles[workspaceFile.Domain.GetConceptID(hl)] = workspaceFile
		}
	}
	mgr.LoadSettings()
	return nil
}

// saveFile saves the file and updates the fileInfo
func (mgr *CrlWorkspaceManager) saveFile(wf *workspaceFile, hl *core.Transaction) error {
	hl.ReadLockElement(wf.Domain)
	if wf.File == nil {
		return errors.New("CrlBrowserEditor.SaveFile called with nil file")
	}
	byteArray, err := mgr.GetUofD().MarshalDomain(wf.Domain, hl)
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
	newFilename := mgr.generateFilename(wf.Domain, hl)
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
func (mgr *CrlWorkspaceManager) SaveWorkspace(hl *core.Transaction) error {
	rootElements := mgr.editor.uOfDManager.UofD.GetRootElements(hl)
	var err error
	for id, el := range rootElements {
		noSaveDomains := mgr.editor.getNoSaveDomains(hl)
		if !el.GetIsCore(hl) && noSaveDomains[el.GetConceptID(hl)] == nil {
			workspaceFile := mgr.workspaceFiles[id]
			if workspaceFile != nil {
				err = mgr.saveFile(workspaceFile, hl)
				if err != nil {
					return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
				}
			} else {
				workspaceFile, err = mgr.newFile(el, hl)
				if err != nil {
					return errors.Wrap(err, "CrlWorkspaceManager.SaveWorkspace failed")
				}
				mgr.workspaceFiles[id] = workspaceFile
				err = mgr.saveFile(workspaceFile, hl)
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
