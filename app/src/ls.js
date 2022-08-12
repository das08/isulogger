export const LS_SECRET_KEY = "secretKey";

export function authHeaders() {
  const secretKey = localStorage.getItem(LS_SECRET_KEY);
  if (secretKey) {
    return {
      'X-Secret-Key': secretKey
    };
  } else {
    return {};
  }
}