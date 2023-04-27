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
}

// onResultChange - Evaluates when an input text changes for the result input field
func onResultChange(text string) {
	node.controller.currentAgent.preferences.Results = util.ParseInt32(text)
}

// onProbabilityChange - Evaluates when an input text changes for the probability input field
func onProbabilityChange(text string) {
	node.controller.currentAgent.preferences.Probabilities = util.ParseInt32(text)
}

// onTemperatureChange - Evaluates when an input text changes for the temperature input field
func onTemperatureChange(text string) {
	node.controller.currentAgent.preferences.Temperature = util.ParseFloat32(text)
}

// onToppChange - Evaluates when an input text changes for the topp input field
func onToppChange(text string) {
	node.controller.currentAgent.preferences.Topp = util.ParseFloat32(text)
}

// onPenaltyChange - Evaluates when an input text changes for the penalty input field
func onPenaltyChange(text string) {
	node.controller.currentAgent.preferences.Penalty = util.ParseFloat32(text)
}

// onFrequencyChange - Evaluates when an input text changes for the frequency penalty input field
func onFrequencyChange(text string) {
	node.controller.currentAgent.preferences.Frequency = util.ParseFloat32(text)
}

// onTemplateChange - Template dropdown selection
func onTemplateChange(option string, index int) {
	if node.controller.currentAgent.preferences.Template != index {
		node.controller.currentAgent.preferences.Template = index
		onNewTopic()
	}
}

// onTypeAccept - Evaluates when an input text matches the field criteria
func onTypeAccept(text string, lastChar rune) bool {
	matched := util.MatchNumber(text)
	return matched
}

// onBack - Button event to return to the main page
func onBack() {
	// Console view
	onConsole()
	// Validate layout forms
	validateRefinementForm()
	// Clean input
	if node.controller.currentAgent.preferences.Mode == "Edit" {
		node.layout.promptArea.SetLabel("Enter your context first: ")
	}
}

// onNewTopic - Define a new conversation button event
func onNewTopic() {
	// Local preferences
	node.controller.currentAgent.preferences.IsNewSession = true
	node.controller.currentAgent.preferences.IsPromptReady = false
	node.controller.currentAgent.preferences.PromptCtx = []string{""}
	node.controller.currentAgent.cachedPrompt = ""
	if node.layout.promptOutput.GetText(true) == "" {
		// Clear console view
		clearConsoleView()
		node.layout.infoOutput.SetText("A new conversation can be started.")
		return
	}

	OnModal()
}

// clearConsoleView - Clean input-output fields
func clearConsoleView() {
	node.layout.infoOutput.SetText("")
	node.layout.metadataOutput.SetText("")
	node.layout.promptOutput.SetText("")
	node.layout.promptArea.SetPlaceholder("Type here...")
	node.layout.promptArea.SetText("", true)
	// Flush training historial
	node.controller.FlushEvents()
}

// onConsole - Console view event
func onConsole() {
	// Console view
	returnToPage(1)
}

// onRefinement - Refinement view event
func onRefinement() {
	// Refinement view
	returnToPage(2)
}

// OnModal - Modal confirmation to export training
func OnModal() {
	// Training modal view
	returnToPage(3)
}

// onProfile - Profile view event
func onProfile() {
	returnToPage(4)
}

// onExportTopic - Export current conversation as a file .txt
func onExportTopic() {
	if node.layout.promptOutput.GetText(true) == "" {
		node.layout.infoOutput.SetText("No converstaion started yet...")
		return
	}
	// Path constructor
	out := util.ConstructTsPathFileTo("export", "txt")
	out.WriteString(node.layout.promptOutput.GetText(true))
}

// onExportTrainedTopic - Export current conversation as a trained model as a .json file
func onExportTrainedTopic() {
	if node.layout.promptOutput.GetText(true) == "" {
		node.layout.infoOutput.SetText("You don't have any interaction to be exported...")
		return
	}
	// Event training
	var event EventManager
	event.ExportTraining(node.controller.events.pool.TrainingSession)
	// Clear console
	clearConsoleView()
	node.layout.infoOutput.SetText("Training session exported, you can continue with a new conversation.")
}

// onChangeRoles - Dropdown from input to change role
func onChangeRoles(option string, optionIndex int) {
	if strings.Contains(option, string(model.User)) {
		node.controller.currentAgent.preferences.Role = model.User
	} else if strings.Contains(option, string(model.Assistant)) {
		node.controller.currentAgent.preferences.Role = model.Assistant
	} else if strings.Contains(option, string(model.System)) {
		node.controller.currentAgent.preferences.Role = model.System
	}
}

