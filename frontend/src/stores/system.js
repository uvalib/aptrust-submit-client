import { defineStore } from 'pinia'
import axios from 'axios'

export const useSystemStore = defineStore('system', {
   state: () => ({
      working: false,
		version: "unknown",
      error: "",
      showError: false,
      toast: {
         error: false,
         summary: "",
         message: "",
         show: false
      }
   }),
   getters: {
   },
   actions: {
      async getConfig() {
         this.working = false
         return axios.get("/config").then(response => {
            this.version = response.data.version
            this.working = false
         }).catch( err => {
            this.setError(  err )
         })
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
