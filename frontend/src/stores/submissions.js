import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useSubmissionsStore = defineStore('submission', {
   state: () => ({
      working: false,
      offset: 0,
      pageSize: 30,
      total: 0,
      searchHits: []
   }),
   getters: {
   },
   actions: {
      getSubmissions() {
         var q = `/api/submissions?start=${this.offset}&limit=${this.pageSize}`
         axios.get(q).then(response => {
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
         this.total = 0
      }
   }
})
