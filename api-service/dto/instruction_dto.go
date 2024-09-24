package dto

type AiInstructionData struct {
	Instruction      string `json:"instruction"`
	InstructionImage string `json:"instruction_image"`
}

type AiInstructionDto struct {
	Instructions []AiInstructionData `json:"instructions"`
}

type InstructionDto struct {
	Instructions []string `json:"instructions" binding:"required"`
}
