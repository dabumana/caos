// Package service section
package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"caos/service/parameters"
	"caos/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var group sync.WaitGroup

// Layout - Recreates the terminal definitions and parameters for a console app
type Layout struct {
	app   *tview.Application
	pages *tview.Pages
	// Flex
	consoleView  *tview.Grid
	affinityView *tview.Grid
	// User form
	refinementInput *tview.Form
	detailsInput    *tview.Form
	// User modal
	modalInput *tview.Modal
	// User input
	promptInput *tview.InputField
	// Details output
	metadataOutput *tview.TextView
	promptOutput   *tview.TextView
	infoOutput     *tview.TextView
}

// OnResultChange - Evaluates when an input text changes for the result input field
func OnResultChange(text string) {
	parameters.Results = util.ParseInt32(text)
}

// OnProbabilityChange - Evaluates when an input text changes for the probability input field
func OnProbabilityChange(text string) {
	parameters.Probabilities = util.ParseInt32(text)
}

// OnTemperatureChange - Evaluates when an input text changes for the temperature input field
func OnTemperatureChange(text string) {
	parameters.Temperature = util.ParseFloat32(text)
}

// OnToppChange - Evaluates when an input text changes for the topp input field
func OnToppChange(text string) {
	parameters.Topp = util.ParseFloat32(text)
}

// OnPenaltyChange - Evaluates when an input text changes for the penalty input field
func OnPenaltyChange(text string) {
	parameters.Penalty = util.ParseFloat32(text)
}

// OnFrequencyChange - Evaluates when an input text changes for the frequency penalty input field
func OnFrequencyChange(text string) {
	parameters.Frequency = util.ParseFloat32(text)
}

// OnTypeAccept - Evaluates when an input text matches the field criteria
func OnTypeAccept(text string, lastChar rune) bool {
	matched := util.MatchNumber(text)
	return matched
}

// OnBack - Button event to return to the main page
func OnBack() {
	// Console view
	node.layout.pages.ShowPage("console")
	node.layout.pages.HidePage("refinement")
	node.layout.pages.HidePage("modal")
	// Validate layout forms
	ValidateRefinementForm()
}

// OnNewTopic - Define a new conversation button event
func OnNewTopic() {
	parameters.IsNewSession = true
	parameters.IsPromptReady = false
	parameters.PromptCtx = []string{""}

	mode := node.layout.detailsInput.GetFormItem(0).(*tview.TextView)
	parameters.Mode = mode.GetText(true)

	engine := node.layout.detailsInput.GetFormItem(1).(*tview.DropDown)
	_, parameters.Engine = engine.GetCurrentOption()

	ValidateRefinementForm()

	if parameters.Mode == "Edit" {
		node.layout.promptInput.SetLabel("Enter your context first: ")
	}

	if parameters.IsTraining {
		OnTrainingTopic()
		return
	}

	CleanConsoleView()
}

// CleanConsoleView - Clean input-output fields
func CleanConsoleView() {
	node.layout.infoOutput.SetText("")
	node.layout.metadataOutput.SetText("")
	node.layout.promptOutput.SetText("")
	node.layout.promptInput.SetPlaceholder("Type here...")
	node.layout.promptInput.SetText("")
}

// OnRefinementTopic - Refinement view button event
func OnRefinementTopic() {
	// Refinement view
	node.layout.pages.HidePage("console")
	node.layout.pages.ShowPage("refinement")
	node.layout.pages.HidePage("modal")
}

// OnTrainingTopic - Modal confirmation to export training
func OnTrainingTopic() {
	// Training modal view
	node.layout.pages.HidePage("console")
	node.layout.pages.HidePage("refinement")
	node.layout.pages.ShowPage("modal")
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
	raw, _ := json.MarshalIndent(parameters.TrainingSessionPool, "", "\u0009")
	out := util.ConstructPathFileTo("training", "json")
	out.WriteString(string(raw))
}

