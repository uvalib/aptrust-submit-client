import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useSubmissionsStore = defineStore('submission', {
   state: () => ({
      working: false,
      loadingBags: false,
      offset: 0,
      query: "",
      filters: [],
      includeAutoApproved: false,
      pageSize: 30,
      sortField: "createdAt",
      sortOrder: "desc",
      total: 0,
      searchHits: [],
      detail: null,
      bags: null
   }),
   getters: {
      filtersAsQueryParam: state => {
         let out = []
         state.filters.forEach( fv => out.push(`${fv.field}=${fv.value}`) )
         return JSON.stringify(out)
      },
      currentStatus: state => {
         if ( !state.detail.status ) return "Unknown"
         return state.detail.status[0].status
      },
      sortOrderInt: state => {
         if (state.sortOrder == "desc") return -1
         return 1
      }
   },
   actions: {
      getSubmissions() {
         this.working = true
         this.detail = null
         this.bags = null
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
         params.push(`sort=${this.sortField}`)
         params.push(`order=${this.sortOrder}`)

         url += params.join("&")
         console.log(url)
         axios.get(url).then(response => {
            this.total = response.data.total
            this.searchHits = response.data.hits
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError( e )
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
         this.working = true
         this.bags = null
         axios.get("/api/submissions/"+identifier).then(response => {
            this.detail = response.data
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError( e )
            this.working = false
         })
      },
      getBags( identifier ) {  
         this.loadingBags = true
         axios.get("/api/submissions/"+identifier+"/bags").then(response => {
            // convert db structure into tree model:
            // key, label, data, children 
            this.bags = []
            response.data.forEach( bag => {
               let bagData = { createdAt: bag.createdAt, status: bag.status }
               let node = { key: bag.id, label: bag.name, data: bagData, children: [], type: "bag" }
               if ( bag.files ) {
                  bag.files.forEach( bf => {
                     node.children.push( {key: bf.id, label: bf.fileName, data: bf, type: "file"})
                  })
               } else if (bag.aptFiles) {
                  bag.aptFiles.forEach( bf => {
                     node.children.push( {key: bf.id, label: bf.fileName, data: bf, type: "file"})
                  })   
               }
               this.bags.push(node)
            })
            this.loadingBags = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError( e )
            this.loadingBags = false
         })
      },
      cancel() {
         axios.post("/api/submissions/"+ this.detail.identifier+"/cancel").then(() => {
            this.detail.status.unshift({createdAt: new Date(), status: "abandoned"})
         }).catch( e => {
            const system = useSystemStore()
            system.setError( e )
         })
      },
      approve( storage ) {
         axios.post("/api/submissions/"+this.detail.identifier+"/approve", {storage: storage}).then(() => {
            this.detail.status.unshift({createdAt: new Date(), status: "submitting"})
         }).catch( e => {
            const system = useSystemStore()
            system.setError( e )
         })
      }
   }
})
