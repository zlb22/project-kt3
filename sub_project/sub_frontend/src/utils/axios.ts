import axios from 'axios'
import type { AxiosInstance } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import md5 from 'md5'
import { getUrlParam } from '@/utils/url'

export class Interceptors {
  instance: AxiosInstance
  constructor() {
    this.instance = axios.create({
      baseURL: '',
      timeout: 10000
    })
    this.init()
  }

  init() {
    this.instance.interceptors.request.use(
      (config) => {
        const authToken = getUrlParam('token') 
        if(authToken) {
          window.localStorage.setItem('AuthToken', authToken) 
        }
        const token = window.localStorage.getItem('AuthToken') || ''
        config.headers['Authorization'] = token
        const body = JSON.stringify(config.data)
        const header = getHeader(body)
        Object.assign(config.headers, {
          'X-Auth-AppId': header.appId,
          'X-Auth-TimeStamp': header.timeStamp,
          'X-Sign': header.sign
        })
        return config
      },
      (err) => {
        console.error('Request interceptor error:', err)
        return Promise.reject(err)
      }
    )
    this.instance.interceptors.response.use(
      (response) => {
        const { status } = response
        if (status === 200) {
          if (response.data.errcode === 60000) {
            ElMessage.error(response.data.errmsg)
           }
          if (response.data.errcode === 20005) {
           ElMessage.error('操作权限不足')
          }
          if (response.data.errcode === 20004) {
            setTimeout(() => {
              window.localStorage.removeItem('AuthToken')
              const currentUrl = window.location.href
              // 跳转到统一登录平台
              window.location.href = `${import.meta.env.VITE_LOGIN_URL}?redirectUrl=${encodeURIComponent(currentUrl)}`
            }, 1000)
          }
          return Promise.resolve(response)
        } else {
          return Promise.reject(response)
        }
      },
      (err) => {
        console.error('Response interceptor error:', err)
        return Promise.reject(err)
      }
    )
  }

  getInterceptors() {
    return this.instance
  }
}

// 创建单例实例
const apiInstance = new Interceptors().getInterceptors()

// 导出单例实例
export { apiInstance }

function getHeader(body: any) {
  const header: any = {
    timeStamp: Math.floor(Date.now() / 1000),
    appId: import.meta.env.VITE_APP_ID
  }
  // 截取参数一部分参与运算:如果请求体参数的长度超过100字节，则只取前100字节的内容；否则，取完整内容
  /*
  if (body?.length > 100) {
    body = body.substring(0, 100)
  }*/
  // 将密钥、时间戳、截取后的body参数按如下顺序组装成字符串
  const str = `${import.meta.env.VITE_SKEY}${header.timeStamp}`
  const sign = md5(str)
  header.sign = sign
  return header
}