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

const OnlineClass: React.FC = () => {
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
          在线微课程场景下智能化测评功能选择
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            在线课程
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              此功能正在开发中，将包含以下特性：
            </Typography>
            <Typography variant="body2" color="text.secondary" component="div" textAlign="left">
              • 在线微课程播放<br/>
              • 学习进度跟踪<br/>
              • 互动练习和测验<br/>
              • 学习效果评估<br/>
              • 个性化学习路径
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

export default OnlineClass;
