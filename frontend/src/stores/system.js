import { defineStore } from 'pinia'
import axios from 'axios'
import { useJwt } from '@vueuse/integrations/useJwt'

export const useSystemStore = defineStore('system', {
   state: () => ({
      authorizing: false,
      signedInUser: "",
      canApproveSubmissions: false,
      userJWT: "",
      requestInterceptor: null,
      responseInterceptor: null,
      working: false,
		version: "unknown",
      error: "",
      showError: false,
      submissionStatuses: [],
      toast: {
         error: false,
         summary: "",
         message: "",
         show: false
      }
   }),
   getters: {
      isSignedIn: state => {
         return state.signedInUser != ""
      },
   },
   actions: {
      async getConfig() {
         this.working = false
         return axios.get("/config").then(response => {
            this.version = response.data.version
            this.submissionStatuses = response.data.submissionStatuses
            this.working = false
         }).catch( err => {
            this.setError(  err )
         })
      },

      authenticate() {
         console.log("AUTHENTICATE")
         this.authorizing = true
         window.location.href = "/authenticate"
      },

      setUserJWT( jwt ) {
         if (jwt == this.jwt || jwt == "" || jwt == null || jwt == "null")  return

         this.authorizing = false
         this.userJWT = jwt
         localStorage.setItem("aptsubmit-client", jwt)

         const { payload } = useJwt(jwt)
         console.log(payload.value)
         this.signedInUser = payload.value.computeID
         this.canApproveSubmissions = payload.value.canApprove
         console.log(`user ${this.signedInUser} signed in`)

         // add interceptor to include jwt bearer
         this.requestInterceptor = axios.interceptors.request.use(config => {
            config.headers['Authorization'] = 'Bearer ' + this.userJWT
            return config
         }, error => {
            return Promise.reject(error)
         })

         // intercept failed responses
         this.responseInterceptor = axios.interceptors.response.use(
            res => res,
            err => {
               console.log(`request ${err.config.url} failed with status ${err.response.status}`)
               console.log(err)
               if (err.config.url.match(/\/authenticate/)) {
                  this.router.push("/forbidden")
               } else {
                  if (err.response && err.response.status == 401) {
                  console.log("REQUEST FAILED WITH 401")
                     this.signOut()
                     system.working = false
                     // this.authenticate() // maybe reauth? prefer just to expire tho
                     return new Promise(() => { })
                  }
               }
               return Promise.reject(err)
            }
         )
      },

      signOut() {
         console.log("SIGNOUT USER")
      
         if ( this.requestInterceptor != null ) {
            console.log("remove existing request interceptor")
            axios.interceptors.request.eject( this.requestInterceptor)
            this.requestInterceptor = null
         }
         if ( this.responseInterceptor != null ) {
            console.log("remove existing response interceptor")
            axios.interceptors.response.eject( this.responseInterceptor)
            this.responseInterceptor = null
         }

         localStorage.removeItem("aptsubmit-client")
         this.$reset()
         this.router.push("/signedout")
      },

      setError( e ) {
         this.error = e
         if (e.response && e.response.data) {
            this.error = e.response.data
         }
         this.showError = true
         this.working = false

         if (e.status && e.status != 406) {
            // 406 is returned on jwt mismatch. do not report this as an error
            this.reportError(e)
         }
      },
      async reportError(data) {
         let ipResp = await fetch("https://api.ipify.org")
         let pubIP = await ipResp.text()
         let err = {
            url: this.router.currentRoute.value.fullPath,
            userAgent: navigator.userAgent,
            error: JSON.stringify(data),
            publicIP: pubIP
         }
         if (err.error == "{}" ) {
            err.error = data.toString()
         }

         // dont report network errors!
         if ( err.error.includes("System error, we regret the inconvenience") ||
              err.error.includes("Network Error") ||
              err.error.includes("status code 401") ||
              err.error.includes("ECONNREFUSED") ) {
            return
         }

         const user = useUserStore()
         if (user.isSignedIn) {
            err.signedIn = true
            err.user = user.signedInUser
         } else {
            err.signedIn = false
         }
         axios.post("/api/error", err)
      },

      toastMessage( summary, message ) {
         this.toast.summary = summary
         this.toast.message = message
         this.toast.show = true
         this.toast.error = false
      },
      toastError( summary, message ) {
         this.toast.summary = summary
         this.toast.message = message
         this.toast.show = true
         this.toast.error = true
         this.reportError(message)
      },
      clearToastMessage() {
         this.toast.summary = ""
         this.toast.message = ""
         this.toast.show = false
         this.toast.error = false
      },
   }
})
