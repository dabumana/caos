// Package service section
package service

import (
	"fmt"
	"strings"
	"sync"

	"caos/model"
	"caos/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var group sync.WaitGroup

// Layout - Recreates the terminal definitions and parameters for a console app
type Layout struct {
	app    *tview.Application
	screen *tcell.Screen
	pages  *tview.Pages
	// Flex
	consoleView  *tview.Grid
	affinityView *tview.Grid
	// User form
	refinementInput *tview.Form
	detailsInput    *tview.Form
	idInput         *tview.Form
	// User modal
	modalInput *tview.Modal
	// User input
	promptArea *tview.TextArea
	// Details output
	metadataOutput *tview.TextView
	promptOutput   *tview.TextView
	infoOutput     *tview.TextView
	// Event manager
	eventManager EventManager
}

// OnResultChange - Evaluates when an input text changes for the result input field
func OnResultChange(text string) {
	node.controller.currentAgent.preferences.Results = util.ParseInt32(text)
}

// OnProbabilityChange - Evaluates when an input text changes for the probability input field
func OnProbabilityChange(text string) {
	node.controller.currentAgent.preferences.Probabilities = util.ParseInt32(text)
}

// OnTemperatureChange - Evaluates when an input text changes for the temperature input field
func OnTemperatureChange(text string) {
	node.controller.currentAgent.preferences.Temperature = util.ParseFloat32(text)
}

// OnToppChange - Evaluates when an input text changes for the topp input field
func OnToppChange(text string) {
	node.controller.currentAgent.preferences.Topp = util.ParseFloat32(text)
}

// OnPenaltyChange - Evaluates when an input text changes for the penalty input field
func OnPenaltyChange(text string) {
	node.controller.currentAgent.preferences.Penalty = util.ParseFloat32(text)
}

// OnFrequencyChange - Evaluates when an input text changes for the frequency penalty input field
func OnFrequencyChange(text string) {
	node.controller.currentAgent.preferences.Frequency = util.ParseFloat32(text)
}

// OnTypeAccept - Evaluates when an input text matches the field criteria
func OnTypeAccept(text string, lastChar rune) bool {
	matched := util.MatchNumber(text)
	return matched
}

// OnBack - Button event to return to the main page
func OnBack() {
	// Console view
	ReturnToPage(1)
	// Validate layout forms
	ValidateRefinementForm()
	// Clean input
	if node.controller.currentAgent.preferences.Mode == "Edit" {
		node.layout.promptArea.SetLabel("Enter your context first: ")
	}
}

// OnNewTopic - Define a new conversation button event
func OnNewTopic() {
	node.controller.currentAgent.preferences.IsNewSession = true
	node.controller.currentAgent.preferences.IsPromptReady = false
	node.controller.currentAgent.preferences.PromptCtx = []string{""}

	if node.controller.currentAgent.preferences.IsTraining {
		OnTrainingTopic()
		return
	}

	ClearConsoleView()
}

// ClearConsoleView - Clean input-output fields
func ClearConsoleView() {
	node.layout.infoOutput.SetText("")
	node.layout.metadataOutput.SetText("")
	node.layout.promptOutput.SetText("")
	node.layout.promptArea.SetPlaceholder("Type here...")
	node.layout.promptArea.SetText("", true)
}

// OnRefinementTopic - Refinement view button event
func OnRefinementTopic() {
	// Refinement view
	ReturnToPage(2)
}

// OnTrainingTopic - Modal confirmation to export training
func OnTrainingTopic() {
	// Training modal view
	ReturnToPage(3)
}

// OnExportTopic - Export current conversation as a file .txt
func OnExportTopic() {
	if node.layout.promptOutput.GetText(true) == "" {
		return
	}

	out := util.ConstructPathFileTo("export", "txt")
	out.WriteString(node.layout.promptOutput.GetText(true))
}

// OnExportTrainedTopic - Export current conversation as a trained model as a .json file
func OnExportTrainedTopic() {
	if node.controller.currentAgent.preferences.IsTraining {
		node.layout.eventManager.SaveTraining()
	}
}

// OnChangeRoles - Dropdown from input to change role
func OnChangeRoles(option string, optionIndex int) {
	if strings.Contains(option, string(model.User)) {
		node.controller.currentAgent.preferences.Role = model.User
	} else if strings.Contains(option, string(model.Assistant)) {
		node.controller.currentAgent.preferences.Role = model.Assistant
	} else if strings.Contains(option, string(model.System)) {
		node.controller.currentAgent.preferences.Role = model.System
	}
}

// OnChangeEngine - Dropdown from input to change engine
func OnChangeEngine(option string, optionIndex int) {
	if node.controller.currentAgent.preferences.IsEditable &&
		!node.controller.currentAgent.preferences.IsPromptReady {
		node.layout.promptArea.SetLabel("Retrieve en editable context sending a prompt to the selected model: ")
	} else {
		node.layout.promptArea.SetLabel("Enter your request: ")
	}

	node.layout.promptArea.SetText("", true)
	mode := node.layout.detailsInput.GetFormItem(0).(*tview.TextView)
	node.controller.currentAgent.preferences.Engine = option

	if strings.Contains(option, "edit") {
		node.controller.currentAgent.preferences.Mode = "Edit"
		node.layout.promptArea.SetLabel("Enter your context first: ")
	} else if strings.Contains(option, "code") {
		node.controller.currentAgent.preferences.Mode = "Code"
	} else if strings.Contains(option, "search") {
		node.controller.currentAgent.preferences.Mode = "Search"
	} else if strings.Contains(option, "insert") {
		node.controller.currentAgent.preferences.Mode = "Insert"
	} else if strings.Contains(option, "instruct") {
		node.controller.currentAgent.preferences.Mode = "Instruct"
	} else if strings.Contains(option, "similarity") {
		node.controller.currentAgent.preferences.Mode = "Similarity"
	} else if strings.Contains(option, "embedding") {
		node.controller.currentAgent.preferences.Mode = "Embedded"
		node.layout.promptArea.SetLabel("Enter the text to search for relatedness: ")
	} else if strings.Contains(option, "turbo") {
		node.controller.currentAgent.preferences.Mode = "Turbo"
	} else {
		node.controller.currentAgent.preferences.Mode = "Text"
	}

	mode.SetText(node.controller.currentAgent.preferences.Mode)
	if node.controller.currentAgent.preferences.Mode != "Edit" {
		OnNewTopic()
	}
}

// OnChangeWords - Dropdown for tokens according to the amount of words
func OnChangeWords(option string, optionIndex int) {
	switch option {
	case "\u0031":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0031")
	case "\u0035\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0033\u0031")
	case "\u0038\u0035":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0036\u0034")
	case "\u0031\u0030\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0037\u0035")
	case "\u0032\u0030\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0031\u0035\u0030")
	case "\u0035\u0030\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0033\u0037\u0035")
	case "\u0031\u0030\u0030\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0037\u0035\u0030")
	case "\u0031\u0035\u0030\u0030":
		node.controller.currentAgent.preferences.MaxTokens = util.ParseInt64("\u0031\u0031\u0032\u0035")
	}
}

