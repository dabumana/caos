package service

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout service - recreates the terminal definitions and parameters for console app
type ServiceLayout struct {
	app   *tview.Application
	pages *tview.Pages
	// Flex
	consoleView *tview.Grid
	// User input
	promptInput *tview.InputField
	// Details output
	metadataOutput *tview.TextView
	promptOutput   *tview.TextView
	infoOutput     *tview.TextView
}

// Define new conversation
func OnNewTopic() {
	mode = "Text"

	Node.Layout.infoOutput.SetText("")
	Node.Layout.metadataOutput.SetText("")
	Node.Layout.promptOutput.SetText("")
	Node.Layout.promptInput.SetPlaceholder("Type here...")
	Node.Layout.promptInput.SetText("")
}

// Dropdown from input to change mode
func OnChangeMode(option string, optionIndex int) {
	switch option {
	case "Edit":
		engine = "text-davinci-edit-001"
		mode = "Edit"
	case "Code":
		engine = "code-davinci-002"
		mode = "Code"
	case "Text":
		engine = "text-davinci-003"
		mode = "Text"
	}

	Node.Agent.currentUser.engineProperties = Node.Agent.currentUser.SetEngineParameters(
		engine,      // "code-davinci-002",
		temperature, // if temperature is used set topp to 1.0
		topp,        // if topp is used set temperature to 1.0
		penalty,     // Penalize from 0 to 1 the repeated tokens
		frequency,   // Frequency  of penalization
	)

	Node.Agent.currentUser.promptProperties = Node.Agent.currentUser.SetRequestParameters(
		promptctx,
		prompt,
		maxtokens,
		results,
		probabilities,
	)
}

// Dropdown from input to change engine
func OnChangeEngine(option string, optionIndex int) {
	switch option {
	case "davinci":
		engine = "text-davinci-003"
		mode = "Text"
	case "babbage":
		engine = "text-babbage-002"
		mode = "Text"
	case "ada":
		engine = "text-ada-001"
		mode = "Text"
	case "curie":
		engine = "text-curie-002"
		mode = "Text"
	case "cushman":
		engine = "code-cushman-001"
		mode = "Code"
	}

	Node.Agent.currentUser.engineProperties = Node.Agent.currentUser.SetEngineParameters(
		engine,      // "code-davinci-002",
		temperature, // if temperature is used set topp to 1.0
		topp,        // if topp is used set temperature to 1.0
		penalty,     // Penalize from 0 to 1 the repeated tokens
		frequency,   // Frequency  of penalization
	)

	Node.Agent.currentUser.promptProperties = Node.Agent.currentUser.SetRequestParameters(
		promptctx,
		prompt,
		maxtokens,
		results,
		probabilities,
	)
}

// Dropdown from input to change results and probabilities
func OnChangeResultsAndProbs(option string, optionIndex int) {
	switch option {
	case "p1":
		probabilities = 1
	case "p2":
		probabilities = 2
	case "p4":
		probabilities = 4
	case "p8":
		probabilities = 8
	case "r1":
		results = 1
	case "r2":
		results = 2
	case "r4":
		results = 4
	case "r8":
		results = 8
	}

	Node.Agent.currentUser.engineProperties = Node.Agent.currentUser.SetEngineParameters(
		engine,      // "code-davinci-002",
		temperature, // if temperature is used set topp to 1.0
		topp,        // if topp is used set temperature to 1.0
		penalty,     // Penalize from 0 to 1 the repeated tokens
		frequency,   // Frequency  of penalization
	)

	Node.Agent.currentUser.promptProperties = Node.Agent.currentUser.SetRequestParameters(
		promptctx,
		prompt,
		maxtokens,
		results,
		probabilities,
	)
}

func OnChangeWords(option string, optionIndex int) {
	switch option {
	case "1":
		maxtokens = 1
	case "50":
		maxtokens = 37
	case "100":
		maxtokens = 75
	case "500":
		maxtokens = 375
	case "1000":
		maxtokens = 750
	case "1500":
		maxtokens = 1125
	}
}

// Checkbox from input
func OnCheck(checked bool) {
	if checked {
		// Global parameters
		temperature = 0.4
		topp = 0.8
		penalty = 0.5
		frequency = 0.5
		// Set engine properties for accurated results
		Node.Agent.currentUser.engineProperties = Node.Agent.currentUser.SetEngineParameters(
			engine,      // "code-davinci-002",
			temperature, // if temperature is used set topp to 1.0
			topp,        // if topp is used set temperature to 1.0
			penalty,     // Penalize from 0 to 1 the repeated tokens
			frequency,   // Frequency  of penalization
		)
	}
}

// Text field from input
func OnTextAccept(textToCheck string, lastChar rune) bool {
	if mode == "Edit" {
		Node.Agent.currentUser.promptProperties = Node.Agent.currentUser.SetRequestParameters(
			promptctx,
			[]string{textToCheck},
			maxtokens,
			results,
			probabilities,
		)
		return true
	} else {
		Node.Agent.currentUser.promptProperties = Node.Agent.currentUser.SetRequestParameters(
			[]string{textToCheck},
			[]string{""},
			maxtokens,
			results,
			probabilities,
		)
		return true
	}
}

