import React, { useState } from 'react';
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
import ThirdPartyApp from '../components/ThirdPartyApp';

const InteractiveDiscussion: React.FC = () => {
  const navigate = useNavigate();
  const [showThirdPartyApp, setShowThirdPartyApp] = useState(false);

  const handleStartDiscussion = () => {
    setShowThirdPartyApp(true);
  };

  const handleBackToSelection = () => {
    setShowThirdPartyApp(false);
  };

  if (showThirdPartyApp) {
    return (
      <ThirdPartyApp 
        onBack={handleBackToSelection}
        appUrl="http://172.24.125.63:8080/creative-solutions/"
        appName="互动讨论"
        description="创意解决方案讨论平台"
        openInNewWindow={false}
        authMethod="url-token"
      />
    );
  }

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
          在线实验场景下智能化测评功能选择
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            在线实验测评
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              此功能正在开发中，将包含以下特性：
            </Typography>
            <Typography variant="body2" color="text.secondary" component="div" textAlign="left">
              • 创意解决方案讨论<br/>
              • 在线协作实验<br/>
              • 问题解决能力评估<br/>
              • 创新思维分析<br/>
              • 实验结果报告
            </Typography>
          </Box>

          <Box className="navigation-buttons" display="flex" flexDirection="column" alignItems="center" gap={2}>
            <Button
              variant="contained"
              color="primary"
              size="large"
              onClick={handleStartDiscussion}
              sx={{ minWidth: 200 }}
            >
              开始在线实验
            </Button>
            <Button
              variant="outlined"
              color="primary"
              onClick={() => navigate('/dashboard')}
              sx={{ minWidth: 200 }}
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
