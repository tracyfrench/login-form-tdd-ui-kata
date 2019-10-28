package login_test

import (
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/SoftDevGang/login-form-tdd-ui-kata/go/internal/login"
)

// UI simulation.
type testingUI struct {
	buttonCalled bool
	buttonText   string
	buttonBounds rl.Rectangle

	textBoxCalled    bool
	textBoxText      string
	textBoxUserInput string // User interaction could stand out from other fields?
}

func (ui *testingUI) TextBox(bounds rl.Rectangle, text string) string {
	ui.textBoxCalled = true
	ui.textBoxText = text
	return ui.textBoxUserInput
}

func (ui *testingUI) Button(bounds rl.Rectangle, text string) bool {
	ui.buttonCalled = true
	ui.buttonText = text
	ui.buttonBounds = bounds
	return false
}

// * There is a "Log in" button in right corner of the dialog.
// right corner is design.

func TestForm_LoginButton(t *testing.T) {
	var form login.Form
	var ui testingUI

	form.Render(&ui)

	if !ui.buttonCalled {
		t.Errorf("not found")
	}
}

func TestForm_LoginButtonText(t *testing.T) {
	var form login.Form
	var ui testingUI

	form.Render(&ui)

	if "Log in" != ui.buttonText {
		t.Errorf("is not \"Log in\"")
	}
}

func TestForm_LoginButtonBounds(t *testing.T) {
	// now we talk about layout. Options:
	// 1. we could abstract Coordinates -> map coordinates
	// 2. we could move coordinates into the UI's responsibility -> more code in UI
	// 3. we ignore coordinates and check manually -> need to update all calls later
	// 4. we use Raylib coordinates and just mock what is hurting us -> use Raylib in tests
	var form login.Form
	var ui testingUI

	form.Render(&ui)

	// first idea: if ui.buttonCoordinate != BottomRight {
	expectedBounds := rl.Rectangle{235, 165, 235 + 110, 165 + 30} // needed to open GIMP
	if ui.buttonBounds != expectedBounds {
		t.Errorf("expected %v, but was %v", expectedBounds, ui.buttonBounds)
	}
	// Peter:
	// I do not feel good because I believe the coordinates will change a lot while working on it.
	// We assert on layout.
	// If we would use 1 screen shot for the layout in the end (Approvals), then we would not need
	// to assert the button, color, text etc. We would only test that button triggers.
}

// * There is a user name field, which is limited to 20 characters.

func TestForm_UserNameField(t *testing.T) {
	var form login.Form
	var ui testingUI

	form.Render(&ui)

	if !ui.textBoxCalled {
		t.Errorf("not found")
	}
}

// ignoring design

func TestForm_UserNameFieldInputKept(t *testing.T) {
	var form login.Form
	var ui testingUI
	ui.textBoxUserInput = "hello"

	form.Render(&ui)

	// we test whether the form keeps input
	if form.UserName != "hello" {
		t.Errorf("not kept")
	}
	// need more tests because UI lib is IM/stateless?
	// Peter surprised that I need to test for keeping the state.
	// Christian says he is used to having the "UI model" of the IM state.
	// Maybe the missing state makes it easier to separate?
}

func TestForm_UserNameFieldInputDisplayed(t *testing.T) {
	var form login.Form
	form.UserName = "username"
	var ui testingUI

	form.Render(&ui)

	// we test whether the next render call receives the new text
	if ui.textBoxText != "username" {
		t.Errorf("not displayed")
	}
}

// test for initial text empty skipped, it is empty.

func TestForm_UserNameFieldInputLimited(t *testing.T) {
	original := strings.Repeat("a", login.UserNameLimit)
	var form login.Form
	form.UserName = original
	var ui testingUI
	ui.textBoxUserInput = form.UserName + "f"

	form.Render(&ui)

	if form.UserName != original {
		t.Errorf("not limited \"%v\"", form.UserName)
	}
}

// * The label "Phone, email or username" is next to the input field.
// * There is a password field, which is limited to 20 characters.
// * The password is either visible as asterisk or bullet signs.
// * The label "Password" is next to the input field.
// * There is a label in a red box above the button(s). It is only visible if there is an error.
