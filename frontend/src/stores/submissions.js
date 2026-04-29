import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useSubmissionsStore = defineStore('submission', {
   state: () => ({
      working: false,
      offset: 0,
      query: "",
      filters: [],
      includeAutoApproved: false,
      pageSize: 30,
      total: 0,
      searchHits: [],
      detail: null
   }),
   getters: {
      filtersAsQueryParam: state => {
         let out = []
         state.filters.forEach( fv => out.push(`${fv.field}=${fv.value}`) )
         return JSON.stringify(out)
      }
   },
   actions: {
      getSubmissions() {
         this.working = true
         this.detail = null
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
         let filterParam = this.filtersAsQueryParam
         if ( filterParam != "") {
            params.push(`filters=${filterParam}`)    
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
         this.detail = null
         this.getSubmissions()
      },
      getSubmissionDetail( identifier ) {
         console.log("load details for "+identifier)
         this.working = true
         axios.get("/api/submissions/"+identifier).then(response => {
            this.detail = response.data
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.error = e
            this.working = false
         })
      }
   }
})
