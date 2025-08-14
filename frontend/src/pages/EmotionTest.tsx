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

const EmotionTest: React.FC = () => {
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
          互动讨论场景下学习与抗挫情绪智能化测评
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            情绪测评
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              此功能正在开发中，将包含以下特性：
            </Typography>
            <Typography variant="body2" color="text.secondary" component="div" textAlign="left">
              • 情绪识别和分析<br/>
              • 抗挫折能力评估<br/>
              • 学习动机测试<br/>
              • 情绪调节能力<br/>
              • 心理健康报告
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

export default EmotionTest;
