import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Alert,
  Box,
  IconButton,
} from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';

const Login: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [captchaCode, setCaptchaCode] = useState('');
  const [captchaData, setCaptchaData] = useState<{ captcha_id: string; image_base64: string } | null>(null);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [captchaLoading, setCaptchaLoading] = useState(false);
  const { login, getCaptcha } = useAuth();
  const navigate = useNavigate();

  const loadCaptcha = async () => {
    setCaptchaLoading(true);
    try {
      const data = await getCaptcha();
      if (data) {
        setCaptchaData(data);
        setCaptchaCode(''); // Clear previous input
      } else {
        setError('获取验证码失败，请重试');
      }
    } catch (err) {
      setError('获取验证码失败，请重试');
    } finally {
      setCaptchaLoading(false);
    }
  };

  useEffect(() => {
    loadCaptcha();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!captchaData) {
      setError('请先获取验证码');
      return;
    }

    if (!captchaCode.trim()) {
      setError('请输入验证码');
      return;
    }

    setIsLoading(true);

    try {
      const success = await login(username, password, captchaData.captcha_id, captchaCode);
      if (success) {
        navigate('/dashboard');
      } else {
        setError('登录失败，请检查用户名、密码或验证码');
        // Refresh captcha on failure
        loadCaptcha();
      }
    } catch (err: any) {
      if (err?.response?.status === 429) {
        setError(err.response.data.detail || '账号已锁定，请稍后再试');
      } else if (err?.response?.status === 400) {
        setError(err.response.data.detail || '验证码错误或过期');
      } else if (err?.response?.status === 401) {
        setError(err.response.data.detail || '用户名或密码错误');
      } else {
        setError('登录失败，请重试');
      }
      // Refresh captcha on any failure
      loadCaptcha();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <Container maxWidth="md">
        <Box textAlign="center" mb={4}>
          <Typography variant="h4" component="h1" className="title-text">
            互动学习场景下创新能力智能化测评工具 V3.0
          </Typography>
        </Box>
        
        <Paper className="login-form" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={3}>
            用户登录
          </Typography>
          
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="用户名"
              variant="outlined"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="form-field"
              required
            />
            
            <TextField
              fullWidth
              label="密码"
              type="password"
              variant="outlined"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="form-field"
              required
            />
            
            {/* Security Notice */}
            <Box sx={{ mt: 1, mb: 1, p: 1, bgcolor: '#fff3cd', border: '1px solid #ffeaa7', borderRadius: 1 }}>
              <Typography variant="body2" color="#856404" sx={{ fontSize: '0.875rem' }}>
                ⚠️ 安全提示：连续登录失败5次将锁定账号15分钟
              </Typography>
            </Box>
            
            {/* Captcha Section */}
            <Box sx={{ mt: 2, mb: 2 }}>
              <Typography variant="body2" sx={{ mb: 1 }}>
                验证码
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <TextField
                  label="请输入验证码"
                  variant="outlined"
                  value={captchaCode}
                  onChange={(e) => setCaptchaCode(e.target.value.toUpperCase())}
                  required
                  sx={{ flex: 1 }}
                  inputProps={{ maxLength: 5 }}
                />
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  {captchaData && (
                    <img
                      src={captchaData.image_base64}
                      alt="验证码"
                      style={{
                        height: '40px',
                        border: '1px solid #ccc',
                        borderRadius: '4px',
                        cursor: 'pointer'
                      }}
                      onClick={loadCaptcha}
                      title="点击刷新验证码"
                    />
                  )}
                  <IconButton
                    onClick={loadCaptcha}
                    disabled={captchaLoading}
                    title="刷新验证码"
                    size="small"
                  >
                    <RefreshIcon />
                  </IconButton>
                </Box>
              </Box>
            </Box>
            
            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={isLoading || captchaLoading || !captchaData}
              className="submit-button"
              sx={{ mt: 2 }}
            >
              {isLoading ? '登录中...' : captchaLoading ? '加载验证码...' : '登录'}
            </Button>
          </form>
          
          {/* 分隔线和注册区域 */}
          <Box sx={{ mt: 3, mb: 2, width: '100%' }}>
            <Box sx={{ 
              display: 'flex', 
              alignItems: 'center', 
              mb: 2,
              '&::before, &::after': {
                content: '""',
                flex: 1,
                height: '1px',
                backgroundColor: '#e0e0e0'
              },
              '&::before': {
                marginRight: '16px'
              },
              '&::after': {
                marginLeft: '16px'
              }
            }}>
              <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.875rem' }}>
                还没有账号？
              </Typography>
            </Box>
            
            <Box textAlign="center">
              <Link to="/register" style={{ textDecoration: 'none' }}>
                <Button 
                  variant="outlined" 
                  color="primary"
                  fullWidth
                  sx={{
                    borderRadius: '8px',
                    padding: '12px 24px',
                    fontSize: '1rem',
                    fontWeight: 500,
                    textTransform: 'none',
                    border: '2px solid',
                    borderColor: 'primary.main',
                    backgroundColor: 'transparent',
                    color: 'primary.main',
                    '&:hover': {
                      backgroundColor: 'primary.main',
                      color: 'white',
                      transform: 'translateY(-1px)',
                      boxShadow: '0 4px 8px rgba(25, 118, 210, 0.3)'
                    },
                    transition: 'all 0.3s ease'
                  }}
                >
                  立即注册
                </Button>
              </Link>
            </Box>
          </Box>
        </Paper>
      </Container>
    </div>
  );
};

export default Login;
