import type { Nullable } from '@/models/extra.ts'

export type User = {
  id: string;
  email: string;
  name: string;
  password: string;
  avatarUrl: string;
  enabledPushNotifications: boolean;
  enabledEmailNotifications: boolean;
  createdAt: Date;
  birthday: Nullable<Date>;
}