// Text key event
func OnTextDone(key tcell.Key) {
	var group sync.WaitGroup
	if key == tcell.KeyEnter {
		if mode == "Edit" {
			group.Add(1)
			go func() {
				Node.Agent.InstructionRequest()
				group.Done()
			}()
		} else {
			group.Add(1)
			go func() {
				Node.Agent.StartRequest()
				mode = "Edit"
				group.Done()
			}()
		}
	}
	Node.Layout.promptInput.SetText("")
	group.Wait()
}

/* Service layout functionality */
func GenerateLayoutContent() {
	// COM
	Node.Layout.promptInput = tview.NewInputField()
	Node.Layout.promptOutput = tview.NewTextView()
	// Metadata
	Node.Layout.metadataOutput = tview.NewTextView()
	// Info
	Node.Layout.infoOutput = tview.NewTextView()
	// Input
	Node.Layout.promptInput.
		SetFieldTextColor(tcell.ColorDarkOrange.TrueColor()).
		SetAcceptanceFunc(OnTextAccept).
		SetDoneFunc(OnTextDone).
		SetLabel("Enter your request: ").
		SetLabelColor(tcell.ColorDarkOrange.TrueColor()).
		SetPlaceholder("Type here...").
		SetFieldBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetFieldTextColor(tcell.ColorWhiteSmoke.TrueColor())
	//Output
	Node.Layout.infoOutput.SetToggleHighlights(true).
		SetLabel("Request: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOrange.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
	Node.Layout.metadataOutput.SetToggleHighlights(true).
		SetLabel("Description: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
	Node.Layout.promptOutput.SetToggleHighlights(true).
		SetLabel("Response: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
}

// Create service layour for terminal session
func InitializeLayout() {
	/* Console View */
	GenerateLayoutContent()
	// help
	helpOutput := tview.NewTextView()
	helpOutput.
		SetText("Press CTRL + C or CMD + Q to exit from the application.\nGo to fullscreen for advanced options.").
		SetTextAlign(tview.AlignRight)
	// Layout
	detailsSection := tview.NewForm()
	metadataSection := tview.NewFlex()
	infoSection := tview.NewFlex()
	comSection := tview.NewFlex()
	// Console section
	detailsSection.
		AddDropDown("Mode", []string{"Edit", "Code", "Text"}, 2, OnChangeMode).
		AddDropDown("Engine", []string{"davinci", "curie", "babbage", "ada", "cushman"}, 0, OnChangeEngine).
		AddDropDown("Results", []string{"r1", "r2", "r4", "r8"}, 0, OnChangeResultsAndProbs).
		AddDropDown("Probabilities", []string{"p1", "p2", "p4", "p8"}, 0, OnChangeResultsAndProbs).
		AddDropDown("Words", []string{"1", "50", "100", "500", "1000", "1500"}, 2, OnChangeWords).
		AddCheckbox("Affinity", false, OnCheck).
		AddButton("New conversation", OnNewTopic).
		SetHorizontal(true).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOrange.TrueColor()).
		SetButtonsAlign(tview.AlignRight)
	metadataSection.
		AddItem(Node.Layout.metadataOutput, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDarkSlateGray.TrueColor()).
		SetTitle("Metadata").
		SetTitleColor(tcell.ColorOrange.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	infoSection.
		AddItem(Node.Layout.infoOutput, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetTitle("Details").
		SetTitleColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	comSection.
		AddItem(Node.Layout.promptOutput, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDarkSlateGray.TrueColor()).
		SetTitle("Prompter").
		SetTitleColor(tcell.ColorOrange.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	// Console grid
	Node.Layout.consoleView = tview.NewGrid().
		SetRows(0, 12).
		SetColumns(0, 2).
		AddItem(metadataSection, 1, 0, 1, 1, 0, 0, false).
		AddItem(infoSection, 1, 1, 1, 4, 0, 0, false).
		AddItem(detailsSection, 2, 0, 1, 5, 50, 0, true).
		AddItem(comSection, 3, 0, 9, 5, 0, 0, false).
		AddItem(Node.Layout.promptInput, 12, 0, 1, 5, 0, 0, true).
		AddItem(helpOutput, 13, 0, 1, 5, 0, 0, false)
	// Console
	Node.Layout.consoleView.SetBorderPadding(0, 0, 9, 9).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor())
	// Window frame
	Node.Layout.pages = tview.NewPages()
	Node.Layout.pages.AddAndSwitchToPage("console", Node.Layout.consoleView, true)
	// Main executor
	Node.Layout.app = tview.NewApplication()
	// App terminal configuration
	Node.Layout.app.
		SetRoot(Node.Layout.pages, true).
		SetFocus(Node.Layout.promptInput).
		EnableMouse(true)
}
