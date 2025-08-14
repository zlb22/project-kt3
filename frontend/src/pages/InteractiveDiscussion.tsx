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

const InteractiveDiscussion: React.FC = () => {
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
          互动讨论场景下智能化测评功能选择
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            互动讨论测评
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              此功能正在开发中，将包含以下特性：
            </Typography>
            <Typography variant="body2" color="text.secondary" component="div" textAlign="left">
              • 实时互动讨论场景<br/>
              • 协作能力评估<br/>
              • 沟通技巧分析<br/>
              • 团队合作表现<br/>
              • 综合能力报告
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

export default InteractiveDiscussion;
