import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Alert,
  Box,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';

const Login: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const success = await login(username, password);
      if (success) {
        navigate('/dashboard');
      } else {
        setError('用户名或密码错误');
      }
    } catch (err) {
      setError('登录失败，请重试');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <Container maxWidth="sm">
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
            
            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={isLoading}
              className="submit-button"
              sx={{ mt: 2 }}
            >
              {isLoading ? '登录中...' : '登录'}
            </Button>
          </form>
          
          <Box textAlign="center" mt={2}>
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 2 }}>
              <Link to="/change-password" style={{ textDecoration: 'none' }}>
                <Button variant="text" color="primary">
                  修改密码
                </Button>
              </Link>
              <Link to="/register" style={{ textDecoration: 'none' }}>
                <Button variant="text" color="primary">
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
