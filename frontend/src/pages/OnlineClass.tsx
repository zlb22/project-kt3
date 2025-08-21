import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Button,
  Typography,
  Box,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Grid,
  Card,
  CardContent,
  CardActions,
} from '@mui/material';
import { ArrowBack, Psychology, Games } from '@mui/icons-material';
import TwentyFourGame from '../components/TwentyFourGame';

const OnlineClass: React.FC = () => {
  const navigate = useNavigate();
  const [showTestDialog, setShowTestDialog] = useState(false);
  const [selectedTest, setSelectedTest] = useState<string | null>(null);

  const handleStartTest = () => {
    setShowTestDialog(true);
  };

  const handleTestSelect = (testType: string) => {
    setSelectedTest(testType);
    setShowTestDialog(false);
  };

  const handleBackToSelection = () => {
    setSelectedTest(null);
  };

  if (selectedTest === '24-point') {
    return <TwentyFourGame onBack={handleBackToSelection} />;
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
          在线实验测试
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" component="h2" textAlign="center" mb={4}>
            在线实验
          </Typography>

          <Box textAlign="center" mb={4}>
            <Typography variant="body1" color="text.secondary" paragraph>
              选择您要参与的实验测试：
            </Typography>
          </Box>

          <Box className="navigation-buttons" display="flex" flexDirection="column" alignItems="center" gap={2}>
            <Button
              variant="contained"
              color="primary"
              size="large"
              onClick={handleStartTest}
              sx={{ minWidth: 200 }}
            >
              开始在线实验测试
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

      {/* 测试选择对话框 */}
      <Dialog 
        open={showTestDialog} 
        onClose={() => setShowTestDialog(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Typography variant="h5" component="h2" textAlign="center">
            在线实验测试
          </Typography>
        </DialogTitle>
        <DialogContent>
          <Typography variant="body1" color="text.secondary" textAlign="center" mb={3}>
            请选择您要参与的测试项目：
          </Typography>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Card elevation={2} sx={{ height: '100%' }}>
                <CardContent>
                  <Box display="flex" alignItems="center" mb={2}>
                    <Games color="primary" sx={{ mr: 1, fontSize: 30 }} />
                    <Typography variant="h6" component="h3">
                      24点游戏
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary">
                    通过观看微课视频学习24点规则，完成10道24点题目。
                    使用加减乘除运算使四个数字结果为24。
                  </Typography>
                </CardContent>
                <CardActions>
                  <Button 
                    variant="contained" 
                    color="primary" 
                    fullWidth
                    onClick={() => handleTestSelect('24-point')}
                  >
                    开始24点测试
                  </Button>
                </CardActions>
              </Card>
            </Grid>
            <Grid item xs={12} md={6}>
              <Card elevation={2} sx={{ height: '100%' }}>
                <CardContent>
                  <Box display="flex" alignItems="center" mb={2}>
                    <Psychology color="primary" sx={{ mr: 1, fontSize: 30 }} />
                    <Typography variant="h6" component="h3">
                      Fool Planning
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary">
                    规划能力测试，评估您的逻辑思维和策略规划能力。
                    （即将推出）
                  </Typography>
                </CardContent>
                <CardActions>
                  <Button 
                    variant="outlined" 
                    color="primary" 
                    fullWidth
                    disabled
                  >
                    敬请期待
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowTestDialog(false)}>
            取消
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default OnlineClass;
