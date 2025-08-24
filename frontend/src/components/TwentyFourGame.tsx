import React, { useState, useEffect, useRef } from 'react';
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
  TextField,
  LinearProgress,
  Card,
  CardContent,
  Grid,
  Chip,
  Alert,
  Stepper,
  Step,
  StepLabel,
} from '@mui/material';
import {
  ArrowBack,
  PlayArrow,
  Pause,
  SkipNext,
  CheckCircle,
  Cancel,
  Timer,
  VideoLibrary,
  Videocam,
  ScreenShare,
  FiberManualRecord,
} from '@mui/icons-material';

interface Question {
  id: number;
  numbers: number[];
  answer?: string;
  attempts: number;
  timeSpent: number;
  skipped: boolean;
}

interface VideoProgress {
  videoId: number;
  title: string;
  duration: number;
  watchedTime: number;
  completed: boolean;
}

interface TwentyFourGameProps {
  onBack: () => void;
}

const TwentyFourGame: React.FC<TwentyFourGameProps> = ({ onBack }) => {
  // 游戏状态
  const [currentPhase, setCurrentPhase] = useState<'video' | 'test' | 'result'>('video');
  const [currentVideoIndex, setCurrentVideoIndex] = useState(0);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [gameStartTime, setGameStartTime] = useState<Date | null>(null);
  const [questionStartTime, setQuestionStartTime] = useState<Date | null>(null);
  
  // 按页面指定上传使用的 MinIO bucket（可按业务需要改为从配置或用户态获取）
  const UPLOAD_BUCKET = 'onlineclass';
  
  // 工具函数：格式化秒为 mm:ss
  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };
  
  // 简单表达式评估器（仅用于演示，不做严格的数字来源校验）
  const evaluateExpression = (expr: string, _numbers: number[]): boolean => {
    try {
      // 仅允许数字、空格、运算符和括号
      if (!/^[-+*/() 0-9.]+$/.test(expr)) return false;
      // 使用 Function 进行计算（注意：生产环境建议替换为更安全的解析器）
      // eslint-disable-next-line no-new-func
      const val = Function(`"use strict"; return (${expr})`)();
      if (typeof val !== 'number' || Number.isNaN(val) || !Number.isFinite(val)) return false;
      return Math.abs(val - 24) < 1e-6;
    } catch {
      return false;
    }
  };
  
  // 视频相关状态
  const [videos] = useState<VideoProgress[]>([
    { videoId: 1, title: '24点游戏基础规则', duration: 120, watchedTime: 0, completed: false },
    { videoId: 2, title: '基本运算技巧', duration: 150, watchedTime: 0, completed: false },
    { videoId: 3, title: '高级解题策略', duration: 180, watchedTime: 0, completed: false },
    { videoId: 4, title: '实战演练示例', duration: 200, watchedTime: 0, completed: false },
  ]);
  const [videoProgress, setVideoProgress] = useState<VideoProgress[]>(videos);
  const [isVideoPlaying, setIsVideoPlaying] = useState(false);
  const videoTimerRef = useRef<NodeJS.Timeout | null>(null);
  
  // 题目相关状态
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentAnswer, setCurrentAnswer] = useState('');
  const [showResult, setShowResult] = useState(false);
  const [testTimeLimit] = useState(30 * 60); // 30分钟
  const [remainingTime, setRemainingTime] = useState(testTimeLimit);
  const testTimerRef = useRef<NodeJS.Timeout | null>(null);
  
  // 结果统计
  const [totalScore, setTotalScore] = useState(0);
  const [totalVideoTime, setTotalVideoTime] = useState(0);
  const [totalTestTime, setTotalTestTime] = useState(0);

  // 视频录制相关状态
  const [isRecording, setIsRecording] = useState(false);
  const [cameraStream, setCameraStream] = useState<MediaStream | null>(null);
  const [screenStream, setScreenStream] = useState<MediaStream | null>(null);
  
  // 摄像头预览引用
  const cameraPreviewRef = useRef<HTMLVideoElement>(null);
  const [cameraRecorder, setCameraRecorder] = useState<MediaRecorder | null>(null);
  const [screenRecorder, setScreenRecorder] = useState<MediaRecorder | null>(null);
  const [cameraChunks, setCameraChunks] = useState<Blob[]>([]);
  const [screenChunks, setScreenChunks] = useState<Blob[]>([]);
  
  // 使用 useRef 来直接存储录制数据，避免异步状态更新问题
  const cameraChunksRef = useRef<Blob[]>([]);
  const screenChunksRef = useRef<Blob[]>([]);

  // 初始化题目
  useEffect(() => {
    const generateQuestions = (): Question[] => {
      const questionList: Question[] = [];
      const questionSets = [
        [1, 1, 8, 8], [2, 2, 10, 10], [3, 3, 8, 8], [4, 4, 6, 6],
        [1, 2, 3, 4], [2, 3, 4, 5], [1, 3, 4, 6], [2, 4, 6, 8],
        [1, 5, 5, 5], [3, 7, 7, 7]
      ];
      
      for (let i = 0; i < 10; i++) {
        questionList.push({
          id: i + 1,
          numbers: questionSets[i] || [Math.floor(Math.random() * 9) + 1, Math.floor(Math.random() * 9) + 1, Math.floor(Math.random() * 9) + 1, Math.floor(Math.random() * 9) + 1],
          attempts: 0,
          timeSpent: 0,
          skipped: false,
        });
      }
      return questionList;
    };
    
    setQuestions(generateQuestions());
  }, []);

  // 视频播放控制
  const handleVideoPlay = () => {
    setIsVideoPlaying(true);
    videoTimerRef.current = setInterval(() => {
      setVideoProgress(prev => {
        const updated = [...prev];
        const current = updated[currentVideoIndex];
        if (current && current.watchedTime < current.duration) {
          current.watchedTime += 1;
          if (current.watchedTime >= current.duration) {
            current.completed = true;
          }
        }
        return updated;
      });
    }, 1000);
  };

  const handleVideoPause = () => {
    setIsVideoPlaying(false);
    if (videoTimerRef.current) {
      clearInterval(videoTimerRef.current);
    }
  };

  const handleNextVideo = () => {
    handleVideoPause();
    if (currentVideoIndex < videos.length - 1) {
      setCurrentVideoIndex(currentVideoIndex + 1);
    } else {
      // 所有视频观看完毕，开始测试
      const totalWatchTime = videoProgress.reduce((sum, video) => sum + video.watchedTime, 0);
      setTotalVideoTime(totalWatchTime);
      setCurrentPhase('test');
      setGameStartTime(new Date());
      setQuestionStartTime(new Date());
      startTestTimer();
      // 自动开始录制
      startRecording();
    }
  };

  const handleSkipAllVideos = () => {
    handleVideoPause();
    // 跳过所有视频，直接开始测试
    setTotalVideoTime(0);
    setCurrentPhase('test');
    setGameStartTime(new Date());
    setQuestionStartTime(new Date());
    startTestTimer();
    // 自动开始录制
    startRecording();
  };

  // 检测浏览器类型
  const getBrowserInfo = () => {
    const userAgent = navigator.userAgent;
    const isSafari = /^((?!chrome|android).)*safari/i.test(userAgent);
    const isChrome = /chrome/i.test(userAgent) && !/edg/i.test(userAgent);
    const isFirefox = /firefox/i.test(userAgent);
    const isEdge = /edg/i.test(userAgent);
    
    return { isSafari, isChrome, isFirefox, isEdge, userAgent };
  };

  // 开始录制功能
  const startRecording = async () => {
    try {
      setIsRecording(true);
      const browserInfo = getBrowserInfo();
      console.log('🌐 浏览器信息:', browserInfo);
      
      let localCameraStream: MediaStream | null = null;
      let localScreenStream: MediaStream | null = null;
      
      // 并行请求权限，避免浏览器差异导致的阻塞问题
      const requestPermissions = async () => {
        const promises: Promise<any>[] = [];
        
        // 摄像头权限请求
        const requestCamera = async (): Promise<MediaStream | null> => {
          try {
            console.log('🔍 正在请求摄像头权限...');
            
            // Safari需要特殊处理权限检查
            if (!browserInfo.isSafari) {
              try {
                const cameraPermission = await navigator.permissions.query({ name: 'camera' as PermissionName });
                console.log('📹 摄像头权限状态:', cameraPermission.state);
                
                if (cameraPermission.state === 'denied') {
                  throw new Error('摄像头权限被拒绝，请在浏览器设置中允许摄像头访问');
                }
              } catch (permError) {
                console.warn('📹 权限查询失败，继续尝试获取流:', permError);
              }
            }
            
            // 尝试多种摄像头配置，从最宽松开始
            let cameraConstraints: MediaStreamConstraints[] = [
              { video: true, audio: false }, // 最宽松的配置
              { video: { facingMode: 'user' }, audio: false }, // 偏好前置摄像头
              { video: { width: 1280, height: 720 }, audio: false }, // 指定分辨率
            ];
            
            for (let i = 0; i < cameraConstraints.length; i++) {
              try {
                console.log(`📹 尝试摄像头配置 ${i + 1}:`, cameraConstraints[i]);
                const stream = await navigator.mediaDevices.getUserMedia(cameraConstraints[i]);
                console.log('✅ 摄像头获取成功！流状态:', stream.active, '轨道数:', stream.getTracks().length);
                
                // 验证流是否真的有效
                const videoTracks = stream.getVideoTracks();
                if (videoTracks.length === 0) {
                  throw new Error('获取到的流没有视频轨道');
                }
                
                console.log('📹 视频轨道信息:', videoTracks.map(t => ({
                  label: t.label,
                  kind: t.kind,
                  readyState: t.readyState,
                  enabled: t.enabled,
                  muted: t.muted
                })));
                
                return stream;
              } catch (constraintError) {
                console.warn(`📹 配置 ${i + 1} 失败:`, constraintError);
                if (i === cameraConstraints.length - 1) {
                  throw constraintError;
                }
              }
            }
            return null;
          } catch (cameraError) {
            console.warn('📹 摄像头权限获取失败:', cameraError);
            return null;
          }
        };
        
        // 屏幕录制权限请求
        const requestScreen = async (): Promise<MediaStream | null> => {
          try {
            console.log('🖥️ 正在请求屏幕录制权限...');
            
            // 检查浏览器是否支持屏幕录制
            if (!navigator.mediaDevices.getDisplayMedia) {
              throw new Error('浏览器不支持屏幕录制功能');
            }
            
            // Safari对屏幕录制的支持较晚，需要特殊处理
            let displayMediaOptions: MediaStreamConstraints;
            
            if (browserInfo.isSafari) {
              // Safari使用更简单的配置
              displayMediaOptions = {
                video: true,
                audio: false // Safari对音频支持不稳定
              };
            } else {
              // Chrome/Firefox/Edge使用完整配置
              displayMediaOptions = {
                video: { 
                  width: 1920, 
                  height: 1080,
                  frameRate: 30
                },
                audio: true
              };
            }
            
            console.log('🖥️ 屏幕录制配置:', displayMediaOptions);
            const stream = await navigator.mediaDevices.getDisplayMedia(displayMediaOptions);
            console.log('✅ 屏幕录制获取成功！流状态:', stream.active, '轨道数:', stream.getTracks().length);
            
            return stream;
          } catch (screenError) {
            console.warn('🖥️ 屏幕录制权限获取失败:', screenError);
            return null;
          }
        };
        
        // 根据浏览器类型决定请求策略
        if (browserInfo.isSafari) {
          // Safari: 先请求摄像头，成功后再请求屏幕录制
          console.log('🍎 Safari浏览器：顺序请求权限');
          localCameraStream = await requestCamera();
          
          // 给Safari一些时间处理摄像头权限
          await new Promise(resolve => setTimeout(resolve, 500));
          
          localScreenStream = await requestScreen();
        } else {
          // Chrome/Firefox/Edge: 并行请求权限
          console.log('🌐 非Safari浏览器：并行请求权限');
          const [cameraResult, screenResult] = await Promise.allSettled([
            requestCamera(),
            requestScreen()
          ]);
          
          localCameraStream = cameraResult.status === 'fulfilled' ? cameraResult.value : null;
          localScreenStream = screenResult.status === 'fulfilled' ? screenResult.value : null;
          
          if (cameraResult.status === 'rejected') {
            console.warn('📹 摄像头请求被拒绝:', cameraResult.reason);
          }
          if (screenResult.status === 'rejected') {
            console.warn('🖥️ 屏幕录制请求被拒绝:', screenResult.reason);
          }
        }
      };
      
      await requestPermissions();
      
      // 设置获取到的流
      if (localCameraStream) {
        setCameraStream(localCameraStream);
        
        // 设置摄像头预览
        if (cameraPreviewRef.current) {
          const v = cameraPreviewRef.current as HTMLVideoElement;
          v.srcObject = localCameraStream;
          const ensurePlay = async () => {
            try { 
              await v.play(); 
              console.log('📹 摄像头预览播放成功');
            } catch (e) { 
              console.warn('📹 预览播放被阻止，等待用户交互后重试', e); 
            }
          };
          if (v.readyState >= 2) {
            ensurePlay();
          } else {
            v.onloadedmetadata = () => { ensurePlay(); };
          }
        }
        
        // 绑定轨道生命周期事件
        (localCameraStream as MediaStream).getTracks().forEach((track: MediaStreamTrack) => {
          track.onended = () => console.warn('📹 摄像头轨道已结束', { kind: track.kind, label: track.label });
        });
        
        console.log('✅ 摄像头权限获取成功');
      }
      
      if (localScreenStream) {
        setScreenStream(localScreenStream);
        
        // 监听屏幕录制停止事件
        const screenVideoTrack = (localScreenStream as MediaStream).getVideoTracks()[0];
        if (screenVideoTrack) {
          screenVideoTrack.addEventListener('ended', () => {
            console.log('🖥️ 用户停止了屏幕分享');
          });
        }
        
        console.log('✅ 屏幕录制权限获取成功');
      }

      // 权限获取结果检查
      const permissionResults = {
        camera: !!localCameraStream,
        screen: !!localScreenStream,
        browser: browserInfo.isSafari ? 'Safari' : browserInfo.isChrome ? 'Chrome' : 'Other'
      };
      
      console.log('🔍 权限获取结果:', permissionResults);
      
      // 如果两个都失败了，显示提示但不阻止测试
      if (!localCameraStream && !localScreenStream) {
        setIsRecording(false);
        const errorMsg = `无法启动录制功能。您可以继续进行测试，但不会记录视频数据。\n\n浏览器: ${permissionResults.browser}\n\n请确保：\n1. 浏览器支持媒体录制功能\n2. 已授权摄像头和屏幕录制权限\n3. 使用HTTPS协议访问（本地开发除外）\n4. ${browserInfo.isSafari ? 'Safari用户请确保系统版本支持屏幕录制' : 'Chrome用户请检查摄像头权限设置'}`;
        alert(errorMsg);
        return;
      }
      
      // 显示获取到的权限信息
      if (localCameraStream && !localScreenStream) {
        console.warn('⚠️ 仅获取到摄像头权限，屏幕录制不可用');
        if (browserInfo.isSafari) {
          alert('Safari检测到仅获取摄像头权限。\n\n可能原因：\n1. Safari版本过低不支持屏幕录制\n2. 系统权限未开启\n3. 用户拒绝了屏幕录制权限\n\n测试将继续进行，但只会录制摄像头画面。');
        }
      } else if (!localCameraStream && localScreenStream) {
        console.warn('⚠️ 仅获取到屏幕录制权限，摄像头不可用');
        if (browserInfo.isChrome) {
          alert('Chrome检测到仅获取屏幕录制权限。\n\n可能原因：\n1. 摄像头被其他应用占用\n2. 用户拒绝了摄像头权限\n3. 摄像头硬件故障\n\n测试将继续进行，但只会录制屏幕画面。');
        }
      } else {
        console.log('✅ 成功获取摄像头和屏幕录制权限');
      }

      // 创建摄像头录制器（如果有摄像头流）
      console.log('🔍 检查摄像头流状态:', {
        cameraStreamExists: !!localCameraStream,
        cameraStreamActive: (localCameraStream as MediaStream | null)?.active,
        cameraStreamTracks: (localCameraStream as MediaStream | null)?.getTracks().length || 0
      });
      
      if (localCameraStream) {
        console.log('🔍 开始创建摄像头录制器...');
        try {
          // 尝试不同的编码格式，从最兼容的开始
          let mimeType = 'video/webm';
          if (MediaRecorder.isTypeSupported('video/webm;codecs=vp8')) {
            mimeType = 'video/webm;codecs=vp8';
          } else if (MediaRecorder.isTypeSupported('video/webm;codecs=vp9')) {
            mimeType = 'video/webm;codecs=vp9';
          } else if (MediaRecorder.isTypeSupported('video/mp4')) {
            mimeType = 'video/mp4';
          }
          
          console.log('📹 摄像头录制使用格式:', mimeType);
          console.log('📹 摄像头流状态:', (localCameraStream as MediaStream).active, '轨道数:', (localCameraStream as MediaStream).getTracks().length);

          // 仅使用视频轨道来进行录制，避免音频轨道导致的不兼容问题
          const videoTracks = (localCameraStream as MediaStream).getVideoTracks();
          const cameraRecordStream = new MediaStream(videoTracks);
                // 监控轨道结束事件，定位流为何变为 inactive
          videoTracks.forEach((track: MediaStreamTrack) => {
            track.onended = () => {
              console.warn('📹 摄像头视频轨道已结束', {
                label: track.label,
                readyState: track.readyState,
                enabled: track.enabled,
                muted: (track as any).muted
              });
            };
          });

          // 尝试不指定编码格式，让浏览器自动选择
          let cameraRecorder: MediaRecorder;
          try {
            // 首先尝试不指定任何选项，让浏览器使用默认设置
            cameraRecorder = new MediaRecorder(cameraRecordStream);
            console.log('📹 使用默认格式创建录制器成功（仅视频轨道）');
          } catch (defaultError) {
            console.warn('📹 默认格式失败，尝试指定格式:', defaultError);
            // 如果默认失败，尝试指定格式
            cameraRecorder = new MediaRecorder(cameraRecordStream, {
              mimeType: mimeType
            });
          }
          
          cameraChunksRef.current = []; // 重置录制数据
          
          cameraRecorder.ondataavailable = (event) => {
            if (event.data.size > 0) {
              cameraChunksRef.current.push(event.data);
              console.log('📹 摄像头数据收集:', event.data.size, '字节，总块数:', cameraChunksRef.current.length);
            } else {
              console.warn('⚠️ 摄像头数据为空，大小:', event.data.size);
            }
          };
          
          cameraRecorder.onstop = () => {
            setCameraChunks([...cameraChunksRef.current]);
            console.log('📹 摄像头录制停止，总数据块:', cameraChunksRef.current.length);
          };
          
          cameraRecorder.onerror = (event: Event) => {
            console.error('📹 摄像头录制错误:', event);
            const errorEvent = event as any;
            if (errorEvent.error) {
              console.error('📹 错误详情:', errorEvent.error);
            }
          };
          
          cameraRecorder.onstart = () => {
            console.log('📹 摄像头录制器启动成功');
          };
          
          cameraRecorder.onpause = () => {
            console.warn('📹 摄像头录制器被暂停');
          };
          
          cameraRecorder.onresume = () => {
            console.log('📹 摄像头录制器恢复录制');
          };

          // 启动录制：不传递 timeslice，让浏览器在 onstop 时统一输出数据（部分浏览器更稳定）
          try {
            cameraRecorder.start();
            setCameraRecorder(cameraRecorder);
            console.log('📹 摄像头录制已开始，格式:', mimeType);
            
            // 添加一个延迟检查，确保录制器真的启动了
            setTimeout(() => {
              console.log('📹 摄像头录制器启动后状态检查:', cameraRecorder.state);
              if (localCameraStream) {
                console.log('📹 摄像头流是否还活跃:', (localCameraStream as MediaStream).active);
                console.log('📹 摄像头轨道状态:', (localCameraStream as MediaStream).getTracks().map(track => ({
                  kind: track.kind,
                  enabled: track.enabled,
                  readyState: track.readyState
                })));
              }
              
              if (cameraRecorder.state === 'inactive') {
                console.error('📹 ❌ 摄像头录制器启动失败！状态仍为 inactive');
                
                // 尝试重新启动录制器
                if (localCameraStream && (localCameraStream as MediaStream).active) {
                  console.log('📹 🔄 尝试重新启动摄像头录制器...');
                  try {
                    cameraRecorder.start();
                  } catch (restartError) {
                    console.error('📹 重启失败:', restartError);
                  }
                }

                // 仍然 inactive 时，尝试使用预览视频元素的 captureStream 进行回退录制
                setTimeout(() => {
                  if (cameraRecorder.state === 'inactive' && cameraPreviewRef.current) {
                    const videoEl: any = cameraPreviewRef.current;
                    const startCanvasFallback = () => {
                      console.warn('📹 准备使用 canvas.captureStream 回退录制');
                      try {
                        const v = cameraPreviewRef.current as HTMLVideoElement;
                        // 若尚未就绪，等到有尺寸后再开始
                        const startWhenReady = () => {
                          const width = v.videoWidth || 1280;
                          const height = v.videoHeight || 720;
                          const canvas = document.createElement('canvas');
                          canvas.width = width;
                          canvas.height = height;
                          const ctx = canvas.getContext('2d');
                          if (!ctx) throw new Error('无法获取canvas上下文');
                          let rafId: number;
                          const draw = () => {
                            try { ctx.drawImage(v, 0, 0, canvas.width, canvas.height); } catch {}
                            rafId = requestAnimationFrame(draw);
                          };
                          draw();
                          const canvasStream = (canvas as any).captureStream(30) as MediaStream;
                          console.warn('📹 使用 canvas.captureStream 回退录制');
                          let canvasRecorder: MediaRecorder;
                          try {
                            canvasRecorder = new MediaRecorder(canvasStream);
                          } catch {
                            canvasRecorder = new MediaRecorder(canvasStream, { mimeType });
                          }
                          cameraChunksRef.current = [];
                          canvasRecorder.ondataavailable = (e) => {
                            if (e.data.size > 0) {
                              cameraChunksRef.current.push(e.data);
                              console.log('📹(canvas) 摄像头数据收集:', e.data.size, '字节，总块数:', cameraChunksRef.current.length);
                            }
                          };
                          canvasRecorder.onstop = () => {
                            cancelAnimationFrame(rafId);
                            setCameraChunks([...cameraChunksRef.current]);
                            console.log('📹(canvas) 摄像头录制停止，总数据块:', cameraChunksRef.current.length);
                          };
                          canvasRecorder.onerror = (ev: Event) => console.error('📹(canvas) 录制错误', ev);
                          canvasRecorder.start(1000);
                          setCameraRecorder(canvasRecorder);
                          console.log('📹(canvas) 摄像头录制已开始');
                        };
                        if ((v as HTMLVideoElement).videoWidth && (v as HTMLVideoElement).videoHeight) {
                          startWhenReady();
                        } else {
                          v.onloadeddata = () => startWhenReady();
                        }
                      } catch (canvasErr) {
                        console.error('📹 使用 canvas 回退失败:', canvasErr);
                      }
                    };
                    if (typeof videoEl.captureStream === 'function') {
                      try {
                        const ensurePreviewPlaying = async () => {
                          const ve = videoEl as HTMLVideoElement;
                          if (ve.readyState < 2) {
                            try { await ve.play(); } catch {}
                          }
                        };
                        ensurePreviewPlaying();
                        const previewStream: MediaStream = videoEl.captureStream(30);
                        const vTracks = previewStream.getVideoTracks();
                        if (!vTracks || vTracks.length === 0) {
                          throw new Error('captureStream 未返回视频轨道');
                        }
                        console.warn('📹 使用预览元素 captureStream 回退录制');
                        let fallbackRecorder: MediaRecorder;
                        try {
                          fallbackRecorder = new MediaRecorder(previewStream);
                        } catch {
                          fallbackRecorder = new MediaRecorder(previewStream, { mimeType });
                        }
                        cameraChunksRef.current = [];
                        fallbackRecorder.ondataavailable = (e) => {
                          if (e.data.size > 0) {
                            cameraChunksRef.current.push(e.data);
                            console.log('📹(fallback) 摄像头数据收集:', e.data.size, '字节，总块数:', cameraChunksRef.current.length);
                          }
                        };
                        fallbackRecorder.onstop = () => {
                          setCameraChunks([...cameraChunksRef.current]);
                          console.log('📹(fallback) 摄像头录制停止，总数据块:', cameraChunksRef.current.length);
                        };
                        fallbackRecorder.onerror = (ev: Event) => console.error('📹(fallback) 录制错误', ev);
                        fallbackRecorder.start(1000);
                        setCameraRecorder(fallbackRecorder);
                        console.log('📹(fallback) 摄像头录制已开始');
                      } catch (capErr) {
                        console.error('📹 使用 captureStream 回退失败:', capErr);
                        // 链接到 canvas 回退
                        startCanvasFallback();
                      }
                    } else {
                      console.warn('📹 浏览器不支持 video.captureStream 回退');
                      startCanvasFallback();
                    }
                  }
                }, 500);
              } else {
                console.log('📹 ✅ 摄像头录制器启动成功！状态:', cameraRecorder.state);
              }
            }, 1000);
            
          } catch (startError) {
            console.error('📹 摄像头录制器启动失败:', startError);
          }
        } catch (error) {
          console.error('📹 摄像头录制器创建失败:', error);
        }
      }

      // 创建屏幕录制器（如果有屏幕流）
      if (localScreenStream) {
        try {
          // 尝试不同的编码格式，从最兼容的开始
          let mimeType = 'video/webm';
          if (MediaRecorder.isTypeSupported('video/webm;codecs=vp8')) {
            mimeType = 'video/webm;codecs=vp8';
          } else if (MediaRecorder.isTypeSupported('video/webm;codecs=vp9')) {
            mimeType = 'video/webm;codecs=vp9';
          } else if (MediaRecorder.isTypeSupported('video/mp4')) {
            mimeType = 'video/mp4';
          }
          
          console.log('🖥️ 屏幕录制使用格式:', mimeType);
          
          const screenRecorder = new MediaRecorder(localScreenStream, {
            mimeType: mimeType
          });
          
          screenChunksRef.current = []; // 重置录制数据
          
          screenRecorder.ondataavailable = (event) => {
            if (event.data.size > 0) {
              screenChunksRef.current.push(event.data);
              console.log('🖥️ 屏幕数据收集:', event.data.size, '字节，总块数:', screenChunksRef.current.length);
            } else {
              console.warn('⚠️ 屏幕数据为空，大小:', event.data.size);
            }
          };
          
          screenRecorder.onstop = () => {
            setScreenChunks([...screenChunksRef.current]);
            console.log('🖥️ 屏幕录制停止，总数据块:', screenChunksRef.current.length);
          };
          
          screenRecorder.onerror = (event) => {
            console.error('🖥️ 屏幕录制错误:', event);
          };
          
          screenRecorder.onstart = () => {
            console.log('🖥️ 屏幕录制器启动成功');
          };

          // 启动录制，使用较短的时间间隔确保数据收集
          screenRecorder.start(500); // 改为500ms间隔
          setScreenRecorder(screenRecorder);
          console.log('🖥️ 屏幕录制已开始，格式:', mimeType);
        } catch (error) {
          console.error('🖥️ 屏幕录制器创建失败:', error);
        }
      }

      // 监听屏幕录制停止事件（用户点击停止分享）
      if (localScreenStream) {
        const screenVideoTrack = (localScreenStream as MediaStream).getVideoTracks()[0];
        if (screenVideoTrack) {
          screenVideoTrack.addEventListener('ended', () => {
            console.log('用户停止了屏幕分享');
            // 可以选择停止整个录制或仅停止屏幕录制
          });
        }
      }

    } catch (error) {
      console.error('录制启动失败:', error);
      setIsRecording(false);
      const errorMessage = error instanceof Error ? error.message : '未知错误';
      alert(`录制启动失败：${errorMessage}\n\n请检查：\n1. 浏览器权限设置\n2. 是否使用支持的浏览器（Chrome、Firefox、Edge等）\n3. 网络连接是否正常`);
    }
  };

  // 停止录制功能
  const stopRecording = () => {
    console.log('🛑 开始停止录制...');
    
    if (cameraRecorder) {
      console.log('📹 摄像头录制器当前状态:', cameraRecorder.state);
      if (cameraRecorder.state !== 'inactive') {
        console.log('📹 正在停止摄像头录制器...');
        cameraRecorder.stop();
      } else {
        console.warn('📹 摄像头录制器已经是 inactive 状态');
      }
    } else {
      console.warn('📹 摄像头录制器不存在');
    }
    
    if (screenRecorder) {
      console.log('🖥️ 屏幕录制器当前状态:', screenRecorder.state);
      if (screenRecorder.state !== 'inactive') {
        console.log('🖥️ 正在停止屏幕录制器...');
        screenRecorder.stop();
      } else {
        console.warn('🖥️ 屏幕录制器已经是 inactive 状态');
      }
    } else {
      console.warn('🖥️ 屏幕录制器不存在');
    }
    
    setIsRecording(false);
    console.log('🛑 录制器停止命令已发送');
    
    // 延迟停止媒体流，给录制器时间收集最后的数据
    setTimeout(() => {
      if (cameraStream) {
        console.log('📹 延迟停止摄像头流...');
        cameraStream.getTracks().forEach(track => track.stop());
      }
      if (screenStream) {
        console.log('🖥️ 延迟停止屏幕流...');
        screenStream.getTracks().forEach(track => track.stop());
      }
      console.log('🛑 媒体流停止完成');
    }, 2000); // 延迟2秒停止媒体流
  };

  // 上传录制的视频到服务器
  const uploadRecordedVideos = async () => {
    const testSessionId = new Date().getTime().toString(); // 使用时间戳作为会话ID
    const uploadPromises: Promise<any>[] = [];
    // 按页面动态选择要上传的 bucket（可按路由/用户/业务动态决定），此处使用常量 UPLOAD_BUCKET
    
    console.log('开始上传视频，录制状态:', {
      isRecording,
      cameraChunksLength: cameraChunks.length,
      screenChunksLength: screenChunks.length,
      cameraChunksRefLength: cameraChunksRef.current.length,
      screenChunksRefLength: screenChunksRef.current.length,
      cameraRecorder: cameraRecorder?.state,
      screenRecorder: screenRecorder?.state,
      cameraStream: !!cameraStream,
      screenStream: !!screenStream
    });
    
    // 详细检查摄像头录制数据
    if (cameraChunksRef.current.length === 0) {
      console.warn('⚠️ 摄像头录制数据为空！可能的原因：');
      console.warn('1. 摄像头权限被拒绝');
      console.warn('2. 摄像头录制器启动失败');
      console.warn('3. 录制过程中没有数据收集');
      console.warn('摄像头流状态:', cameraStream ? '存在' : '不存在');
      console.warn('摄像头录制器状态:', cameraRecorder?.state || '未创建');
    } else {
      console.log('✅ 摄像头录制数据正常，块数:', cameraChunksRef.current.length);
    }
    
    try {
      // 上传摄像头录制 - 使用 useRef 中的数据
      if (cameraChunksRef.current.length > 0) {
        const cameraBlob = new Blob(cameraChunksRef.current, { type: 'video/webm' });
        const cameraFormData = new FormData();
        cameraFormData.append('file', cameraBlob, `camera_${testSessionId}.webm`);
        cameraFormData.append('video_type', 'camera');
        cameraFormData.append('test_session_id', testSessionId);
        cameraFormData.append('bucket', UPLOAD_BUCKET);
        
        const cameraUpload = fetch('/api/upload/video', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: cameraFormData
        });
        
        uploadPromises.push(cameraUpload);
      }
      
      // 上传屏幕录制 - 使用 useRef 中的数据
      if (screenChunksRef.current.length > 0) {
        const screenBlob = new Blob(screenChunksRef.current, { type: 'video/webm' });
        const screenFormData = new FormData();
        screenFormData.append('file', screenBlob, `screen_${testSessionId}.webm`);
        screenFormData.append('video_type', 'screen');
        screenFormData.append('test_session_id', testSessionId);
        screenFormData.append('bucket', UPLOAD_BUCKET);
        
        const screenUpload = fetch('/api/upload/video', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: screenFormData
        });
        
        uploadPromises.push(screenUpload);
      }
      
      if (uploadPromises.length === 0) {
        console.log('没有录制数据需要上传');
        return { success: true, results: [], sessionId: testSessionId, message: '没有录制数据' };
      }
      
      // 等待所有上传完成
      const results = await Promise.all(uploadPromises);
      
      // 检查上传结果
      const uploadResults = await Promise.all(
        results.map(response => response.json())
      );
      
      console.log('视频上传成功:', uploadResults);
      return { success: true, results: uploadResults, sessionId: testSessionId };
      
    } catch (error) {
      console.error('视频上传失败:', error);
      return { success: false, error: error };
    }
  };

  // 测试计时器
  const startTestTimer = () => {
    testTimerRef.current = setInterval(() => {
      setRemainingTime(prev => {
        if (prev <= 1) {
          handleTestComplete();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  // 提交答案
  const handleSubmitAnswer = () => {
    if (!currentAnswer.trim()) return;
    
    const currentTime = new Date();
    const timeSpent = questionStartTime ? Math.floor((currentTime.getTime() - questionStartTime.getTime()) / 1000) : 0;
    
    setQuestions(prev => {
      const updated = [...prev];
      const current = updated[currentQuestionIndex];
      if (current) {
        current.answer = currentAnswer;
        current.attempts += 1;
        current.timeSpent = timeSpent;
        
        // 简单验证答案（这里可以后续完善验证逻辑）
        if (evaluateExpression(currentAnswer, current.numbers)) {
          setTotalScore(prev => prev + 1);
        }
      }
      return updated;
    });
    
    moveToNextQuestion();
  };

  // 跳过题目
  const handleSkipQuestion = () => {
    const currentTime = new Date();
    const timeSpent = questionStartTime ? Math.floor((currentTime.getTime() - questionStartTime.getTime()) / 1000) : 0;
    
    setQuestions(prev => {
      const updated = [...prev];
      const current = updated[currentQuestionIndex];
      if (current) {
        current.skipped = true;
        current.timeSpent = timeSpent;
      }
      return updated;
    });
    
    moveToNextQuestion();
  };

  const moveToNextQuestion = () => {
    setCurrentAnswer('');
    if (currentQuestionIndex < questions.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
      setQuestionStartTime(new Date());
    } else {
      handleTestComplete();
    }
  };

  const handleTestComplete = async () => {
    if (testTimerRef.current) {
      clearInterval(testTimerRef.current);
    }
    const endTime = new Date();
    const testDuration = gameStartTime ? Math.floor((endTime.getTime() - gameStartTime.getTime()) / 1000) : 0;
    setTotalTestTime(testDuration);
    
    // 停止录制并切到结果页
    stopRecording();
    setCurrentPhase('result');
    
    // 延迟上传，确保录制器完全停止并收集到所有数据
    setTimeout(async () => {
      console.log('🔍 延迟检查录制数据状态:');
      console.log('📹 摄像头数据块数:', cameraChunksRef.current.length);
      console.log('🖥️ 屏幕数据块数:', screenChunksRef.current.length);
      console.log('📹 摄像头录制器状态:', cameraRecorder?.state);
      console.log('🖥️ 屏幕录制器状态:', screenRecorder?.state);
      
      const uploadResult = await uploadRecordedVideos();
      if (uploadResult.success) {
        console.log('视频上传成功，会话ID:', uploadResult.sessionId);
        const videosMeta = (uploadResult.results || []).map((r: any) => ({
          type: r.video_type,
          bucket: r.bucket,
          object_name: r.object_name,
          file_path: r.file_path,
          file_size: r.file_size,
          content_type: r.content_type,
          upload_time: r.upload_time,
          test_session_id: r.test_session_id,
        }));
        await saveTestResults(videosMeta, uploadResult.sessionId);
      } else {
        console.error('视频上传失败:', uploadResult.error);
        // 即便视频失败，也提交测试结果（无视频引用）
        await saveTestResults([], undefined);
      }
    }, 2000);
  };

  // 保存测试结果到服务器（带视频元信息）
  const saveTestResults = async (videosMeta?: any[], sessionId?: string) => {
    const results: any = {
      totalScore,
      totalVideoTime,
      totalTestTime,
      questions: questions.map(q => ({
        id: q.id,
        numbers: q.numbers,
        answer: q.answer,
        attempts: q.attempts,
        timeSpent: q.timeSpent,
        skipped: q.skipped,
      })),
      timestamp: new Date().toISOString(),
    };
    if (sessionId) results.test_session_id = sessionId;
    if (videosMeta && videosMeta.length > 0) results.videos = videosMeta;
    
    try {
      const response = await fetch('/api/tests/24point/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(results)
      });
      if (response.ok) {
        const result = await response.json();
        console.log('测试结果上传成功:', result);
      } else {
        console.error('测试结果上传失败:', response.statusText);
      }
    } catch (error) {
      console.error('测试结果上传失败:', error);
    }
    // 本地备份
    localStorage.setItem('twentyFourGameResults', JSON.stringify(results));
  };

  // 清理定时器和录制资源
  useEffect(() => {
    return () => {
      if (videoTimerRef.current) clearInterval(videoTimerRef.current);
      if (testTimerRef.current) clearInterval(testTimerRef.current);
      
      // 清理录制资源
      if (cameraStream) {
        cameraStream.getTracks().forEach(track => track.stop());
      }
      if (screenStream) {
        screenStream.getTracks().forEach(track => track.stop());
      }
      if (cameraRecorder && cameraRecorder.state !== 'inactive') {
        cameraRecorder.stop();
      }
      if (screenRecorder && screenRecorder.state !== 'inactive') {
        screenRecorder.stop();
      }
    };
  }, [cameraStream, screenStream, cameraRecorder, screenRecorder]);

  // 渲染视频学习阶段
  if (currentPhase === 'video') {
    const currentVideo = videoProgress[currentVideoIndex];
    const progress = currentVideo ? (currentVideo.watchedTime / currentVideo.duration) * 100 : 0;
    
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
            24点游戏 - 微课学习
          </Typography>

          <Paper className="button-container" elevation={3}>
            <Box mb={3}>
              <Stepper activeStep={currentVideoIndex} alternativeLabel>
                {videos.map((video, index) => (
                  <Step key={video.videoId}>
                    <StepLabel
                      StepIconComponent={() => (
                        <VideoLibrary color={index <= currentVideoIndex ? 'primary' : 'disabled'} />
                      )}
                    >
                      {video.title}
                    </StepLabel>
                  </Step>
                ))}
              </Stepper>
            </Box>

            {currentVideo && (
              <Card elevation={2} sx={{ mb: 3 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    {currentVideo.title}
                  </Typography>
                  
                  {/* 模拟视频播放器 */}
                  <Box 
                    sx={{ 
                      width: '100%', 
                      height: 300, 
                      backgroundColor: '#000', 
                      display: 'flex', 
                      alignItems: 'center', 
                      justifyContent: 'center',
                      mb: 2,
                      borderRadius: 1,
                    }}
                  >
                    <Typography variant="h6" color="white">
                      {isVideoPlaying ? '视频播放中...' : '点击播放按钮开始学习'}
                    </Typography>
                  </Box>

                  <LinearProgress 
                    variant="determinate" 
                    value={progress} 
                    sx={{ mb: 2, height: 8, borderRadius: 4 }}
                  />
                  
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="body2" color="text.secondary">
                      {formatTime(currentVideo.watchedTime)} / {formatTime(currentVideo.duration)}
                    </Typography>
                    <Box>
                      <IconButton 
                        onClick={isVideoPlaying ? handleVideoPause : handleVideoPlay}
                        color="primary"
                        disabled={currentVideo.completed}
                      >
                        {isVideoPlaying ? <Pause /> : <PlayArrow />}
                      </IconButton>
                      <IconButton 
                        onClick={handleNextVideo}
                        color="primary"
                        disabled={!currentVideo.completed && currentVideoIndex < videos.length - 1}
                      >
                        <SkipNext />
                      </IconButton>
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            )}

            <Alert severity="info" sx={{ mb: 2 }}>
              请认真观看所有微课视频，了解24点游戏的规则和技巧。观看完成后将开始正式测试。
            </Alert>

            <Box textAlign="center" mt={3}>
              <Button
                variant="outlined"
                color="secondary"
                onClick={handleSkipAllVideos}
                sx={{ mr: 2 }}
              >
                跳过所有视频，直接开始测试
              </Button>
              <Button
                variant="text"
                color="primary"
                onClick={onBack}
              >
                返回选择
              </Button>
            </Box>
          </Paper>
        </Container>
      </div>
    );
  }

  // 渲染测试阶段
  if (currentPhase === 'test') {
    const currentQuestion = questions[currentQuestionIndex];
    
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
            24点游戏测试
          </Typography>

          <Paper className="button-container" elevation={3}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h6">
                题目 {currentQuestionIndex + 1} / {questions.length}
              </Typography>
              <Box display="flex" alignItems="center" gap={2}>
                {isRecording && (
                  <Box display="flex" alignItems="center" sx={{ color: 'error.main' }}>
                    <FiberManualRecord sx={{ mr: 0.5, fontSize: 16 }} />
                    <Videocam sx={{ mr: 0.5, fontSize: 20 }} />
                    <ScreenShare sx={{ mr: 1, fontSize: 20 }} />
                    <Typography variant="body2" color="error">
                      录制中
                    </Typography>
                  </Box>
                )}
                <Box display="flex" alignItems="center">
                  <Timer sx={{ mr: 1 }} />
                  <Typography variant="h6" color={remainingTime < 300 ? 'error' : 'primary'}>
                    {formatTime(remainingTime)}
                  </Typography>
                </Box>
              </Box>
            </Box>

            {/* 隐藏的摄像头预览元素，用于录制但不显示 */}
            <video
              ref={cameraPreviewRef}
              autoPlay
              muted
              style={{ display: 'none' }}
            />

            <LinearProgress 
              variant="determinate" 
              value={((currentQuestionIndex) / questions.length) * 100} 
              sx={{ mb: 3, height: 8, borderRadius: 4 }}
            />

            {currentQuestion && (
              <Card elevation={2} sx={{ mb: 3 }}>
                <CardContent>
                  <Typography variant="h5" gutterBottom textAlign="center">
                    使用以下四个数字，通过加减乘除运算得到24：
                  </Typography>
                  
                  <Box display="flex" justifyContent="center" gap={2} mb={3}>
                    {currentQuestion.numbers.map((num, index) => (
                      <Chip 
                        key={index}
                        label={num}
                        color="primary"
                        sx={{ fontSize: '1.5rem', padding: '20px 15px', height: '60px' }}
                      />
                    ))}
                  </Box>

                  <TextField
                    fullWidth
                    label="请输入您的答案（如：(1+2+3)*6）"
                    value={currentAnswer}
                    onChange={(e) => setCurrentAnswer(e.target.value)}
                    placeholder="例如：(8-1-1)*8/8+8*3"
                    sx={{ mb: 3 }}
                    onKeyPress={(e) => {
                      if (e.key === 'Enter') {
                        handleSubmitAnswer();
                      }
                    }}
                  />

                  <Box display="flex" justifyContent="center" gap={2}>
                    <Button
                      variant="contained"
                      color="primary"
                      onClick={handleSubmitAnswer}
                      disabled={!currentAnswer.trim()}
                    >
                      提交答案
                    </Button>
                    <Button
                      variant="outlined"
                      color="secondary"
                      onClick={handleSkipQuestion}
                    >
                      跳过此题
                    </Button>
                    {/* 调试按钮 */}
                    {/* <Button
                      variant="outlined"
                      size="small"
                      onClick={() => {
                        console.log('🔍 录制状态调试信息:');
                        console.log('摄像头流:', cameraStream ? '✅存在' : '❌不存在');
                        console.log('屏幕流:', screenStream ? '✅存在' : '❌不存在');
                        console.log('摄像头录制器:', cameraRecorder?.state || '❌未创建');
                        console.log('屏幕录制器:', screenRecorder?.state || '❌未创建');
                        console.log('摄像头数据块数:', cameraChunksRef.current.length);
                        console.log('屏幕数据块数:', screenChunksRef.current.length);
                        console.log('录制状态:', isRecording ? '✅录制中' : '❌未录制');
                      }}
                    >
                      🔍 调试录制状态
                    </Button> */}
                  </Box>
                </CardContent>
              </Card>
            )}

            <Alert severity="warning" sx={{ mb: 2 }}>
              提示：每个数字必须且只能使用一次，可以使用括号改变运算顺序。
            </Alert>

            {isRecording && (
              <Alert severity="success" sx={{ mb: 2 }}>
                <Typography variant="body2">
                  📹 正在录制您的答题过程...
                </Typography>
              </Alert>
            )}
          </Paper>
        </Container>
      </div>
    );
  }

  // 渲染结果阶段
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
          24点游戏 - 测试结果
        </Typography>

        <Paper className="button-container" elevation={3}>
          <Typography variant="h5" textAlign="center" mb={2}>
            测试完成！
          </Typography>
          
          <Alert severity="success" sx={{ mb: 4 }}>
            <Typography variant="body1">
              🎥 录制已完成！摄像头和屏幕录制视频已自动上传到服务器。
            </Typography>
            <Typography variant="body2" color="text.secondary" mt={1}>
              视频文件已安全保存在服务器端，用于后续的数据分析和研究。
            </Typography>
          </Alert>

          <Grid container spacing={3} mb={4}>
            <Grid item xs={12} md={3}>
              <Card elevation={2}>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h4" color="primary">
                    {totalScore}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    总得分 / 10
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card elevation={2}>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h6" color="primary">
                    {formatTime(totalTestTime)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    答题用时
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card elevation={2}>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h6" color="primary">
                    {formatTime(totalVideoTime)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    视频学习用时
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card elevation={2}>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h6" color="primary">
                    {questions.reduce((sum, q) => sum + q.attempts, 0)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    总尝试次数
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>

          <Typography variant="h6" mb={2}>题目详情：</Typography>
          {questions.map((question, index) => (
            <Card key={question.id} elevation={1} sx={{ mb: 2 }}>
              <CardContent>
                <Grid container alignItems="center" spacing={2}>
                  <Grid item xs={12} sm={3}>
                    <Typography variant="body1">
                      题目 {index + 1}: {question.numbers.join(', ')}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={3}>
                    <Typography variant="body2" color="text.secondary">
                      {question.skipped ? '已跳过' : `答案: ${question.answer || '未作答'}`}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={2}>
                    <Typography variant="body2" color="text.secondary">
                      用时: {formatTime(question.timeSpent)}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={2}>
                    <Typography variant="body2" color="text.secondary">
                      尝试: {question.attempts} 次
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={2}>
                    <Chip 
                      label={question.skipped ? '跳过' : (question.answer && evaluateExpression(question.answer, question.numbers) ? '正确' : '错误')}
                      color={question.skipped ? 'default' : (question.answer && evaluateExpression(question.answer, question.numbers) ? 'success' : 'error')}
                      size="small"
                    />
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          ))}

          <Box textAlign="center" mt={4}>
            <Button
              variant="contained"
              color="primary"
              onClick={onBack}
              size="large"
            >
              返回测试选择
            </Button>
          </Box>
        </Paper>
      </Container>
    </div>
  );
};

export default TwentyFourGame;
