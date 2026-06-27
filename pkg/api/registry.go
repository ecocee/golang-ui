package api

import (
	"fmt"
	"reflect"

	"github.com/webview/webview_go"
)

// Registry manages the Go-to-WebView IPC bridge.
type Registry struct {
	w        webview.WebView
	services map[string]interface{}
}

// NewRegistry creates a new API registry for binding Go services to the WebView.
func NewRegistry(w webview.WebView) *Registry {
	return &Registry{
		w:        w,
		services: make(map[string]interface{}),
	}
}

// Register binds all exported methods of a struct to the WebView.
// For example, Register("Auth", &AuthService{}) binds window.Auth_Login().
func (r *Registry) Register(name string, service interface{}) error {
	r.services[name] = service
	val := reflect.ValueOf(service)
	typ := reflect.TypeOf(service)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		// We use an underscore because w.Bind registers functions on the global window object.
		bindName := fmt.Sprintf("%s_%s", name, method.Name)

		methodFunc := val.Method(i).Interface()
		if err := r.w.Bind(bindName, methodFunc); err != nil {
			return fmt.Errorf("failed to bind %s: %w", bindName, err)
		}
	}
	return nil
}
