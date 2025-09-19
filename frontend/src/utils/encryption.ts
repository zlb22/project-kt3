// RSA encryption utilities using Web Crypto API

/**
 * Convert PEM public key to CryptoKey
 */
async function importRSAPublicKey(pemKey: string): Promise<CryptoKey> {
  // Remove PEM headers and decode base64
  const pemHeader = "-----BEGIN PUBLIC KEY-----";
  const pemFooter = "-----END PUBLIC KEY-----";
  const pemContents = pemKey.replace(pemHeader, "").replace(pemFooter, "").replace(/\s/g, "");
  
  // Convert base64 to ArrayBuffer
  const binaryDer = atob(pemContents);
  const bytes = new Uint8Array(binaryDer.length);
  for (let i = 0; i < binaryDer.length; i++) {
    bytes[i] = binaryDer.charCodeAt(i);
  }

  // Import the key
  return await window.crypto.subtle.importKey(
    "spki",
    bytes.buffer,
    {
      name: "RSA-OAEP",
      hash: "SHA-256",
    },
    false,
    ["encrypt"]
  );
}

/**
 * Encrypt password using RSA public key
 */
export async function encryptPassword(password: string, publicKeyPem: string): Promise<string> {
  try {
    // Import the public key
    const publicKey = await importRSAPublicKey(publicKeyPem);
    
    // Convert password to ArrayBuffer
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    
    // Encrypt the data
    const encrypted = await window.crypto.subtle.encrypt(
      {
        name: "RSA-OAEP",
      },
      publicKey,
      data
    );
    
    // Convert to base64
    const encryptedArray = new Uint8Array(encrypted);
    const encryptedBase64 = btoa(String.fromCharCode.apply(null, Array.from(encryptedArray)));
    
    return encryptedBase64;
  } catch (error) {
    console.error('Password encryption failed:', error);
    throw new Error('Password encryption failed');
  }
}

/**
 * Get public key from backend
 */
export async function getPublicKey(): Promise<string> {
  try {
    const response = await fetch('/api/auth/public-key');
    const data = await response.json();
    return data.public_key;
  } catch (error) {
    console.error('Failed to get public key:', error);
    throw new Error('Failed to get public key');
  }
}
