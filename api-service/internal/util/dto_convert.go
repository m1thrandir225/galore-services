package util

import (
	"fmt"

	"github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/dto"
)

func ConvertPromptsToAiInstructionDto(prompts []dto.PromptInstruction, imageRequests []db.GenerateImageRequest) (dto.AiInstructionDto, error) {
	if len(prompts) != (len(imageRequests) - 1) {
		return dto.AiInstructionDto{}, fmt.Errorf("the two given arrays are not of the same size")
	}
	instructions := make([]dto.AiInstructionData, len(prompts))
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
			return dto.AiInstructionDto{}, fmt.Errorf("the image url is nil")
		}

		instruction := dto.AiInstructionData{
			Instruction:      prompt.Instruction,
			InstructionImage: *imageUrl,
		}
		instructions[i] = instruction
	}

	return dto.AiInstructionDto{Instructions: instructions}, nil
}
