export type Session = {
  id: string;
  email: string;
  refreshToken: string;
  userAgent: string;
  clientIp: string;
  isBlocked: boolean;
  expiresAt: Date;
  createdAt: Date;
}
