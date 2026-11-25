// bcrypt rejects passwords longer than 72 bytes (bytes, not characters)
const MAX_PASSWORD_BYTES = 72;

export function isPasswordTooLong(password: string) {
  return new TextEncoder().encode(password).length > MAX_PASSWORD_BYTES;
}
