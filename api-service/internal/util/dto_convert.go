package util

import (
	"fmt"

	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/db/types"
)

func ConvertPromptsToAiInstructionDto(prompts []types.PromptInstruction, imageRequests []db.GenerateImageRequest) (types.AiInstructionDto, error) {
	if len(prompts) != (len(imageRequests) - 1) {
		return types.AiInstructionDto{}, fmt.Errorf("the two given arrays are not of the same size")
	}
	instructions := make([]types.AiInstructionData, len(prompts))
	localRequests := imageRequests
	for i, img := range localRequests {
		if img.IsMain {
			localRequests = append(localRequests[:i], localRequests[i+1:]...)
			break
		}
	}

	for i, prompt := range prompts {
		imageUrl := localRequests[i].ImageUrl
		if imageUrl == nil {
			return types.AiInstructionDto{}, fmt.Errorf("the image url is nil")
		}

		instruction := types.AiInstructionData{
			Instruction:      prompt.Instruction,
			InstructionImage: *imageUrl,
		}
		instructions[i] = instruction
	}

	return types.AiInstructionDto{Instructions: instructions}, nil
}
