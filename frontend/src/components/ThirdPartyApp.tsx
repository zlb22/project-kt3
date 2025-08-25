import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Button,
  Typography,
  Box,
  IconButton,
  CircularProgress,
} from '@mui/material';
import { ArrowBack } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';

interface ThirdPartyAppProps {
  onBack: () => void;
  appUrl: string;
  appName: string;
  description: string;
  openInNewWindow?: boolean;
  authMethod?: 'url-token' | 'postmessage' | 'none';
}

const ThirdPartyApp: React.FC<ThirdPartyAppProps> = ({
  onBack,
  appUrl,
  appName,
  description,
  openInNewWindow = false,
  authMethod = 'url-token'
}) => {
  const { token } = useAuth();
  const [loading, setLoading] = useState(true);
  const [iframeUrl, setIframeUrl] = useState('');

  useEffect(() => {
    if (authMethod === 'url-token' && token) {
      // 构建带token的URL
      const separator = appUrl.includes('?') ? '&' : '?';
      const urlWithToken = `${appUrl}${separator}token=${token}`;
      setIframeUrl(urlWithToken);
    } else {
      setIframeUrl(appUrl);
    }
    setLoading(false);
  }, [appUrl, token, authMethod]);

  const handleOpenInNewWindow = () => {
    if (iframeUrl) {
      window.open(iframeUrl, '_blank');
    }
  };

  const handleIframeLoad = () => {
    setLoading(false);
  };

  if (openInNewWindow) {
    return (
      <div className="page-container">
        <IconButton
          className="back-button"
          onClick={onBack}
          sx={{ color: 'primary.main' }}
        >
          <ArrowBack />
        </IconButton>

        <Container maxWidth="lg">
          <Typography variant="h4" component="h1" className="title-text">
            {appName}
          </Typography>

          <Paper className="button-container" elevation={3}>
            <Typography variant="h5" component="h2" textAlign="center" mb={4}>
              {description}
            </Typography>

            <Box textAlign="center" mb={4}>
              <Typography variant="body1" color="text.secondary" paragraph>
                点击下方按钮在新窗口中打开应用：
              </Typography>
            </Box>

            <Box className="navigation-buttons" display="flex" flexDirection="column" alignItems="center" gap={2}>
              <Button
                variant="contained"
                color="primary"
                size="large"
                onClick={handleOpenInNewWindow}
                sx={{ minWidth: 200 }}
              >
                打开 {appName}
              </Button>
              <Button
                variant="outlined"
                color="primary"
                onClick={onBack}
                sx={{ minWidth: 200 }}
              >
                返回
              </Button>
            </Box>
          </Paper>
        </Container>
      </div>
    );
  }

  return (
    <div className="page-container">
      <IconButton
        className="back-button"
        onClick={onBack}
        sx={{ color: 'primary.main' }}
      >
        <ArrowBack />
      </IconButton>

      <Container maxWidth="lg" sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Typography variant="h4" component="h1" className="title-text" sx={{ mb: 2 }}>
          {appName}
        </Typography>

        <Paper 
          elevation={3} 
          sx={{ 
            flex: 1, 
            display: 'flex', 
            flexDirection: 'column',
            position: 'relative',
            minHeight: '600px'
          }}
        >
          {loading && (
            <Box 
              display="flex" 
              justifyContent="center" 
              alignItems="center" 
              position="absolute"
              top="50%"
              left="50%"
              sx={{ transform: 'translate(-50%, -50%)' }}
            >
              <CircularProgress />
              <Typography variant="body1" sx={{ ml: 2 }}>
                正在加载 {appName}...
              </Typography>
            </Box>
          )}
          
          {iframeUrl && (
            <iframe
              src={iframeUrl}
              style={{
                width: '100%',
                height: '100%',
                border: 'none',
                borderRadius: '4px',
                opacity: loading ? 0 : 1,
                transition: 'opacity 0.3s ease'
              }}
              onLoad={handleIframeLoad}
              title={appName}
              sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
            />
          )}
        </Paper>
      </Container>
    </div>
  );
};

export default ThirdPartyApp;
