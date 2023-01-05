package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"caos/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var group sync.WaitGroup

// Layout service - recreates the terminal definitions and parameters for console app
type ServiceLayout struct {
	app   *tview.Application
	pages *tview.Pages
	// Flex
	consoleView  *tview.Grid
	affinityView *tview.Grid
	// User form
	refinementInput *tview.Form
	// User input
	promptInput *tview.InputField
	// Details output
	metadataOutput *tview.TextView
	promptOutput   *tview.TextView
	infoOutput     *tview.TextView
}

func OnResultChange(text string) {
	results = util.ParseInt(text)
}

func OnProbabilityChange(text string) {
	probabilities = util.ParseInt(text)
}

func OnTemperatureChange(text string) {
	temperature = util.ParseFloat32(text)
}

func OnToppChange(text string) {
	topp = util.ParseFloat32(text)
}

func OnPenaltyChange(text string) {
	penalty = util.ParseFloat32(text)
}

func OnFrequencyPenaltyChange(text string) {
	frequency = util.ParseFloat32(text)
}

func OnTypeAccept(text string, lastChar rune) bool {
	matched := util.MatchNumber(text)
	return matched
}

func OnBack() {
	// Console view
	Node.Layout.pages.ShowPage("console")
	Node.Layout.pages.HidePage("refinement")
	ValidateRefinementForm()
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

func OnRefinementTopic() {
	// Refinement view
	Node.Layout.pages.HidePage("console")
	Node.Layout.pages.ShowPage("refinement")
}

func OnExportTopic() {
	if Node.Layout.promptOutput.GetText(true) == "" {
		return
	}

	unix_milliseconds := fmt.Sprint(time.Now().UnixMilli())
	ts_file := fmt.Sprintf("prompt-%s.txt", unix_milliseconds)
	var dir string
	if dir, e := os.Getwd(); e != nil {
		fmt.Printf("e: %v\n", e)
		fmt.Printf("dir: %v\n", dir)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/export", dir)); os.IsNotExist(err) {
		os.Mkdir("export", 0755)
	}

	path := filepath.Join(dir, "export", ts_file)
	out, err := os.Create(path)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	out.WriteString(Node.Layout.promptOutput.GetText(true))
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
}

func OnChangeWords(option string, optionIndex int) {
	switch option {
	case "1":
		maxtokens = 1
	case "50":
		maxtokens = 37
	case "85":
		maxtokens = 64
	case "100":
		maxtokens = 75
	case "200":
		maxtokens = 150
	case "500":
		maxtokens = 375
	case "1000":
		maxtokens = 750
	case "1500":
		maxtokens = 1125
	}
}

// Text field from input
func OnTextAccept(textToCheck string, lastChar rune) bool {
	if isLoading {
		Node.Layout.promptInput.SetText("...")
		return false
	}

	Node.Agent.currentUser.engineProperties = Node.Agent.currentUser.SetEngineParameters(
		engine,      // "code-davinci-002",
		temperature, // if temperature is used set topp to 1.0
		topp,        // if topp is used set temperature to 1.0
		penalty,     // Penalize from 0 to 1 the repeated tokens
		frequency,   // Frequency  of penalization
	)

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
	if key == tcell.KeyEnter && !isLoading {
		if mode == "Edit" {
			group.Add(1)
			go func() {
				isLoading = true
				Node.Agent.InstructionRequest()
				Node.Layout.promptInput.SetText("")
				group.Done()
			}()
		} else {
			group.Add(1)
			go func() {
				isLoading = true
				Node.Agent.StartRequest()
				Node.Layout.promptInput.SetText("")
				group.Done()
				if isEditable {
					mode = "Edit"
				}
			}()
		}
		Node.Layout.promptInput.SetText("...")
	} else {
		Node.Layout.promptInput.SetText("...")
	}
}

func OnEditChecked(state bool) {
	if state {
		isEditable = true
	} else {
		isEditable = false
		mode = "Text"
	}
}

func OnConversationChecked(state bool) {
	if state {
		isConversational = true
	} else {
		isConversational = false
	}
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

func ValidateRefinementForm() {
	// Default Values
	resultInput := Node.Layout.refinementInput.GetFormItem(0).(*tview.InputField)
	probabilityInput := Node.Layout.refinementInput.GetFormItem(1).(*tview.InputField)
	temperatureInput := Node.Layout.refinementInput.GetFormItem(2).(*tview.InputField)
	toppInput := Node.Layout.refinementInput.GetFormItem(3).(*tview.InputField)
	penaltyInput := Node.Layout.refinementInput.GetFormItem(4).(*tview.InputField)
	frecuencyInput := Node.Layout.refinementInput.GetFormItem(5).(*tview.InputField)

	if !util.MatchNumber(resultInput.GetText()) {
		resultInput.SetText("1")
	}

	if !util.MatchNumber(probabilityInput.GetText()) {
		probabilityInput.SetText("1")
	}

	if !util.MatchNumber(toppInput.GetText()) {
		temperatureInput.SetText("1.0")
	}

	if !util.MatchNumber(toppInput.GetText()) {
		toppInput.SetText("0.4")
	}

	if !util.MatchNumber(penaltyInput.GetText()) {
		penaltyInput.SetText("0.5")
	}

	if !util.MatchNumber(frecuencyInput.GetText()) {
		frecuencyInput.SetText("0.5")
	}
}

func CreateConsoleView() bool {
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
		AddDropDown("Words", []string{"1", "50", "85", "100", "200", "500", "1000", "1500"}, 4, OnChangeWords).
		AddButton("Affinity", OnRefinementTopic).
		AddButton("New conversation", OnNewTopic).
		AddButton("Export conversation", OnExportTopic).
		SetHorizontal(true).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
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
	Node.Layout.consoleView.
		SetBorderPadding(0, 0, 9, 9).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor())
	// Validate view
	if Node.Layout.consoleView != nil {
		return true
	} else {
		return false
	}
}

func CreateRefinementView() bool {
	// Layout
	affinitySection := tview.NewForm()
	// Affinity section
	affinitySection.
		AddInputField("Results", fmt.Sprintf("%v", results), 5, OnTypeAccept, OnResultChange).
		AddInputField("Probabilities", fmt.Sprintf("%v", probabilities), 5, OnTypeAccept, OnProbabilityChange).
		AddInputField("Temperature [0.0 / 1.0]", fmt.Sprintf("%v", temperature), 5, OnTypeAccept, OnTemperatureChange).
		AddInputField("Topp [0.0 / 1.0]", fmt.Sprintf("%v", topp), 5, OnTypeAccept, OnToppChange).
		AddInputField("Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", penalty), 5, OnTypeAccept, OnPenaltyChange).
		AddInputField("Frecuency Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", frequency), 5, OnTypeAccept, OnFrequencyPenaltyChange).
		AddCheckbox("Edit mode (edit and improve the previous reponse)", true, OnEditChecked).
		AddCheckbox("Conversational mode", false, OnConversationChecked).
		AddButton("Back to chat", OnBack).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignCenter).
		SetTitle("Improve your search criteria: ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorDarkCyan.TrueColor()).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(2, 1, 9, 9)
	// Refinement form
	Node.Layout.refinementInput = affinitySection
	// Affinity grid
	Node.Layout.affinityView = tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 1).
		AddItem(affinitySection, 0, 0, 1, 1, 0, 0, true)
	// Affinity
	Node.Layout.affinityView.
		SetBorderPadding(15, 15, 20, 20).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor())
	// Validate view
	if Node.Layout.affinityView != nil {
		return true
	} else {
		return false
	}
}

// Create service layour for terminal session
func InitializeLayout() {
	/* Layout content */
	GenerateLayoutContent()
	// Create views
	CreateConsoleView()
	CreateRefinementView()
	// Window frame
	Node.Layout.pages = tview.NewPages()
	Node.Layout.pages.
		AddAndSwitchToPage("console", Node.Layout.consoleView, true).
		AddAndSwitchToPage("refinement", Node.Layout.affinityView, true)
	// Main executor
	Node.Layout.app = tview.NewApplication()
	// App terminal configuration
	Node.Layout.app.
		SetRoot(Node.Layout.pages, true).
		SetFocus(Node.Layout.promptInput).
		EnableMouse(true)
	// Console view
	Node.Layout.pages.ShowPage("console")
	Node.Layout.pages.HidePage("refinement")
}
