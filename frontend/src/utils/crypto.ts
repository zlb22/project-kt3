import bcrypt from 'bcryptjs';

/**
 * Hash password on frontend before transmission
 * Uses bcrypt with salt for secure password hashing
 */
export const hashPassword = async (password: string): Promise<string> => {
  const saltRounds = 12;
  return await bcrypt.hash(password, saltRounds);
};

/**
 * Verify password against hash (for client-side validation if needed)
 */
export const verifyPassword = async (password: string, hash: string): Promise<boolean> => {
  return await bcrypt.compare(password, hash);
};

/**
 * Generate a secure random salt
 */
export const generateSalt = async (): Promise<string> => {
  return await bcrypt.genSalt(12);
};