// OnTextChange - Text field from input
func OnTextChange(textToCheck string, lastChar rune) bool {
	if node.controller.currentAgent.preferences.IsLoading {
		return false
	}

	textToCheck = strings.ReplaceAll(textToCheck, "\u000D", "\u0020")

	node.controller.currentAgent.engineProperties = node.controller.currentAgent.SetEngineParameters(
		node.controller.currentAgent.id,
		node.controller.currentAgent.preferences.Engine, // "text-davinci-003",
		node.controller.currentAgent.preferences.Role,
		node.controller.currentAgent.preferences.Temperature, // if temperature is used set topp to 1.0
		node.controller.currentAgent.preferences.Topp,        // if topp is used set temperature to 1.0
		node.controller.currentAgent.preferences.Penalty,     // Penalize from 0 to 1 the repeated tokens
		node.controller.currentAgent.preferences.Frequency,   // Frequency  of penalization
	)

	if !node.controller.currentAgent.preferences.IsPromptReady {
		node.controller.currentAgent.promptProperties = node.controller.currentAgent.SetPromptParameters(
			[]string{textToCheck},
			[]string{""},
			int(node.controller.currentAgent.preferences.MaxTokens),
			int(node.controller.currentAgent.preferences.Results),
			int(node.controller.currentAgent.preferences.Probabilities),
		)
		node.controller.currentAgent.preferences.PromptCtx = []string{textToCheck}
	} else {
		if node.controller.currentAgent.preferences.IsPromptReady &&
			node.controller.currentAgent.preferences.Mode == "Edit" {
			node.controller.currentAgent.promptProperties = node.controller.currentAgent.SetPromptParameters(
				node.controller.currentAgent.preferences.PromptCtx,
				[]string{textToCheck},
				int(node.controller.currentAgent.preferences.MaxTokens),
				int(node.controller.currentAgent.preferences.Results),
				int(node.controller.currentAgent.preferences.Probabilities),
			)
		}
	}

	return true
}

