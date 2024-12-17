export type AiInstructionData = {
  instruction: string;
  instructionImage: string;
};

export type AiInstructionDto = {
  instructions: AiInstructionData[];
};

export type InstructionDto = {
  instructions: string[];
};
