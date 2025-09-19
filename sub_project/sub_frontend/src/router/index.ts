import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// Global auth guard: require token either in query (?token=) or in localStorage (AuthToken)
function getJwtExp(token: string | null | undefined): number | null {
  try {
    if (!token) return null
    const parts = token.split('.')
    if (parts.length < 2) return null
    const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))
    const exp = typeof payload?.exp === 'number' ? payload.exp : null
    return exp
  } catch {
    return null
  }
}

router.beforeEach((to, _from, next) => {
  // Token in query
  const q = to.query || {}
  let queryToken: string | undefined
  const raw = q.token as any
  if (typeof raw === 'string') queryToken = raw
  else if (Array.isArray(raw) && raw.length > 0) queryToken = raw[0]

  // Validate JWT exp (in seconds) for query token only (strict: must come via URL)
  const now = Math.floor(Date.now() / 1000)
  const queryExp = getJwtExp(queryToken)

  // Allow only when query token exists and is not expired (or has no exp)
  if (queryToken && (!queryExp || queryExp > now)) {
    next()
    return
  }

  // No valid token: redirect to unified login with return URL
  const loginUrl = (import.meta as any).env?.VITE_LOGIN_URL
  if (loginUrl) {
    const target = window.location.origin + to.fullPath
    const redirectUrl = encodeURIComponent(target)
    window.location.href = `${loginUrl}?redirectUrl=${redirectUrl}`
    return
  }

  // No token: redirect to unified login with return URL
  // Fallback: stay on page (avoid infinite loop) if no login URL configured
  next()
})

export default router