// OnTextAccept - Text key event
func OnTextAccept(key tcell.Key) {
	if node.controller.currentAgent.preferences.IsLoading {
		return
	}

	if node.controller.currentAgent.preferences.IsNewSession {
		node.controller.currentAgent.preferences.IsNewSession = false
	}

	err := func() {
		node.layout.infoOutput.SetText("\nUncheck edit mode in affinity preferences and try again...")
		OnNewTopic()
	}

	if key == tcell.KeyCtrlSpace && !node.controller.currentAgent.preferences.IsLoading {
		if node.controller.currentAgent.preferences.Mode == "Edit" {
			if node.controller.currentAgent.preferences.IsPromptReady {
				group.Add(1)
				go func() {
					node.controller.currentAgent.preferences.IsLoading = true
					defer group.Done()
					node.controller.EditRequest()
				}()
			}
		} else if node.controller.currentAgent.preferences.Mode == "Embedded" {
			if !node.controller.currentAgent.preferences.IsEditable {
				group.Add(1)
				go func() {
					node.controller.currentAgent.preferences.IsLoading = true
					defer group.Done()
					node.controller.EmbeddingRequest()
				}()
			} else {
				err()
			}
		} else if node.controller.currentAgent.preferences.Mode == "Turbo" {
			group.Add(1)
			go func() {
				node.controller.currentAgent.preferences.IsLoading = true
				defer group.Done()
				node.controller.ChatCompletionRequest()
			}()
		} else {
			group.Add(1)
			go func() {
				node.controller.currentAgent.preferences.IsLoading = true
				defer group.Done()
				node.controller.CompletionRequest()
			}()
		}

		group.Wait()
		if node.controller.currentAgent.preferences.IsEditable ||
			(node.controller.currentAgent.preferences.Mode == "Edit" &&
				node.controller.currentAgent.preferences.PromptCtx != nil) {

			node.controller.currentAgent.preferences.Engine = "text-davinci-edit-001"
			node.controller.currentAgent.preferences.IsPromptReady = true

			engine := node.layout.detailsInput.GetFormItem(1).(*tview.DropDown)
			engine.SetCurrentOption(ValidateSelector(node.controller.currentAgent.preferences.Engine))

			node.layout.promptArea.SetLabel("Enter your request: ")
		}
	}
}

// OnEditChecked - Editable mode activation
func OnEditChecked(state bool) {
	node.controller.currentAgent.preferences.IsEditable = state
	if node.controller.currentAgent.preferences.IsEditable {
		node.layout.promptArea.SetLabel("Retrieve en editable context sending a prompt from the selected model: ")
	} else {
		node.layout.promptArea.SetLabel("Enter your request: ")
	}
	// new topic event
	OnNewTopic()
}

