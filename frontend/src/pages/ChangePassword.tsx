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
  IconButton,
} from '@mui/material';
import { ArrowBack } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';

const ChangePassword: React.FC = () => {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { changePassword, user } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (newPassword !== confirmPassword) {
      setError('新密码和确认密码不匹配');
      return;
    }

    if (newPassword.length < 1) {
      setError('密码不能为空');
      return;
    }

    setIsLoading(true);

    try {
      const success = await changePassword(oldPassword, newPassword);
      if (success) {
        setSuccess('密码修改成功！');
        setTimeout(() => {
          navigate(user ? '/dashboard' : '/login');
        }, 2000);
      } else {
        setError('密码修改失败，请检查旧密码是否正确');
      }
    } catch (err) {
      setError('密码修改失败，请重试');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <IconButton
        className="back-button"
        onClick={() => navigate(user ? '/dashboard' : '/login')}
        sx={{ color: 'primary.main' }}
      >
        <ArrowBack />
      </IconButton>

      <Container maxWidth="sm">
        <Box textAlign="center" mb={4}>
          <Typography variant="h4" component="h1" className="title-text">
            修改密码
          </Typography>
        </Box>
        
        <Paper className="login-form" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={3}>
            {user ? '修改密码' : '重置密码'}
          </Typography>
          
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}
          
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="旧密码"
              type="password"
              variant="outlined"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              className="form-field"
              required
            />
            
            <TextField
              fullWidth
              label="新密码"
              type="password"
              variant="outlined"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="form-field"
              required
            />

            <TextField
              fullWidth
              label="确认新密码"
              type="password"
              variant="outlined"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
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
              {isLoading ? '修改中...' : '修改密码'}
            </Button>
          </form>
          
          {!user && (
            <Box textAlign="center" mt={2}>
              <Link to="/login" style={{ textDecoration: 'none' }}>
                <Button variant="text" color="primary">
                  返回登录
                </Button>
              </Link>
            </Box>
          )}
        </Paper>
      </Container>
    </div>
  );
};

export default ChangePassword;
