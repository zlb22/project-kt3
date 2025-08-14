import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Button,
  Typography,
  Grid,
  Box,
  AppBar,
  Toolbar,
  IconButton,
} from '@mui/material';
import { ExitToApp } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import axios from 'axios';

interface TestOption {
  id: string;
  name: string;
  description: string;
}

const Dashboard: React.FC = () => {
  const [tests, setTests] = useState<TestOption[]>([]);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchTests = async () => {
      try {
        const response = await axios.get('/api/tests/available');
        setTests(response.data.tests);
      } catch (error) {
        console.error('Failed to fetch tests:', error);
      }
    };

    fetchTests();
  }, []);

  const handleTestClick = (testId: string) => {
    switch (testId) {
      case 'problem_solving':
        navigate('/problem-solving');
        break;
      case 'interactive_discussion':
        navigate('/interactive-discussion');
        break;
      case 'aut_test':
        navigate('/aut-test');
        break;
      case 'emotion_test':
        navigate('/emotion-test');
        break;
      case 'online_class':
        navigate('/online-class');
        break;
      default:
        console.log(`Test ${testId} not implemented yet`);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <>
      <AppBar position="static" sx={{ backgroundColor: 'rgba(122, 98, 245, 0.9)' }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            欢迎, {user?.username}
          </Typography>
          <IconButton color="inherit" onClick={handleLogout}>
            <ExitToApp />
          </IconButton>
        </Toolbar>
      </AppBar>

      <div className="page-container">
        <Container maxWidth="lg">
          <Typography variant="h4" component="h1" className="title-text">
            互动学习场景下青少年创新能力智能化测评工具 V2.0
          </Typography>

          <Paper className="button-container" elevation={3}>
            <Typography variant="h5" component="h2" textAlign="center" mb={4}>
              测评功能选择
            </Typography>

            <Grid container spacing={3}>
              {tests.map((test) => (
                <Grid item xs={12} sm={6} md={4} key={test.id}>
                  <div className="bordered-box">
                    <Button
                      variant="contained"
                      className="test-button"
                      onClick={() => handleTestClick(test.id)}
                      sx={{
                        backgroundColor: 'rgba(122, 98, 245, 0.1)',
                        color: 'rgb(122, 98, 245)',
                        border: '1px solid #ffffff',
                        '&:hover': {
                          backgroundColor: 'rgba(122, 98, 245, 0.2)',
                        },
                      }}
                    >
                      <Box textAlign="center">
                        <Typography variant="h6" component="div" mb={1}>
                          {test.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {test.description}
                        </Typography>
                      </Box>
                    </Button>
                  </div>
                </Grid>
              ))}
            </Grid>

            <Box textAlign="center" mt={4}>
              <Button
                variant="contained"
                color="secondary"
                size="large"
                onClick={() => navigate('/change-password')}
                sx={{
                  borderRadius: '69px',
                  px: 4,
                  py: 2,
                }}
              >
                修改密码
              </Button>
            </Box>
          </Paper>
        </Container>
      </div>
    </>
  );
};

export default Dashboard;