// OnConversationChecked - Conversation mode for friendly responses
func OnConversationChecked(state bool) {
	node.controller.currentAgent.preferences.IsConversational = state
	if node.controller.currentAgent.preferences.IsConversational {
		node.controller.currentAgent.preferences.Mode = "Text"
		node.controller.currentAgent.preferences.Engine = "text-davinci-003"
	}
}

// OnTrainingChecked - Training mode to store the current conversation
func OnTrainingChecked(state bool) {
	node.controller.currentAgent.preferences.IsTraining = state
}

//OnStreamingChecked - Streaming mode
func OnStreamingChecked(state bool) {
	node.controller.currentAgent.preferences.IsPromptStreaming = state
}

// ValidateSelector - Validate the selected engine
func ValidateSelector(engine string) int {
	var index int
	for i := range node.controller.currentAgent.preferences.Models {
		if engine == node.controller.currentAgent.preferences.Models[i] {
			index = i
		}
	}
	return index
}

// ValidateRefinementForm - Service layout functionality
func ValidateRefinementForm() {
	// Default Values
	resultInput := node.layout.refinementInput.GetFormItem(0).(*tview.InputField)
	probabilityInput := node.layout.refinementInput.GetFormItem(1).(*tview.InputField)
	temperatureInput := node.layout.refinementInput.GetFormItem(2).(*tview.InputField)
	toppInput := node.layout.refinementInput.GetFormItem(3).(*tview.InputField)
	penaltyInput := node.layout.refinementInput.GetFormItem(4).(*tview.InputField)
	frequencyInput := node.layout.refinementInput.GetFormItem(5).(*tview.InputField)

	if !util.MatchNumber(resultInput.GetText()) {
		resultInput.SetText("\u0031")
	}

	if !util.MatchNumber(probabilityInput.GetText()) {
		probabilityInput.SetText("\u0031")
	}

	if !util.MatchNumber(toppInput.GetText()) {
		temperatureInput.SetText("\u0031\u002e\u0030")
	}

	if !util.MatchNumber(toppInput.GetText()) {
		toppInput.SetText("\u0030\u002e\u0034")
	}

	if !util.MatchNumber(penaltyInput.GetText()) {
		penaltyInput.SetText("\u0030\u002e\u0035")
	}

	if !util.MatchNumber(frequencyInput.GetText()) {
		frequencyInput.SetText("\u0030\u002e\u0035")
	}
}

// ReturnToPage - Switch to page according to their index
func ReturnToPage(index int) {
	switch index {
	case 1:
		node.layout.pages.ShowPage("console")
		node.layout.pages.HidePage("refinement")
		node.layout.pages.HidePage("training")
		node.layout.pages.HidePage("id")
	case 2:
		node.layout.pages.HidePage("console")
		node.layout.pages.ShowPage("refinement")
		node.layout.pages.HidePage("training")
		node.layout.pages.HidePage("id")
	case 3:
		node.layout.pages.HidePage("console")
		node.layout.pages.HidePage("refinement")
		node.layout.pages.ShowPage("training")
		node.layout.pages.HidePage("id")
	case 4:
		node.layout.pages.HidePage("console")
		node.layout.pages.HidePage("refinement")
		node.layout.pages.HidePage("training")
		node.layout.pages.ShowPage("id")
	default:
		node.layout.pages.ShowPage("console")
		node.layout.pages.HidePage("refinement")
		node.layout.pages.HidePage("training")
		node.layout.pages.HidePage("id")
	}
	// Clear console
	ClearConsoleView()
}

