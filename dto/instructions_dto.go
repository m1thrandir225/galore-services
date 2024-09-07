package dto

type AiInstructionData struct {
	Instruction      string `json:"instruction"`
	InstructionImage string `json:"instruction_image"`
}

type AiInstructionDto struct {
	Instructions []AiInstructionData `json:"instructions"`
}

type InstructionData struct {
	Instruction string `json:"instruction" binding:"required"`
}

type InstructionDto struct {
	Instructions []InstructionData `json:"instructions" binding:"required"`
}