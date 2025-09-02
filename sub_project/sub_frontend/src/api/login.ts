import { apiInstance } from '@/utils/axios'


export const login = (body?: any) => {
  return apiInstance.post('/web/keti3/student/login', body)
}

