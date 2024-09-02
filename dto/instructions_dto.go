package dto

type InstructionData struct {
	Instruction      string `json:"instruction"`
	InstructionImage string `json:"instruction_image"`
}

type InstructionDto struct {
	Instructions []InstructionData `json:"instructions"`
}