// OnChangeEngine - Dropdown from input to change engine
func OnChangeEngine(option string, optionIndex int) {
	if parameters.IsEditable && !parameters.IsPromptReady {
		node.layout.promptInput.SetLabel("Retrieve en editable context sending a prompt from the selected model: ")
	} else {
		node.layout.promptInput.SetLabel("Enter your request: ")
	}
	mode := node.layout.detailsInput.GetFormItem(0).(*tview.TextView)
	parameters.Engine = option

	if strings.Contains(option, "edit") {
		parameters.Mode = "Edit"
		if parameters.PromptCtx != nil &&
			parameters.PromptCtx[0] == "" {
			parameters.IsPromptReady = false
		}
		node.layout.promptInput.SetLabel("Enter your context first: ")
	} else if strings.Contains(option, "code") {
		parameters.Mode = "Code"
	} else if strings.Contains(option, "search") {
		parameters.Mode = "Search"
	} else if strings.Contains(option, "insert") {
		parameters.Mode = "Insert"
	} else if strings.Contains(option, "instruct") {
		parameters.Mode = "Instruct"
	} else if strings.Contains(option, "similarity") {
		parameters.Mode = "Similarity"
	} else if strings.Contains(option, "embedding") {
		parameters.Mode = "Embedded"
	} else if strings.Contains(option, "zero") {
		parameters.Mode = "Predicted"
	} else {
		parameters.Mode = "Text"
	}

	mode.SetText(parameters.Mode)
}

// OnChangeWords - Dropdown for tokens according to the amount of words
func OnChangeWords(option string, optionIndex int) {
	switch option {
	case "1":
		parameters.MaxTokens = util.ParseInt64("\u0031")
	case "50":
		parameters.MaxTokens = util.ParseInt64("\u0033\u0031")
	case "85":
		parameters.MaxTokens = util.ParseInt64("\u0036\u0034")
	case "100":
		parameters.MaxTokens = util.ParseInt64("\u0037\u0035")
	case "200":
		parameters.MaxTokens = util.ParseInt64("\u0031\u0035\u0030")
	case "500":
		parameters.MaxTokens = util.ParseInt64("\u0033\u0037\u0035")
	case "1000":
		parameters.MaxTokens = util.ParseInt64("\u0037\u0035\u0030")
	case "1500":
		parameters.MaxTokens = util.ParseInt64("\u0031\u0031\u0032\u0035")
	}
}

// OnTextAccept - Text field from input
func OnTextAccept(textToCheck string, lastChar rune) bool {
	if parameters.IsLoading {
		return false
	}

	textToCheck = strings.ReplaceAll(textToCheck, "\u000D", "\u0020")

	if parameters.Mode == "Edit" {
		if !parameters.IsPromptReady {
			parameters.PromptCtx = []string{textToCheck}
			node.agent.currentAgent.promptProperties = node.agent.currentAgent.SetPromptParameters(
				[]string{textToCheck},
				[]string{""},
				int(parameters.MaxTokens),
				int(parameters.Results),
				int(parameters.Probabilities),
			)
		} else {
			node.agent.currentAgent.promptProperties = node.agent.currentAgent.SetPromptParameters(
				parameters.PromptCtx,
				[]string{textToCheck},
				int(parameters.MaxTokens),
				int(parameters.Results),
				int(parameters.Probabilities),
			)
		}
	} else if parameters.Mode == "Predicted" {
		node.agent.currentAgent.predictProperties = node.agent.currentAgent.SetPredictionParameters(
			[]string{textToCheck},
		)
	} else {
		node.agent.currentAgent.promptProperties = node.agent.currentAgent.SetPromptParameters(
			[]string{textToCheck},
			[]string{""},
			int(parameters.MaxTokens),
			int(parameters.Results),
			int(parameters.Probabilities),
		)
	}

	node.agent.currentAgent.engineProperties = node.agent.currentAgent.SetEngineParameters(
		parameters.Engine,      // "text-davinci-003",
		parameters.Temperature, // if temperature is used set topp to 1.0
		parameters.Topp,        // if topp is used set temperature to 1.0
		parameters.Penalty,     // Penalize from 0 to 1 the repeated tokens
		parameters.Frequency,   // Frequency  of penalization
	)

	return true
}

