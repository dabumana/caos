package service

import "caos/model"

func Constructor() {
	node.controller.currentAgent.templateProperties = node.controller.currentAgent.SetTemplateParameters("", model.ChainPrompt{}, .0)
}

/*
func (c *Chain) onStageAssemble()       {}
func (c *Chain) onStageImplementation() {}
func (c *Chain) onStageValidation()     {}
*/
