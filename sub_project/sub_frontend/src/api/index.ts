import { apiInstance } from '@/utils/axios'

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
  return apiInstance.post('/********', body)
}



