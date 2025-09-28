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
    // 直接设置三个主要功能
    const mainTests: TestOption[] = [
      {
        id: 'problem_solving',
        name: '问题解决',
        description: '创新问题解决能力测评'
      },
      {
        id: 'interactive_discussion',
        name: '在线实验',
        description: '在线实验能力评估'
      },
      {
        id: 'online_class',
        name: '在线微课',
        description: '在线微课学习平台'
      }
    ];
    setTests(mainTests);
  }, []);

  const handleTestClick = (testId: string) => {
    switch (testId) {
      case 'problem_solving':
        navigate('/problem-solving');
        break;
      case 'interactive_discussion':
        // Direct approach: create JWT token and construct URL directly
        try {
          const token = localStorage.getItem('token');
          if (!token) {
            console.error('No authentication token found');
            navigate('/login');
            return;
          }

          // Create a new JWT token for sub-frontend access
          console.log('Creating sub-frontend access token...');
          
          axios.post('/api/auth/create-sub-token', {}, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          }).then(response => {
            const subToken = response.data.token;
            console.log('Sub-frontend token created');
            
            // Construct URL to sub-frontend
            let subFrontendUrl;
              const currentOrigin = window.location.origin;
              subFrontendUrl = `${currentOrigin}/topic-three/?token=${subToken}`;
            console.log('Redirecting to:', subFrontendUrl);
            
            window.location.href = subFrontendUrl;
          }).catch(error => {
            console.error('Failed to create sub-frontend token:', error);
            if (error.response?.status === 401) {
              navigate('/login');
            }
          });
        } catch (e) {
          console.error('Redirect to online experiment failed:', e);
        }
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
            互动学习场景下青少年创新能力智能化测评工具 V3.0
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
          </Paper>
        </Container>
      </div>
    </>
  );
};

export default Dashboard;
