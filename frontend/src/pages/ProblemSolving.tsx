import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Button,
  Typography,
  Box,
  IconButton,
} from '@mui/material';
import { ArrowBack } from '@mui/icons-material';

const ProblemSolving: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="page-container">
      <IconButton
        className="back-button"
        onClick={() => navigate('/dashboard')}
        sx={{ color: 'primary.main' }}
      >
        <ArrowBack />
      </IconButton>

      <Container maxWidth="lg">
        <Typography variant="h4" component="h1" className="title-text">
          在线实验场景下青少年创新能力智能化测评工具 V2.0
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            问题解决能力测评
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              此功能正在开发中，将包含以下特性：
            </Typography>
            <Typography variant="body2" color="text.secondary" component="div" textAlign="left">
              • 在线实验场景模拟<br/>
              • 问题解决过程记录<br/>
              • 创新思维能力评估<br/>
              • 实时反馈和指导<br/>
              • 详细的能力报告
            </Typography>
          </Box>

          <Box className="navigation-buttons">
            <Button
              variant="contained"
              color="primary"
              onClick={() => navigate('/dashboard')}
            >
              返回主页
            </Button>
          </Box>
        </Paper>
      </Container>
    </div>
  );
};

export default ProblemSolving;