// GenerateLayoutContent - Layout content for console view
func GenerateLayoutContent() {
	// COM
	node.layout.promptOutput = tview.NewTextView()
	node.layout.promptArea = tview.NewTextArea()
	// Metadata
	node.layout.metadataOutput = tview.NewTextView()
	// Info
	node.layout.infoOutput = tview.NewTextView()
	//Output
	node.layout.infoOutput.
		SetToggleHighlights(true).
		SetLabel("Request: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOrange.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
	node.layout.metadataOutput.
		SetToggleHighlights(true).
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
	node.layout.promptOutput.
		SetToggleHighlights(true).
		SetLabel("Response: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true)
	// Input
	node.layout.promptArea.
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkOrange).
		SetTitle("Input").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorDarkCyan)
	// List models
	node.controller.ListModels()
	// Add types
	node.controller.currentAgent.preferences.Roles = append(node.controller.currentAgent.preferences.Roles, string(model.User))
	node.controller.currentAgent.preferences.Roles = append(node.controller.currentAgent.preferences.Roles, string(model.Assistant))
	node.controller.currentAgent.preferences.Roles = append(node.controller.currentAgent.preferences.Roles, string(model.System))
}

// CreateConsoleView - Create console view page
func CreateConsoleView() bool {
	// help
	helpOutput := tview.NewTextView()
	helpOutput.
		SetText("Press CTRL+SPACE or CMD+SPACE to send the prompt.\nPress CTRL+C or CMD+Q to exit from the application.\nGo to fullscreen for advanced options.").
		SetTextAlign(tview.AlignRight)
	// Layout
	node.layout.detailsInput = tview.NewForm()
	metadataSection := tview.NewFlex()
	infoSection := tview.NewFlex()
	comSection := tview.NewFlex()
	// Console section
	node.layout.promptArea.
		SetBorderPadding(1, 2, 2, 4)
	node.layout.detailsInput.
		AddTextView("Mode", "", 12, 2, true, false).
		AddDropDown("Engine", node.controller.currentAgent.preferences.Models, 11, OnChangeEngine).
		AddDropDown("Role", node.controller.currentAgent.preferences.Roles, 1, OnChangeRoles).
		AddDropDown("Words", []string{"\u0031", "\u0035\u0030", "\u0038\u0035", "\u0031\u0030\u0030", "\u0032\u0030\u0030", "\u0035\u0030\u0030", "\u0031\u0030\u0030\u0030", "\u0031\u0035\u0030\u0030"}, 4, OnChangeWords).
		AddButton("Affinity", OnRefinementTopic).
		AddButton("New conversation", OnNewTopic).
		AddButton("Export conversation", OnExportTopic).
		AddButton("Export training", OnExportTrainedTopic).
		SetHorizontal(true).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignRight)
	metadataSection.
		AddItem(node.layout.metadataOutput, 0, 1, false).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkSlateGray.TrueColor()).
		SetBorderPadding(1, 2, 2, 4).
		SetTitle("Metadata").
		SetTitleColor(tcell.ColorOrange.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	infoSection.
		AddItem(node.layout.infoOutput, 0, 1, false).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(1, 2, 2, 4).
		SetTitle("Details").
		SetTitleColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	comSection.
		AddItem(node.layout.promptOutput, 0, 1, false).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkCyan.TrueColor()).
		SetBorderPadding(1, 2, 2, 4).
		SetTitle("Prompter").
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	// Console grid
	node.layout.consoleView = tview.NewGrid().
		SetRows(0, 12).
		SetColumns(0, 2).
		AddItem(metadataSection, 1, 0, 1, 1, 0, 0, false).
		AddItem(infoSection, 1, 1, 1, 4, 0, 0, false).
		AddItem(node.layout.detailsInput, 2, 0, 1, 5, 50, 0, true).
		AddItem(comSection, 3, 0, 8, 5, 0, 0, false).
		AddItem(node.layout.promptArea, 11, 0, 2, 5, 0, 0, true).
		AddItem(helpOutput, 13, 0, 1, 5, 0, 0, false)
	// Key event
	_ = node.layout.promptArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			if OnTextChange(node.layout.promptArea.GetText(), rune(event.Key())) {
				OnTextAccept(event.Key())
				node.layout.promptArea.SetText("", true)
				return nil
			}
		}
		return event
	})
	// Console
	node.layout.consoleView.
		SetBorderPadding(0, 0, 9, 9).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor())
	// Validate view
	return node.layout.consoleView != nil
}

