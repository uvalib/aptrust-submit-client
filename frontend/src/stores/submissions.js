import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useSubmissionsStore = defineStore('submission', {
   state: () => ({
      working: false,
      offset: 0,
      query: "",
      includeAutoApproved: false,
      pageSize: 30,
      total: 0,
      searchHits: []
   }),
   getters: {
   },
   actions: {
      getSubmissions() {
         var url = `/api/submissions?`
         var params = [] 
         if ( this.query != "" ) {
            params.push(`q=${this.query}`)
         }
         if ( this.includeAutoApproved ) {
            params.push("includeauto=1")
         } else {
            params.push("includeauto=0")    
         }
         params.push(`start=${this.offset}`)
         params.push(`limit=${this.pageSize}`)

         url += params.join("&")
         axios.get(url).then(response => {
            this.total = response.data.total
            this.searchHits = response.data.hits
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.error = e
            this.working = false
         })

      },
      resetSearch() {
         this.page = 1
         this.pageSize = 30 
         this.searchHits = []
         this.includeAutoApproved = false 
         this.query = ""
         this.total = 0
         this.getSubmissions()
      }
   }
})
