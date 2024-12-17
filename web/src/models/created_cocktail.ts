import type { IngredientDto } from '@/models/dto/ingredient_dto.ts'
import type { AiInstructionDto } from '@/models/dto/instruction_dto.ts'

export type CreatedCocktail = {
  id: string;
  name: string;
  image: string;
  ingredients: IngredientDto;
  instructions: AiInstructionDto;
  description: string;
  userId: string;
  embedding: number[];
  createdAt: Date;
}