// onChangeEngine - Dropdown from input to change engine
func onChangeEngine(option string, optionIndex int) {
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
	} else if strings.Contains(option, "code") || strings.Contains(option, "curie") {
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
	} else if strings.Contains(option, "text") {
		node.controller.currentAgent.preferences.Mode = "Text"
	} else if strings.Contains(option, "zero") {
		node.controller.currentAgent.preferences.Mode = "Predicted"
		node.layout.promptArea.SetLabel("Enter the text that you want to analyze for AI plagiarism: ")
	} else {
		node.controller.currentAgent.preferences.Mode = "NOT_SUPPORTED"
	}

	mode.SetText(node.controller.currentAgent.preferences.Mode)
	if node.controller.currentAgent.preferences.Mode != "Edit" {
		onNewTopic()
	}
}

// onChangeCharacter - Dropdown for tokens according to the amount of words
func onChangeCharacter(option string, optionIndex int) {
	node.controller.currentAgent.preferences.MaxTokens = util.CalcTokenPrompt(int(util.ParseInt64(option)))
}

// onTextChange - Text field from input
func onTextChange(textToCheck string, lastChar rune) bool {
	if node.controller.currentAgent.preferences.IsLoading {
		return false
	}

	textToCheck = strings.ReplaceAll(textToCheck, "\u000D", "\u0020")

	input := []string{textToCheck}
	ctx := []string{""}

	if !node.controller.currentAgent.preferences.IsPromptReady &&
		node.controller.currentAgent.preferences.Mode == "Predicted" {
		node.controller.currentAgent.predictProperties = node.controller.currentAgent.SetPredictionParameters(
			input,
		)
	}

	if !node.controller.currentAgent.preferences.IsPromptReady &&
		node.controller.currentAgent.preferences.Mode != "Predicted" {
		node.controller.currentAgent.preferences.PromptCtx = input
	} else if node.controller.currentAgent.preferences.IsPromptReady &&
		node.controller.currentAgent.preferences.Mode == "Edit" {
		input = node.controller.currentAgent.preferences.PromptCtx
		ctx = []string{textToCheck}
	}

	node.controller.currentAgent.engineProperties = node.controller.currentAgent.SetEngineParameters(
		node.controller.currentAgent.id,
		node.controller.currentAgent.preferences.Engine, // "text-davinci-003",
		node.controller.currentAgent.preferences.Role,
		node.controller.currentAgent.preferences.Temperature, // if temperature is used set topp to 1.0
		node.controller.currentAgent.preferences.Topp,        // if topp is used set temperature to 1.0
		node.controller.currentAgent.preferences.Penalty,     // Penalize from 0 to 1 the repeated tokens
		node.controller.currentAgent.preferences.Frequency,   // Frequency  of penalization
	)

	node.controller.currentAgent.promptProperties = node.controller.currentAgent.SetPromptParameters(
		input,
		ctx,
		int(node.controller.currentAgent.preferences.MaxTokens),
		int(node.controller.currentAgent.preferences.Results),
		int(node.controller.currentAgent.preferences.Probabilities),
	)

	return true
}

