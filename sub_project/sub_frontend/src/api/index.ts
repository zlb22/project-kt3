import { apiInstance } from '@/utils/axios'
import type { AxiosProgressEvent } from 'axios'

/**
 * 获取配置
 * @param body 
 * @returns 
 */
export const getConfig = (body?: any) => {
  return apiInstance.post('/web/keti3/config/list', body)
}

/**
 * 保存数据
 * @param body 
 * @returns 
 */
export const saveData = (body?: any) => {
  return apiInstance.post('/web/keti3/log/save', body)
}

/**
 * 获取OSS认证
 * @param body 
 * @returns 
 */
export const getOssAuth = (body?: any) => {
  return apiInstance.post('/web/keti3/oss/auth', body)
}

/**
 * 通过后端代理上传（避免浏览器直传 MinIO 的 CORS 问题）
 */
export const uploadToServer = (form: FormData, onUploadProgress?: (e: AxiosProgressEvent) => void) => {
  // Do NOT set Content-Type manually; let Axios add the correct boundary
  // Kept for backward compatibility without uid; prefer using the overload with uid
  return apiInstance.post('/web/keti3/oss/upload', form, { onUploadProgress })
}

// Overload with uid: send X-User-Id so backend can partition by user
export const uploadToServerWithUid = (form: FormData, uid: number | string, onUploadProgress?: (e: AxiosProgressEvent) => void) => {
  const userId = String(uid ?? '')
  return apiInstance.post('/web/keti3/oss/upload', form, {
    headers: { 'X-User-Id': userId },
    onUploadProgress
  })
}



