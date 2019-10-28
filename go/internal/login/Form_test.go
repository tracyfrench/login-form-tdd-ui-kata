package login_test

import (
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/SoftDevGang/login-form-tdd-ui-kata/go/internal/login"
)

// UI simulation.
type testingUI struct {
	buttonCalled map[string]bool
	buttonText   map[string]string
	buttonBounds map[string]rl.Rectangle

	textBoxCalled    map[string]bool
	textBoxText      map[string]string
	textBoxUserInput map[string]string // User interaction could stand out from other fields?
}

func newTestingUI() *testingUI {
	ui := &testingUI{
		buttonCalled:     make(map[string]bool),
		buttonText:       make(map[string]string),
		buttonBounds:     make(map[string]rl.Rectangle),
		textBoxCalled:    make(map[string]bool),
		textBoxText:      make(map[string]string),
		textBoxUserInput: make(map[string]string),
	}
	return ui
}

func (ui *testingUI) TextBox(id string, bounds rl.Rectangle, text string) string {
	ui.textBoxCalled[id] = true
	ui.textBoxText[id] = text
	return ui.textBoxUserInput[id]
}

func (ui *testingUI) Label(id string, bounds rl.Rectangle, text string) {
}

func (ui *testingUI) Button(id string, bounds rl.Rectangle, text string) bool {
	ui.buttonCalled[id] = true
	ui.buttonText[id] = text
	ui.buttonBounds[id] = bounds
	return false
}

// ***** There is a "Log in" button in right corner of the dialog. *****
// right corner is design.

func TestForm_LoginButton(t *testing.T) {
	var form login.Form
	ui := newTestingUI()

	form.Render(ui)

	if !ui.buttonCalled["login"] {
		t.Errorf("not found")
	}
}

func TestForm_LoginButtonText(t *testing.T) {
	var form login.Form
	ui := newTestingUI()

	form.Render(ui)

	if "Log in" != ui.buttonText["login"] {
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
	ui := newTestingUI()

	form.Render(ui)

	// first idea: if ui.buttonCoordinate != BottomRight {
	expectedBounds := rl.Rectangle{235, 165, 235 + 110, 165 + 30} // needed to open GIMP
	if ui.buttonBounds["login"] != expectedBounds {
		t.Errorf("expected %v, but was %v", expectedBounds, ui.buttonBounds)
	}
	// Peter:
	// I do not feel good because I believe the coordinates will change a lot while working on it.
	// We assert on layout.
	// If we would use 1 screen shot for the layout in the end (Approvals), then we would not need
	// to assert the button, color, text etc. We would only test that button triggers.
}

// ***** There is a user name field, which is limited to 20 characters. *****

func TestForm_UserNameField(t *testing.T) {
	var form login.Form
	ui := newTestingUI()

	form.Render(ui)

	if !ui.textBoxCalled["username"] {
		t.Errorf("not found")
	}
}

// ignoring design

func TestForm_UserNameFieldInputKept(t *testing.T) {
	var form login.Form
	ui := newTestingUI()
	ui.textBoxUserInput["username"] = "Christian"

	form.Render(ui)

	// we test whether the form keeps input
	if form.UserName != "Christian" {
		t.Errorf("not kept")
	}
	// need more tests because UI lib is IM/stateless?
	// Peter surprised that I need to test for keeping the state.
	// Christian says he is used to having the "UI model" of the IM state.
	// Maybe the missing state makes it easier to separate?
}

func TestForm_UserNameFieldInputDisplayed(t *testing.T) {
	var form login.Form
	form.UserName = "Peter"
	ui := newTestingUI()

	form.Render(ui)

	// we test whether the next render call receives the new text
	if ui.textBoxText["username"] != "Peter" {
		t.Errorf("not displayed")
	}
}

// test for initial text empty skipped, it is empty.

func TestForm_UserNameFieldInputLimited(t *testing.T) {
	original := strings.Repeat("a", login.UserNameLimit)
	var form login.Form
	form.UserName = original
	ui := newTestingUI()
	ui.textBoxUserInput["username"] = form.UserName + "f"

	form.Render(ui)

	if form.UserName != original {
		t.Errorf("not limited \"%v\"", form.UserName)
	}
}

// ***** The label "Phone, email or username" is next to the input field. *****

// Do we need to test the label? It is only style?
// Typical style = colours, positions, pixel perfect stuff.
// Does it matter, we test everything which we can easily?
// Screen shot does not say what is wrong. The message is the image diff. We need to show the diff.
//   Christian likes images. Peter thinks it's heavy weight.

// What about the label?
// Peter would write a unit test.
// Christian would do all presentation stuff with screen shot like a schema validation.
// -> just wrote the code.

// btw. the test UI does not force any design.

// ***** There is a password field, which is limited to 20 characters. *****

func TestForm_PasswordField(t *testing.T) {
	var form login.Form
	ui := newTestingUI()

	form.Render(ui)

	if !ui.textBoxCalled["password"] {
		t.Errorf("not found")
	}
	// now we need an ID of the field because we have 2.
	// * use order of invocations -> can NOT reorder statements on refactor.
	// * add ID to our API (only for tests). IDs are used in other library imgui as well.
}

// skip tests for state, similar to 1st field
// skip tests for limit, similar to 1st field

// ***** The password is either visible as asterisk or bullet signs. *****

func TestForm_PasswordFieldIsDisplayedMasked(t *testing.T) {
	t.SkipNow()

	var form login.Form
	form.Password = "secret"
	ui := newTestingUI()

	form.Render(ui)

	// this does not work, RayGui cannot do this because we need to change the rendering.
	// need to create your own function to use Password function. Will never see "*", only in screen shot.
	if ui.textBoxText["password"] != "*****" {
		t.Errorf("not masked")
	}
	// ImGui can do it with ImGuiInputTextFlags_Password = 1 << 15 ... Password mode, display all characters as '*'
}

// ***** The label "Password" is next to the input field. *****

// skip tests, similar to 1st label

// ***** There is a label in a red box above the button(s). It is only visible if there is an error. *****