// onTextAccept - Text key event
func onTextAccept(key tcell.Key) {
	if node.controller.currentAgent.preferences.IsLoading {
		return
	}

	if node.controller.currentAgent.preferences.IsNewSession {
		node.controller.currentAgent.preferences.IsNewSession = false
	}

	tokenConsumption := util.MatchToken(node.controller.currentAgent.preferences.PromptCtx)
	node.controller.currentAgent.preferences.MaxTokens += tokenConsumption
	node.controller.currentAgent.promptProperties.MaxTokens = int(node.controller.currentAgent.preferences.MaxTokens)

	if key == tcell.KeyCtrlSpace &&
		!node.controller.currentAgent.preferences.IsLoading {
		group.Add(1)
		go func() {
			defer group.Done()
			if node.controller.currentAgent.preferences.Mode == "Edit" &&
				node.controller.currentAgent.preferences.IsPromptReady {
				node.controller.EditRequest()
			} else if node.controller.currentAgent.preferences.Mode == "Embedded" &&
				!node.controller.currentAgent.preferences.IsEditable {
				node.controller.EmbeddingRequest()
			} else if node.controller.currentAgent.preferences.Mode == "Predicted" &&
				!node.controller.currentAgent.preferences.IsEditable {
				node.controller.PredictableRequest()
			} else if node.controller.currentAgent.preferences.Mode == "Turbo" {
				node.controller.ChatCompletionRequest()
			} else {
				if node.controller.currentAgent.preferences.Mode != "Edit" &&
					!node.controller.currentAgent.preferences.IsPromptReady {
					node.controller.CompletionRequest()
				}
			}
			node.controller.currentAgent.preferences.IsLoading = false
		}()

		group.Wait()
		node.controller.currentAgent.cachedPrompt = node.layout.promptOutput.GetText(true)

		if node.controller.currentAgent.preferences.IsEditable ||
			(node.controller.currentAgent.preferences.Mode == "Edit" &&
				node.controller.currentAgent.preferences.PromptCtx != nil) {
			node.controller.currentAgent.preferences.Engine = "text-davinci-edit-001"
			node.controller.currentAgent.preferences.IsPromptReady = true
			engine := node.layout.detailsInput.GetFormItem(1).(*tview.DropDown)
			engine.SetCurrentOption(validateSelector(node.controller.currentAgent.preferences.Engine))
			node.layout.promptArea.SetLabel("Enter your request: ")
		}
	}
}

// onEditChecked - Editable mode activation
func onEditChecked(state bool) {
	node.controller.currentAgent.preferences.IsEditable = state
	if node.controller.currentAgent.preferences.IsEditable {
		node.layout.promptArea.SetLabel("Retrieve en editable context sending a prompt from the selected model: ")
	} else {
		node.layout.promptArea.SetLabel("Enter your request: ")
	}
	// new topic event
	onNewTopic()
}

// onStreamingChecked - Streaming mode
func onStreamingChecked(state bool) {
	node.controller.currentAgent.preferences.IsPromptStreaming = state
}

// validateSelector - Validate the selected engine
func validateSelector(engine string) int {
	var index int
	for i := range node.controller.currentAgent.preferences.Models {
		if engine == node.controller.currentAgent.preferences.Models[i] {
			index = i
		}
	}
	return index
}

// validateRefinementForm - Service layout functionality
func validateRefinementForm() {
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

// returnToPage - Switch to page according to their index
func returnToPage(index int) {
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
	}
}

