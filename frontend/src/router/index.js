import { createRouter, createWebHistory } from 'vue-router'
import { useCookies } from "vue3-cookies"
import { useToast } from "primevue/usetoast"
import { useSystemStore } from '@/stores/system'

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
      {
         path: '/',
         name: 'browse',
         component: () => import('../views/BrowseSubmissions.vue'),
      },
       {
         path: '/submissions/:id',
         name: 'details',
         component: () => import('../views/SubmissionDetail.vue'),
      },
      {
         path: '/signedout',
         name: "signedout",
         component: () => import('../views/SignedOut.vue'),
      },
      {
         path: '/forbidden',
         name: "forbidden",
         component: () => import('../views/Forbidden.vue'),
      },
      {
         path: '/:pathMatch(.*)*',
         name: "not_found",
         component: () => import('../views/NotFound.vue'),
      }
   ],
})

router.beforeEach(async (to) => {
   console.log("BEFORE ROUTE " + to.path)
   const system = useSystemStore()
   const { cookies } = useCookies()
   const noAuthRoutes = ["not_found", "forbidden", "expired", "signedout"]
   useToast().removeAllGroups()

   // the /signedin endpoint called after authorization. it has no page itself; it just
   // processes the authorization response and redirects to the next page (or forbidden)
   if (to.path == '/signedin') {
      const jwtStr = cookies.get("aptsubmit-client")
      console.log(`GRANTED [${jwtStr}]`)
      system.setUserJWT(jwtStr)
      return "/"
   }

   // for all other routes, pull the existing jwt from storage from storage and set in the  system.
   // depending upon the page requested, this token may or may not be used.
   const jwtStr = localStorage.getItem('aptsubmit-client')
   system.setUserJWT(jwtStr)

   if ( noAuthRoutes.includes(to.name)) {
      console.log("NOT A PROTECTED PAGE")
   } else {
      if ( system.isSignedIn == false) {
         system.authenticate()
         return false   // cancel the original navigation
      }
   }
})

export default router