// OnTextDone - Text key event
func OnTextDone(key tcell.Key) {
	if parameters.IsLoading {
		return
	}

	if parameters.IsNewSession {
		parameters.IsNewSession = false
	}

	if key == tcell.KeyEnter && !parameters.IsLoading {
		if parameters.Mode == "Edit" {
			if !parameters.IsPromptReady &&
				node.layout.promptInput.GetText() != "" {
				node.layout.promptInput.SetText("")
				parameters.IsPromptReady = true
			} else {
				if parameters.IsPromptReady {
					group.Add(1)
					go func() {
						parameters.IsLoading = true
						defer group.Done()
						defer node.agent.EditRequest()
					}()
				}
			}
		} else if parameters.Mode == "Embedded" {
			group.Add(1)
			go func() {
				parameters.IsLoading = true
				defer group.Done()
				defer node.agent.EmbeddingRequest()
			}()
		} else if parameters.Mode == "Predicted" {
			if !parameters.IsEditable {
				group.Add(1)
				go func() {
					parameters.IsLoading = true
					defer group.Done()
					defer node.agent.PredictableRequest()
				}()
			} else {
				node.layout.infoOutput.SetText("\nUncheck edit mode in affinity preferences and try again...")
				node.layout.promptInput.SetText("")
			}
		} else {
			group.Add(1)
			go func() {
				parameters.IsLoading = true
				defer group.Done()
				defer node.agent.CompletionRequest()
			}()
		}
	}

	if parameters.IsEditable &&
		parameters.Mode != "Edit" {
		parameters.Mode = "Edit"
		parameters.Engine = "text-davinci-edit-001"
		parameters.IsPromptReady = true
	}

	if parameters.IsPromptReady {
		node.layout.promptInput.SetLabel("Enter your request: ")
	}

	group.Wait()
}

// OnEditChecked - Editable mode activation
func OnEditChecked(state bool) {
	parameters.IsEditable = state
	if parameters.IsEditable {
		node.layout.promptInput.SetLabel("Retrieve en editable context sending a prompt from the selected model: ")
	} else {
		node.layout.promptInput.SetLabel("Enter your request: ")
	}
}

// OnConversationChecked - Conversation mode for friendly responses
func OnConversationChecked(state bool) {
	parameters.IsConversational = state
	if parameters.IsConversational {
		parameters.Mode = "Text"
		parameters.Engine = "text-davinci-003"
	}
}

