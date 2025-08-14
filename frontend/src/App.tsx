import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import ProblemSolving from './pages/ProblemSolving';
import InteractiveDiscussion from './pages/InteractiveDiscussion';
import AutTest from './pages/AutTest';
import EmotionTest from './pages/EmotionTest';
import OnlineClass from './pages/OnlineClass';
import ChangePassword from './pages/ChangePassword';
import './App.css';

const theme = createTheme({
  palette: {
    primary: {
      main: '#7A62F5',
    },
    secondary: {
      main: '#DAD3FF',
    },
    background: {
      default: 'rgb(228, 235, 253)',
    },
  },
  typography: {
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
    h4: {
      fontWeight: 700,
      color: '#333333',
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: '20px',
          padding: '12px 24px',
          fontSize: '1.1rem',
          fontWeight: 500,
          textTransform: 'none',
          boxShadow: '0 4px 24.4px rgba(40, 5, 217, 0.15)',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-2px)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: '20px',
          boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
        },
      },
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <div className="App">
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/change-password" element={<ChangePassword />} />
              <Route
                path="/dashboard"
                element={
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/problem-solving"
                element={
                  <ProtectedRoute>
                    <ProblemSolving />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/interactive-discussion"
                element={
                  <ProtectedRoute>
                    <InteractiveDiscussion />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/aut-test"
                element={
                  <ProtectedRoute>
                    <AutTest />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/emotion-test"
                element={
                  <ProtectedRoute>
                    <EmotionTest />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/online-class"
                element={
                  <ProtectedRoute>
                    <OnlineClass />
                  </ProtectedRoute>
                }
              />
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </div>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
