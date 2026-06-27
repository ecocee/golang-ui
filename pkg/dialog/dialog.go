package dialog

import (
	"github.com/sqweek/dialog"
)

// Service provides native OS dialogs (file picking, alerts, messages).
// It can be directly registered into the Glyra API Bridge so the frontend can trigger OS dialogs.
type Service struct{}

// MessageReq defines the parameters for a native message box.
type MessageReq struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"` // "info", "error"
}

// ShowMessage displays a native OS alert dialog.
func (s *Service) ShowMessage(req MessageReq) (string, error) {
	b := dialog.Message("%s", req.Message).Title(req.Title)
	if req.Type == "error" {
		b.Error()
	} else if req.Type == "info" {
		b.Info()
	} else {
		// Default to yes/no for other types just as a fallback
		b.YesNo()
	}
	return "ok", nil
}

// FilePickerReq defines the parameters for an Open/Save file dialog.
type FilePickerReq struct {
	Title   string `json:"title"`
	StartDir string `json:"start_dir"`
	FilterDesc string `json:"filter_desc"`
	FilterExt string `json:"filter_ext"` // e.g. "*.txt", "*.png"
	SaveMode bool `json:"save_mode"`
}

// PickFile opens a native OS file chooser and returns the absolute path selected by the user.
func (s *Service) PickFile(req FilePickerReq) (string, error) {
	b := dialog.File().Title(req.Title)
	if req.StartDir != "" {
		b = b.SetStartDir(req.StartDir)
	}
	if req.FilterDesc != "" && req.FilterExt != "" {
		b = b.Filter(req.FilterDesc, req.FilterExt)
	}
	
	var path string
	var err error
	if req.SaveMode {
		path, err = b.Save()
	} else {
		path, err = b.Load()
	}
	return path, err
}

// PickDirectory opens a native OS folder chooser and returns the absolute path selected.
func (s *Service) PickDirectory(req FilePickerReq) (string, error) {
	b := dialog.Directory().Title(req.Title)
	if req.StartDir != "" {
		b = b.SetStartDir(req.StartDir)
	}
	path, err := b.Browse()
	return path, err
}