// CreateRefinementView - Creates refinement page view
func CreateRefinementView() bool {
	// Layout
	affinitySection := tview.NewForm()
	// Affinity section
	affinitySection.
		AddInputField("Results", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Results), 5, OnTypeAccept, OnResultChange).
		AddInputField("Probabilities", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Probabilities), 5, OnTypeAccept, OnProbabilityChange).
		AddInputField("Temperature [0.0 / 1.0]", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Temperature), 5, OnTypeAccept, OnTemperatureChange).
		AddInputField("Topp [0.0 / 1.0]", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Topp), 5, OnTypeAccept, OnToppChange).
		AddInputField("Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Penalty), 5, OnTypeAccept, OnPenaltyChange).
		AddInputField("Frequency Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Frequency), 5, OnTypeAccept, OnFrequencyChange).
		AddCheckbox("Edit mode (edit and improve the previous response)", false, OnEditChecked).
		AddCheckbox("Conversational mode (on Text mode only)", false, OnConversationChecked).
		AddCheckbox("Training mode", false, OnTrainingChecked).
		AddCheckbox("Streaming mode", true, OnStreamingChecked).
		AddButton("Back to chat", OnBack).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetTitle("Improve your search criteria: ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorDarkOrange.TrueColor()).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(1, 1, 18, 18)
	// Refinement form
	node.layout.refinementInput = affinitySection
	// Affinity grid
	node.layout.affinityView = tview.NewGrid().
		SetBorders(true)
	// Affinity
	node.layout.affinityView.
		SetSize(1, 3, 25, 55).
		AddItem(affinitySection, 0, 0, 1, 2, 0, 0, true).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(12, 12, 15, 15)
	// Validate view
	return node.layout.affinityView != nil
}

// CreateTrainingModalView - Create modal view for training mode
func CreateTrainingModalView() {
	// Modal layout
	node.layout.modalInput = tview.NewModal()
	// Modal section
	node.layout.modalInput.
		SetText("Do you want to export the current conversation? Press Ok to export it or Cancel to start a new conversation.").
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				OnExportTrainedTopic()
			}
			ReturnToPage(1)
		})
}

// CreateIDModalView - Create modal view for training mode
func CreateIDModalView() {
	// Form layout
	node.layout.idInput = tview.NewForm()
	// Form section
	node.layout.idInput.
		AddInputField("Enter your name: ", node.controller.currentAgent.id, 12, func(textToCheck string, lastChar rune) bool {
			node.controller.currentAgent.id = textToCheck
			return true
		}, nil).
		AddButton("Save", func() {
			ReturnToPage(2)
		}).
		AddButton("Cancel", func() {
			ReturnToPage(1)
		}).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkCyan.TrueColor()).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetTitleColor(tcell.ColorDarkOrange.TrueColor()).
		SetBorderPadding(25, 25, 30, 30).
		SetTitleAlign(tview.AlignCenter)
}

// InitializeLayout - Create service layout for terminal session
func InitializeLayout() {
	/* Layout content */
	GenerateLayoutContent()
	// Create views
	CreateConsoleView()
	CreateRefinementView()
	CreateTrainingModalView()
	CreateIDModalView()
	// Window frame
	node.layout.pages = tview.NewPages()
	node.layout.pages.
		AddAndSwitchToPage("console", node.layout.consoleView, true).
		AddAndSwitchToPage("refinement", node.layout.affinityView, true).
		AddAndSwitchToPage("training", node.layout.modalInput, true).
		AddAndSwitchToPage("id", node.layout.idInput, true)
	// Main executor
	node.layout.app = tview.NewApplication()
	// Inline
	node.controller.currentAgent.preferences.InlineText = make(chan string)
	// Main screen
	node.layout.screen = new(tcell.Screen)
	node.layout.app.SetScreen(*node.layout.screen)
	// App terminal configuration
	node.layout.app.
		SetRoot(node.layout.pages, true).
		SetFocus(node.layout.promptArea).
		EnableMouse(true)
	// Console view
	ReturnToPage(4)
	// Validate forms
	ValidateRefinementForm()
}