// generateLayoutContent - Layout content for console view
func generateLayoutContent() {
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

// createConsoleSections - Create sections for console view
func createConsoleSections() (*tview.Flex, *tview.Flex, *tview.Flex) {
	// Layout
	metadataSection := tview.NewFlex()
	infoSection := tview.NewFlex()
	comSection := tview.NewFlex()
	// Initialize sections
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
	return metadataSection, infoSection, comSection
}

// createConsoleView - Create console view page
func createConsoleView() bool {
	// help
	helpOutput := tview.NewTextView()
	helpOutput.
		SetText("Press CTRL+SPACE or CMD+SPACE to send the prompt.\nPress CTRL+C or CMD+Q to exit from the application.\nGo to fullscreen for advanced options.").
		SetTextAlign(tview.AlignRight)
	// Layout
	node.layout.detailsInput = tview.NewForm()
	// Console section
	node.layout.promptArea.
		SetBorderPadding(1, 2, 2, 4)
	node.layout.detailsInput.
		AddTextView("Mode", "", 15, 2, true, false).
		AddDropDown("Engine", node.controller.currentAgent.preferences.Models, 11, onChangeEngine).
		AddDropDown("Role", node.controller.currentAgent.preferences.Roles, 1, onChangeRoles).
		AddDropDown("Characters",
			[]string{
				"\u0032\u0035\u0030",
				"\u0035\u0030\u0030",
				"\u0031\u0030\u0030\u0030",
				"\u0032\u0035\u0030\u0030",
				"\u0035\u0030\u0030\u0030",
				"\u0037\u0035\u0030\u0030",
				"\u0031\u0030\u0030\u0030\u0030",
				"\u0031\u0035\u0030\u0030\u0030",
				"\u0032\u0030\u0030\u0030\u0030"},
			3,
			onChangeCharacter).
		AddDropDown("Template", node.controller.currentAgent.templateID, 0, onTemplateChange).
		AddButton("Affinity", onRefinement).
		AddButton("New conversation", onNewTopic).
		AddButton("Export conversation", onExportTopic).
		AddButton("Export training", onExportTrainedTopic).
		AddButton("Profile", onProfile).
		SetHorizontal(true).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignRight)
	// Create sections
	metadataSection, infoSection, comSection := createConsoleSections()
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
			if onTextChange(node.layout.promptArea.GetText(), rune(event.Key())) {
				onTextAccept(event.Key())
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

// createRefinementView - Creates refinement page view
func createRefinementView() bool {
	// Layout
	affinitySection := tview.NewForm()
	// Affinity section
	affinitySection.
		AddInputField("Results: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Results), 5, onTypeAccept, onResultChange).
		AddInputField("Probabilities: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Probabilities), 5, onTypeAccept, onProbabilityChange).
		AddInputField("Temperature [0.0 / 1.0]: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Temperature), 5, onTypeAccept, onTemperatureChange).
		AddInputField("Topp [0.0 / 1.0]: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Topp), 5, onTypeAccept, onToppChange).
		AddInputField("Penalty [-2.0 / 2.0]: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Penalty), 5, onTypeAccept, onPenaltyChange).
		AddInputField("Frequency Penalty [-2.0 / 2.0]: ", fmt.Sprintf("%v", node.controller.currentAgent.preferences.Frequency), 5, onTypeAccept, onFrequencyChange).
		AddCheckbox("Edit mode (edit and improve the previous response)", false, onEditChecked).
		AddCheckbox("Streaming mode (on Text and Turbo mode only)", true, onStreamingChecked).
		AddButton("Back to chat", onBack).
		SetFieldBackgroundColor(tcell.ColorDarkGrey.TrueColor()).
		SetButtonBackgroundColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(tcell.ColorDarkCyan.TrueColor()).
		SetTitle("Improve your search criteria: ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorDarkOrange.TrueColor()).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(5, 5, 25, 25)
	// Refinement form
	node.layout.refinementInput = affinitySection
	// Affinity grid
	node.layout.affinityView = tview.NewGrid()
	// Affinity
	node.layout.affinityView.
		SetSize(1, 3, 0, 0).
		AddItem(affinitySection, 0, 0, 1, 3, 0, 0, true).
		SetBorder(true).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetBorderColor(tcell.ColorDarkSlateGrey.TrueColor()).
		SetTitleColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetBorderPadding(12, 12, 24, 24)
	// Validate view
	return node.layout.affinityView != nil
}

// createModalView - Create modal view for training mode
func createModalView() {
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
				onExportTrainedTopic()
			}

			if buttonLabel == "Cancel" {
				clearConsoleView()
			}

			onConsole()
		})
}

// createIDModalView - Create modal view for training mode
func createIDModalView() {
	// Form layout
	node.layout.idInput = tview.NewForm()
	// Form section
	node.layout.idInput.
		AddInputField("Enter your name: ", node.controller.currentAgent.id, 30, func(textToCheck string, lastChar rune) bool {
			node.controller.currentAgent.id = textToCheck
			return true
		}, nil).
		AddInputField("API key: ", node.controller.currentAgent.key[0], 60, func(textToCheck string, lastChar rune) bool {
			node.controller.currentAgent.key[0] = textToCheck
			return true
		}, nil).
		AddInputField("API key [GPT-Zero]: ", node.controller.currentAgent.key[1], 60, func(textToCheck string, lastChar rune) bool {
			node.controller.currentAgent.key[1] = textToCheck
			return true
		}, nil).
		AddButton("Save", func() {
			node.controller.currentAgent.client, node.controller.currentAgent.exClient = node.controller.currentAgent.Connect()
			node.controller.currentAgent.SaveKeys()
			onRefinement()
		}).
		AddButton("Cancel", func() {
			onConsole()
		}).
		SetLabelColor(tcell.ColorDarkOliveGreen.TrueColor()).
		SetButtonsAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkCyan.TrueColor()).
		SetTitle(" C A O S - Conversational Assistant for OpenAI Services ").
		SetTitleColor(tcell.ColorDarkOrange.TrueColor()).
		SetBorderPadding(10, 10, 35, 35).
		SetTitleAlign(tview.AlignCenter)

}

// InitializeLayout - Create service layout for terminal session
func InitializeLayout() {
	/* Layout content */
	generateLayoutContent()
	// Create views
	createConsoleView()
	createRefinementView()
	createModalView()
	createIDModalView()
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
	// Initial view
	onProfile()
	// Validate forms
	validateRefinementForm()
}
