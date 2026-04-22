import { createRouter, createWebHistory } from 'vue-router'
import { useToast } from "primevue/usetoast"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
         path: '/',
         alias: '/dashboard',
         name: 'dashboard',
         component: () => import('../views/DashboardView.vue'),
      },
  ],
})

router.beforeEach( async (to) => {
   console.log("BEFORE ROUTE "+to.path)
    useToast().removeAllGroups()
})

export default router
