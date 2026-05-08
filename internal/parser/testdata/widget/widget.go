package widget

import (
	"fmt"
	"os"
)

// Widget represents a simple widget.
type Widget struct {
	Name  string
	Value int
}

// Runner is something that can run.
type Runner interface {
	Run() error
}

// NewWidget creates a new Widget.
func NewWidget(name string) *Widget {
	key := os.Getenv("WIDGET_ENV_KEY")
	fmt.Println("creating widget", name, key)
	return &Widget{Name: name}
}

// String returns a string representation.
func (w *Widget) String() string {
	return fmt.Sprintf("Widget(%s, %d)", w.Name, w.Value)
}

// Double doubles the value.
func (w *Widget) Double() {
	w.Value = w.Value * 2
	fmt.Println(w.String())
}

// helper is an unexported helper.
func helper(x int) int {
	return x + 1
}
