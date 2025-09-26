import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import axios from 'axios';
import { encryptPassword, getPublicKey } from '../utils/encryption';

// Configure axios base URL for API requests
const API_BASE_URL = process.env.NODE_ENV === 'production' 
  ? window.location.origin 
  : 'https://172.24.130.213:8443';

axios.defaults.baseURL = API_BASE_URL;

interface User {
  username: string;
  school: string;
  student_id: string;
  grade: string;
  is_active: boolean;
  created_at?: string;
}

interface RegisterData {
  username: string;
  school: string;
  student_id: string;
  grade: string;
  password: string;
}

interface CaptchaData {
  captcha_id: string;
  image_base64: string;
  expires_in: number;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string, captchaId: string, captchaCode: string) => Promise<boolean>;
  register: (userData: RegisterData) => Promise<boolean>;
  logout: () => void;
  getCaptcha: () => Promise<CaptchaData | null>;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [isLoading, setIsLoading] = useState(true);

  // Configure axios defaults
  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete axios.defaults.headers.common['Authorization'];
    }
  }, [token]);

  // Check if user is authenticated on app load
  useEffect(() => {
    const checkAuth = async () => {
      if (token) {
        try {
          const response = await axios.get('/api/auth/me');
          setUser(response.data);
        } catch (error) {
          // Token is invalid, remove it
          localStorage.removeItem('token');
          setToken(null);
          setUser(null);
        }
      }
      setIsLoading(false);
    };

    checkAuth();
  }, [token]);

  const login = async (username: string, password: string, captchaId: string, captchaCode: string): Promise<boolean> => {
    try {
      // Get public key and encrypt password
      const publicKey = await getPublicKey();
      const encryptedPassword = await encryptPassword(password, publicKey);
      
      const response = await axios.post('/api/auth/login', {
        username,
        password: encryptedPassword,
        captcha_id: captchaId,
        captcha_code: captchaCode,
      });

      const { access_token } = response.data;
      localStorage.setItem('token', access_token);
      setToken(access_token);
      
      // Fetch user data
      const userResponse = await axios.get('/api/auth/me', {
        headers: { Authorization: `Bearer ${access_token}` }
      });
      setUser(userResponse.data);
      
      return true;
    } catch (error) {
      console.error('Login failed:', error);
      return false;
    }
  };

  const getCaptcha = async (): Promise<CaptchaData | null> => {
    try {
      const response = await axios.get('/api/auth/captcha');
      return response.data;
    } catch (error) {
      console.error('Failed to get captcha:', error);
      return null;
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    delete axios.defaults.headers.common['Authorization'];
  };

  const register = async (userData: RegisterData): Promise<boolean> => {
    try {
      // Get public key and encrypt password
      const publicKey = await getPublicKey();
      const encryptedPassword = await encryptPassword(userData.password, publicKey);
      const secureUserData = {
        ...userData,
        password: encryptedPassword
      };
      
      const response = await axios.post('/api/auth/register', secureUserData);
      console.log('Registration response:', response.data);
      
      // 检查后端返回的success字段
      return response.data.success === true;
    } catch (error) {
      console.error('Registration failed:', error);
      if (axios.isAxiosError(error)) {
        console.error('Error response:', error.response?.data);
        console.error('Error status:', error.response?.status);
      }
      return false;
    }
  };



  const value: AuthContextType = {
    user,
    token,
    login,
    register,
    logout,
    getCaptcha,
    isLoading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
