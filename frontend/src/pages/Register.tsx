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
  Grid,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';

const Register: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    school: '',
    student_id: '',
    grade: '',
    password: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setIsLoading(true);

    // 验证密码
    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致');
      setIsLoading(false);
      return;
    }

    if (formData.password.length < 6) {
      setError('密码长度至少6位');
      setIsLoading(false);
      return;
    }

    // 验证必填字段
    if (!formData.username || !formData.school || !formData.student_id || !formData.grade) {
      setError('请填写所有必填字段');
      setIsLoading(false);
      return;
    }

    try {
      console.log('Attempting registration with data:', {
        username: formData.username,
        school: formData.school,
        student_id: formData.student_id,
        grade: formData.grade,
        password: '***'
      });
      
      const success = await register({
        username: formData.username,
        school: formData.school,
        student_id: formData.student_id,
        grade: formData.grade,
        password: formData.password
      });
      
      if (success) {
        setSuccess('注册成功！正在跳转到登录页面...');
        setTimeout(() => {
          navigate('/login');
        }, 2000);
      } else {
        setError('注册失败，请检查网络连接或联系管理员');
      }
    } catch (err) {
      console.error('Registration error:', err);
      setError('注册失败，请检查网络连接或联系管理员');
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
            用户注册
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
            <Grid container spacing={2}>
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="用户名"
                  name="username"
                  variant="outlined"
                  value={formData.username}
                  onChange={handleChange}
                  className="form-field"
                  required
                />
              </Grid>
              
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="学校"
                  name="school"
                  variant="outlined"
                  value={formData.school}
                  onChange={handleChange}
                  className="form-field"
                  required
                />
              </Grid>
              
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="学号"
                  name="student_id"
                  variant="outlined"
                  value={formData.student_id}
                  onChange={handleChange}
                  className="form-field"
                  required
                />
              </Grid>
              
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="年级"
                  name="grade"
                  variant="outlined"
                  value={formData.grade}
                  onChange={handleChange}
                  className="form-field"
                  required
                  placeholder="如：大一、大二、大三、大四"
                />
              </Grid>
              
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="密码"
                  name="password"
                  type="password"
                  variant="outlined"
                  value={formData.password}
                  onChange={handleChange}
                  className="form-field"
                  required
                  helperText="密码长度至少6位"
                />
              </Grid>
              
              <Grid item xs={12} sm={12}>
                <TextField
                  fullWidth
                  label="确认密码"
                  name="confirmPassword"
                  type="password"
                  variant="outlined"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  className="form-field"
                  required
                />
              </Grid>
            </Grid>
            
            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={isLoading}
              className="submit-button"
              sx={{ mt: 2 }}
            >
              {isLoading ? '注册中...' : '注册'}
            </Button>
          </form>
          
          <Box mt={2} textAlign="center">
            <Typography variant="body2" color="textSecondary" mb={1}>
              已有账号？
            </Typography>
            <Link to="/login" style={{ textDecoration: 'none' }}>
              <Button variant="text" color="primary" size="medium">
                立即登录
              </Button>
            </Link>
          </Box>
        </Paper>
      </Container>
    </div>
  );
};

export default Register;
