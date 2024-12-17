import type { IngredientDto } from '@/models/dto/ingredient_dto.ts'

export type Cocktail = {
  id: string;
  name: string;
  isAlcoholic: boolean | null;
  glass: string;
  image: string;
  instructions: string;
  ingredients: IngredientDto;
  embedding: number[]; // Assuming Vector is a number array
  createdAt: Date;
}