// OnTrainingChecked - Training mode to store the current conversation
func OnTrainingChecked(state bool) {
	parameters.IsTraining = state
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

// GenerateLayoutContent - Layout content for console view
func GenerateLayoutContent() {
	// COM
	node.layout.promptInput = tview.NewInputField()
	node.layout.promptOutput = tview.NewTextView()
	// Metadata
	node.layout.metadataOutput = tview.NewTextView()
	// Info
	node.layout.infoOutput = tview.NewTextView()
	// Input
	node.layout.promptInput.
		SetFieldTextColor(tcell.ColorDarkOrange.TrueColor()).
		SetAcceptanceFunc(OnTextAccept).
		SetDoneFunc(OnTextDone).
		SetLabel("Enter your request: ").
		SetLabelColor(tcell.ColorDarkOrange.TrueColor()).
		SetPlaceholder("Type here...").
		SetFieldBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetFieldTextColor(tcell.ColorWhiteSmoke.TrueColor())
	//Output
	node.layout.infoOutput.
		SetToggleHighlights(true).
		SetLabel("Request: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOrange.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true).
		SetBorderPadding(0, 1, 1, 3)
	node.layout.metadataOutput.
		SetToggleHighlights(true).
		SetLabel("Description: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true).
		SetBorderPadding(0, 1, 1, 3)
	node.layout.promptOutput.
		SetToggleHighlights(true).
		SetLabel("Response: ").
		SetScrollable(true).
		ScrollToEnd().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetRegions(true).
		SetDynamicColors(true).
		SetBorderPadding(0, 3, 1, 1)
	// Engine availables
	parameters.Models = append(parameters.Models, "zero-gpt")
	node.agent.ListModels()
}

// CreateConsoleView - Create console view page
func CreateConsoleView() bool {
	// help
	helpOutput := tview.NewTextView()
	helpOutput.
		SetText("Press CTRL + C or CMD + Q to exit from the application.\nGo to fullscreen for advanced options.").
		SetTextAlign(tview.AlignRight)
	// Layout
	node.layout.detailsInput = tview.NewForm()
	metadataSection := tview.NewFlex()
	infoSection := tview.NewFlex()
	comSection := tview.NewFlex()
	// Console section
	node.layout.detailsInput.
		AddTextView("Mode", "", 10, 2, true, false).
		AddDropDown("Engine", parameters.Models, 11, OnChangeEngine).
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
		SetBackgroundColor(tcell.ColorDarkSlateGray.TrueColor()).
		SetTitle("Metadata").
		SetTitleColor(tcell.ColorOrange.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	infoSection.
		AddItem(node.layout.infoOutput, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetTitle("Details").
		SetTitleColor(tcell.ColorDarkTurquoise.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	comSection.
		AddItem(node.layout.promptOutput, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDarkSlateGray.TrueColor()).
		SetTitle("Prompter").
		SetTitleColor(tcell.ColorOrange.TrueColor()).
		SetTitleAlign(tview.AlignLeft)
	// Console grid
	node.layout.consoleView = tview.NewGrid().
		SetRows(0, 12).
		SetColumns(0, 2).
		AddItem(metadataSection, 1, 0, 1, 1, 0, 0, false).
		AddItem(infoSection, 1, 1, 1, 4, 0, 0, false).
		AddItem(node.layout.detailsInput, 2, 0, 1, 5, 50, 0, true).
		AddItem(comSection, 3, 0, 9, 5, 0, 0, false).
		AddItem(node.layout.promptInput, 12, 0, 1, 5, 0, 0, true).
		AddItem(helpOutput, 13, 0, 1, 5, 0, 0, false)
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
		AddInputField("Results", fmt.Sprintf("%v", parameters.Results), 5, OnTypeAccept, OnResultChange).
		AddInputField("Probabilities", fmt.Sprintf("%v", parameters.Probabilities), 5, OnTypeAccept, OnProbabilityChange).
		AddInputField("Temperature [0.0 / 1.0]", fmt.Sprintf("%v", parameters.Temperature), 5, OnTypeAccept, OnTemperatureChange).
		AddInputField("Topp [0.0 / 1.0]", fmt.Sprintf("%v", parameters.Topp), 5, OnTypeAccept, OnToppChange).
		AddInputField("Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", parameters.Penalty), 5, OnTypeAccept, OnPenaltyChange).
		AddInputField("Frequency Penalty [-2.0 / 2.0]", fmt.Sprintf("%v", parameters.Frequency), 5, OnTypeAccept, OnFrequencyChange).
		AddCheckbox("Edit mode (edit and improve the previous response)", false, OnEditChecked).
		AddCheckbox("Conversational mode (on Text mode only)", false, OnConversationChecked).
		AddCheckbox("Training mode", false, OnTrainingChecked).
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
		SetBorderPadding(1, 0, 9, 9)
	// Refinement form
	node.layout.refinementInput = affinitySection
	// Affinity grid
	node.layout.affinityView = tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 1).
		AddItem(affinitySection, 0, 0, 1, 1, 0, 0, true)
	// Affinity
	node.layout.affinityView.
		SetBorderPadding(15, 15, 20, 20).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor())
	// Validate view
	return node.layout.affinityView != nil
}

// CreateModalView - Create modal view for training mode
func CreateModalView() {
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

			node.layout.pages.ShowPage("console")
			node.layout.pages.HidePage("refinement")
			node.layout.pages.HidePage("modal")

			CleanConsoleView()
		})
}

// InitializeLayout - Create service layout for terminal session
func InitializeLayout() {
	/* Layout content */
	GenerateLayoutContent()
	// Create views
	CreateConsoleView()
	CreateRefinementView()
	CreateModalView()
	// Window frame
	node.layout.pages = tview.NewPages()
	node.layout.pages.
		AddAndSwitchToPage("console", node.layout.consoleView, true).
		AddAndSwitchToPage("refinement", node.layout.affinityView, true).
		AddAndSwitchToPage("modal", node.layout.modalInput, true)
	// Main executor
	node.layout.app = tview.NewApplication()
	// App terminal configuration
	node.layout.app.
		SetRoot(node.layout.pages, true).
		SetFocus(node.layout.promptInput).
		EnableMouse(true)
	// Console view
	node.layout.pages.ShowPage("console")
	node.layout.pages.HidePage("refinement")
	node.layout.pages.HidePage("modal")
	// Validate forms
	ValidateRefinementForm()
}
